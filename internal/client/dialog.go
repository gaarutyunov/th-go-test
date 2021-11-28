package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type MsgDialog struct {
	conn     *http.Client
	PersonID string
	Choice   string
}

const serverURL = "http://localhost:8080/"

func NewMsgDialog(connection *http.Client) *MsgDialog {
	return &MsgDialog{
		conn: connection,
	}
}

func (d *MsgDialog) Identify() {
	fmt.Print("Name yourself: ")
	fmt.Scanln(&d.PersonID)
}

func (d *MsgDialog) Choose() {
	menu := "\n" +
		"[1] Add new message\n" +
		"[2] Get all messages\n" +
		"[Q] Quit\n" +
		"Choose: "

	fmt.Printf(menu)
	fmt.Scanln(&d.Choice)

	switch d.Choice {
	case "1":
		fmt.Printf("Enter a message: ")
		var text string
		fmt.Scanln(&text)
		if err := d.putMessage(text); err != nil {
			fmt.Printf("%s\n\n", err.Error())
		}
	case "2":
		msgs, err := d.getAllMesages()
		if err != nil {
			fmt.Printf("%s\n\n", err.Error())
		}
		for idx, msg := range msgs {
			fmt.Printf("[%2d]: %s\n", idx, msg)
		}
	}
}

var ErrInsufficientStorage = errors.New("insufficient storage")
var ErrServerDown = errors.New("server down")

func (d *MsgDialog) putMessage(text string) error {
	type JSONRequest struct {
		PersonID string `json:"person_id"`
		Message  string `json:"message"`
	}

	jr := &JSONRequest{
		PersonID: d.PersonID,
		Message:  text,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(jr)
	req, err := http.NewRequest(http.MethodPut, serverURL, payloadBuf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := d.conn.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 507:
		return ErrInsufficientStorage
	case 200:
		return nil
	default:
		return ErrServerDown
	}
}

func (d *MsgDialog) getAllMesages() ([]string, error) {
	// TODO communicate with server, return messages, handle errors
	type JSONRequest struct {
		PersonID string `json:"person_id"`
	}

	jr := &JSONRequest{
		PersonID: d.PersonID,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(jr)
	req, err := http.NewRequest(http.MethodGet, serverURL, payloadBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := d.conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	type JSONResponse struct {
		Messages []string `json:"messages"`
	}

	var jw JSONResponse
	if err := json.NewDecoder(res.Body).Decode(&jw); err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 200:
		return jw.Messages, nil
	default:
		return nil, ErrServerDown
	}
}
