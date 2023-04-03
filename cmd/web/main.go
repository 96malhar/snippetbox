package main

import (
	"database/sql"
	"flag"
	"github.com/96malhar/snippetbox/internal/store"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const (
	psqlDriver = "postgres"
	psqlInfo   = "host=localhost port=5432 user=web password=malhar123 sslmode=disable dbname=snippetbox"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", psqlInfo, "Postgres data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(psqlDriver, *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippetStore:  store.NewSnippetStore(db),
		templateCache: templateCache,
		formDecoder:   form.NewDecoder(),
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given driverName and DSN.
func openDB(driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
