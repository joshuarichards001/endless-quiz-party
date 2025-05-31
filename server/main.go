package main

import (
	"log"
	"net/http"
)

func main() {
	loadEnv()
	http.HandleFunc("/ws", websocketHandler)
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
