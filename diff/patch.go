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

	// Split content into lines
	lines := strings.Split(string(originalContent), "\n")

	// Parse the diff header to get line numbers
	header := diffLines[0]
	var startLine, removedLines, addedLines int
	fmt.Sscanf(header, "@@ -%d,%d +%d,%d @@", &startLine, &removedLines, &startLine, &addedLines)
	startLine-- // Convert to 0-based index

	// Create new content
	newLines := make([]string, 0, len(lines)+addedLines-removedLines)
	newLines = append(newLines, lines[:startLine]...)

	// Apply the changes
	diffIndex := 1
	for diffIndex < len(diffLines) {
		line := diffLines[diffIndex]

		// Check if we've gone past the end of the original file
		if startLine >= len(lines) {
			return fmt.Errorf("patch extends beyond end of file at line %d", startLine+1)
		}

		switch {
		case strings.HasPrefix(line, "-"):
			// Skip removed line in original content
			startLine++
		case strings.HasPrefix(line, "+"):
			// Add new line
			newLines = append(newLines, line[1:])
		default:
			// Make sure we're not past the end of the file
			if startLine < len(lines) {
				// Copy unchanged line
				newLines = append(newLines, lines[startLine])
			}
			startLine++
		}
		diffIndex++
	}

	// Add remaining lines
	if startLine < len(lines) {
		newLines = append(newLines, lines[startLine:]...)
	}

	// Write back to file
	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(filename, []byte(newContent), 0644)
}
