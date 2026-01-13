package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/princetheprogrammerbtw/gitsynq/internal/ui"
	"github.com/spf13/cobra"
)

var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "⚓ Manage sync hooks",
	Long:  `Configure scripts to run before or after synchronization operations.`,
	Run:   runHooks,
}

func init() {
	// Hooks don't have subcommands yet, just a status check
}

func runHooks(cmd *cobra.Command, args []string) {
	printBanner()
	ui.Cyan.Println("\n⚓ Sync Hooks Status")

	hooksDir := ".gitsync-hooks"
	os.MkdirAll(hooksDir, 0755)

	hooks := []string{"pre-push", "post-push", "pre-pull", "post-pull"}

	fmt.Printf("\n%-15s %-10s %-30s\n", "HOOK", "STATUS", "PATH")
	fmt.Println(strings.Repeat("-", 60))

	for _, hook := range hooks {
		path := filepath.Join(hooksDir, hook)
		status := "∅ Not set"
		if _, err := os.Stat(path); err == nil {
			status = "✅ Active"
		}
		fmt.Printf("%-15s %-10s %-30s\n", hook, status, path)
	}

	ui.Yellow.Println("\n💡 To create a hook, simply create an executable script in .gitsync-hooks/")
	fmt.Println("   Example: echo '#!/bin/bash\necho \"Syncing!\"' > .gitsync-hooks/pre-push")
	fmt.Println("   chmod +x .gitsync-hooks/pre-push")
}

func executeHook(name string) error {
	path := filepath.Join(".gitsync-hooks", name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Hook doesn't exist, which is fine
	}

	ui.Cyan.Printf("⚓ Running hook: %s...\n", name)
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
