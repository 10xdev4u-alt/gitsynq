package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "👀 Watch for changes and auto-sync",
	Long:  `Automatically run 'gitsync push' whenever files in the repository are changed.`, 
	Run:   runWatch,
}

func runWatch(cmd *cobra.Command, args []string) {
	printBanner()
	green.Println("\n👀 Watching for changes... (Press Ctrl+C to stop)")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		red.Printf("❌ Failed to create watcher: %v\n", err)
		os.Exit(1)
	}
	defer watcher.Close()

	done := make(chan bool)
	
	// Debounce timer
	var timer *time.Timer
	const delay = 2 * time.Second

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(delay, func() {
						cyan.Printf("\n🔄 Change detected in %s. Syncing...\n", event.Name)
						// We need a way to auto-commit or just push if commits exist
						// For now, let's just trigger a normal push
						runPush(cmd, args)
						green.Println("\n👀 Still watching...")
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				red.Printf("❌ Watcher error: %v\n", err)
			}
		}
	}()

	// Add directories to watch
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip hidden directories like .git
			if len(path) > 1 && path[0] == '.' && path != "./" {
				return filepath.SkipDir
			}
			if strings.Contains(path, ".gitsync-bundles") {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		red.Printf("❌ Failed to walk directory: %v\n", err)
		os.Exit(1)
	}

	<-done
}
