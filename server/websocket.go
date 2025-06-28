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
			// return origin == "https://endlessquiz.party"
			return true
		}
	},
}

func websocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	if !rateLimiter.AllowConnection(r.RemoteAddr) {
		http.Error(w, "Too many connections from this IP", http.StatusTooManyRequests)
		log.Printf("Websocket - Connection rejected due to rate limiting: %s", r.RemoteAddr)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	rateLimiter.AddConnection(r.RemoteAddr)

	username := GenerateRandomUsername()
	client := NewClient(hub, conn, username, r.RemoteAddr)
	hub.Register <- client

	go client.ReadPump()
	go client.WritePump()

	log.Printf("Websocket - New client connected: %s", client.Name)
}
