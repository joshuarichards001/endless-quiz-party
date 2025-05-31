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
	Answer   int      `json:"answer"` // index of correct option
}

// Question represents a quiz question
var questions = []QuizQuestion{
	{
		Question: "What is the capital of France?",
		Options:  []string{"Paris", "London", "Berlin", "Madrid"},
		Answer:   0,
	},
	{
		Question: "What is 2 + 2?",
		Options:  []string{"3", "4", "5", "6"},
		Answer:   1,
	},
	{
		Question: "What is the largest planet in our solar system?",
		Options:  []string{"Earth", "Mars", "Jupiter", "Saturn"},
		Answer:   2,
	},
	{
		Question: "Who wrote 'To Kill a Mockingbird'?",
		Options:  []string{"Harper Lee", "Mark Twain", "Jane Austen", "Ernest Hemingway"},
		Answer:   0,
	},
	{
		Question: "What is the boiling point of water in Celsius?",
		Options:  []string{"90", "100", "110", "120"},
		Answer:   1,
	},
	{
		Question: "What is the chemical symbol for gold?",
		Options:  []string{"Au", "Ag", "Gd", "Go"},
		Answer:   0,
	},
	{
		Question: "Who painted the Mona Lisa?",
		Options:  []string{"Vincent van Gogh", "Pablo Picasso", "Leonardo da Vinci", "Claude Monet"},
		Answer:   2,
	},
	{
		Question: "What is the smallest prime number?",
		Options:  []string{"0", "1", "2", "3"},
		Answer:   2,
	},
	{
		Question: "What year did the Titanic sink?",
		Options:  []string{"1910", "1912", "1914", "1920"},
		Answer:   1,
	},
	{
		Question: "What is the square root of 64?",
		Options:  []string{"6", "7", "8", "9"},
		Answer:   2,
	},
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// We'll keep a simple in-memory list of users
var users = make(map[string]User)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	questionIndex := 0
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Send the first question immediately
	questionMsg := map[string]interface{}{
		"type":     "question",
		"question": questions[questionIndex].Question,
		"options":  questions[questionIndex].Options,
	}
	if err := conn.WriteJSON(questionMsg); err != nil {
		log.Println("Error writing message:", err)
		return
	}

	for range ticker.C {
		questionIndex = (questionIndex + 1) % len(questions)
		questionMsg := map[string]interface{}{
			"type":     "question",
			"question": questions[questionIndex].Question,
			"options":  questions[questionIndex].Options,
		}
		if err := conn.WriteJSON(questionMsg); err != nil {
			log.Println("Error writing message:", err)
			break
		}
		log.Printf("Sent question: %s", questions[questionIndex].Question)
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