package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"stash.io/internal/models"
)

// home handler function with byte slice string
func (a *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := a.snippets.Latest()
	if err != nil {
		a.ServerError(w, r, err)
		return
	}

	data := a.newTemplateData(r)
	data.Snippets = snippets

	a.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (a *Application) stashView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			a.ServerError(w, r, err)
		}
		return
	}

	data := a.newTemplateData(r)
	data.Snippet = snippet

	a.render(w, r, http.StatusOK, "view.tmpl.html", data)
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
