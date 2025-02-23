package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("caught signal", zap.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks", zap.String("addr", srv.Addr))

		app.stopWorkers()
		shutdownError <- nil
	}()

	app.startWorkers()

	app.logger.Info("starting server", zap.String("addr", srv.Addr))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", zap.String("addr", srv.Addr))

	return nil
}
