package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/PuerkitoBio/ghost/handlers"
	"github.com/PuerkitoBio/ghost/templates"
	_ "github.com/PuerkitoBio/ghost/templates/amber"
	_ "github.com/PuerkitoBio/ghost/templates/gotpl"
	"github.com/howeyc/fsnotify"
)

const (
	publicDir        = "./public/"
	templateDir      = "./templates/"
	port             = 9000
	tmplRefreshDelay = 2 * time.Second
	faviconCache     = 1 * time.Minute
)

var (
	faviconPath   = path.Join(publicDir, "favicon.ico")
	staticHandler = http.StripPrefix("/public/", http.FileServer(http.Dir(publicDir)))

	// Protect the templates
	tmplLock sync.RWMutex
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// If a template matches the path, minus the starting /
	tmplLock.RLock()
	err := templates.Render(r.URL.Path[1:], w, nil)
	tmplLock.RUnlock()
	if err != nil {
		w.Header().Del("Content-Type")
		staticHandler.ServeHTTP(w, r)
	}
}

func refreshTemplates() {
	tmplLock.Lock()
	defer tmplLock.Unlock()
	if err := templates.CompileDir(templateDir); err != nil {
		log.Fatal(err)
	}
}

func watchDir(watcher *fsnotify.Watcher) {
	var delay <-chan time.Time
	for {
		select {
		case <-watcher.Event:
			if delay == nil {
				delay = time.After(tmplRefreshDelay)
			}
		case err, ok := <-watcher.Error:
			log.Printf("watch error: %s\n", err)
			if !ok {
				return
			}
		case <-delay:
			refreshTemplates()
			delay = nil
		}
	}
}

func watchRecursive(w *fsnotify.Watcher) {
	filepath.Walk(templateDir, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			if err := w.Watch(path); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
}

func main() {
	log.SetFlags(0)

	// Initial compilation of the templates, then activate the watcher
	refreshTemplates()
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go watchDir(w)
	watchRecursive(w)
	defer w.Close()

	h := handlers.FaviconHandler(
		handlers.PanicHandler(
			handlers.LogHandler(
				handlers.GZIPHandlerFunc(
					serveTemplate,
					nil),
				handlers.NewLogOptions(nil, handlers.Lshort)),
			nil),
		faviconPath,
		faviconCache)

	// Assign the combined handler to the server.
	http.Handle("/", h)

	// Start it up.
	log.Printf("server listening on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
