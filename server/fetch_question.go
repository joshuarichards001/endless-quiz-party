package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	client := &http.Client{}
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

	var apiResp struct {
		ResponseCode int `json:"response_code"`
		Results      []struct {
			Question         string   `json:"question"`
			CorrectAnswer    string   `json:"correct_answer"`
			IncorrectAnswers []string `json:"incorrect_answers"`
		} `json:"results"`
	}
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("could not parse question JSON: %v", err)
	}

	result := apiResp.Results[0]
	options := append([]string{result.CorrectAnswer}, result.IncorrectAnswers...)

	return &Question{
		Question: result.Question,
		Options:  options,
		Answer:   0,
	}, nil
}
