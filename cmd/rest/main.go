package main

import (
	"context"
	"database/sql"
	"os"
	"sync"

	"github.com/ahmad-abuziad/transaction-processor/internal/cache"
	"github.com/ahmad-abuziad/transaction-processor/internal/data"
	"github.com/ahmad-abuziad/transaction-processor/internal/metrics"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type application struct {
	logger  *zap.Logger
	models  data.Models
	errors  httpErrors
	wg      sync.WaitGroup
	cache   cache.Cache
	metrics metrics.Metrics
}

func main() {
	// zap logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// mysql db
	db, err := openDB(os.Getenv("DSN"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// redis client
	redisClient, err := newRedisClient(os.Getenv("REDIS_ADDR"))
	if err != nil {
		logger.Error("unable to connect to cache, app will run without caching", zap.String("error", err.Error()))
	}

	// metrics
	m := metrics.NewMetrics()

	// app
	app := &application{
		logger:  logger,
		models:  data.NewModels(db),
		errors:  newHTTPErrors(logger),
		cache:   cache.NewCache(redisClient, m),
		metrics: m,
	}

	// serve
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

func newRedisClient(redisAddr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
