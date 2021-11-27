package server

import (
	"fmt"
	"net/http"

	"th-go-test/pkg/msgstore"
)

type MsgHandler struct {
	storage *msgstore.MsgStore
}

func NewMsgHandler(storage *msgstore.MsgStore) *MsgHandler {
	return &MsgHandler{
		storage: storage,
	}
}

func (h *MsgHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getMessages(w, r)
	case http.MethodPut:
		h.addMessage(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *MsgHandler) addMessage(w http.ResponseWriter, r *http.Request) {
	// TODO json unmarshall, message add, handle errors
	fmt.Fprintf(w, "ADD, %v, http: %v", r.URL.Path, r.TLS == nil)
}

func (h *MsgHandler) getMessages(w http.ResponseWriter, r *http.Request) {
	// TODO messages get, json marshall, handle errors
	fmt.Fprintf(w, "GET, %v, http: %v", r.URL.Path, r.TLS == nil)
}
