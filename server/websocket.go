package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		environment := os.Getenv("ENVIRONMENT")
		if environment == "development" {
			return true
		} else {
			return true
		}
	},
}

func websocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	username := GenerateRandomUsername()
	client := NewClient(hub, conn, username)
	hub.Register <- client

	go client.ReadPump()
	go client.WritePump()

	log.Printf("Websocket - New client connected: %s", client.Name)
}
