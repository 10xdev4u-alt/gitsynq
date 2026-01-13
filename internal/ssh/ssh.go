// Package ssh provides a high-level wrapper around SSH and SFTP for file transfers
// and remote command execution.
package ssh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"github.com/princetheprogrammerbtw/gitsynq/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Client is a wrapper around ssh.Client and sftp.Client.
type Client struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

// NewClient creates and connects a new SSH and SFTP client using the provided configuration.
// It tries multiple authentication methods: explicitly provided SSH key, default SSH keys
// (~/.ssh/id_rsa, etc.), and SSH agent.
func NewClient(cfg config.ServerConfig) (*Client, error) {
	var authMethods []ssh.AuthMethod

	// Method 1: Explicit SSH key from config
	if cfg.SSHKeyPath != "" {
		key, err := os.ReadFile(cfg.SSHKeyPath)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
	}

	// Method 2: Default SSH keys in home directory
	homeDir, _ := os.UserHomeDir()
	keyPaths := []string{
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_ecdsa"),
	}

	for _, keyPath := range keyPaths {
		key, err := os.ReadFile(keyPath)
		if err != nil {
			continue
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			continue
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// Method 3: SSH Agent
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		authMethods = append(authMethods, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no SSH authentication methods found (checked config, defaults, and agent)")
	}

	sshConfig := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Support host key verification
		Timeout:         30 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("failed to initialize SFTP: %w", err)
	}

	return &Client{
		sshClient:  sshClient,
		sftpClient: sftpClient,
	}, nil
}

// Close closes both SSH and SFTP connections gracefully.
func (c *Client) Close() error {
	var errs []error
	if c.sftpClient != nil {
		if err := c.sftpClient.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if c.sshClient != nil {
		if err := c.sshClient.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing connection: %v", errs)
	}
	return nil
}

// Run executes a command on the remote server and returns its combined stdout and stderr.
// It supports cancellation via the provided context.
func (c *Client) Run(ctx context.Context, command string) (string, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Handle cancellation
	done := make(chan error, 1)
	go func() {
		done <- session.Run(command)
	}()

	select {
	case <-ctx.Done():
		session.Signal(ssh.SIGKILL)
		return "", ctx.Err()
	case err := <-done:
		if err != nil {
			return stdout.String() + stderr.String(), fmt.Errorf("command failed: %w", err)
		}
	}

	return stdout.String(), nil
}

// ProgressFunc is a callback for reporting transfer progress.
type ProgressFunc func(current, total int64)

// Upload transfers a local file to the remote server via SFTP with optional progress reporting.
func (c *Client) Upload(localPath, remotePath string, onProgress ProgressFunc) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()

	info, err := localFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat local file: %w", err)
	}
	total := info.Size()

	// Ensure remote directory exists
	remoteDir := filepath.Dir(remotePath)
	if err := c.sftpClient.MkdirAll(remoteDir); err != nil {
		return fmt.Errorf("failed to create remote directory %s: %w", remoteDir, err)
	}

	remoteFile, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer remoteFile.Close()

	if onProgress == nil {
		_, err = io.Copy(remoteFile, localFile)
	} else {
		pw := &progressWriter{
			writer:     remoteFile,
			total:      total,
			onProgress: onProgress,
		}
		_, err = io.Copy(pw, localFile)
	}

	if err != nil {
		return fmt.Errorf("failed to upload data: %w", err)
	}

	return nil
}

// Download transfers a remote file to the local machine via SFTP with optional progress reporting.
func (c *Client) Download(remotePath, localPath string, onProgress ProgressFunc) error {
	remoteFile, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer remoteFile.Close()

	info, err := remoteFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat remote file: %w", err)
	}
	total := info.Size()

	// Ensure local directory exists
	localDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create local directory %s: %w", localDir, err)
	}

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	if onProgress == nil {
		_, err = io.Copy(localFile, remoteFile)
	} else {
		pw := &progressWriter{
			writer:     localFile,
			total:      total,
			onProgress: onProgress,
		}
		_, err = io.Copy(pw, remoteFile)
	}

	if err != nil {
		return fmt.Errorf("failed to download data: %w", err)
	}

	return nil
}

type progressWriter struct {
	writer     io.Writer
	current    int64
	total      int64
	onProgress ProgressFunc
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	pw.current += int64(n)
	pw.onProgress(pw.current, pw.total)
	return n, err
}
