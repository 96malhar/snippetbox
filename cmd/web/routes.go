package main

import (
	"github.com/96malhar/snippetbox/ui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {

	r := chi.NewRouter()
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	r.Method(http.MethodGet, "/static/*", fileServer)
	r.Get("/ping", ping)

	standardMiddlewares := []func(handler http.Handler) http.Handler{
		app.recoverPanic, middleware.StripSlashes, app.logRequest, middleware.GetHead, secureHeaders,
	}

	r.Group(func(r chi.Router) {
		r.Use(standardMiddlewares...)
		r.Use(app.sessionManager.LoadAndSave, app.authenticate)
		r.Get("/", app.home)
		r.Get("/about", app.about)
		r.Get("/snippet/view/{id}", app.snippetView)
		r.Get("/user/signup", app.userSignup)
		r.Post("/user/signup", app.userSignupPost)
		r.Get("/user/login", app.userLogin)
		r.Post("/user/login", app.userLoginPost)
	})

	r.Group(func(r chi.Router) {
		r.Use(standardMiddlewares...)
		r.Use(app.sessionManager.LoadAndSave, app.authenticate, app.requireAuthentication)
		r.Get("/snippet/create", app.snippetCreate)
		r.Post("/snippet/create", app.snippetCreatePost)
		r.Post("/user/logout", app.userLogoutPost)
	})

	return r
}
