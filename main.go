package main

import (
	"aaai/anthropic"
	"aaai/diff"
	"aaai/prompt"
	"fmt"
	"os"
	"path/filepath"
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
		//fmt.Printf("DEBUG: line='%s', err='%v'\n", line, err)
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		input := strings.TrimSpace(line)

		if input == "." {
			joined := strings.Join(buffer, "\\n")

			// Open history file in append mode
			historyFile, err := os.OpenFile(".aaai.input.history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				historyFile.WriteString(joined + "\n")
				historyFile.Close()
			}

			// Process the command
			joined = strings.Join(buffer, "\n")
			p := prompt.MakePrompt(joined, fcs)
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

			// Create tests directory if it doesn't exist
			testsDir := "tests"
			if err := os.MkdirAll(testsDir, 0755); err != nil {
				fmt.Printf("Error creating tests directory: %v\n", err)
				continue
			}

			for k, v := range m {
				// Get original file content
				origContent, err := os.ReadFile(filepath.Join(dir, k))
				if err != nil {
					fmt.Printf("Error reading original file %s: %v\n", k, err)
					continue
				}

				// Write original file to tests/file.orig
				origPath := filepath.Join(testsDir, k+".orig")
				if err := os.MkdirAll(filepath.Dir(origPath), 0755); err != nil {
					fmt.Printf("Error creating directory for %s: %v\n", origPath, err)
					continue
				}
				os.WriteFile(origPath, origContent, 0644)
				// Write diff to tests/file.diff
				os.WriteFile(filepath.Join(testsDir, k+".diff"), []byte(v), 0644)
				diff.ApplyPatch(dir+"/"+k, v)
			}
			buffer = []string{}
		} else {
			buffer = append(buffer, input)
		}

	}
}
