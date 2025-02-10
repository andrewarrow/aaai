package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Story struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	By          string `json:"by"`
	TimeCreated int64  `json:"time"`
}

func fetchStoriesSync() ([]Story, error) {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/newstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var storyIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&storyIDs); err != nil {
		return nil, err
	}

	stories := make([]Story, 30)
	for i := 0; i < 30; i++ {
		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyIDs[i])
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var story Story
		if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
			return nil, err
		}
		stories[i] = story
	}

	return stories[0:20], nil
}
