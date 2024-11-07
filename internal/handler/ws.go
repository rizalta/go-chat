package handler

import (
	"go-chat/internal/ws"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub}
}

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		log.Println("no id")
		return
	}
	username, ok := r.Context().Value("username").(string)
	if !ok {
		log.Println("no username")
		return
	}

	client := ws.NewClient(ws.ClientOption{
		UserID:   userID,
		Username: username,
		Conn:     conn,
		Hub:      h.hub,
	})

	go client.ReadPump()
	go client.WritePump()
}
