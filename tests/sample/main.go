package main

import (
	"fmt"
	"sample/bar"
	"sample/foo"
)

func main() {
	fmt.Println("sample")
	foo.Foo()
	bar.Bar()
}
