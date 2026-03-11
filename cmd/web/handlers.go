package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// home handler function with byte slice string
func (a *Application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.ServerError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		a.ServerError(w, r, err)
	}
}

func (a *Application) stashView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display specific stash with id %d...", id)
}

func (a *Application) stashCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display Form for new stash..."))
}

func (a *Application) stashCreatePost(w http.ResponseWriter, r *http.Request) {
	//Dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/stash/view/%d", id), http.StatusSeeOther)
}
