package main

import (
	"aaai/anthropic"
	"bufio"
	"fmt"
	"os"
	"strings"
)

const DELIMETER = "---------"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("./aaai [dir]")
		return
	}
	dir := os.Args[1]

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set ANTHROPIC_API_KEY environment variable")
		return
	}

	client := anthropic.NewClient(apiKey)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		oneOrMoreFiles := AssembleFiles(dir)
		fmt.Print("> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "quit" {
			break
		}

		file, _ := os.OpenFile(".aaai.input.history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		file.Write([]byte(input + "\n"))

		fmt.Println(oneOrMoreFiles)
		s, err := client.Complete(input, oneOrMoreFiles)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println("\n================\n")

		files := strings.Split(s, DELIMETER)
		fmt.Println(len(files))
		for _, f := range files {
			//fmt.Printf("%d: %v\n", i, []byte(f))
			lines := strings.Split(f, "\n")
			j := 0
			if strings.TrimSpace(lines[0]) == "" {
				j++
			}
			newFile := dir + "/" + lines[j]
			content := strings.Join(lines[j+1:], "\n")
			fmt.Println(newFile, len(content))

			os.Remove(newFile)
			os.WriteFile(newFile, []byte(content), 0644)
		}

	}
}

func AssembleFiles(dir string) string {
	files, _ := os.ReadDir(dir)
	buffer := []string{}
	first := true
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".go") == false {
			continue
		}
		if first == false {
			buffer = append(buffer, DELIMETER+"\n")
		}
		topLine := file.Name() + "\n"
		buffer = append(buffer, topLine)
		goFile, _ := os.ReadFile(dir + "/" + name)
		buffer = append(buffer, string(goFile)+"\n")
		first = false
	}
	return strings.Join(buffer, "")
}
