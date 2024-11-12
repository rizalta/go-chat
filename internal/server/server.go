package server

import (
	"context"
	"fmt"
	"go-chat/internal/database"
	"go-chat/internal/ws"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port int
	hub  *ws.Hub

	userRepo    *database.UserRepo
	messageRepo *database.MessageRepo
}

type ServerParams struct {
	UserRepo    *database.UserRepo
	MessageRepo *database.MessageRepo
	Hub         *ws.Hub
}

func NewServer(ctx context.Context, params ServerParams) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:        port,
		messageRepo: params.MessageRepo,
		userRepo:    params.UserRepo,
		hub:         params.Hub,
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
