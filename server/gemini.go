package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// GeminiAPIResponse is the structure for the response from the Gemini API
type GeminiAPIResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// cleanCodeBlock removes Markdown code block markers from a string.
func cleanCodeBlock(s string) string {
	re := regexp.MustCompile("(?s)```[a-zA-Z]*\\n(.*)\\n```")
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		return matches[1]
	}
	// Remove single-line backticks if present
	return strings.Trim(s, "`")
}

// getGeminiAPIKey loads the Gemini API key from the environment
func getGeminiAPIKey() (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set in environment")
	}
	return apiKey, nil
}

// fetchGeminiQuestion fetches a new question from the Gemini API
func fetchGeminiQuestion() (*QuizQuestion, error) {
	apiKey, err := getGeminiAPIKey()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	systemPrompt := "You are a game show host for a general knowledge quiz show similar to 'Who Wants to Be a Millionaire'.\nYour role is to generate engaging and challenging multiple-choice questions for the contestants.\nEach question should come with four possible answers, only one of which is correct.\nThe questions can span a range of categories, and subcategories.\nEach reply should be in a json format, like this:\n{\n\"Question\": \"question\",\n\"Options\": [\"Option 1\", \"Option 2\", \"Option 3\", \"Option 4\"],\n\"Answer\": \"Option x\"\n}\nYou only reply with one question each time.\nIf you use apostrophes, make sure to escape them with a backslash, like this: \\'.\n"
	userPrompt := "Now generate a new question"

	fullPrompt := systemPrompt + "\n" + userPrompt

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": fullPrompt}}},
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var geminiResp GeminiAPIResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		log.Printf("Failed to unmarshal Gemini API response: %v\nRaw body: %s", err, string(body))
		return nil, fmt.Errorf("invalid Gemini API response: %v", err)
	}
	if len(geminiResp.Candidates) == 0 ||
		len(geminiResp.Candidates[0].Content.Parts) == 0 {
		log.Printf("Gemini API response missing candidates or parts. Raw body: %s", string(body))
		return nil, fmt.Errorf("invalid Gemini API response: missing candidates or parts")
	}

	// The response text should be a JSON object for a question
	questionText := geminiResp.Candidates[0].Content.Parts[0].Text
	questionText = cleanCodeBlock(questionText)
	var q struct {
		Question string   `json:"Question"`
		Options  []string `json:"Options"`
		Answer   string   `json:"Answer"`
	}
	err = json.Unmarshal([]byte(questionText), &q)
	if err != nil {
		return nil, fmt.Errorf("could not parse question JSON: %v", err)
	}

	// Find the answer index
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
