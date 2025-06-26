package main

import (
	"log"
	"net/http"
)

var rateLimiter *RateLimiter

func main() {
	loadEnv()

	rateLimiter = NewRateLimiter()
	hub := NewHub(nil)
	quizManager := NewQuizManager(hub, fetchQuestion)
	hub.QuizManager = quizManager

	go hub.Run()
	go quizManager.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(hub, w, r)
	})

	log.Println("Main - Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
