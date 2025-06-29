package main

import (
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type messageWindow struct {
	count       int
	windowStart time.Time
}

type RateLimiter struct {
	mu          sync.RWMutex
	connections map[string]int
	messages    map[string]*messageWindow

	maxConnectionsPerIP int
	maxMessagesPerMin   int
	windowDuration      time.Duration
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		connections:         make(map[string]int),
		messages:            make(map[string]*messageWindow),
		maxConnectionsPerIP: 5,
		maxMessagesPerMin:   32,
		windowDuration:      time.Minute,
	}
}

func (rl *RateLimiter) AllowConnection(remoteAddr string) bool {
	if isDevelopmentMode() {
		return true
	}

	ip := getIPFromAddr(remoteAddr)
	if ip == "" {
		return false
	}

	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return rl.connections[ip] < rl.maxConnectionsPerIP
}

func (rl *RateLimiter) AddConnection(remoteAddr string) {
	ip := getIPFromAddr(remoteAddr)
	if ip == "" {
		return
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.connections[ip]++
}

func (rl *RateLimiter) RemoveConnection(remoteAddr string) {
	ip := getIPFromAddr(remoteAddr)
	if ip == "" {
		return
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.connections[ip] > 0 {
		rl.connections[ip]--
		if rl.connections[ip] == 0 {
			delete(rl.connections, ip)
		}
	}
}

func (rl *RateLimiter) AllowMessage(remoteAddr string) bool {
	if isDevelopmentMode() {
		return true
	}

	ip := getIPFromAddr(remoteAddr)
	if ip == "" {
		return false
	}

	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	window := rl.messages[ip]

	// Reset window if expired or doesn't exist
	if window == nil || now.Sub(window.windowStart) >= rl.windowDuration {
		rl.messages[ip] = &messageWindow{
			count:       1,
			windowStart: now,
		}
		return true
	}

	// Check if over limit
	if window.count >= rl.maxMessagesPerMin {
		return false
	}

	// Allow message
	window.count++
	return true
}

func getIPFromAddr(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		if net.ParseIP(remoteAddr) != nil {
			return remoteAddr
		}
		return ""
	}

	log.Print("Remote address:", remoteAddr)
	log.Print("IP:", ip)

	return ip
}

func isDevelopmentMode() bool {
	return os.Getenv("ENVIRONMENT") == "development"
}
