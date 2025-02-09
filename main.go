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
		goFile, _ := os.ReadFile(dir + "/main.go")
		fmt.Print("> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "quit" {
			break
		}

		s, err := client.Complete(input, string(goFile))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		os.Remove(file)
		os.WriteFile(file, []byte(s), 0644)
		fmt.Println("")

	}
}
