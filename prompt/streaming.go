package prompt

import (
	"encoding/json"
	"fmt"
	"strings"
)

type StreamParser struct {
	buffer strings.Builder
}

func NewStreamParser() *StreamParser {
	return &StreamParser{}
}

func (p *StreamParser) ProcessLine(line string) error {

	var m map[string]any

	json.Unmarshal([]byte(line), &m)
	if m["type"] == "content_block_delta" {

		d := m["delta"].(map[string]any)
		s := d["text"].(string)
		fmt.Print(s)
		p.buffer.WriteString(s)
	}

	return nil
}

func (p *StreamParser) ProcessLineAsString(s string) error {

	fmt.Print(s)
	p.buffer.WriteString(s)
	return nil
}

// GetResult returns the final concatenated text
func (p *StreamParser) Result() string {
	return p.buffer.String()
}
