package persist

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)

// Lock локальный mutex, чтобы не писать и не читать одновременно
var lock sync.Mutex

// Marshal is a function that marshals the object into an io.Reader.
// By default, it uses the JSON marshaller.
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Save пишет из переменной v в файл по указанному пути.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(v)
}

// Load читает из файла по указанному пути в переменную v.
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}
