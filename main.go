package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"revai/config"
	"time"
)

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model     string          `json:"model"`
	Messages  []OpenAIMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	apiKey := flag.String("key", "", "OpenAI API key")
	file := flag.String("file", "", "File to analyze (optional)")
	configPath := flag.String("config", "config/config.json", "Path to configuration file")
	flag.Parse()

	cfg := config.DefaultConfig()
	err := cfg.LoadFromFile(*configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if *apiKey != "" {
		cfg.AI.ApiKey = *apiKey
	}

	var diff string

	if *file != "" {
		diff, err = getFileDiff(*file)
	} else {
		diff, err = getGitDiff()
	}

	if err != nil {
		fmt.Printf("Error getting diff: %v\n", err)
		os.Exit(1)
	}

	review, err := getCodeReview(diff, cfg)
	if err != nil {
		fmt.Printf("Error getting code review: %v\n", err)
		os.Exit(1)
	}

	curDate := time.Now()
	fileName := fmt.Sprintf("cr%04d%02d%02d_%02d%02d.md", curDate.Year(), curDate.Month(), curDate.Day(), curDate.Hour(), curDate.Minute())
	err = os.WriteFile(fileName, []byte(review), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func getFileDiff(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getCodeReview(diff string, cfg config.Config) (string, error) {
	if diff == "" {
		return "No diff found", nil
	}
	prompt := cfg.Prompt
	prompt = fmt.Sprintf(`%s Diff:%s`, prompt, diff)

	requestBody := OpenAIRequest{
		Model: cfg.AI.Model,
		Messages: []OpenAIMessage{
			{
				Role:    cfg.Role,
				Content: prompt,
			},
		},
		MaxTokens: 2000,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", cfg.AI.ApiBase, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.AI.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response OpenAIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
