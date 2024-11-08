package ws

import (
	"context"
	"go-chat/internal/database"
	"go-chat/internal/domain"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan domain.Message
	mutex      sync.Mutex
	repo       *database.MessageRepo
}

func NewHub(db *redis.Client) *Hub {
	repo := database.NewMessageRepo(db)
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan domain.Message),
		mutex:      sync.Mutex{},
		repo:       repo,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.userID]; ok {
				close(client.send)
				delete(h.clients, client.userID)
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.Lock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.userID)
				}
			}
			h.repo.AddMessage(ctx, message)
			h.mutex.Unlock()
		}
	}
}
