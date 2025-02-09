package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Story struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	By          string `json:"by"`
	TimeCreated int64  `json:"time"`
}

type model struct {
	stories  []Story
	cursor   int
	selected *Story
	err      error
	loading  bool
}

func initialModel() model {
	return model{
		stories: []Story{},
		loading: true,
	}
}

func (m model) Init() tea.Cmd {
	return fetchStories
}

func fetchStories() tea.Msg {
	// First, fetch the list of newest story IDs
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/newstories.json")
	if err != nil {
		return errMsg{err}
	}
	defer resp.Body.Close()

	var storyIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&storyIDs); err != nil {
		return errMsg{err}
	}

	// Fetch first 30 stories
	stories := make([]Story, 0, 30)
	for i := 0; i < 30 && i < len(storyIDs); i++ {
		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyIDs[i])
		resp, err := http.Get(url)
		if err != nil {
			continue
		}

		var story Story
		if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		stories = append(stories, story)
		time.Sleep(100 * time.Millisecond) // Be nice to the API
	}

	return storiesMsg(stories)
}

type storiesMsg []Story
type errMsg struct{ error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.stories)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.stories) > 0 {
				m.selected = &m.stories[m.cursor]
			}
		}

	case storiesMsg:
		m.stories = msg
		m.loading = false

	case errMsg:
		m.err = msg.error
		m.loading = false
	}

	return m, nil
}

func (m model) View() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#383838")).
		PaddingLeft(4)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#00FF00")).
		PaddingLeft(4)

	if m.loading {
		return "Loading stories...\n"
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	var s string
	s += "Newest Hacker News Stories\n\n"

	for i, story := range m.stories {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			timeAgo := time.Since(time.Unix(story.TimeCreated, 0)).Round(time.Minute)
			storyLine := fmt.Sprintf("%s %s\n   by %s | %v ago | %d points\n",
				cursor,
				story.Title,
				story.By,
				timeAgo,
				story.Score,
			)
			s += selectedStyle.Render(storyLine)
		} else {
			storyLine := fmt.Sprintf("%s %s\n", cursor, story.Title)
			s += style.Render(storyLine)
		}
	}

	if m.selected != nil {
		s += "\nSelected story URL:\n"
		s += m.selected.URL + "\n"
	}

	s += "\nPress q to quit, up/down to navigate, enter to select\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
