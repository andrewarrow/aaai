package bar

import "fmt"

const (
	message = "Bar package: executing Bar function"
)

// Bar prints a message indicating the bar package execution
func Bar() {
	fmt.Println(message)
}
