package main

import (
	"context"
	"log/slog"
	"nova/internal/handler"
	mapstorage "nova/internal/storage/map"
	"nova/internal/tcp"
	"os"
)

// TODO: move to configuration module (or package)
var (
	addr = "localhost:6379"
)

func main() {
	storage := mapstorage.New(context.Background())

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	handler := handler.NewHandler(log, storage)

	srv := &tcp.Server{
		Addr:    addr,
		Handler: handler,
	}

	srv.ListenAndServe()
}
