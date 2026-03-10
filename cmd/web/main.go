package main

import (
	"log"
	"net/http"
)

func main() {
	//initialize a new servemux
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /stash/view/{id}", stashView)
	mux.HandleFunc("GET /stash/create", stashCreate)
	mux.HandleFunc("POST /stash/create", stashCreatePost)

	log.Print("Server starting on: 8000")

	//Start web server
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
