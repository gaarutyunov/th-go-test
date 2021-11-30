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
	// Unbuffered (blocking) channel
	mp.messages = make(chan Message)

	return mp
}

// AddMessage Добавляет сообщение в прокси-очередь сообщений с сохранением владельца, ждёт отправки на сервер
func (mp *MsgProxy) AddMessage(text, owner string) {
	// Write-only channel in parallel goroutine
	go func(c chan<- Message) {
		c <- Message{
			PersonID: owner,
			Message:  text,
		}
	}(mp.messages)

	mp.sender()
}

// sender Функция-отправитель
func (mp *MsgProxy) sender() {
	// Fetch a message
	select {
	case payload, ok := <-mp.messages:
		if !ok {
			fmt.Printf("proxy: error, channel closed")
			break
		}

		// Repeatable part, until server done
	sloop:
		for {
			payloadBuf := new(bytes.Buffer)
			json.NewEncoder(payloadBuf).Encode(&payload)

			req, err := http.NewRequest(http.MethodPut, serverURL, payloadBuf)
			if err != nil {
				fmt.Printf("proxy: error making a request (%s)\n", err.Error())
				return
			}
			req.Header.Add("Content-Type", "application/json")

			res, err := mp.connector.Do(req)
			if err != nil {
				countdown("Server down, repeat in", 5*time.Second)
				continue sloop
			}

			// Wee need StatusCode, so close Body stream
			res.Body.Close()

			switch res.StatusCode {
			case 507:
				countdown("Insufficient storage, repeat in", 5*time.Second)
				time.Sleep(5 * time.Second)
			case 200:
				fmt.Printf("Done!\n")
				break sloop
			default:
				fmt.Printf("Unknown error...%s\n", res.Status)
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

	// Buffered (non blocking) channel
	msgs := make(chan string, 10)

	// Repeatable part, until server done
rloop:
	for {
		payloadBuf := bytes.NewBuffer(payload)

		req, err := http.NewRequest(http.MethodGet, serverURL, payloadBuf)
		if err != nil {
			fmt.Printf("proxy: error making a request (%s)", err.Error())
			return nil
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := mp.connector.Do(req)
		if err != nil {
			countdown("Server down, repeat in", 5*time.Second)
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
			fmt.Printf("Unknown error...%s\n", res.Status)
			break rloop
		}
	}

	// No writes allowed for now
	close(msgs)

	// Read-only channel in return
	return msgs
}

// countdown Обратный отсчёт по секундам на той же строке
func countdown(prompt string, d time.Duration) {
	for range time.Tick(1 * time.Second) {
		fmt.Printf("\r%s %s", prompt, d)
		if d.Milliseconds() <= 0 {
			break
		}
		d -= time.Second
	}

	fmt.Printf("\r%s\r", bytes.Repeat([]byte{' '}, 40))
}
