package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func getTriviaApiUrl() (string, error) {
	apiKey := os.Getenv("TRIVIA_API_URL")
	if apiKey == "" {
		return "", fmt.Errorf("TRIVIA_API_URL not set in environment")
	}
	return apiKey, nil
}

func fetchQuestion(hub *Hub) (*Question, error) {
	triviaApiUrl, err := getTriviaApiUrl()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", triviaApiUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var apiResp []struct {
		CorrectAnswer    string   `json:"correctAnswer"`
		IncorrectAnswers []string `json:"incorrectAnswers"`
		Question         struct {
			Text string `json:"text"`
		} `json:"question"`
	}
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("could not parse question JSON: %v", err)
	}

	if len(apiResp) == 0 {
		return nil, fmt.Errorf("no results found in trivia API response")
	}

	result := apiResp[0]

	decodedQuestion := html.UnescapeString(result.Question.Text)
	decodedCorrectAnswer := html.UnescapeString(result.CorrectAnswer)
	decodedIncorrectAnswers := make([]string, len(result.IncorrectAnswers))
	for i, incorrectAnswer := range result.IncorrectAnswers {
		decodedIncorrectAnswers[i] = html.UnescapeString(incorrectAnswer)
	}

	options := append([]string{decodedCorrectAnswer}, decodedIncorrectAnswers...)

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	correctIndex := 0
	for i, opt := range options {
		if opt == decodedCorrectAnswer {
			correctIndex = i
			break
		}
	}

	return &Question{
		Question: decodedQuestion,
		Options:  options,
		Answer:   correctIndex,
	}, nil
}
