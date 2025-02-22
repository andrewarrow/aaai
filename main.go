package main

import (
	"aaai/anthropic"
	"aaai/prompt"
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

	p := prompt.MakePrompt()
	fmt.Println(dir, p)
}
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
		s, err := client.Complete(input, oneOrMoreFiles)
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
}

func ApplyPatch(dir string) {
	cmd := exec.Command("git", "apply", "--whitespace=fix", "diff.patch")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	fmt.Printf("%s\n", output)
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
			buffer = append(buffer, "\n")
		}
		topLine := file.Name() + "\n"
		buffer = append(buffer, topLine)
		goFile, _ := os.ReadFile(dir + "/" + name)
		buffer = append(buffer, string(goFile)+"\n")
		first = false
	}
	return strings.Join(buffer, "")
}
