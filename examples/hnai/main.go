package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/go-rod/rod"
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
				// Fetch HTML content using rod
				browser := rod.New().MustConnect()
				defer browser.MustClose()
				
				page := browser.MustPage(story.URL).MustWaitStable()
				
				// Wait for network to be idle and DOM to be ready
				page.MustWaitNavigation()
				
				// Wait a bit longer for any delayed JavaScript execution
				err = page.WaitIdle(5000)
				if err != nil {
					fmt.Printf("Error waiting for page to stabilize: %v\n", err)
					return
				}
				
				html, err := page.HTML()
				if err != nil {
					fmt.Printf("Error fetching HTML: %v\n", err)
					return
				}
				fmt.Printf("HTML Content:\n%s\n", html)
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
