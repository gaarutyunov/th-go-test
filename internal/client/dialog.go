package client

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"th-go-test/pkg/msgproxy"
)

type MsgDialog struct {
	connector *http.Client
	proxy     *msgproxy.MsgProxy
	scanner   *bufio.Scanner
	PersonID  string
	Choice    string
}

func NewMsgDialog(connector *http.Client) *MsgDialog {
	px := msgproxy.New(connector)
	in := bufio.NewScanner(os.Stdin)

	return &MsgDialog{
		connector: connector,
		proxy:     px,
		scanner:   in,
	}
}

func (d *MsgDialog) Identify() {
	fmt.Print("Name yourself: ")
	d.PersonID = d.input()
}

func (d *MsgDialog) Choose() {
	menu := "\n" +
		"[1] Add new message\n" +
		"[2] Get all messages\n" +
		"[Q] Quit\n" +
		"Choose: "

	fmt.Printf(menu)
	d.Choice = d.input()

	switch d.Choice {
	case "1":
		fmt.Printf("Enter a message: ")
		text := d.input()
		d.proxy.AddMessage(text, d.PersonID)
	case "2":
		d.proxy.GetMessages(d.PersonID)
	}
}

func (d *MsgDialog) input() string {
	d.scanner.Scan()
	return d.scanner.Text()
}
