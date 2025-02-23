package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-rod/rod"
)

func main() {
	s, _ := fetchStoriesSync()

	switch len(os.Args) {
	case 1:
		// No arguments - show story list
		showStoryList(s)
	case 2:
		// One argument - check if it's "screenshots"
		if os.Args[1] == "screenshots" {
			takeAllScreenshots(s)
		} else {
			fmt.Println("Usage:")
			fmt.Println("  No args: Show story list")
			fmt.Println("  screenshots: Take screenshots of all stories")
			fmt.Println("  story <ID>: Show and screenshot specific story")
		}
	case 3:
		// Two arguments - must be "story" and ID
		if os.Args[1] == "story" {
			handleSingleStory(s, os.Args[2])
		} else {
			fmt.Println("First argument must be 'story' followed by ID")
		}
	default:
		fmt.Println("Invalid number of arguments")
	}
}

func showStoryList(stories []Story) {
	for _, item := range stories {
		fmt.Printf("ID: %d\nTitle: %s\n\n", item.ID, item.Title)
	}
}

func handleSingleStory(stories []Story, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("Invalid story ID: %s\n", idStr)
		return
	}

	for _, story := range stories {
		if story.ID == id {
			fmt.Printf("URL: %s\n", story.URL)
			screenshot := captureScreenshot(story.URL)
			if screenshot != nil {
				saveScreenshot(screenshot, "screenshot.png")
			}
			return
		}
	}
	fmt.Printf("Story with ID %d not found\n", id)
}

func takeAllScreenshots(stories []Story) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	for _, story := range stories {
		fmt.Printf("Taking screenshot for story %d\n", story.ID)
		screenshot := captureScreenshot(story.URL)
		if screenshot != nil {
			filename := fmt.Sprintf("%d.jpg", story.ID)
			saveScreenshot(screenshot, filename)
		}
	}
}

func captureScreenshot(url string) []byte {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(url).MustWaitStable()
	page.MustWaitNavigation()

	err := page.WaitIdle(1000)
	if err != nil {
		fmt.Printf("Error waiting for page to stabilize: %v\n", err)
		return nil
	}
	screenshot, _ := page.Screenshot(false, nil)
	return screenshot
}

func saveScreenshot(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0644)
}
