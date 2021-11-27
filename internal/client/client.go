package client

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Start() error {
	// Placeholder init
	time.Sleep(10 * time.Millisecond)

	d := NewMsgDialog(c.httpClient)

	// Dialogue loop
dloop:
	for {
		switch {
		case d.PersonID == "":
			d.Identify()
		case d.Choice != "Q":
			d.Choose()
		default:
			break dloop
		}
	}

	return ErrClientClosed
}

var ErrClientClosed = errors.New("client: quit")

func (c *Client) Stop(ctx context.Context) error {
	c.httpClient.CloseIdleConnections()

	//WARNING: ctx not used because of simple net/http client implementation, so here is a placeholder
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// if ctx has not expired yet, we're all done
		return nil
	}
}
