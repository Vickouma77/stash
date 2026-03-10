package main

import (
	"log"
	"net/http"
)

// home handler function with byte slice string
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("stash starts here"))
}

func main() {
	//initialize a new servemux
	mux := http.NewServeMux()
	//register home function as the handler for this url "/"
	mux.HandleFunc("/", home)

	log.Print("Server starting on: 8000")

	//Start web server
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
