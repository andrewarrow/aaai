package groq

import (
	"encoding/json"
	"fmt"
)

// ChatCompletionChunk represents the JSON data
type ChatCompletionChunk struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"` // Use int64 for Unix timestamp
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []string `json:"choices"`
}

type Delta struct {
	Content string `json:"content"`
}

func firstChoice(data []byte) string {
	var chunk ChatCompletionChunk
	if err := json.Unmarshal(data, &chunk); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ""
	}

	for _, choice := range chunk.Choices {
		return choice
	}
	return ""
}
