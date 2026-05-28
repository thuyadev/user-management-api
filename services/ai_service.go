package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"user-management-api/utils"
)

type AIService interface {
	GenerateProductDescription(name string, category string) (string, error)
	SuggestCategoryName(keywords string) (string, error)
}

type aiService struct {
	cfg    *utils.Config
	client *http.Client
}

func NewAIService(cfg *utils.Config) AIService {
	return &aiService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func (s *aiService) GenerateProductDescription(name string, category string) (string, error) {
	prompt := fmt.Sprintf(
		"Write a concise, professional product description (2-3 sentences) for a product named '%s' in the '%s' category. Return only the description text.",
		name, category,
	)
	return s.complete(prompt)
}

func (s *aiService) SuggestCategoryName(keywords string) (string, error) {
	prompt := fmt.Sprintf(
		"Suggest a single short category name (2-4 words max) for products related to: %s. Return only the category name, nothing else.",
		keywords,
	)
	return s.complete(prompt)
}

func (s *aiService) complete(prompt string) (string, error) {
	if !s.cfg.AIEnabled || s.cfg.AIAPIKey == "" {
		return s.fallback(prompt), nil
	}

	body, err := json.Marshal(chatRequest{
		Model: s.cfg.AIModel,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, s.cfg.AIAPIURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.AIAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return s.fallback(prompt), nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.fallback(prompt), nil
	}

	if resp.StatusCode != http.StatusOK {
		return s.fallback(prompt), nil
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return s.fallback(prompt), nil
	}

	if len(chatResp.Choices) == 0 {
		return s.fallback(prompt), nil
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}

func (s *aiService) fallback(prompt string) string {
	if strings.Contains(prompt, "category name") {
		return "General Merchandise"
	}
	return "A high-quality product designed to meet your needs with excellent value and reliable performance."
}
