package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/PuerkitoBio/ghost/handlers"
	"github.com/PuerkitoBio/ghost/templates"
	_ "github.com/PuerkitoBio/ghost/templates/amber"
	_ "github.com/PuerkitoBio/ghost/templates/gotpl"
	"github.com/howeyc/fsnotify"
)

const (
	publicDir   = "./public/"
	templateDir = "./templates/"
	port        = 9000
)

var (
	faviconPath   = path.Join(publicDir, "favicon.ico")
	faviconCache  = 1 * time.Minute
	staticHandler = http.StripPrefix("/public/", http.FileServer(http.Dir(publicDir)))
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// If a template matches the path, minus the starting /
	log.Printf("in serveTemplate, URL=%s\n", r.URL)
	err := templates.Render(r.URL.Path[1:], w, nil)
	if err != nil {
		w.Header().Del("Content-Type")
		staticHandler.ServeHTTP(w, r)
	}
}

func main() {
	log.SetFlags(0)

	if err := templates.CompileDir(templateDir); err != nil {
		panic(err)
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	if err := w.Watch(templateDir); err != nil {
		panic(err)
	}

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
