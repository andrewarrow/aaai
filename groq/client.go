package groq

import (
	"aaai/prompt"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	DefaultAPIEndpoint = "https://api.groq.com/openai/v1/chat/completions"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

type CompletionRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
	Stream    bool      `json:"stream,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Choices []Choice       `json:"choices"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
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

func (c *Client) Complete(promptString string) (string, error) {
	req := CompletionRequest{
		Model:  "mixtral-8x7b-32768", // Replace with the appropriate model name for Groq
		Stream: true,
		Messages: []Message{
			{
				Role:    "user",
				Content: promptString,
			},
		},
		MaxTokens: 8192,
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
	request.Header.Set("Authorization", "Bearer "+c.APIKey)

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	reader := bufio.NewReader(response.Body)
	parser := prompt.NewStreamParser()
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading stream: %w", err)
		}

		if len(line) == 0 {
			continue
		}

		if !bytes.HasPrefix(line, []byte("data: ")) {
			continue
		}

		data := bytes.TrimPrefix(line, []byte("data: "))

		if string(data) == "[DONE]\n" {
			break
		}

		parser.ProcessLine(firstChoice(data))
	}
	return parser.Result(), nil
}
