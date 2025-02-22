package prompt

import (
	"os"
	"strings"
)

func AssembleFiles(dir string) []FileContent {
	files, _ := os.ReadDir(dir)
	buffer := []FileContent{}
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".go") == false {
			continue
		}
		fc := FileContent{}
		fc.Filename = file.Name()
		goFile, _ := os.ReadFile(dir + "/" + name)
		fc.Content = string(goFile)
		buffer = append(buffer, fc)
	}
	return buffer
}
