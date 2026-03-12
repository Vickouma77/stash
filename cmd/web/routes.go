package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *Application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", a.home)
	mux.HandleFunc("GET /stash/view/{id}", a.stashView)
	mux.HandleFunc("GET /stash/create", a.stashCreate)
	mux.HandleFunc("POST /stash/create", a.stashCreatePost)

	standard := alice.New(a.recoverPanic, a.logRequest, commonHandler)

	return standard.Then(mux)
}
