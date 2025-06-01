package main

import (
	"log"
	"net/http"
)

func main() {
	loadEnv()

	hub := NewHub(nil)
	quizManager := NewQuizManager(hub, fetchQuestion)
	hub.QuizManager = quizManager

	go hub.Run()
	go quizManager.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(hub, w, r)
	})

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
