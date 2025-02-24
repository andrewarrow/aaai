package prompt

import (
	"strings"
)

func ParseDiffs(response string) map[string]string {
	diffs := make(map[string]string)

	sections := strings.Split(response, "```diff")
	sections2 := strings.Split(response, "```go")
	if len(sections) == 0 {
		sections = sections2
	}

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

				// Handle a/ and b/ prefixes
				if strings.HasPrefix(filename, "a/") || strings.HasPrefix(filename, "b/") {
					filename = filename[2:]
				}

				// Handle filenames that start with a single forward slash
				if strings.HasPrefix(filename, "/") {
					filename = filename[1:]
				}

				// Keep the full path but clean up any "./" prefixes
				filename = strings.TrimPrefix(filename, "./")
			}

			newLines = append(newLines, line)
		}

		if filename != "" {
			diffs[filename] = strings.Join(newLines, "\n")
		}
	}

	return diffs
}
