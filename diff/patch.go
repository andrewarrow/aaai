package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func ProcesssDiffs(dir, jsonData string) {
	// Parse the JSON
	var changes []Change
	err := json.Unmarshal([]byte(jsonData), &changes)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Process each change
	for _, change := range changes {
		// Read the original file
		content, err := os.ReadFile(dir + "/" + change.File)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", change.File, err)
			continue
		}

		lines := strings.Split(string(content), "\n")

		// Apply each range change
		for _, r := range change.Ranges {
			// Verify the "before" content matches
			actualBefore := lines[r.Start : r.End+1]
			if !compareLines(actualBefore, r.Before) {
				fmt.Printf("Warning: Content mismatch in file %s at lines %d-%d\n",
					change.File, r.Start+1, r.End+1)
				continue
			}

			// Create new content
			newLines := make([]string, 0, len(lines)+(len(r.After)-len(r.Before)))
			newLines = append(newLines, lines[:r.Start]...)
			newLines = append(newLines, r.After...)
			newLines = append(newLines, lines[r.End+1:]...)
			lines = newLines
		}

		// Write the modified content back to the file
		newContent := strings.Join(lines, "\n")
		err = os.WriteFile(dir+"/"+change.File, []byte(newContent), 0644)
		if err != nil {
			fmt.Printf("Error writing file %s: %v\n", change.File, err)
			continue
		}

		fmt.Printf("Successfully updated file: %s\n", change.File)
	}
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
