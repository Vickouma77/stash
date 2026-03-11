package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"stash.io/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	//Command-line flag 'addr', a default value of :8000
	addr := flag.String("addr", ":8000", "HTTP network address")
	//Command-line flag 'dsn', MYSQL data source name
	dsn := flag.String("mysql", "root:rootpassword@/stash?parseTime=true", "MYSQL data source name")
	flag.Parse()

	//initialize a new structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	//Connection pool
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	//Initialize a new instance of Application struct
	app := &Application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("Starting server", slog.String("addr", ":8000"))

	//Start web server
	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
