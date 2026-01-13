package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/schollz/progressbar/v3"
	"github.com/princetheprogrammerbtw/gitsynq/internal/bundle"
	"github.com/princetheprogrammerbtw/gitsynq/internal/config"
	"github.com/princetheprogrammerbtw/gitsynq/internal/ssh"
	"github.com/princetheprogrammerbtw/gitsynq/internal/ui"
	"github.com/princetheprogrammerbtw/gitsynq/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	fullPush   bool
	includeAll bool
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "📤 Push repository to remote server",
	Long: `Create a Git bundle and transfer it to the remote server.
	
Examples:
  gitsync push           # Push new commits only
  gitsync push --full    # Push entire repository
  gitsync push --all     # Include all branches`,
	Run: runPush,
}

func init() {
	pushCmd.Flags().BoolVarP(&fullPush, "full", "f", false, "Push entire repository (not just new commits)")
	pushCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "Include all branches")
}

func runPush(cmd *cobra.Command, args []string) {
	printBanner()
	ui.Green.Println("\n📤 Pushing to Remote Server")

	// Run pre-push hook
	if err := executeHook("pre-push"); err != nil {
		ui.Red.Printf("❌ pre-push hook failed: %v\n", err)
		os.Exit(1)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		ui.Red.Printf("❌ Error loading config: %v\n", err)
		ui.Yellow.Println("💡 Run 'gitsync init' first!")
		os.Exit(1)
	}

	// Start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	// Step 1: Create bundle
	s.Suffix = " Creating Git bundle..."
	s.Start()

	timestamp := time.Now().Format("20060102-150405")
	bundleName := fmt.Sprintf("%s-%s.bundle", cfg.Project.Name, timestamp)
	bundlePath := filepath.Join(cfg.Bundle.Directory, bundleName)

	var bundleErr error
	if fullPush {
		bundleErr = bundle.CreateFull(bundlePath)
	} else {
		bundleErr = bundle.CreateIncremental(bundlePath, cfg.Project.Branch)
	}

	s.Stop()

	if bundleErr != nil {
		// If incremental fails, try full
		if !fullPush {
			ui.Yellow.Println("⚠️  Incremental push failed. Attempting full bundle...")
			s.Suffix = " Creating full bundle..."
			s.Start()
			bundleErr = bundle.CreateFull(bundlePath)
			s.Stop()
		}

		if bundleErr != nil {
			ui.Red.Printf("❌ Error creating bundle: %v\n", bundleErr)
			os.Exit(1)
		}
	}

	ui.Green.Println("✅ Bundle created:", bundleName)

	// Get bundle size
	info, _ := os.Stat(bundlePath)
	ui.Cyan.Printf("📦 Bundle size: %s\n", utils.FormatBytes(info.Size()))

	// Step 2: Transfer to server
	s.Suffix = fmt.Sprintf(" Transferring to %s@%s...", cfg.Server.User, cfg.Server.Host)
	s.Start()

	client, err := ssh.NewClient(cfg.Server)
	if err != nil {
		s.Stop()
		ui.Red.Printf("❌ SSH connection failed: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	remoteBundlePath := filepath.Join(cfg.Server.RemotePath, bundleName)
	
	bar := progressbar.DefaultBytes(
		info.Size(),
		"🚀 Uploading bundle",
	)

	err = client.Upload(bundlePath, remoteBundlePath, func(current, total int64) {
		bar.Set64(current)
	})

	if err != nil {
		ui.Red.Printf("❌ Upload failed: %v\n", err)
		os.Exit(1)
	}

	ui.Green.Println("\n✅ Bundle transferred successfully!")

	// Step 3: Setup/Update repo on server
	s.Suffix = " Setting up repository on server..."
	s.Start()

	remoteRepoPath := filepath.Join(cfg.Server.RemotePath, cfg.Project.Name)
	setupScript := generateSetupScript(remoteBundlePath, remoteRepoPath, cfg.Project.Branch)

	output, err := client.Run(setupScript)
	s.Stop()

	if err != nil {
		ui.Red.Printf("❌ Remote setup failed: %v\n", err)
		if verbose {
			fmt.Println("Output:", output)
		}
		os.Exit(1)
	}

	// Success!
	printPushSuccess(cfg, bundleName)

	// Run post-push hook
	_ = executeHook("post-push")
}

func generateSetupScript(bundlePath, repoPath, branch string) string {
	return fmt.Sprintf(`
		set -e
		
		BUNDLE_PATH="%s"
		REPO_PATH="%s"
		BRANCH="%s"
		
		if [ ! -d "$REPO_PATH/.git" ]; then
			echo "📂 Cloning from bundle..."
			git clone "$BUNDLE_PATH" "$REPO_PATH"
			cd "$REPO_PATH"
			git checkout "$BRANCH" 2>/dev/null || git checkout -b "$BRANCH"
		else
			echo "🔄 Updating existing repository..."
			cd "$REPO_PATH"
			
			# Add bundle as remote temporarily
			git remote add bundle "$BUNDLE_PATH" 2>/dev/null || git remote set-url bundle "$BUNDLE_PATH"
			
			# Fetch and merge
			git fetch bundle
			git merge bundle/"$BRANCH" --no-edit 2>/dev/null || git merge bundle/master --no-edit 2>/dev/null || true
			
			# Cleanup
			git remote remove bundle 2>/dev/null || true
		fi
	`, bundlePath, repoPath, branch)
}

func printPushSuccess(cfg *config.Config, bundleName string) {
	ui.Green.Println("\n" + strings.Repeat("═", 50))
	ui.Green.Println("          🎉 PUSH SUCCESSFUL! 🎉")
	ui.Green.Println(strings.Repeat("═", 50))

	ui.Cyan.Printf(`
📦 Bundle:    %s
🖥️  Server:    %s@%s
📂 Path:      %s/%s

`, bundleName, cfg.Server.User, cfg.Server.Host, cfg.Server.RemotePath, cfg.Project.Name)

	ui.Yellow.Println("🔜 Next steps on server:")
	fmt.Printf("   ssh %s@%s\n", cfg.Server.User, cfg.Server.Host)
	fmt.Printf("   cd %s/%s\n", cfg.Server.RemotePath, cfg.Project.Name)
	fmt.Println("   # Start coding! 🚀")
}