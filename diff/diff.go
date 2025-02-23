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
	// Get the first few context lines from the hunk
	contextLines := make([]string, 0)
	for _, line := range hunk.Lines {
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "-") {
			trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, " "), "-"))
			if trimmed != "" {
				contextLines = append(contextLines, trimmed)
				if len(contextLines) == 3 {
					break
				}
			}
		}
	}

	// Look for matching context
	for i := 0; i < len(lines); i++ {
		matches := 0
		for j, ctx := range contextLines {
			if i+j >= len(lines) {
				break
			}
			if strings.TrimSpace(lines[i+j]) == ctx {
				matches++
			}
		}
		if matches == len(contextLines) {
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

		before := result[:pos]
		var newContent []string
		skipLines := 0

		for _, line := range hunk.Lines {
			switch {
			case strings.HasPrefix(line, " "):
				newContent = append(newContent, line[1:])
				skipLines++
			case strings.HasPrefix(line, "+"):
				newContent = append(newContent, line[1:])
			case strings.HasPrefix(line, "-"):
				skipLines++
			}
		}

		after := result[pos+skipLines:]
		result = append(before, append(newContent, after...)...)
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
