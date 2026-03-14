package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *Application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(a.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(a.home))
	mux.Handle("GET /stash/view/{id}", dynamic.ThenFunc(a.stashView))
	mux.Handle("GET /stash/create", dynamic.ThenFunc(a.stashCreate))
	mux.Handle("POST /stash/create", dynamic.ThenFunc(a.stashCreatePost))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(a.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(a.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(a.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(a.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(a.userLogoutPost))

	standard := alice.New(a.recoverPanic, a.logRequest, commonHandler)

	return standard.Then(mux)
}
