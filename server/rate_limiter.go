package main

import (
	"net"
	"os"
	"sync"
	"time"
)

type RateLimiter struct {
	mu             sync.RWMutex
	connections    map[string]int
	messageTimings map[string][]time.Time
	lastCleanup    time.Time

	maxConnectionsPerIP int
	maxMessagesPerMin   int
	cleanupInterval     time.Duration
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		connections:         make(map[string]int),
		messageTimings:      make(map[string][]time.Time),
		lastCleanup:         time.Now(),
		maxConnectionsPerIP: 5,
		maxMessagesPerMin:   32,
		cleanupInterval:     time.Minute,
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

	rl.mu.Lock()
	defer rl.mu.Unlock()

	connectionCount := rl.connections[ip]
	return connectionCount < rl.maxConnectionsPerIP
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
	}
	if rl.connections[ip] == 0 {
		delete(rl.connections, ip)
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
	oneMinuteAgo := now.Add(-time.Minute)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Clean up old message timings
	timings := rl.messageTimings[ip]
	validTimings := make([]time.Time, 0, len(timings))
	for _, timing := range timings {
		if timing.After(oneMinuteAgo) {
			validTimings = append(validTimings, timing)
		}
	}

	// Check if under limit
	if len(validTimings) >= rl.maxMessagesPerMin {
		rl.messageTimings[ip] = validTimings
		return false
	}

	// Add current message timing
	validTimings = append(validTimings, now)
	rl.messageTimings[ip] = validTimings

	// Periodic cleanup
	if now.Sub(rl.lastCleanup) > rl.cleanupInterval {
		go rl.cleanup(now)
		rl.lastCleanup = now
	}

	return true
}

func (rl *RateLimiter) cleanup(now time.Time) {
	oneMinuteAgo := now.Add(-time.Minute)

	for ip, timings := range rl.messageTimings {
		validTimings := make([]time.Time, 0, len(timings))
		for _, timing := range timings {
			if timing.After(oneMinuteAgo) {
				validTimings = append(validTimings, timing)
			}
		}

		if len(validTimings) == 0 {
			delete(rl.messageTimings, ip)
		} else {
			rl.messageTimings[ip] = validTimings
		}
	}
}

func getIPFromAddr(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		if net.ParseIP(remoteAddr) != nil {
			return remoteAddr
		}
		return ""
	}
	return ip
}

func isDevelopmentMode() bool {
	return os.Getenv("ENVIRONMENT") == "development"
}
