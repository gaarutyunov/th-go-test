package server

import (
	"os"
	"testing"
)

const TestFilePath = "./test.json"

func TestNew(t *testing.T) {
	srv := NewServer()
	if srv.httpServer == nil {
		t.Errorf("no *http.server, want *http.server")
	}
}

func TestSave(t *testing.T) {
	srv := NewServer()

	_ = srv.storage.AddMessage("TestMsg2", "User1")
	_ = srv.storage.AddMessage("TestMsg3", "User2")
	_ = srv.storage.AddMessage("TestMsg4", "User1")
	_ = srv.storage.AddMessage("TestMsg5", "User2")
	_ = srv.storage.AddMessage("TestMsg6", "User3")

	if err := srv.storage.SaveMessages(TestFilePath); err != nil {
		t.Errorf("error saving messages, want no errors")
	}

	if length := srv.storage.Length(); length != 0 {
		t.Errorf("error got messages length %d, want 0", length)
	}
}

func TestLoad(t *testing.T) {
	srv := NewServer()

	if err := srv.storage.LoadMessages(TestFilePath); err != nil {
		t.Errorf("error loading messages, want no errors")
	}

	if length := srv.storage.Length(); length != 5 {
		t.Errorf("error got messages length %d, want 5", length)
	}
}

func TestClean(t *testing.T) {
	f, err := os.Open(TestFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("error %s not found", TestFilePath)
	}
	_ = f.Close()

	if err := os.Remove(TestFilePath); err != nil {
		t.Fatalf("error removing %s", TestFilePath)
	}
}
