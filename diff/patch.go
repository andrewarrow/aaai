package diff

import (
	"fmt"
	"os"
	"strings"
)

func HandleDiffs(filename string, diffLines []string) error {
	// Read the original file
	originalContent, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	// Split content into lines, preserving empty lines
	lines := strings.Split(string(originalContent), "\n")

	// Parse the diff header to get line numbers
	// Format: @@ -3,6 +3,9 @@ means:
	// - starts at line 3, removes 6 lines
	// + starts at line 3, adds 9 lines
	header := diffLines[0]
	var oldStart, oldCount, newStart, newCount int
	fmt.Sscanf(header, "@@ -%d,%d +%d,%d @@", &oldStart, &oldCount, &newStart, &newCount)

	// Convert to 0-based index
	startLine := oldStart - 1

	// Create new content
	newLines := make([]string, 0)

	// Add all lines before the change
	newLines = append(newLines, lines[:startLine]...)

	// Track position in original file
	currentLine := startLine

	// Apply the changes
	for i := 1; i < len(diffLines); i++ {
		line := diffLines[i]

		switch {
		case strings.HasPrefix(line, "-"):
			// Skip the line in original content
			currentLine++

		case strings.HasPrefix(line, "+"):
			// Add new line (without the + prefix)
			newLines = append(newLines, line[1:])

		case line == "\\ No newline at end of file":
			// Ignore this line
			continue

		default:
			// Context line - copy from original and advance
			if currentLine < len(lines) {
				newLines = append(newLines, line)
			}
			currentLine++
		}
	}

	// Add remaining lines from the original file
	if currentLine < len(lines) {
		newLines = append(newLines, lines[currentLine:]...)
	}

	// Write back to file
	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(filename, []byte(newContent), 0644)
}
