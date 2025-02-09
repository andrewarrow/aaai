package main

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	strings []string
	cursor  int
	width   int
	height  int
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func initialModel() model {
	strings := make([]string, 20)
	for i := range strings {
		strings[i] = generateRandomString(10)
	}
	return model{
		strings: strings,
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
			if m.cursor < len(m.strings)-1 {
				m.cursor++
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
		Background(lipgloss.Color("5"))

	content := strings.Builder{}
	for i, str := range m.strings {
		if i == m.cursor {
			content.WriteString(selectedStyle.Render(str))
		} else {
			content.WriteString(style.Render(str))
		}
		content.WriteString("\n")
	}

	s.WriteString(tableStyle.Render(content.String()))
	s.WriteString("\nPress q to quit")
	return s.String()
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
