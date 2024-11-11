package ws

import (
	"bytes"
	"context"
	"fmt"
	"go-chat/cmd/web/components"
	"go-chat/internal/domain"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	userID   string
	username string
	conn     *websocket.Conn
	hub      *Hub
	mu       sync.Mutex
	send     chan domain.Message
	closed   bool
}

type ClientOption struct {
	UserID   string
	Username string
	Conn     *websocket.Conn
	Hub      *Hub
}

func NewClient(ctx context.Context, opts ClientOption) *Client {
	if opts.Conn == nil || opts.Hub == nil {
		return nil
	}
	c := &Client{
		userID:   opts.UserID,
		username: opts.Username,
		conn:     opts.Conn,
		hub:      opts.Hub,
		mu:       sync.Mutex{},
		send:     make(chan domain.Message, 256),
		closed:   false,
	}

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPingHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	c.hub.register <- c

	go c.loadHistoricalMessages(ctx)

	return c
}

func (c *Client) loadHistoricalMessages(ctx context.Context) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	messages, err := c.hub.repo.GetAllMessages(ctx)
	if err != nil {
		log.Printf("Error loading historical messages, %v", err)
		return
	}

	for _, m := range messages {
		select {
		case c.send <- *m:
		default:
			log.Printf("Failed to send historical message, channel full")
			return
		}
	}
}

type WSMessage struct {
	Message string      `json:"message"`
	Headers interface{} `json:"HEADERS"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.Close()
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
				log.Printf("read error: %v", err)
			}
			break
		}

		message := strings.TrimSpace(msg.Message)
		if message == "" {
			continue
		}

		if len(message) > maxMessageSize {
			log.Printf("Message too large from user %s", c.userID)
			continue
		}

		newMsg := domain.Message{
			UserID:    c.userID,
			Username:  c.username,
			Content:   msg.Message,
			TimeStamp: time.Now(),
		}

		select {
		case c.hub.broadcast <- newMsg:
		default:
			log.Printf("failed to broadcast message, channel full")
		}
	}
}

func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.writeMessage(websocket.CloseMessage, []byte{})
				return
			}
			buf := &bytes.Buffer{}
			err := components.Message(msg.Content, msg.Username, msg.TimeStamp.Format(time.RFC3339), c.userID == msg.UserID).
				Render(ctx, buf)
			if err != nil {
				log.Printf("error rendering message: %v", err)
			}
			err = c.conn.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				log.Printf("error writing message from user %s: %v", c.userID, err)
				return
			}
		case <-ticker.C:
			if err := c.writeMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) writeMessage(messageType int, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("connection closed")
	}

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	return c.conn.WriteMessage(messageType, payload)
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		c.closed = true
		c.hub.unregister <- c
		close(c.send)
		c.conn.Close()
	}
}
