package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type MsgDialog struct {
	conn  *http.Client
	input *bufio.Reader

	PersonID string
	Choice   string
}

func NewMsgDialog(connection *http.Client) *MsgDialog {
	return &MsgDialog{
		conn:  connection,
		input: bufio.NewReader(os.Stdin),
	}
}

//TODO Clear очищает экран, понадобится?
//func (d *MsgDialog) Clear() {
//	switch runtime.GOOS {
//	case "linux":
//		cmd := exec.Command("clear")
//		cmd.Stdout = os.Stdout
//		cmd.Run()
//	case "windows":
//		cmd := exec.Command("cmd", "/c", "cls")
//		cmd.Stdout = os.Stdout
//		cmd.Run()
//	}
//}

func (d *MsgDialog) Identify() {
	fmt.Print("Name yourself: ")
	d.PersonID, _ = d.input.ReadString('\n')
}

func (d *MsgDialog) Choose() {
	menu := "\n" +
		"[1] Add new message\n" +
		"[2] Get all messages\n" +
		"[Q] Quit\n" +
		"Choose: "

	fmt.Printf(menu)
	if str, _ := d.input.ReadString('\n'); len(str) > 0 {
		d.Choice = strings.ToUpper(string(str[0]))
	} else {
		d.Choice = ""
	}

	switch d.Choice {
	case "1":
		// TODO
		fmt.Println("< 1")
	case "2":
		// TODO
		fmt.Println("< 2")
	}
}

func (d *MsgDialog) communicate(message string) ([]byte, error) {
	// Testing
	res, err := http.Get("http://localhost:8080/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	return body, nil
}

func (d *MsgDialog) putMessage(text string) error {
	// TODO store in dialog queue, communicate with server, handle errors
	return nil
}

func (d *MsgDialog) getAllMesages() ([]string, error) {
	// TODO communicate with server, return messages, handle errors
	var messages []string

	return messages, nil
}
