package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) > 1 {
		id, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Please provide a valid story ID")
			return
		}
		story, err := fetchStoryByID(id)
		if err != nil {
			fmt.Println("Error fetching story:", err)
			return
		}
		if story.URL != "" {
			fmt.Printf("URL for story %d: %s\n", id, story.URL)
		} else {
			fmt.Printf("Story %d has no URL\n", id)
		}
		return
	}

	s, _ := fetchStoriesSync()
	for _, item := range s {
		fmt.Println(item.ID)
		fmt.Println(item.Title)
	}
}
