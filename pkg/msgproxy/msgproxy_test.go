package msgproxy

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	// New
	mp := New(&http.Client{})

	if l := len(mp.messages); l > 0 {
		t.Errorf("got len(messages)=%d, want 0", l)
	}
}
