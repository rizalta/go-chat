package server

import (
	"context"
	"fmt"
	"go-chat/internal/ws"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Server struct {
	port int
	hub  *ws.Hub

	db *redis.Client
}

func NewServer(ctx context.Context, db *redis.Client, hub *ws.Hub) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		hub:  hub,

		db: db,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(ctx),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
