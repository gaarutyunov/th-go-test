package msgproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// serverURL постоянный URL для обращения на сервер по HTTP
const serverURL = "http://localhost:8080/"

// Message Прототип сообщения
type Message struct {
	PersonID string `json:"person_id"`
	Message  string `json:"message"`
}

// MsgProxy Простая прокси-очередь сообщений в памяти, самостоятельно связывается с сервером, клиент ждёт отправки
type MsgProxy struct {
	connector *http.Client
	messages  chan Message
}

// New Конструктор экземпляра прокси-очереди сообщений
func New(connector *http.Client) *MsgProxy {
	mp := &MsgProxy{
		connector: connector,
	}
	mp.messages = make(chan Message, 1)

	return mp
}

// AddMessage Добавляет сообщение в прокси-очередь сообщений с сохранением владельца, ждёт отправки на сервер
func (mp *MsgProxy) AddMessage(text, owner string) {
	mp.messages <- Message{
		PersonID: owner,
		Message:  text,
	}
	mp.sender()
}

// sender Функция-отправитель
func (mp *MsgProxy) sender() {
	// Fetch a message
	select {
	case msg, ok := <-mp.messages:
		if !ok {
			fmt.Printf("proxy: error, channel closed")
			break
		}

		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(&msg)

		req, err := http.NewRequest(http.MethodPut, serverURL, payloadBuf)
		if err != nil {
			fmt.Printf("proxy: error making a request (%s)\n", err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")

		// Repeatable part, until done
	sloop:
		for {
			res, err := mp.connector.Do(req)
			if err != nil {
				fmt.Printf("Server down, repeat in 5s...\n")
				fmt.Printf("FIXME: %s\n", err.Error())
				time.Sleep(5 * time.Second)
				continue sloop
			}

			// Wee need StatusCode, so close Body stream
			res.Body.Close()

			switch res.StatusCode {
			case 507:
				// FIXME body size = 0 after
				fmt.Printf("Insufficient storage, repeat in 10s...\n")
				time.Sleep(10 * time.Second)
			case 200:
				fmt.Printf("Done!\n")
				break sloop
			default:
				fmt.Printf("Unknown error...\n")
				break sloop
			}
		}
	}
}

// GetMessages Ждёт успешной связи с сервером, затем получает все сообщения для указанного владельца
func (mp *MsgProxy) GetMessages(owner string) {
	// On read-only channel will do until empty
	for msg := range mp.receiver(owner) {
		fmt.Printf("* %s\n", msg)
	}
}

// receiver Функция-получатель
func (mp *MsgProxy) receiver(owner string) <-chan string {
	payload := []byte(fmt.Sprintf(`{"person_id":"%s"}`, owner))

	req, err := http.NewRequest(http.MethodGet, serverURL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("proxy: error making a request (%s)", err.Error())
		return nil
	}
	req.Header.Add("Content-Type", "application/json")

	// Buffered channel so will not block
	msgs := make(chan string, 10)

	// Repeatable part, until done
rloop:
	for {
		res, err := mp.connector.Do(req)
		if err != nil {
			fmt.Printf("Server down, repeat in 5s...\n")
			fmt.Printf("FIXME: %s\n", err.Error())
			time.Sleep(5 * time.Second)
			continue rloop
		}

		type JSONResponse struct {
			Messages []string `json:"messages"`
		}

		var payout JSONResponse
		if err := json.NewDecoder(res.Body).Decode(&payout); err != nil {
			fmt.Printf("proxy: json decode error (%s)\n", err.Error())
			break rloop
		}

		// Close Body stream
		res.Body.Close()

		switch res.StatusCode {
		case 200:
			for _, msg := range payout.Messages {
				msgs <- msg
			}
			fmt.Printf("Done!\n")
			break rloop
		default:
			fmt.Printf("Unknown error...\n")
			break rloop
		}
	}

	// No writes allowed for now
	close(msgs)

	return msgs
}
