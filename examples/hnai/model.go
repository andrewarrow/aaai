package main

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	stories []Story
	cursor  int
	width   int
	height  int
	content string
}

func initialModel() model {
	stories, err := fetchStoriesSync()
	if err != nil {
		stories = []Story{}
	}
	return model{
		stories: stories,
		cursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
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
				story := m.stories[m.cursor]
				if story.URL != "" {
					resp, err := http.Get(story.URL)
					if err == nil {
						defer resp.Body.Close()
						doc, err := html.Parse(resp.Body)
						if err == nil {
							m.content = extractText(doc)
						}
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}

	if m.content != "" {
		contentStyle := lipgloss.NewStyle().
			Width((m.width*80)/100 - 4).
			Height((m.height*80)/100 - 2)
		return contentStyle.Render(m.content)
	}

	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width((m.width*80)/100 - 4).
		Height((m.height*80)/100 - 2)

	style := lipgloss.NewStyle().
		Width((m.width*80)/100 - 6).
		Height(1)

	selectedStyle := lipgloss.NewStyle().
		Width((m.width*80)/100 - 6).
		Height(1).
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("15"))

	content := strings.Builder{}
	for i, story := range m.stories {
		if i == m.cursor {
			content.WriteString(selectedStyle.Render(story.Title))
		} else {
			content.WriteString(style.Render(story.Title))
		}
		content.WriteString("\n")
	}

	s.WriteString(tableStyle.Render(content.String()))
	s.WriteString("\nPress q to quit")
	return s.String()
}
