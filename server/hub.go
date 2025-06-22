package main

import (
	"encoding/json"
	"sort"
	"sync/atomic"
	"time"
)

type UserAnswer struct {
	Client *Client
	Answer int
}

type ProcessResultsRequest struct {
	Answer         int
	Votes          map[int]int
	RevealDuration int
}

type Hub struct {
	Clients        map[*Client]bool
	Broadcast      chan []byte
	Register       chan *Client
	Unregister     chan *Client
	ProcessAnswer  chan UserAnswer
	ProcessResults chan ProcessResultsRequest
	QuizManager    *QuizManager
	UserCount      int32
}

func NewHub(quizManager *QuizManager) *Hub {
	return &Hub{
		Clients:        make(map[*Client]bool),
		Broadcast:      make(chan []byte, 512),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		ProcessAnswer:  make(chan UserAnswer),
		ProcessResults: make(chan ProcessResultsRequest),
		QuizManager:    quizManager,
		UserCount:      0,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}

		case client := <-h.Register:
			h.Clients[client] = true
			atomic.AddInt32(&h.UserCount, 1)

			welcomeMessage := WelcomeMessage{
				Type:     MessageTypeWelcome,
				Username: client.Name,
			}
			welcomeBytes, _ := json.Marshal(welcomeMessage)
			client.Send <- welcomeBytes

			if h.QuizManager.IsQuestionActive {
				questionMessage := QuestionMessage{
					Type:      "question",
					Question:  h.QuizManager.CurrentQuestion.Question,
					Options:   h.QuizManager.CurrentQuestion.Options,
					TimeLeft:  int(time.Until(h.QuizManager.QuestionEndTime).Seconds()),
					UserCount: int(atomic.LoadInt32(&h.UserCount)),
				}
				messageBytes, _ := json.Marshal(questionMessage)
				client.Send <- messageBytes
			} else {
				answerResultMsg := AnswerResultMessage{
					Type:              "answer_result",
					CorrectAnswer:     h.QuizManager.CurrentQuestion.Answer,
					Votes:             h.QuizManager.CurrentVotes,
					YourAnswerCorrect: false,
					CurrentStreak:     0,
					UserCount:         int(atomic.LoadInt32(&h.UserCount)),
				}
				messageBytes, _ := json.Marshal(answerResultMsg)
				client.Send <- messageBytes
			}

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				close(client.Send)
				delete(h.Clients, client)
				atomic.AddInt32(&h.UserCount, -1)
			}

		case answer := <-h.ProcessAnswer:
			answer.Client.CurrentAnswer = answer.Answer
			h.QuizManager.RecordVote(answer.Answer)

		case results := <-h.ProcessResults:
			for client := range h.Clients {
				if client.CurrentAnswer == results.Answer {
					client.Streak++
				} else {
					client.Streak = 0
				}
			}

			userCount := int(atomic.LoadInt32(&h.UserCount))
			leaderboard := h.getLeaderboard()

			for client := range h.Clients {
				answerResultMsg := AnswerResultMessage{
					Type:              "answer_result",
					CorrectAnswer:     results.Answer,
					Votes:             results.Votes,
					YourAnswerCorrect: client.CurrentAnswer == results.Answer,
					CurrentStreak:     client.Streak,
					UserCount:         userCount,
					Leaderboard:       leaderboard,
				}

				messageBytes, err := json.Marshal(answerResultMsg)
				if err != nil {
					continue
				}

				select {
				case client.Send <- messageBytes:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
				client.CurrentAnswer = -1
			}
		}
	}
}

func (h *Hub) broadcastUserCount() {
	countMsg := UserCountUpdateMessage{Type: "user_count", Count: int(atomic.LoadInt32(&h.UserCount))}
	msgBytes, _ := json.Marshal(countMsg)
	h.Broadcast <- msgBytes
}

func (h *Hub) BroadcastMessage(message []byte) {
	h.Broadcast <- message
}

func (h *Hub) getLeaderboard() []LeaderboardEntry {
	var leaderboard []LeaderboardEntry
	clientStreaks := make([]*Client, 0, len(h.Clients))
	
	for client := range h.Clients {
		clientStreaks = append(clientStreaks, client)
	}

	sort.Slice(clientStreaks, func(i, j int) bool {
		return clientStreaks[i].Streak > clientStreaks[j].Streak
	})

	maxEntries := min(len(clientStreaks), 5)

	for i := range maxEntries {
		leaderboard = append(leaderboard, LeaderboardEntry{
			Username: clientStreaks[i].Name,
			Streak:   clientStreaks[i].Streak,
			Rank:     i + 1,
		})
	}

	return leaderboard
}
