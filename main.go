package main

import (
	"aaai/anthropic"
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
	file := dir + "/main.go"

	for {
		goFile, _ := os.ReadFile(file)
		fmt.Print("> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "quit" {
			break
		}

		s, err := client.Complete(input, "main.go\n"+string(goFile))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println("")

		files := strings.Split(s, "---------")
		fmt.Println(len(files))
		for _, f := range files {
			//fmt.Printf("%d: %v\n", i, []byte(f))
			lines := strings.Split(f, "\n")
			j := 0
			if strings.TrimSpace(lines[0]) == "" {
				j++
			}
			newFile := dir + "/" + lines[j]
			content := strings.Join(lines[1:], "\n")
			fmt.Println(newFile, len(content))

			os.Remove(newFile)
			os.WriteFile(newFile, []byte(content), 0644)
		}

	}
}
