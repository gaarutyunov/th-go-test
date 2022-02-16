package server

import (
	"context"
	"log"
	"net/http"
	"th-go-test/pkg/msgstore"
	"time"
)

type Server struct {
	httpServer *http.Server
	storage    *msgstore.MsgStore
}

const (
	serverAddr = "0.0.0.0:8080"
	dataPath   = "./server.json"
)

func NewServer() *Server {
	storage := msgstore.New()
	handler := NewMsgHandler(storage)
	server := &http.Server{
		Addr:           serverAddr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		httpServer: server,
		storage:    storage,
	}
}

func (s *Server) Start() error {
	if n, err := s.storage.LoadMessages(dataPath); err != nil {
		log.Printf("Error loading messages: %s", err.Error())
	} else {
		log.Printf("Loaded %d messages", n)
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if n, err := s.storage.SaveMessages(dataPath); err != nil {
		log.Printf("Error saving messages: %s", err.Error())
	} else {
		log.Printf("Saved %d messages", n)
	}

	return s.httpServer.Shutdown(ctx)
}
