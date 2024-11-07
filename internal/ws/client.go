package ws

import (
	"bytes"
	"context"
	"fmt"
	"go-chat/cmd/web/components"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	userID   string
	username string
	conn     *websocket.Conn
	hub      *Hub
	send     chan Message
}

type ClientOption struct {
	UserID   string
	Username string
	Conn     *websocket.Conn
	Hub      *Hub
}

func NewClient(opts ClientOption) *Client {
	c := &Client{
		userID:   opts.UserID,
		username: opts.Username,
		conn:     opts.Conn,
		hub:      opts.Hub,
		send:     make(chan Message),
	}
	opts.Hub.register <- c

	return c
}

type WSMessage struct {
	Message string      `json:"message"`
	Headers interface{} `json:"HEADERS"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg WSMessage
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- Message{
			UserID:    c.userID,
			Username:  c.username,
			Content:   msg.Message,
			TimeStamp: time.Now(),
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		msg, ok := <-c.send
		if !ok {
			log.Println("write error")
			return
		}
		fmt.Println(msg)

		var buf bytes.Buffer
		components.Message(msg.Content, msg.Username, msg.TimeStamp.Format(time.RFC3339)).
			Render(context.Background(), &buf)
		err := c.conn.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			return
		}
	}
}
