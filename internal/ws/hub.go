package ws

import (
	"context"
	"go-chat/internal/database"
	"go-chat/internal/domain"
	"log"
	"sync"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan domain.Message
	mutex      sync.RWMutex
	repo       *database.MessageRepo
	done       chan struct{}
}

func NewHub(repo *database.MessageRepo) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan domain.Message, 256),
		mutex:      sync.RWMutex{},
		repo:       repo,
		done:       make(chan struct{}),
	}
}

func (h *Hub) Run(ctx context.Context) {
	defer h.shutdown()

	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client
			h.mutex.Unlock()
			log.Printf("client registered: %s", client.userID)
		case client := <-h.unregister:
			h.mutex.Lock()
			if existingClient, ok := h.clients[client.userID]; ok {
				existingClient.Close()
				delete(h.clients, client.userID)
				log.Printf("client unregistered: %s", client.userID)
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			err := h.repo.AddMessage(ctx, message)
			if err != nil {
				log.Printf("error saving message %v", err)
				continue
			}
			h.mutex.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					go func(c *Client) {
						c.Close()
						h.unregister <- c
					}(client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) shutdown() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, client := range h.clients {
		client.Close()
	}

	h.clients = make(map[string]*Client)

	close(h.done)
	close(h.broadcast)
	close(h.register)
	close(h.unregister)
}
