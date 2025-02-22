package diff

import (
	"fmt"
	"os"
	"strings"
)

func HandleDiffs(filename string, diffLines []string) error {
	originalContent, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	lines := strings.Split(string(originalContent), "\n")

	header := diffLines[0]
	var startLine, removedLines, addedLines int
	fmt.Sscanf(header, "@@ -%d,%d +%d,%d @@", &startLine, &removedLines, &startLine, &addedLines)
	startLine-- // Convert to 0-based index

	newLines := make([]string, 0, len(lines)+addedLines-removedLines)
	newLines = append(newLines, lines[:startLine]...)

	diffIndex := 1
	for diffIndex < len(diffLines) {
		line := diffLines[diffIndex]
		switch {
		case strings.HasPrefix(line, "-"):
			// Skip removed line in original content
			startLine++
		case strings.HasPrefix(line, "+"):
			// Add new line
			newLines = append(newLines, line[1:])
		default:
			// Copy unchanged line
			newLines = append(newLines, lines[startLine])
			startLine++
		}
		diffIndex++
	}

	newLines = append(newLines, lines[startLine:]...)

	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(filename, []byte(newContent), 0644)
}

/*
func foo() {
	for filename, diff := range diffs {
		fmt.Printf("Applying patch to %s...\n", filename)
		if err := applyPatch(filename, diff); err != nil {
			fmt.Printf("Error applying patch to %s: %v\n", filename, err)
			continue
		}
		fmt.Printf("Successfully patched %s\n", filename)
	}
}*/
