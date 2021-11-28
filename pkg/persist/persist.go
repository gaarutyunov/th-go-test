package persist

import (
	"encoding/json"
	"os"
	"sync"
)

// Lock локальный mutex, чтобы не писать и не читать одновременно
var lock sync.Mutex

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
