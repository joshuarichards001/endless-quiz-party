package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Utils - error loading .env file")
	}
}

func getIPFromRequest(r *http.Request) string {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		log.Println("Utils - IP extracted from X-Real-Ip header:", ip)
		return ip
	}

	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		log.Println("Utils - IP extracted from X-Forwarded-For header:", strings.Split(ip, ",")[0])
		return strings.Split(ip, ",")[0]
	}

	log.Println("Utils - IP extracted from RemoteAddr:", r.RemoteAddr)
	return r.RemoteAddr
}
