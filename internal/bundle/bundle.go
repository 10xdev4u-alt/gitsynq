package bundle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateFull creates a bundle with entire repository
func CreateFull(outputPath string) error {
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	cmd := exec.Command("git", "bundle", "create", outputPath, "--all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

// CreateIncremental creates a bundle with only new commits
func CreateIncremental(outputPath, branch string) error {
	// First, check if we have a remote tracking branch
	cmd := exec.Command("git", "rev-parse", "--verify", "origin/"+branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("no tracking branch found, use --full")
	}

	bundleCmd := exec.Command("git", "bundle", "create", outputPath,
		fmt.Sprintf("origin/%s..%s", branch, branch))

	output, err := bundleCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	if info.Size() < 100 {
		os.Remove(outputPath)
		return fmt.Errorf("no new commits to bundle")
	}

	return nil
}

// Merge merges a bundle into the current repository
func Merge(bundlePath, branch string) error {
	verifyCmd := exec.Command("git", "bundle", "verify", bundlePath)
	if err := verifyCmd.Run(); err != nil {
		return fmt.Errorf("invalid bundle: %v", err)
	}

	exec.Command("git", "remote", "remove", "bundle").Run()
	addCmd := exec.Command("git", "remote", "add", "bundle", bundlePath)
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add bundle remote: %v", err)
	}

	fetchCmd := exec.Command("git", "fetch", "bundle")
	if output, err := fetchCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("fetch failed: %s", string(output))
	}

	mergeCmd := exec.Command("git", "merge", "bundle/"+branch, "--no-edit")
	if output, err := mergeCmd.CombinedOutput(); err != nil {
		mergeCmd = exec.Command("git", "merge", "bundle/master", "--no-edit")
		if output2, err2 := mergeCmd.CombinedOutput(); err2 != nil {
			return fmt.Errorf("merge failed: %s %s", string(output), string(output2))
		}
	}

	exec.Command("git", "remote", "remove", "bundle").Run()

	return nil
}

// PushToOrigin pushes to the origin remote
func PushToOrigin(branch string) error {
	cmd := exec.Command("git", "push", "origin", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

// ShowRecentCommits displays recent commits
func ShowRecentCommits(n int) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", n), "--oneline", "--graph")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
