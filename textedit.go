package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	content   []string
	cursorX   int
	virtualX  int
	cursorY   int
	height    int
	width     int
	pageStart int
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:

		switch key := msg.String(); key {
		// application control
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		// navigation
		case "up":
			if m.cursorY > 0 {
				m.cursorY--
			}
		case "pgup":
			m.cursorY = max(0, m.cursorY-m.height/2)

		case "down":
			if m.cursorY < len(m.content)-1 {
				m.cursorY++
			}
		case "pgdown":
			m.cursorY = min(len(m.content)-1, m.cursorY+m.height/2)

		case "home":
			m.cursorX = 0
			m.virtualX = m.cursorX

		case "end":
			m.cursorX = len(m.content[m.cursorY]) - 1
			m.virtualX = m.cursorX

		case "left":
			if m.cursorX > 0 {
				m.cursorX--
				m.virtualX = m.cursorX
			}
		case "right":
			if m.cursorX < len(m.content[m.cursorY])-1 {
				m.cursorX++
				m.virtualX = m.cursorX
			}

		// text editing
		case "enter":
			// TODO: break line
			m.content = slices.Insert(m.content, m.cursorY, "")
		default:
			m.content = slices.Insert(m.content, m.cursorY, key)
		}
		max_x := max(len(m.content[m.cursorY])-1, 0)
		m.cursorX = min(max_x, max(m.cursorX, m.virtualX))

		// set up the pageStart
		if m.cursorY < m.pageStart {
			m.pageStart = m.cursorY
		} else if m.cursorY >= m.pageStart+m.height-2 {
			m.pageStart = max(0, m.cursorY-(m.height-2))
		}

	}
	return m, nil
}

func (m model) View() string {
	var sb strings.Builder

	for y := m.pageStart; y < m.pageStart+m.height-1; y++ {
		line := "\033[90m -\033[0m"
		if y < len(m.content) {
			line = m.content[y]
		}
		if y == m.cursorY && len(line) == 0 {
			sb.WriteString(ansiInvertRune(' '))
		} else {
			for x, c := range line {
				if x == m.cursorX && y == m.cursorY {
					sb.WriteString(ansiInvertRune(c))
				} else {
					sb.WriteRune(c)
				}
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func ansiInvertRune(c rune) string {
	return fmt.Sprintf("\033[07m%c\033[27m", c)
}

func initModelWithText(s string) model {
	return model{
		content:  strings.Split(s, "\n"),
		cursorY:  0,
		virtualX: 0,
		cursorX:  0,
	}
}

func main() {
	p := tea.NewProgram(initModelWithText("Hello!\n\nPress ctrl+q to quit."))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
