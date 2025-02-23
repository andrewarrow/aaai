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
	apiKey = os.Getenv("DEEPSEEK")
	if apiKey == "" {
		fmt.Println("Please set DEEPSEEK environment variable")
		return
	}
	apiKey = os.Getenv("GROQ")
	if apiKey == "" {
		fmt.Println("Please set GROQ environment variable")
		return
	}
	apiKey = os.Getenv("ANTHROPIC_API_KEY")

	client := anthropic.NewClient(apiKey)
	//client := deepseek.NewClient(apiKey)
	//client := groq.NewClient(apiKey)

	rl, _ := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     ".aaai.input.history",
		InterruptPrompt: "^C",
		EOFPrompt:       "quit",
	})

	buffer := []string{}
	for {
		fcs := prompt.AssembleFiles(dir)
		fmt.Print("> ")

		line, err := rl.Readline()
		fmt.Printf("DEBUG: line='%s', err='%v'\n", line, err)
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		input := strings.TrimSpace(line)
		if input == "quit" || input == "exit" {
			break
		}
		if input == "." {
			p := prompt.MakePrompt(strings.Join(buffer, "\n"), fcs)
			s, err := client.Complete(p)
			fmt.Println(err)
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
			fmt.Println(s)
			m := prompt.ParseDiffs(s)
			fmt.Println("===")
			fmt.Println(m)
			fmt.Println("===")
			for k, v := range m {
				diff.ApplyPatch(dir+"/"+k, v)
			}
			buffer = []string{}
		} else {
			buffer = append(buffer, input)
		}

	}
}
