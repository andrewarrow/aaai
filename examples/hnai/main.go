package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	s, _ := fetchStoriesSync()
	
	if len(os.Args) > 1 {
		id, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Invalid story ID: %s\n", os.Args[1])
			return
		}
		for _, story := range s {
			if story.ID == id {
				fmt.Printf("URL: %s\n", story.URL)
				return
			}
		}
		fmt.Printf("Story with ID %d not found\n", id)
	} else {
		for _, item := range s {
			fmt.Printf("ID: %d\nTitle: %s\n\n", item.ID, item.Title)
		}
	}
}
