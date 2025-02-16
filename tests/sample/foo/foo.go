package foo

import "fmt"

func Foo() {
	const (
		message = "Foo package: executing Foo function"
	)

	// Foo prints a message indicating the foo package execution
	fmt.Println(message)
}
