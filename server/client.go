package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Hub           *Hub
	Conn          *websocket.Conn
	Send          chan []byte
	Name          string
	CurrentAnswer int
	Streak        int
}

func NewClient(hub *Hub, conn *websocket.Conn, name string) *Client {
	return &Client{
		Hub:           hub,
		Conn:          conn,
		Send:          make(chan []byte, 256),
		Name:          name,
		CurrentAnswer: -1,
		Streak:        0,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var incomingMessage SubmitAnswerMessage
		if err := json.Unmarshal(message, &incomingMessage); err != nil {
			log.Println("Client - Error unmarshaling message:", err, "Raw message:", string(message))
		}

		if incomingMessage.Type == MessageTypeSubmitAnswer {
			userAnswer := UserAnswer{
				Client: c,
				Answer: incomingMessage.Answer,
			}
			c.Hub.ProcessAnswer <- userAnswer
		} else {
			log.Println("Client - Received unknown message type after successful unmarshal:", incomingMessage.Type)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Client - Error writing message:", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Client - Error writing ping message:", err)
				return
			}
		}
	}
}
