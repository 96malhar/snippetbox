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

	r.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)

		r.Get("/", app.home)
		r.Get("/snippet/view/{id}", app.snippetView)
		r.Get("/user/signup", app.userSignup)
		r.Post("/user/signup", app.userSignupPost)
		r.Get("/user/login", app.userLogin)
		r.Post("/user/login", app.userLoginPost)
	})

	r.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave, app.requireAuthentication)

		r.Get("/snippet/create", app.snippetCreate)
		r.Post("/snippet/create", app.snippetCreatePost)
		r.Post("/user/logout", app.userLogoutPost)
	})

	return r
}
