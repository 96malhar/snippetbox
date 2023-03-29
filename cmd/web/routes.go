package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(app.recoverPanic, middleware.StripSlashes, app.logRequest, middleware.GetHead, secureHeaders)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Method(http.MethodGet, "/static/*", http.StripPrefix("/static", fileServer))

	r.Get("/", app.home)
	r.Get("/snippet/view/{id}", app.snippetView)
	r.Post("/snippet/create", app.snippetCreate)
	return r
}
