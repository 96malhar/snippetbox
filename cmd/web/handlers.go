package main

import (
	"errors"
	"fmt"
	"github.com/96malhar/snippetbox/internal/store"
	"github.com/96malhar/snippetbox/internal/validation"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type snippetCreateForm struct {
	Title       string            `form:"title"`
	Content     string            `form:"content"`
	Expires     int               `form:"expires"`
	FieldErrors map[string]string `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippetStore.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippetStore.Get(id)
	if err != nil {
		if errors.Is(err, store.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fieldErrors := checkFields(
		fieldEntry{validation.NotBlank(form.Title), "title", "This field cannot be blank"},
		fieldEntry{validation.NotBlank(form.Content), "content", "This field cannot be blank"},
		fieldEntry{validation.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long"},
		fieldEntry{validation.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365"},
	)

	if len(fieldErrors) > 0 {
		data := app.newTemplateData(r)
		form.FieldErrors = fieldErrors
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippetStore.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
