package main

import (
	"net/http"
)

//Server Error.Writing log entry at Error level,then sends generic 500 internal server error
func (a *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri = r.URL.RequestURI()
	)

	a.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//Client Error
func (a *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}