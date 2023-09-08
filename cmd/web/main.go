package main

import (
	"crypto/tls"
	"database/sql"
	"github.com/96malhar/snippetbox/internal/store"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"log/slog"
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
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", "addr", srv.Addr)
	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.logger.Error(err.Error())
	os.Exit(1)
}

func newApplication() *application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:         logger,
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
