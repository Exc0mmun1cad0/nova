package main

import (
	"context"
	"nova/internal/handler"
	mapstorage "nova/internal/storage/map"
	"nova/internal/tcp"

	"go.uber.org/zap"
)

// TODO: move to configuration module (or package)
var (
	addr = "localhost:6379"
)

func main() {

	storage := mapstorage.New(context.Background())

	log := setupLogger()

	handler := handler.NewHandler(log, storage)

	srv := &tcp.Server{
		Addr:    addr,
		Handler: handler,
	}

	srv.ListenAndServe()
}

// TODO: add more options
func setupLogger() *zap.Logger {
	logger := zap.Must(zap.NewProduction())
	return logger
}
