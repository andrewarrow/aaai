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
						Text: "Please review this Go code and suggest improvements: " + document,
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

type Request struct {
	Messages           []Message    `json:"messages"`
	Model              string       `json:"model"`
	Prompt             string       `json:"prompt"`
	ParentMessageUUID  string       `json:"parent_message_uuid"`
	Timezone           string       `json:"timezone"`
	PersonalizedStyles []Style      `json:"personalized_styles"`
	Tools              []Tool       `json:"tools"`
	Attachments        []Attachment `json:"attachments"`
	Files              []any        `json:"files"`
	SyncSources        []any        `json:"sync_sources"`
	RenderingMode      string       `json:"rendering_mode"`
	MaxTokens          int          `json:"max_tokens"`
}

type Style struct {
	Name      string `json:"name"`
	Prompt    string `json:"prompt"`
	Summary   string `json:"summary"`
	IsDefault bool   `json:"isDefault"`
	Type      string `json:"type"`
	Key       string `json:"key"`
}

type Tool struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Attachment struct {
	FileName         string `json:"file_name"`
	FileType         string `json:"file_type"`
	FileSize         int    `json:"file_size"`
	ExtractedContent string `json:"extracted_content"`
}

/*
func makeRequest(goFile string) []byte {
	req := Request{
		Model:  "claude-3-5-sonnet-20241022",
		Prompt: "what is this file?",
		//ParentMessageUUID: "00000000-0000-4000-8000-000000000000",
		Timezone: "America/Los_Angeles",
		PersonalizedStyles: []Style{
			{
				Name:      "Normal",
				Prompt:    "Normal",
				Summary:   "Default responses from Claude",
				IsDefault: true,
				Type:      "default",
				Key:       "Default",
			},
		},
		Tools: []Tool{
			{Type: "text_editor_20250124", Name: "str_replace_editor"},
		},
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: "Please review this Go code and suggest improvements:",
					},
				},
			},
		},
		MaxTokens: 1024,
		Attachments: []Attachment{
			{
				FileName:         "main.go",
				FileType:         "text/plain",
				FileSize:         643,
				ExtractedContent: goFile,
			},
		},
		Files:         []interface{}{},
		SyncSources:   []interface{}{},
		RenderingMode: "messages",
	}

	data, _ := json.Marshal(req)
	return data

}*/
