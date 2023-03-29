package main

import (
	"bytes"
	"fmt"
	"github.com/96malhar/snippetbox/internal/store"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	snippetStore interface {
		Insert(title string, content string, expirationDays int) (int, error)
		Get(id int) (*store.Snippet, error)
		Latest() ([]*store.Snippet, error)
	}
	templateCache map[string]*template.Template
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
