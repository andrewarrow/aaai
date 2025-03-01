package diff

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ApplyPatch(fileOrig, fileDiff string) {
	linesOrig, _ := readLines(fileOrig)
	if linesOrig == nil {
		linesOrig = []string{} // Initialize empty slice for new files
	}

	linesDiff, err := readLinesFromString(fileDiff)
	if err != nil {
		fmt.Printf("Error reading diff: %v\n", err)
		os.Exit(1)
	}

	hunks := parseHunks(linesDiff)
	updatedLines := applyHunks(linesOrig, hunks)

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
			start-- // 0-based
			length, _ := strconv.Atoi(match[2])
			newStart, _ := strconv.Atoi(match[3])
			newStart-- // 0-based
			newLength, _ := strconv.Atoi(match[4])
			currentHunk = &Hunk{StartLine: start, Length: length, NewStart: newStart, NewLength: newLength}
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
	if hunk.Length == 0 {
		pos := hunk.StartLine + 1
		if pos >= 0 && pos <= len(lines) {
			return pos
		}
		return -1
	}

	var contextLines []string
	for _, line := range hunk.Lines {
		trimmedLine := strings.TrimRight(line, "\n")
		if strings.HasPrefix(trimmedLine, " ") || strings.HasPrefix(trimmedLine, "-") {
			contextLines = append(contextLines, trimmedLine[1:])
		}
	}
	if len(contextLines) == 0 {
		return -1
	}

	// Check expected position
	i := hunk.StartLine
	if i >= 0 && i <= len(lines)-len(contextLines) {
		matches := true
		for j, ctx := range contextLines {
			if i+j >= len(lines) || strings.TrimRight(lines[i+j], "\n") != ctx {
				matches = false
				break
			}
		}
		if matches {
			return i
		}
	}

	// Search nearby
	startSearch := hunk.StartLine - 3
	if startSearch < 0 {
		startSearch = 0
	}
	endSearch := hunk.StartLine + 3
	if endSearch > len(lines) {
		endSearch = len(lines)
	}
	for i := startSearch; i <= endSearch-len(contextLines); i++ {
		matches := true
		for j, ctx := range contextLines {
			if i+j >= len(lines) || strings.TrimRight(lines[i+j], "\n") != ctx {
				matches = false
				break
			}
		}
		if matches {
			return i
		}
	}

	// Search entire file
	for i := 0; i <= len(lines)-len(contextLines); i++ {
		matches := true
		for j, ctx := range contextLines {
			if strings.TrimRight(lines[i+j], "\n") != ctx {
				matches = false
				break
			}
		}
		if matches {
			return i
		}
	}
	return -1
}

func applyHunks(original []string, hunks []Hunk) []string {
	result := make([]string, len(original))
	copy(result, original)

	for i, hunk := range hunks {
		pos := findHunkPosition(result, hunk)
		if pos == -1 {
			fmt.Printf("Warning: Hunk %d skipped (context not found at line %d)\n", i+1, hunk.StartLine+1)
			continue
		}

		end := pos + hunk.Length
		if end > len(result) {
			end = len(result)
		}
		before := result[:pos]
		after := result[end:]

		var newLines []string
		for _, line := range hunk.Lines {
			if strings.HasPrefix(line, "+") || strings.HasPrefix(line, " ") {
				newLine := line[1:]
				newLine = strings.TrimRight(newLine, "\n") + "\n"
				newLines = append(newLines, newLine)
			}
		}

		result = append(before, append(newLines, after...)...)
	}
	return result
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
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating directory: %v", err)
	}

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
