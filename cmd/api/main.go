package main

import (
	"fmt"
	"go-chat/internal/database"
	"go-chat/internal/server"
	"go-chat/internal/ws"
)

func main() {
	db := database.New()
	defer db.Close()

	hub := ws.NewHub()
	go hub.Run()

	server := server.NewServer(db, hub)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
