package main

import (
	"flag"
	"log"
	"net/http"
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
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippet", app.ShowSnippet)
	mux.HandleFunc("/snippet/new", app.NewSnippet)

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So we need to dereference the pointer (i.e.
	// prefix it with the * symbol) before we use it as the path for our static file
	// server.
	fileServer := http.FileServer(http.Dir(*staticDir))

	// Use the mux.Handle() function to register the file server as the
	// handler for all URL paths that start with "/static/". For matching
	// paths, we strip the "/static" prefix before the request reaches the file
	// server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Again, we dereference the addr variable and use it as the network address
	// to listen on. Notice that we also use the log.Printf() function to interpolate
	// the correct address in the log message.
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, app.Routes())
	log.Fatal(err)
}
