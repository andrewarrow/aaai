package main

import (
	"aaai/anthropic"
	"aaai/diff"
	"aaai/prompt"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("./aaai [dir]")
		return
	}
	dir := os.Args[1]
	_ = dir

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set ANTHROPIC_API_KEY environment variable")
		return
	}

	fcs := prompt.AssembleFiles(dir)

	if false {
		p := prompt.MakePrompt(fcs)
		client := anthropic.NewClient(apiKey)
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

/*
func maini2() {
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
		s, err := client.Complete2(input, oneOrMoreFiles)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println("\n================\n")

		newFile := dir + "/diff.patch"
		os.Remove(newFile)
		os.WriteFile(newFile, []byte(s+"\n"), 0644)
		ApplyPatch(dir)

	}
} */
