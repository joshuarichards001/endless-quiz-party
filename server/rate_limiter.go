package main

import (
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

func (rl *RateLimiter) AllowConnection(ip string) bool {
	if isDevelopmentMode() {
		return true
	}

	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return rl.connections[ip] < rl.maxConnectionsPerIP
}

func (rl *RateLimiter) AddConnection(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.connections[ip]++
}

func (rl *RateLimiter) RemoveConnection(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.connections[ip] > 0 {
		rl.connections[ip]--
		if rl.connections[ip] == 0 {
			delete(rl.connections, ip)
		}
	}
}

func (rl *RateLimiter) AllowMessage(ip string) bool {
	if isDevelopmentMode() {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

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

func isDevelopmentMode() bool {
	return os.Getenv("ENVIRONMENT") == "development"
}
