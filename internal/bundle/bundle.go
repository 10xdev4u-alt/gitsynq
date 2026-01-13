// Package bundle handles Git bundle operations like creation, verification, and merging.
package bundle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateFull creates a Git bundle containing the entire repository history.
// It includes all branches and tags.
func CreateFull(outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create bundle directory: %w", err)
	}

	cmd := exec.Command("git", "bundle", "create", outputPath, "--all")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git bundle create failed: %v: %s", err, string(output))
	}
	return nil
}

// CreateIncremental creates a Git bundle containing only the commits that exist on the
// specified branch but not on its remote tracking counterpart (origin/branch).
func CreateIncremental(outputPath, branch string) error {
	// First, check if we have a remote tracking branch
	checkCmd := exec.Command("git", "rev-parse", "--verify", "origin/"+branch)
	if err := checkCmd.Run(); err != nil {
		return fmt.Errorf("no tracking branch found for %s, a full push is required", branch)
	}

	bundleCmd := exec.Command("git", "bundle", "create", outputPath,
		fmt.Sprintf("origin/%s..%s", branch, branch))

	if output, err := bundleCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git bundle incremental create failed: %v: %s", err, string(output))
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("failed to stat created bundle: %w", err)
	}

	// Minimal bundle size check (empty bundles are around 100 bytes)
	if info.Size() < 100 {
		os.Remove(outputPath)
		return fmt.Errorf("no new commits found to bundle")
	}

	return nil
}

// Merge takes a path to a Git bundle and merges the changes into the local repository
// on the specified branch.
func Merge(bundlePath, branch string) error {
	// Verify bundle
	verifyCmd := exec.Command("git", "bundle", "verify", bundlePath)
	if output, err := verifyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("invalid or incompatible bundle: %v: %s", err, string(output))
	}

	// Ensure any old bundle remote is removed
	_ = exec.Command("git", "remote", "remove", "bundle").Run()

	addCmd := exec.Command("git", "remote", "add", "bundle", bundlePath)
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add bundle as remote: %w", err)
	}
	defer func() {
		_ = exec.Command("git", "remote", "remove", "bundle").Run()
	}()

	// Fetch from bundle
	fetchCmd := exec.Command("git", "fetch", "bundle")
	if output, err := fetchCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to fetch from bundle: %s", string(output))
	}

	// Merge
	mergeCmd := exec.Command("git", "merge", "bundle/"+branch, "--no-edit")
	if output, err := mergeCmd.CombinedOutput(); err != nil {
		// Fallback to master if main branch merge fails (common transition issue)
		mergeCmdMaster := exec.Command("git", "merge", "bundle/master", "--no-edit")
		if output2, err2 := mergeCmdMaster.CombinedOutput(); err2 != nil {
			return fmt.Errorf("merge failed (tried %s and master): %s", branch, string(output)+string(output2))
		}
	}

	return nil
}

// PushToOrigin pushes the specified branch to the 'origin' remote.
func PushToOrigin(branch string) error {
	cmd := exec.Command("git", "push", "origin", branch)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push origin failed: %v: %s", err, string(output))
	}
	return nil
}

// ShowRecentCommits prints the last n commits to stdout using a pretty graph format.
func ShowRecentCommits(n int) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", n), "--oneline", "--graph", "--color")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
