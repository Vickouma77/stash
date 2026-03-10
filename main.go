package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// home handler function with byte slice string
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("stash starts here"))
}

func stashView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	msg := fmt.Sprintf("Display specific stash with id %d...", id)
	w.Write([]byte(msg))
}

func stashCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display Form for new stash..."))
}

func stashCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Save a new stash..."))
}

func main() {
	//initialize a new servemux
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /stash/view/{id}", stashView)
	mux.HandleFunc("GET /stash/create", stashCreate)
	mux.HandleFunc("POST /stash/create", stashCreatePost)

	log.Print("Server starting on: 8000")

	//Start web server
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
