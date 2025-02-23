package groq

import (
	"encoding/json"
	"fmt"
)

// ChatCompletionChunk represents the JSON data
type ChatCompletionChunk struct {
	ID                string        `json:"id"`
	Object            string        `json:"object"`
	Created           int64         `json:"created"` // Use int64 for Unix timestamp
	Model             string        `json:"model"`
	SystemFingerprint string        `json:"system_fingerprint"`
	Choices           []ParseChoice `json:"choices"`
}

type ParseChoice struct {
	Index        int             `json:"index"`
	Delta        Delta           `json:"delta"`
	LogProbs     json.RawMessage `json:"logprobs"`      // Use json.RawMessage for null or not-present values
	FinishReason json.RawMessage `json:"finish_reason"` // Use json.RawMessage for null or not-present values
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

	s := chunk.Choices[0].Delta.Content
	return s
}
