package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ProcessDiffs(dir, jsonData string) error {
	// Parse the JSON
	var changes []Change
	err := json.Unmarshal([]byte(sanitizeJSON(jsonData)), &changes)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	// Process each change
	for _, change := range changes {
		if err := processFileChange(dir, change); err != nil {
			return fmt.Errorf("error processing file %s: %w", change.File, err)
		}
		fmt.Printf("Successfully updated file: %s\n", change.File)
	}
	return nil
}

func processFileChange(dir string, change Change) error {
	// Read the original file
	content, err := os.ReadFile(dir + "/" + change.File)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Process each range change
	for _, r := range change.Ranges {
		// Validate range bounds
		if err := validateRange(r, len(lines)); err != nil {
			return err
		}

		// Create new content with proper capacity
		newLines := make([]string, 0, len(lines)+(len(r.After)-len(r.Before)))

		// Append lines before the change
		newLines = append(newLines, lines[:r.Start]...)

		// Append the new lines
		newLines = append(newLines, r.After...)

		// Append remaining lines after the change, checking bounds
		if r.End+1 <= len(lines) {
			newLines = append(newLines, lines[r.End+1:]...)
		}

		lines = newLines
	}

	// Write the modified content back to the file
	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(dir+"/"+change.File, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func validateRange(r Range, lineCount int) error {
	if r.Start < 0 {
		return fmt.Errorf("invalid start line: %d", r.Start)
	}
	if r.End >= lineCount {
		return fmt.Errorf("invalid end line: %d (file has %d lines)", r.End, lineCount)
	}
	if r.Start > r.End {
		return fmt.Errorf("start line (%d) is after end line (%d)", r.Start, r.End)
	}
	return nil
}

// compareLines compares two string slices for equality, ignoring trailing whitespace
func compareLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if strings.TrimRight(a[i], " \t\r\n") != strings.TrimRight(b[i], " \t\r\n") {
			return false
		}
	}
	return true
}

func sanitizeJSON(input string) string {
	// Replace literal tabs with \t escape sequence
	sanitized := strings.ReplaceAll(input, "\t", "\\t")
	if !strings.HasPrefix(sanitized, "[") {
		return "[" + sanitized + "]"
	}
	return sanitized
}
