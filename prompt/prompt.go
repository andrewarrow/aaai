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
		SystemPrompt: `You are a skilled programmer helping edit code, using unified diffs.
Follow the indentation and style of the existing code.
Keep line length to 80 characters or less unless other conventions override.
Update all imports needed by your changes.
List each file using the CodeFence and the string diff after the fence.
Always include 100% of the diffs for a single file in just one file between CodeFences.
For example do not list 2 ranges of diffs for foo.txt and then a CodeFence and then one
more diff for foo.txt. Instead all 3 diffs should be together for foo.txt file.
Make sure to list +++ and the filename and --- and the filename at start of each diff.
Use unified diff format with 3 lines of context.` + request,
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
		fmt.Println(file.Filename)
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

func MakePrompt(request string, files []FileContent) string {
	pm := NewPromptManager(request)

	pm.Files = files

	prompt := pm.BuildPrompt(request)
	return prompt
}
