package main

import (
	"fmt"
	"go-chat/internal/database"
	"go-chat/internal/server"
)

func main() {
	db := database.New()
	defer db.Close()
	server := server.NewServer(db)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
