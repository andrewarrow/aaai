package prompt

import (
	"strings"
)

func ParseDiffs(response string) map[string]string {
	diffs := make(map[string]string)

	sections := strings.Split(response, "```diff")

	for _, section := range sections[1:] { // Skip first section (pre-diff text)
		parts := strings.SplitN(section, "```", 2)
		if len(parts) < 1 {
			continue
		}
		diffContent := strings.TrimSpace(parts[0])

		newLines := []string{}
		lines := strings.Split(diffContent, "\n")
		var filename string
		for _, line := range lines {
			if strings.HasPrefix(line, "+++ ") {
				filename = strings.TrimPrefix(line, "+++ ")
				filename = strings.TrimSpace(filename)
				tokens := strings.Split(filename, "/")
				filename = tokens[len(tokens)-1]
			}
			newLines = append(newLines, line)
		}

		if filename != "" {
			diffs[filename] = strings.Join(newLines, "\n")
		}
	}

	return diffs
}
