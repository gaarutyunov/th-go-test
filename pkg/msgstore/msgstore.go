package msgstore

import (
	"fmt"
	"sync"
)

// MsgLimit Лимит для сообщений в хранилище
const MsgLimit = 10

// Message Прототип сообщения
type Message struct {
	PersonID string `json:"person_id"`
	Message  string `json:"message"`
}

// MsgStore Простое хранилище сообщений в памяти, безопасно для конкурентного доступа
type MsgStore struct {
	sync.Mutex

	messages chan Message
}

// New Конструктор экземпляра хранилище сообщений
func New() *MsgStore {
	ms := &MsgStore{}
	ms.messages = make(chan Message, MsgLimit)

	return ms
}

// AddMessage Добавляет сообщение в хранилище сообщений с сохранением владельца, есть лимит сообщений
func (ms *MsgStore) AddMessage(text, owner string) error {
	ms.Lock()
	defer ms.Unlock()

	if len(ms.messages) >= cap(ms.messages) {
		return fmt.Errorf("messages limit reached, no more space")
	}

	msg := Message{
		PersonID: owner,
		Message:  text,
	}

	ms.messages <- msg

	return nil
}

// GetMessages Получает все сообщения из хранилища для указанного владельца и удаляет их из хранилища
func (ms *MsgStore) GetMessages(owner string) []Message {
	ms.Lock()
	defer ms.Unlock()

	close(ms.messages)
	t := make(chan Message, cap(ms.messages))

	var messages []Message
	for msg := range ms.messages {
		if owner == msg.PersonID {
			messages = append(messages, msg)
		} else {
			t <- msg
		}
	}

	ms.messages = t

	return messages
}

func (ms *MsgStore) Length() int {
	return len(ms.messages)
}
