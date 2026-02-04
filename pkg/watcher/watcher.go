package watcher

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileWatcher watches a file for changes.
type FileWatcher struct {
	watcher *fsnotify.Watcher
	stop    chan struct{}

	mu       sync.Mutex
	filePath string
	callback func()
}

// NewFileWatcher creates a new FileWatcher.
func NewFileWatcher() (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		watcher: w,
		stop:    make(chan struct{}),
	}

	// Start the single event loop
	go fw.startLoop()

	return fw, nil
}

// Watch starts watching the specified file.
// If it was already watching a file, it stops watching the old one.
// The callback is executed when the file changes.
func (fw *FileWatcher) Watch(path string, callback func()) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// Remove old watch if exists
	if fw.filePath != "" {
		_ = fw.watcher.Remove(fw.filePath)
	}

	// Clean path just in case
	cleanPath := filepath.Clean(path)

	err := fw.watcher.Add(cleanPath)
	if err != nil {
		return err
	}

	fw.filePath = cleanPath
	fw.callback = callback
	return nil
}

func (fw *FileWatcher) startLoop() {
	var lastEvent time.Time
	debounceDuration := 100 * time.Millisecond

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// We care about Write or Rename
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Rename) || event.Has(fsnotify.Chmod) {
				if time.Since(lastEvent) < debounceDuration {
					continue
				}
				lastEvent = time.Now()

				log.Printf("File changed event: %s (Op: %s)", event.Name, event.Op)

				fw.mu.Lock()
				cb := fw.callback
				currentPath := fw.filePath
				fw.mu.Unlock()

				if cb != nil {
					// Check if valid? fsnotify might send removal event.
					// We blindly trigger.
					cb()

					// Re-add logic for atomic saves (Rename/Remove)
					if event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
						time.Sleep(50 * time.Millisecond)
						// We need to re-add the SAME path
						// Note: This needs to be done carefully to avoid deadlock if we used Start/Stop logic
						// But here we access fw.watcher directly.
						_ = fw.watcher.Add(currentPath)
					}
				}
			}
		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		case <-fw.stop:
			return
		}
	}
}

// Close stops the watcher.
func (fw *FileWatcher) Close() {
	close(fw.stop)
	fw.watcher.Close()
}
