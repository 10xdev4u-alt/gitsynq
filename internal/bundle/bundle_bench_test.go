package bundle

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func BenchmarkCreateFull(b *testing.B) {
	dir := b.TempDir()
	
	// Setup repo once
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	cmd.Run()
	
	err := os.WriteFile(filepath.Join(dir, "large_file.txt"), make([]byte, 1024*1024), 0644)
	if err != nil {
		b.Fatal(err)
	}
	
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "large commit").Run()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bundlePath := filepath.Join(dir, "bench.bundle")
		
		oldWd, _ := os.Getwd()
		os.Chdir(dir)
		CreateFull(bundlePath)
		os.Chdir(oldWd)
		
		os.Remove(bundlePath)
	}
}
