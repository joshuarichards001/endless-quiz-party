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
		return ip
	}

	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}

	return r.RemoteAddr
}
