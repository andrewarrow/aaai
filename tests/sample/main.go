package main

import (
	"fmt"
	"sample/bar"
	"sample/foo"
)

const (
	mainMessage = "Starting sample application"
)

// main is the entry point of the application
func main() {
	fmt.Println(mainMessage)
	foo.Foo()
	bar.Bar()
}
