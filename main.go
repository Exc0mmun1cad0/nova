package main

import (
	"context"
	"nova/internal/handler"
	mapstorage "nova/internal/storage/map"
	"nova/internal/tcp"
	"nova/pkg/logger"

	"go.uber.org/zap"
)

// TODO: move to configuration module (or package)
var (
	addr = "localhost:6379"
)

func main() {
	log := logger.Setup()

	log.Info("starting nova")

	log.Info("initializing storage")
	storage := mapstorage.New(context.Background())

	srv, err := tcp.NewServer(
		addr,
		handler.NewHandler(storage),
		log,
	)
	if err != nil {
		log.Panic("failed to init tcp server", zap.Error(err))
	}

	// TODO: wrap in MustRun() function
	srv.ListenAndServe()
}
