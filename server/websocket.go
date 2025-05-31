package main

import (
	"log"
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

type QuizQuestion struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = make(map[string]User)

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

	// Send the first question immediately
	quizQ, err := fetchGeminiQuestion()
	if err == nil {
		questionMsg := map[string]interface{}{
			"type":     "question",
			"question": quizQ.Question,
			"options":  quizQ.Options,
		}
		conn.WriteJSON(questionMsg)
	}

	for range ticker.C {
		quizQ, err := fetchGeminiQuestion()
		if err != nil {
			log.Println("Error fetching Gemini question:", err)
			continue
		}
		questionMsg := map[string]interface{}{
			"type":     "question",
			"question": quizQ.Question,
			"options":  quizQ.Options,
		}
		if err := conn.WriteJSON(questionMsg); err != nil {
			log.Println("Error writing message:", err)
			break
		}
		log.Printf("Sent question: %s", quizQ.Question)
	}
}
