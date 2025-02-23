package prompt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func NewPromptManager(request string) *PromptManager {
	return &PromptManager{
		SystemPrompt: `You are a skilled programmer helping edit code. Reply only in json.
Characters from U+0000 through U+001F must be escaped in the json.
Follow the indentation and style of the existing code.
Keep line length to 80 characters or less unless other conventions override.
Update all imports needed by your changes.
Use this json format to send back diffs

[
{
  "file": "file1.go",
  "ranges": [
    {
      "s": 42,     // start line
      "e": 42,     // end line
      "b": ["func oldName(x, y int) bool {"],  // before
      "a": ["func newName(x, y int) int {"]   // after
    },
    {
      "s": 47,
      "e": 49,
      "b": [
        "  if (x > 0) {",
        "    return true;",
        "  }"
      ],
      "a": [
        "  if (x > 0 && y < 100) {",
        "    return x * y;",
        "  }"
      ]
    }
  ]
},
{
  "file": "file2.go",
  "ranges": [
    {
      "s": 22,     // start line
      "e": 22,     // end line
      "b": ["func oldName2(x, y int) bool {"],  // before
      "a": ["func newName2(x, y int) bool{"]   // after
    }
  ]
}]
. ` + request,
		CodeFence: "```",
	}
}

func (pm *PromptManager) AddFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	pm.Files = append(pm.Files, FileContent{
		Filename: filename,
		Content:  string(content),
	})
	return nil
}

func (pm *PromptManager) BuildPrompt(userRequest string) string {
	var buf bytes.Buffer

	buf.WriteString(pm.SystemPrompt)
	buf.WriteString("\n\n")

	for _, file := range pm.Files {
		ext := filepath.Ext(file.Filename)
		lang := strings.TrimPrefix(ext, ".")
		if lang == "" {
			lang = "text"
		}

		buf.WriteString(file.Filename)
		buf.WriteString("\n")
		buf.WriteString(pm.CodeFence)
		buf.WriteString(lang)
		buf.WriteString("\n")
		buf.WriteString(file.Content)
		if !strings.HasSuffix(file.Content, "\n") {
			buf.WriteString("\n")
		}
		buf.WriteString(pm.CodeFence)
		buf.WriteString("\n\n")
	}

	buf.WriteString(userRequest)
	return buf.String()
}

func (pm *PromptManager) ParseDiffs(response string) map[string][]string {
	diffs := make(map[string][]string)

	sections := strings.Split(response, "```diff")

	for _, section := range sections[1:] { // Skip first section (pre-diff text)
		parts := strings.SplitN(section, "```", 2)
		if len(parts) < 1 {
			continue
		}
		diffContent := parts[0]

		lines := strings.Split(diffContent, "\n")
		var filename string
		for _, line := range lines {
			if strings.HasPrefix(line, "+++ ") {
				filename = strings.TrimPrefix(line, "+++ ")
				filename = strings.TrimSpace(filename)
				break
			}
		}

		if filename != "" {
			diffs[filename] = lines
		}
	}

	return diffs
}

func MakePrompt(request string, files []FileContent) string {
	pm := NewPromptManager(request)

	pm.Files = files

	prompt := pm.BuildPrompt(request)
	return prompt
}
