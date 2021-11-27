package server

import (
	"context"
	"net/http"
	"time"

	"th-go-test/pkg/msgstore"
)

type Server struct {
	httpServer *http.Server
	storage    *msgstore.MsgStore
}

func NewServer() *Server {
	storage := msgstore.New()
	handler := NewMsgHandler(storage)
	server := &http.Server{
		Addr:           "0.0.0.0:8080",
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
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
