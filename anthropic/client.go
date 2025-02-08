package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	DefaultAPIEndpoint = "https://api.anthropic.com/v1/messages"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

type CompletionRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	Document *Document `json:"document,omitempty"`
}

type Document struct {
	Source Source `json:"source"`
}

type Source struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type CompletionResponse struct {
	Content string         `json:"content"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) Complete(prompt, document string) (string, error) {
	req := CompletionRequest{
		Model: "claude-3-5-sonnet-20241022",
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: "I'm going to give you some instructions and a .go source file. I want you to make the changes to the go file and return the entire file with your changes. Include nothing else in your reply. " + prompt + ": " + document,
					},
				},
			},
		},
		MaxTokens: 1024,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	request, err := http.NewRequest("POST", DefaultAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.APIKey)
	request.Header.Set("anthropic-version", "2023-06-01")

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	fmt.Println(string(body))

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	contents := result["content"].([]any)
	content := contents[0].(map[string]any)
	t := content["text"].(string)
	fmt.Println(t)

	return t, nil
}
