package main

import (
	"net/http"
	"strconv"
)

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w) // Use the app.NotFound() helper.
		return
	}

//Fetch a slice of the latest snippets from the database.
snippets, err := app.Database.LatestSnippets()
if err != nil {
	app.ServerError(w,err)
	return
}

	app.RenderHTML(w, "home.page.html", &HTMLData{
		Snippets: snippets,
	})

}

func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	snippet, err := app.Database.GetSnippet(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if snippet == nil {
		app.NotFound(w)
		return
	}

	// Render the show.page.html template, passing in the snippet data wrapped in
	// our HTMLData struct.
	app.RenderHTML(w, "show.page.html", &HTMLData{
		Snippet: snippet,
	})
}

// Add a placeholder NewSnippet handler function.
func (app *App) NewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the new snippet form..."))
}
