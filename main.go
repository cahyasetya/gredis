package main

import (
	"context"

	"os"
	"os/signal"
	"syscall"

	"github.com/cahyasetya/gredis/logger" // Import the logger package
	"github.com/cahyasetya/gredis/server"
	"github.com/cahyasetya/gredis/storage"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	defer logger.Log.Sync()

	storage.InitStorage()

	// Create and start the server
	srv := server.New(server.Config{Port: 6379})
	ctx := context.Background()
	if err := srv.Start(ctx); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}

	// Wait for a signal to shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Shutdown the server
	logger.Log.Info("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Error during server shutdown", zap.Error(err))
	}
}
