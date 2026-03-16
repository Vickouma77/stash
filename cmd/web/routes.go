package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *Application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// dynamic includes session management so handlers can read and write session data.
	dynamic := alice.New(a.sessionManager.LoadAndSave, noSurf)

	mux.Handle("GET /{$}", dynamic.ThenFunc(a.home))
	mux.Handle("GET /stash/view/{id}", dynamic.ThenFunc(a.stashView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(a.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(a.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(a.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(a.userLoginPost))

	// protected restricts routes to authenticated users only.
	protected := dynamic.Append(a.requireAuthentication)

	mux.Handle("GET /stash/create", protected.ThenFunc(a.stashCreate))
	mux.Handle("POST /stash/create", protected.ThenFunc(a.stashCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(a.userLogoutPost))

	standard := alice.New(a.recoverPanic, a.logRequest, commonHandler)

	return standard.Then(mux)
}
