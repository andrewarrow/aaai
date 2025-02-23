package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Change represents a modification to a file
type Change struct {
	File   string  `json:"file"`
	Ranges []Range `json:"ranges"`
}

// Range represents a specific change within a file
type Range struct {
	Start  int      `json:"s"` // Start line number (1-based)
	End    int      `json:"e"` // End line number (1-based)
	Before []string `json:"b"` // Lines before the change
	After  []string `json:"a"` // Lines after the change
}

func ProcessDiffs(dir, jsonData string) error {
	var changes []Change
	if err := json.Unmarshal([]byte(sanitizeJSON(jsonData)), &changes); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	for _, change := range changes {
		if err := applyChange(dir, change); err != nil {
			return fmt.Errorf("error processing %s: %w", change.File, err)
		}
		fmt.Printf("Successfully updated file: %s\n", change.File)
	}
	return nil
}

func applyChange(dir string, change Change) error {
	// Read file content
	content, err := os.ReadFile(dir + "/" + change.File)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Split into lines and remove trailing empty line if present
	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Sort ranges in descending order to apply from bottom to top
	sort.Slice(change.Ranges, func(i, j int) bool {
		return change.Ranges[i].Start > change.Ranges[j].Start
	})

	// Apply each range change
	for _, r := range change.Ranges {
		start := r.Start - 1 // Convert to 0-based
		end := r.End - 1

		// Validate range
		if err := validateRange(start, end, len(lines)); err != nil {
			return err
		}

		// Verify the "before" content matches
		if !matchesBeforeLines(lines[start:end+1], r.Before) {
			return fmt.Errorf("content mismatch at lines %d-%d", r.Start, r.End)
		}

		// Replace the range with new content
		lines = append(append(lines[:start], r.After...), lines[end+1:]...)
	}

	// Ensure trailing newline
	lines = append(lines, "")

	// Write back to file
	return os.WriteFile(dir+"/"+change.File, []byte(strings.Join(lines, "\n")), 0644)
}

func matchesBeforeLines(fileLines, beforeLines []string) bool {
	if len(fileLines) != len(beforeLines) {
		return false
	}
	for i := range fileLines {
		if fileLines[i] != beforeLines[i] {
			return false
		}
	}
	return true
}

func validateRange(start, end, lineCount int) error {
	if start < 0 {
		return fmt.Errorf("invalid start line: %d", start+1)
	}
	if end >= lineCount {
		return fmt.Errorf("invalid end line: %d (file has %d lines)", end+1, lineCount)
	}
	if start > end {
		return fmt.Errorf("start line (%d) is after end line (%d)", start+1, end+1)
	}
	return nil
}

func sanitizeJSON(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\t", "\\t")
	if !strings.HasPrefix(input, "[") {
		input = "[" + input
	}
	if !strings.HasSuffix(input, "]") {
		input = input + "]"
	}
	return input
}
