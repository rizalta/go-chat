package main

import (
	"context"
	"fmt"
	"go-chat/internal/database"
	"go-chat/internal/server"
	"go-chat/internal/ws"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := database.New()
	defer db.Close()

	ctx := context.Background()

	hub := ws.NewHub(db)
	go hub.Run(ctx)

	server := server.NewServer(db, hub)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
