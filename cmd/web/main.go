package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/pootwaddle/snippetbox1/pkg/models"
)

func main() {
	// Define command-line flags for the network address and location of the static
	// files directory.
	addr := flag.String("addr", ":4000", "HTTP network address")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
	dsn := flag.String("dsn", "sb:12345@/snippetbox?parseTime=true", "MySql DSN")

	flag.Parse()

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate connect() function below. We pass connect() the DSN
	// from the command-line flag.
	db := connect(*dsn)

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new instance of App containing the dependencies.
	app := &App{
		Database:  &models.Database{db},
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

// The connect() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
