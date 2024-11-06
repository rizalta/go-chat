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
	userID string
	conn   *websocket.Conn
	hub    *Hub
	send   chan Message
}

func NewClient(userID string, conn *websocket.Conn, hub *Hub) *Client {
	c := &Client{
		userID: userID,
		conn:   conn,
		hub:    hub,
		send:   make(chan Message),
	}
	hub.register <- c

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
			Username:  c.userID[:7],
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
		components.Message(msg.Username, msg.Content).Render(context.Background(), &buf)
		err := c.conn.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			return
		}
	}
}
