package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func HandleDiffs(filePath, diffContent string) {

	originalContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	dmp := diffmatchpatch.New()

	diffs := convertUnifiedDiffToDiffs(string(originalContent), diffContent, dmp)
	if len(diffs) == 0 {
		fmt.Println("No valid diff found")
		os.Exit(1)
	}

	patches := dmp.PatchMake(string(originalContent), diffs)
	newContent, applied := dmp.PatchApply(patches, string(originalContent))

	for i, wasApplied := range applied {
		if !wasApplied {
			fmt.Printf("Warning: Patch %d could not be applied\n", i+1)
		}
	}

	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File updated successfully!")
}

func convertUnifiedDiffToDiffs(originalText, unifiedDiff string, dmp *diffmatchpatch.DiffMatchPatch) []diffmatchpatch.Diff {
	var diffs []diffmatchpatch.Diff
	lines := strings.Split(unifiedDiff, "\n")

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		switch line[0] {
		case '+':
			// Added line
			diffs = append(diffs, diffmatchpatch.Diff{
				Type: diffmatchpatch.DiffInsert,
				Text: line[1:] + "\n",
			})
		case '-':
			// Removed line
			diffs = append(diffs, diffmatchpatch.Diff{
				Type: diffmatchpatch.DiffDelete,
				Text: line[1:] + "\n",
			})
		case '@':
			// Hunk header - skip
			continue
		case ' ':
			// Context line
			diffs = append(diffs, diffmatchpatch.Diff{
				Type: diffmatchpatch.DiffEqual,
				Text: line[1:] + "\n",
			})
		default:
			// Headers or other lines - skip
			continue
		}
	}

	return diffs
}

func ParseDiffs(response string) map[string][]string {
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
