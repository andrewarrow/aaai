package main

import (
	"fmt"
)

func main() {
	s, _ := fetchStoriesSync()
	for _, item := range s {
		fmt.Println(item.ID)
		fmt.Println(item.Title)
	}
}
