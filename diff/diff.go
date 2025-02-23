package diff

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func ApplyPatch(fileOrig, fileDiff string) {

	linesOrig, err := readLines(fileOrig)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", fileOrig, err)
		os.Exit(1)
	}

	linesDiff, err := readLinesFromString(fileDiff)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", fileDiff, err)
		os.Exit(1)
	}

	// Parse unified diff into hunks
	hunks := parseHunks(linesDiff)

	// Apply all hunks at once
	updatedLines := applyHunks(linesOrig, hunks)

	// Write back to original file
	err = writeLines(fileOrig, updatedLines)
	if err != nil {
		fmt.Printf("Error writing to %s: %v\n", fileOrig, err)
		os.Exit(1)
	}
}

type Hunk struct {
	StartLine int
	Length    int
	NewStart  int
	NewLength int
	Lines     []string
}

func parseHunks(diffLines []string) []Hunk {
	var hunks []Hunk
	var currentHunk *Hunk

	hunkHeader := regexp.MustCompile(`@@ -(\d+),(\d+) \+(\d+),(\d+) @@`)

	for _, line := range diffLines {
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			continue
		}

		if match := hunkHeader.FindStringSubmatch(line); match != nil {
			if currentHunk != nil {
				hunks = append(hunks, *currentHunk)
			}
			start, _ := strconv.Atoi(match[1])
			length, _ := strconv.Atoi(match[2])
			newStart, _ := strconv.Atoi(match[3])
			newLength, _ := strconv.Atoi(match[4])
			currentHunk = &Hunk{
				StartLine: start,
				Length:    length,
				NewStart:  newStart,
				NewLength: newLength,
			}
			continue
		}

		if currentHunk != nil {
			currentHunk.Lines = append(currentHunk.Lines, line)
		}
	}

	if currentHunk != nil {
		hunks = append(hunks, *currentHunk)
	}

	return hunks
}

func findHunkPosition(lines []string, hunk Hunk) int {
	// Get context lines from the hunk (leading lines that start with space or minus)
	var contextLines []string
	for _, line := range hunk.Lines {
		trimmedLine := strings.TrimRight(line, "\n")
		if strings.HasPrefix(trimmedLine, " ") {
			contextLines = append(contextLines, trimmedLine[1:])
		} else if strings.HasPrefix(trimmedLine, "-") {
			contextLines = append(contextLines, trimmedLine[1:])
		}
		if len(contextLines) == 3 { // Get first 3 context lines
			break
		}
	}

	// Check each position in the file
	for i := 0; i < len(lines); i++ {
		matches := 0
		for j, ctx := range contextLines {
			if i+j >= len(lines) {
				break
			}
			fileLine := strings.TrimRight(lines[i+j], "\n")
			if fileLine == ctx {
				matches++
			}
		}
		if matches == len(contextLines) {
			// Found matching context at position i
			return i
		}
	}
	return -1
}

func applyHunks(original []string, hunks []Hunk) []string {
	result := make([]string, len(original))
	copy(result, original)

	for _, hunk := range hunks {
		pos := findHunkPosition(result, hunk)
		if pos == -1 {
			continue
		}

		updated := make([]string, 0, len(result))
		updated = append(updated, result[:pos]...)

		currentPos := pos

		for _, line := range hunk.Lines {
			switch {
			case strings.HasPrefix(line, " "):
				if currentPos < len(result) {
					updated = append(updated, result[currentPos])
					currentPos++
				}
			case strings.HasPrefix(line, "+"):
				updated = append(updated, line[1:])
			case strings.HasPrefix(line, "-"):
				currentPos++
			}
		}

		// Only append remaining content if we haven't exceeded bounds
		if currentPos < len(result) {
			updated = append(updated, result[currentPos:]...)
		}

		result = updated
	}

	return result
}

func findLastNonDeleted(linesOrig []string, linesUpdated []string) int {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(strings.Join(linesOrig, "\n"), strings.Join(linesUpdated, "\n"), false)

	pos := 0
	lastNonDeleted := 0
	origText := strings.Join(linesOrig, "\n")

	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			pos += len(diff.Text)
			lastNonDeleted = strings.Count(origText[:pos], "\n") + 1
		case diffmatchpatch.DiffDelete:
			pos += len(diff.Text)
		}
	}

	return lastNonDeleted
}

func assertNewlines(lines []string) error {
	if len(lines) == 0 {
		return nil
	}
	for i, line := range lines[:len(lines)-1] {
		if len(line) == 0 || !strings.HasSuffix(line, "\n") {
			return fmt.Errorf("line %d does not end with newline", i)
		}
	}
	return nil
}

func readLinesFromString(content string) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")
	}
	return lines, scanner.Err()
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")
	}
	return lines, scanner.Err()
}

func writeLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
