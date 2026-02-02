package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AIProvider defines the interface for AI services
type AIProvider interface {
	GenerateCompletion(ctx context.Context, prompt string, model string) (string, error)
}

// OpenRouterProvider implements AIProvider for OpenRouter
type OpenRouterProvider struct {
	apiKey     string
	httpClient *http.Client
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(apiKey string) *OpenRouterProvider {
	return &OpenRouterProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// GenerateCompletion generates text using OpenRouter
func (p *OpenRouterProvider) GenerateCompletion(ctx context.Context, prompt string, model string) (string, error) {
	if model == "" {
		model = "anthropic/claude-3.5-sonnet" // Default model for OpenRouter
	}

	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/boobachad/clankerloop")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenRouter API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenRouter")
	}

	return response.Choices[0].Message.Content, nil
}

// GeminiProvider implements AIProvider for Google Gemini
type GeminiProvider struct {
	apiKey     string
	httpClient *http.Client
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// GenerateCompletion generates text using Gemini
func (p *GeminiProvider) GenerateCompletion(ctx context.Context, prompt string, model string) (string, error) {
	if model == "" {
		model = "gemini-1.5-pro-latest" // Default Gemini model
	}

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, p.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// AIService provides AI capabilities using configured provider
type AIService struct {
	provider AIProvider
}

// NewAIService creates a new AI service
func NewAIService(providerType, openRouterKey, geminiKey string) (*AIService, error) {
	var provider AIProvider

	switch providerType {
	case "openrouter":
		provider = NewOpenRouterProvider(openRouterKey)
	case "gemini":
		provider = NewGeminiProvider(geminiKey)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", providerType)
	}

	return &AIService{provider: provider}, nil
}

// GenerateText generates text using the configured AI provider
func (s *AIService) GenerateText(ctx context.Context, prompt string, model string) (string, error) {
	return s.provider.GenerateCompletion(ctx, prompt, model)
}
