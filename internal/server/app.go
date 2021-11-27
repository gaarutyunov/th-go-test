package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	srv := NewServer()

	// Goroutine with a server
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %s", err.Error())
		}
	}()

	log.Printf("Server started\n")

	// Graceful shutdown on OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	// Shutdown server w/ 5s timeout
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("failed to stop server: %s", err.Error())
	}

	log.Printf("Server stopped\n")
}
