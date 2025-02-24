package prompt

import (
	"os"
	"path/filepath"
	"strings"
)

const maxFiles = 90

func AssembleFiles(dir string) []FileContent {
	buffer := []FileContent{}

	// Walk through directory recursively
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Check if we've reached the max files limit
		if len(buffer) >= maxFiles {
			return filepath.SkipDir
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .go files
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// Get relative path from root dir
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return nil
		}

		fc := FileContent{}
		fc.Filename = relPath
		goFile, _ := os.ReadFile(path)
		fc.Content = string(goFile)
		buffer = append(buffer, fc)
		return nil
	})

	return buffer
}
