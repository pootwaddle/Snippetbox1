package main

import (
	"flag"
	"log"
	"net/http"

	"snippetbox/pkg/models"
)

func main() {
	// Define command-line flags for the network address and location of the static
	// files directory.
	addr := flag.String("addr", ":4000", "HTTP network address")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")

	flag.Parse()

	// Initialize a new instance of App containing the dependencies.
	app := &App{
		Database:  &models.Database{},
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
	}

	// Again, we dereference the addr variable and use it as the network address
	// to listen on. Notice that we also use the log.Printf() function to interpolate
	// the correct address in the log message.
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, app.Routes())
	log.Fatal(err)
}
