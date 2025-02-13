package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	content   []string
	fileName  string
	cursorX   int
	virtualX  int
	cursorY   int
	height    int
	width     int
	pageStart int
}

func (m model) Init() tea.Cmd {
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

		case "ctrl+s":
			err := m.writeFile()
			if err != nil {
				panic(err)
			}
			return m, nil

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

func initModelWithFile(fileName string) model {
	m := model{}
	if len(fileName) > 0 {
		m.fileName = fileName
		m.content = strings.Split(readFile(fileName), "\n")
	} else {
		// TODO make sure the file doesn't exist already
		m.fileName = "untitled.txt"
		m.content = []string{""}
	}
	return m
}

func readFile(name string) string {
	var contents, err = os.ReadFile(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		} else {
			// TODO find other errors to handle
			return err.Error()
		}

	}
	return string(contents)
}

func (m model) writeFile() error {
	fileContent := strings.Join(m.content, "\n")
	return os.WriteFile(m.fileName, []byte(fileContent), 0644)
}

func main() {
	var fileName string
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	p := tea.NewProgram(initModelWithFile(fileName))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
