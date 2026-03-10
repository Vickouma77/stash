package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type Application struct {
	logger *slog.Logger
}

func main() {
	//Command-line flag 'addr', a default value of :8000
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	//initialize a new structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	//Initialize a new instance of Application struct
	app := &Application{
		logger: logger,
	}

	//initialize a new servemux
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /stash/view/{id}", app.stashView)
	mux.HandleFunc("GET /stash/create", app.stashCreate)
	mux.HandleFunc("POST /stash/create", app.stashCreatePost)

	logger.Info("Starting server", slog.String("addr", ":8000"))

	//Start web server
	err := http.ListenAndServe(*addr, mux)

	logger.Error(err.Error())
	os.Exit(1)
}
