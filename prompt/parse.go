package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func HandleDiffs() {

	filePath := os.Args[1]
	diffContent := os.Args[2]

	// Read original file
	originalContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create diff-match-patch instance
	dmp := diffmatchpatch.New()

	// Convert unified diff to diffs
	// First, we need to extract the actual changes from the unified diff
	diffs := convertUnifiedDiffToDiffs(string(originalContent), diffContent, dmp)
	if len(diffs) == 0 {
		fmt.Println("No valid diff found")
		os.Exit(1)
	}

	// Create and apply patches
	patches := dmp.PatchMake(string(originalContent), diffs)
	newContent, applied := dmp.PatchApply(patches, string(originalContent))

	// Verify all patches were applied
	for i, wasApplied := range applied {
		if !wasApplied {
			fmt.Printf("Warning: Patch %d could not be applied\n", i+1)
		}
	}

	// Write the new content back to file
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
