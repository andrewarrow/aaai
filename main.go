package main

import (
	"aaai/anthropic"
	"aaai/diff"
	"aaai/prompt"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
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

	rl, _ := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     ".aaai.input.history",
		InterruptPrompt: "^C",
		EOFPrompt:       "quit",
	})

	for {
		fcs := prompt.AssembleFiles(dir)
		fmt.Print("> ")

		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		input := strings.TrimSpace(line)
		if input == "quit" || input == "exit" {
			break
		}

		p := prompt.MakePrompt(input, fcs)
		s, _ := client.Complete(p)
		diff.ProcesssDiffs(dir, s)
	}
}
