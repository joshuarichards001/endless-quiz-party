package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		environment := os.Getenv("ENVIRONMENT")
		if environment == "development" {
			return origin == "http://localhost:3000"
		} else {
			return origin == "https://endlessquiz.party"
		}
	},
}

type QuizQuestion struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"`
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		quiz, err := fetchQuestion()
		if err != nil {
			log.Println("Error fetching Gemini question:", err)
			continue
		}
		questionMsg := map[string]interface{}{
			"type":     "question",
			"question": quiz.Question,
			"options":  quiz.Options,
		}
		if err := conn.WriteJSON(questionMsg); err != nil {
			log.Println("Error writing message:", err)
			break
		}
		log.Printf("Sent question: %s", quiz.Question)
	}
}
