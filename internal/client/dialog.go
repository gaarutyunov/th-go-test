package client

import (
	"fmt"
	"net/http"
	"th-go-test/pkg/msgproxy"
)

type MsgDialog struct {
	connector *http.Client
	proxy     *msgproxy.MsgProxy
	PersonID  string
	Choice    string
}

func NewMsgDialog(connector *http.Client) *MsgDialog {
	px := msgproxy.New(connector)
	return &MsgDialog{
		connector: connector,
		proxy:     px,
	}
}

func (d *MsgDialog) Identify() {
	fmt.Print("Name yourself: ")

	// TODO fix for Macs
	fmt.Scanln(&d.PersonID)
}

func (d *MsgDialog) Choose() {
	menu := "\n" +
		"[1] Add new message\n" +
		"[2] Get all messages\n" +
		"[Q] Quit\n" +
		"Choose: "

	fmt.Printf(menu)

	// TODO fix for Macs
	fmt.Scanln(&d.Choice)

	switch d.Choice {
	case "1":
		fmt.Printf("Enter a message: ")
		// TODO fix for Macs
		var text string
		fmt.Scanln(&text)
		d.proxy.AddMessage(text, d.PersonID)
	case "2":
		d.proxy.GetMessages(d.PersonID)
	}
}
