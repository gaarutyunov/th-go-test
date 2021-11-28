package server

import "testing"

func TestNew(t *testing.T) {
	srv := NewServer()
	if srv.httpServer == nil {
		t.Errorf("no *http.server, want *http.server")
	}
}
