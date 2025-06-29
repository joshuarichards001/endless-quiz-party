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
		origin := r.Header.Get("Origin")
		environment := os.Getenv("ENVIRONMENT")
		if environment == "development" {
			return true
		} else {
			return origin == "https://endlessquiz.party"
		}
	},
}

func websocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ip := getIPFromRequest(r)

	if !rateLimiter.AllowConnection(ip) {
		http.Error(w, "Too many connections from this IP", http.StatusTooManyRequests)
		log.Printf("Websocket - Connection rejected due to rate limiting: %s", ip)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	rateLimiter.AddConnection(ip)
	username := GenerateRandomUsername()
	client := NewClient(hub, conn, username, ip)
	hub.Register <- client

	go client.ReadPump()
	go client.WritePump()

	log.Printf("Websocket - New client connected: %s from IP: %s", username, ip)
}
