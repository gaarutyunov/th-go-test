package msgstore

import (
	"fmt"
	"sync"
)

// MsgLimit Лимит для сообщений в хранилище
const MsgLimit = 10

// Message Прототип сообщения
type Message struct {
	ID      int    `json:"id"`
	OwnerID string `json:"owner_id"`
	Message string `json:"message"`
}

// MsgStore Простое хранилище сообщений в памяти, безопасно для конкурентного доступа
// TODO по ТЗ переделать на chan?
type MsgStore struct {
	sync.Mutex

	messages map[int]Message
	nextID   int
}

// New Конструктор экземпляра хранилище сообщений
func New() *MsgStore {
	ms := &MsgStore{}
	ms.messages = make(map[int]Message)
	ms.nextID = 0

	return ms
}

// AddMessage Добавляет сообщение в хранилище сообщений с сохранением владельца, есть лимит сообщений
func (ms *MsgStore) AddMessage(text, owner string) error {
	ms.Lock()
	defer ms.Unlock()

	if len(ms.messages) >= MsgLimit {
		return fmt.Errorf("messages limit reached, no more space")
	}

	msg := Message{
		ID:      ms.nextID,
		OwnerID: owner,
		Message: text,
	}

	ms.messages[ms.nextID] = msg
	ms.nextID++

	return nil
}

// GetMessages Получает все сообщения из хранилища для указанного владельца и удаляет их из хранилища
func (ms *MsgStore) GetMessages(owner string) []Message {
	ms.Lock()
	defer ms.Unlock()

	var messages []Message

	for idx, msg := range ms.messages {
		if owner == msg.OwnerID {
			messages = append(messages, msg)
			delete(ms.messages, idx)
		}
	}

	return messages
}
