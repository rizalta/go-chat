package main

import (
	"context"
	"go-chat/internal/database"
	"go-chat/internal/server"
	"go-chat/internal/ws"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	db, err := database.New(ctx)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	messageRepo := database.NewMessageRepo(db)
	userRepo := database.NewUserRepo(db)
	hub := ws.NewHub(messageRepo)
	go hub.Run(ctx)

	server := server.NewServer(ctx, server.ServerParams{
		UserRepo:    userRepo,
		MessageRepo: messageRepo,
		Hub:         hub,
	})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gracefully..")
		cancel()
		timeout, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelTimeout()

		if err := server.Shutdown(timeout); err != nil {
			log.Println("error during shutdown:", err)
		}
	}()

	log.Println("Runnning server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start server: %s", err)
	}

	<-ctx.Done()
	log.Println("Server stopped")
}
