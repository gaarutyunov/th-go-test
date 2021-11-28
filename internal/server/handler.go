package server

import (
	"encoding/json"
	"log"
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

func (h *MsgHandler) getMessages(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "400 Bad Request", 400)
		return
	}

	type JSONRequest struct {
		PersonID string `json:"person_id"`
	}

	var jr JSONRequest
	if err := json.NewDecoder(r.Body).Decode(&jr); err != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	if jr.PersonID == "" {
		http.Error(w, "401 Unauthorized", 401)
		return
	}

	type JSONResponse struct {
		Messages []string `json:"messages"`
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	jw := JSONResponse{}
	for _, msg := range h.storage.GetMessages(jr.PersonID) {
		jw.Messages = append(jw.Messages, msg.Message)
	}

	if err := json.NewEncoder(w).Encode(jw); err != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	log.Printf("MsgStore size is %d", h.storage.Length())
}

func (h *MsgHandler) addMessage(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "400 Bad Request", 400)
		return
	}

	type JSONRequest struct {
		PersonID string `json:"person_id"`
		Message  string `json:"message"`
	}

	var jr JSONRequest
	if err := json.NewDecoder(r.Body).Decode(&jr); err != nil {
		http.Error(w, "500 Internal Server Error", 500)
		return
	}

	if jr.PersonID == "" {
		http.Error(w, "401 Unauthorized", 401)
		return
	}

	if err := h.storage.AddMessage(jr.Message, jr.PersonID); err != nil {
		http.Error(w, "507 Insufficient Storage", 507)
		return
	}

	// Response
	http.Error(w, "200 OK", 200)

	log.Printf("MsgStore size is %d", h.storage.Length())
}
