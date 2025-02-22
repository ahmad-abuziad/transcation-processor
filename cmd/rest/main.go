package main

import (
	"database/sql"
	"log/slog"
	"os"
	"sync"

	"github.com/ahmad-abuziad/transaction-processor/internal/data"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *slog.Logger
	models data.Models
	errors httpErrors
	wg     sync.WaitGroup
}

func main() {

	// use Zap or Logrus
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(os.Getenv("DSN"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		models: data.NewModels(db),
		errors: newHTTPErrors(logger),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
