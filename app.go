package main

import (
	"context"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/aki-kuramoto/sosaykey/pkg/markdown"
	"github.com/aki-kuramoto/sosaykey/pkg/watcher"
)

// App struct
type App struct {
	ctx context.Context
	watcher *watcher.FileWatcher
}

// NewApp creates a new App application struct
func NewApp() *App {
	w, err := watcher.NewFileWatcher()
	if err != nil {
		// In production we should handle this better, but for now panic or log?
		// Since we can't return error here easily without changing signature used in main...
		// Let's print and leave nil, checking for nil later.
		println("Error creating watcher:", err.Error())
		return &App{}
	}
	return &App{
		watcher: w,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called at termination
func (a *App) shutdown(ctx context.Context) {
	if a.watcher != nil {
		a.watcher.Close()
	}
}

// LoadFile reads the file at the given path, parses it as markdown, and returns the HTML string.
func (a *App) LoadFile(path string) (string, error) {
	// Start watching the file
	if a.watcher != nil {
		err := a.watcher.Watch(path, func() {
			runtime.EventsEmit(a.ctx, "file:changed", path)
		})
		if err != nil {
			fmt.Printf("Error watching file: %v\n", err)
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	parser := markdown.NewParser(string(content))
	doc := parser.Parse()
	html := doc.ToHTML()

	return html, nil
}
