package msgstore

import (
	"testing"
)

func TestNewAddGet(t *testing.T) {
	// New
	ms := New()

	if l := len(ms.messages); l > 0 {
		t.Errorf("got len(messages)=%d, want 0", l)
	}

	// AddMessage
	if err := ms.AddMessage("TestMsg1", "User1"); err != nil {
		t.Errorf("error on AddMessage, want no error")
	}

	if l := len(ms.messages); l != 1 {
		t.Errorf("got len(messages)=%d, want 1", l)
	}

	// AddMessage up to 10 items
	_ = ms.AddMessage("TestMsg2", "User1")
	_ = ms.AddMessage("TestMsg3", "User2")
	_ = ms.AddMessage("TestMsg4", "User1")
	_ = ms.AddMessage("TestMsg5", "User2")
	_ = ms.AddMessage("TestMsg6", "User3")
	_ = ms.AddMessage("TestMsg7", "User2")
	_ = ms.AddMessage("TestMsg8", "User1")
	_ = ms.AddMessage("TestMsg9", "User3")
	_ = ms.AddMessage("TestMsg10", "User2")

	if l := len(ms.messages); l != 10 {
		t.Errorf("got len(messages)=%d, want 10", l)
	}

	// AddMessage more than 10 items
	if err := ms.AddMessage("TestMsg11", "User3"); err == nil {
		t.Errorf("no error on AddMessage, want an error")
	}

	// GetMessages for User1, must remain 6 items
	arr1 := ms.GetMessages("User1")
	if l := len(arr1); l != 4 {
		t.Errorf("got len(arr1)=%d, want 4", l)
	}
	if l := len(ms.messages); l != 6 {
		t.Errorf("got len(messages)=%d, want 6", l)
	}

	// GetMessages for User2, must remain 2 items
	arr2 := ms.GetMessages("User2")
	if l := len(arr2); l != 4 {
		t.Errorf("got len(arr2)=%d, want 4", l)
	}
	if l := len(ms.messages); l != 2 {
		t.Errorf("got len(messages)=%d, want 2", l)
	}

	// GetMessages for User3, must remain 0 items
	arr3 := ms.GetMessages("User3")
	if l := len(arr3); l != 2 {
		t.Errorf("got len(arr3)=%d, want 2", l)
	}
	if l := len(ms.messages); l != 0 {
		t.Errorf("got len(messages)=%d, want 0", l)
	}
}
