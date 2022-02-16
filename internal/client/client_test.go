package client

import (
	"testing"
)

func TestNew(t *testing.T) {
	cln := NewClient()
	if cln.httpClient == nil {
		t.Errorf("no *http.client, want *http.client")
	}
}
