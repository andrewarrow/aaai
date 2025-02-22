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

		if input == "quit" {
			break
		}

		//file, _ := os.OpenFile(".aaai.input.history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		//file.Write([]byte(input + "\n"))
		p := prompt.MakePrompt(input, fcs)
		s, _ := client.Complete(p)
		m := prompt.ParseDiffs(s)
		for k, v := range m {
			fmt.Println(k)
			fmt.Println(v)
			fmt.Println("")
			err := diff.HandleDiffs(dir+"/"+k, v)
			fmt.Println(err)
			fmt.Println("")
		}
	}
}
