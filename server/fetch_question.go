package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func cleanCodeBlock(s string) string {
	re := regexp.MustCompile("(?s)```[a-zA-Z]*\\n(.*)\\n```")
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	return strings.Trim(s, "`")
}

func getOpenAIAPIKey() (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set in environment")
	}
	return apiKey, nil
}

func fetchQuestion() (*QuizQuestion, error) {
	apiKey, err := getOpenAIAPIKey()
	if err != nil {
		return nil, err
	}

	var catKeys []string
	for k := range categories {
		catKeys = append(catKeys, k)
	}
	category := catKeys[rand.Intn(len(catKeys))]
	subcats := categories[category]
	subcategory := subcats[rand.Intn(len(subcats))]

	url := "https://api.openai.com/v1/chat/completions"

	systemPromptBytes, err := os.ReadFile("system_prompt.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read system_prompt.txt: %v", err)
	}
	systemPrompt := string(systemPromptBytes)
	userPrompt := fmt.Sprintf("Generate a new question in the category '%s' and sub-category '%s'.", category, subcategory)

	requestBody := map[string]interface{}{
		"model": "gpt-4.1-mini",
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.7,
	}

	jsonBody, _ := json.Marshal(requestBody)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(body, &openAIResp)
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAI API response: %v", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("invalid OpenAI API response: missing question")
	}

	questionText := cleanCodeBlock(openAIResp.Choices[0].Message.Content)

	var q struct {
		Question string   `json:"Question"`
		Options  []string `json:"Options"`
		Answer   string   `json:"Answer"`
	}
	err = json.Unmarshal([]byte(questionText), &q)
	if err != nil {
		return nil, fmt.Errorf("could not parse question JSON: %v", err)
	}

	answerIdx := -1
	for i, opt := range q.Options {
		if opt == q.Answer {
			answerIdx = i
			break
		}
	}
	if answerIdx == -1 {
		return nil, fmt.Errorf("answer not found in options")
	}

	return &QuizQuestion{
		Question: q.Question,
		Options:  q.Options,
		Answer:   answerIdx,
	}, nil
}
