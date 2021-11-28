package persist

import (
	"os"
	"reflect"
	"testing"
)

const TestFilePath = "./test.json"

type TestStruct struct {
	Field1 string
	Field2 int
	Field3 []float32
}

func NewTestStruct() *TestStruct {
	return &TestStruct{
		Field1: "value",
		Field2: 123,
		Field3: []float32{4.5, 6.7, 8.9},
	}
}

func TestSave(t *testing.T) {
	obj := NewTestStruct()

	if err := Save(TestFilePath, obj); err != nil {
		t.Fatalf("error saving v, want no error")
	}
}

func TestLoad(t *testing.T) {
	var obj TestStruct

	if err := Load(TestFilePath, &obj); err != nil {
		t.Fatalf("error loading v, want no error")
	}

	obj2 := NewTestStruct()

	if reflect.DeepEqual(obj, obj2) {
		t.Fatalf("error loaded v not match, want match")
	}
}

func TestClean(t *testing.T) {
	f, err := os.Open(TestFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("error test.json not found")
	}
	_ = f.Close()

	if err := os.Remove(TestFilePath); err != nil {
		t.Fatalf("error removing test.json")
	}
}
