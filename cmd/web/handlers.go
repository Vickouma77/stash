package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// home handler function with byte slice string
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	w.Write([]byte("stash starts here"))
}

func stashView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display specific stash with id %d...", id)
}

func stashCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display Form for new stash..."))
}

func stashCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Save a new stash..."))
}