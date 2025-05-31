package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type StockPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Time   string  `json:"time"`
}

func fetchAppleStockPrice() StockPrice {
	return StockPrice{
		Symbol: "AAPL",
		Price:  150.00 + rand.Float64()*(200.00-150.00),
		Time:   time.Now().Format(time.RFC3339),
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		price := fetchAppleStockPrice()
		message, err := json.Marshal(price)
		if err != nil {
			log.Println("Error marshalling json:", err)
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			break
		}
		log.Printf("Sent: %s", message)
	}
}

func main() {
	http.HandleFunc("/ws", websocketHandler)
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}