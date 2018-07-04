package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"snippetbox/pkg/models"
)

// Define a new HTMLData struct to act as a wrapper for the dynamic data we want
// to pass to our templates. For now this just contains the snippet data that we
// want to display, which has the underling type *models.Snippet.
type HTMLData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func (app *App) RenderHTML(w http.ResponseWriter, page string, data *HTMLData) {
	files := []string{
		filepath.Join(app.HTMLDir, "base.html"),
		filepath.Join(app.HTMLDir, page),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err) // Use the new app.ServerError() helper.
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.ServerError(w, err) // Use the new app.ServerError() helper.
	}
}
