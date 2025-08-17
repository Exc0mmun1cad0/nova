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
	log := setupLogger()

	log.Info("starting nova")

	log.Info("initializing storage")
	storage := mapstorage.New(context.Background())

	handler := handler.NewHandler(log, storage)
	srv, err := tcp.NewServer(addr, handler, log)
	if err != nil {
		log.Panic("failed to init tcp server", zap.Error(err))
	}

	// TODO: wrap in MustRun() function
	srv.ListenAndServe()
}

// TODO: add more options
func setupLogger() *zap.Logger {
	logger := zap.Must(zap.NewProduction())
	return logger
}
