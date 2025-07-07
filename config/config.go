package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Role         string   `json:"role"`
	Prompt       string   `json:"prompt"`
	AI           AIConfig `json:"ai"`
	ExcludeDirs  []string `json:"excludeDirs"`
	ExcludeFiles []string `json:"excludeFiles"`
}

type AIConfig struct {
	Model     string `json:"model"`
	Provider  string `json:"provider"`
	ApiKey    string `json:"apiKey"`
	ApiBase   string `json:"apiBase"`
	MaxTokens int    `json:"maxTokens,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		Role:   "assistant",
		Prompt: "Проведи строгое код-ревью. Будь конструктивным.",
		AI: AIConfig{
			Model:    "",
			Provider: "openai",
			ApiKey:   "",
			ApiBase:  "",
		},
		ExcludeDirs:  []string{},
		ExcludeFiles: []string{},
	}
}

func (c *Config) IsValid() bool {
	return !(c.AI.Model == "" || c.AI.ApiKey == "" || c.AI.ApiBase == "" || c.Role == "" || c.Prompt == "")
}

func (c *Config) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}
