package client

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cln := NewClient()

	// Graceful shutdown on user's Quit
	usrQuit := make(chan os.Signal, 1)

	// Goroutine with a client
	go func(c chan os.Signal) {
		if err := cln.Start(); err != nil && err != ErrClientClosed {
			log.Fatalf("failed to start client: %s", err.Error())
		}
		c <- syscall.SIGQUIT
	}(usrQuit)

	log.Printf("Client started\n")

	// Graceful shutdown on OS signals
	sysQuit := make(chan os.Signal, 1)
	signal.Notify(sysQuit, syscall.SIGTERM, syscall.SIGINT)

	// Fan-in: wait for any quits
	select {
	case <-sysQuit:
	case <-usrQuit:
	}

	// Shutdown client w/ 5s timeout
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := cln.Stop(ctx); err != nil {
		log.Fatalf("failed to stop client: %s", err.Error())
	}

	log.Printf("Client stopped\n")
}
