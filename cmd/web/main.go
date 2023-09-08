package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/96malhar/snippetbox/internal/store"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	app := newApplication()

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":4000"
	}
	srv := &http.Server{
		Addr:         serverPort,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Printf("Starting server on %s", serverPort)
	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.errorLog.Fatal(err)
}

func newApplication() *application {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	debug := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	if *debug {
		infoLog.Print("Running in debug mode...")
	}

	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		debug:          *debug,
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippetStore:   store.NewSnippetStore(db),
		userStore:      store.NewUserStore(db),
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	return app
}

func openDB() (*sql.DB, error) {
	dsn := os.Getenv("SNIPPETBOX_DB_DSN")
	if dsn == "" {
		panic("SNIPPETBOX_DB_DSN environment variable not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
