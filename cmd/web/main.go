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

	logger.Info("Starting server", slog.String("addr", ":8000"))

	//Start web server
	err := http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}
