package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	content   []string
	cursor_x  int
	virtual_x int
	cursor_y  int
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor_y > 0 {
				m.cursor_y--
			}

		case "down":
			if m.cursor_y < len(m.content)-1 {
				m.cursor_y++
			}

		case "left":
			if m.cursor_x > 0 {
				m.cursor_x--
				m.virtual_x = m.cursor_x
			}
		case "right":
			if m.cursor_x < len(m.content[m.cursor_y])-1 {
				m.cursor_x++
				m.virtual_x = m.cursor_x
			}
		}
		max_x := max(len(m.content[m.cursor_y])-1, 0)
		m.cursor_x = min(max_x, max(m.cursor_x, m.virtual_x))

	}
	return m, nil
}

func (m model) View() string {
	var sb strings.Builder
	for y, line := range m.content {
		if y == m.cursor_y && len(line) == 0 {
			sb.WriteString(">\n")
		} else {
			for x, c := range line {
				if x == m.cursor_x && y == m.cursor_y {
					sb.WriteRune('>')
				} else {
					sb.WriteRune(c)
				}
			}
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func initModelWithText(s string) model {
	return model{
		content:   strings.Split(s, "\n"),
		cursor_y:  0,
		virtual_x: 0,
		cursor_x:  0,
	}
}

func main() {
	p := tea.NewProgram(initModelWithText("Hello!\n\nPress q to quit."))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
