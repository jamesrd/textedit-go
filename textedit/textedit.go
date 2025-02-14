package textedit

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type state struct {
	editor    Model
	content   []string
	fileName  string
	cursorX   int
	virtualX  int
	cursorY   int
	height    int
	width     int
	pageStart int
}

func (m state) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.editor.MoveCursor(-1, 0)
		case "right":
			m.editor.MoveCursor(1, 0)

		// text editing
		case "enter":
			// TODO: break line
			m.content = slices.Insert(m.content, m.cursorY, "")
		case "backspace":
			// TODO implement removing previous character
			m.content = slices.Delete(m.content, m.cursorY, m.cursorY+1)
			if len(m.content) == 0 {
				m.content = append(m.content, "")
			}
			m.cursorY = min(m.cursorY, len(m.content)-1)
		case "delete":
			//TODO implement removing current character
			m.content = slices.Delete(m.content, m.cursorY, m.cursorY+1)
			if len(m.content) == 0 {
				m.content = append(m.content, "")
			}
			m.cursorY = min(m.cursorY, len(m.content)-1)
		case "esc":

		default:
			// TODO actually insert the character
			if len(key) == 1 {
				m.content = slices.Insert(m.content, m.cursorY, key)
			} else if key == "tab" {
				m.content = slices.Insert(m.content, m.cursorY, "\t")
			}
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

func (m state) View() string {
	var sb strings.Builder
	ed := &m.editor

	for i := 0; i < len(ed.content); i++ {
		b := ed.content[i]
		if i == ed.index {
			switch b {
			case '\n':
				sb.WriteString(ansiInvertRune(' '))
				sb.WriteRune('\n')
			default:
				sb.WriteString(ansiInvertByte(b))
			}
		} else {
			sb.WriteByte(b)
		}
	}
	if ed.index == len(ed.content) {
		sb.WriteString(ansiInvertRune(' '))
	}

	return sb.String()
}

func (m state) View2() string {
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

func ansiInvertByte(c byte) string {
	return fmt.Sprintf("\033[07m%c\033[27m", c)
}

func InitModelWithFile(fileName string) state {
	m := state{}
	if len(fileName) > 0 {
		m.fileName = fileName
		content, err := readFile(fileName)
		if err != nil {
			panic(err)
		}
		m.editor = Model{
			content: content,
		}
		m.content = strings.Split(string(content), "\n")
	} else {
		// TODO make sure the file doesn't exist already
		m.fileName = "untitled.txt"
		m.content = []string{""}
	}
	return m
}

func readFile(name string) ([]byte, error) {
	var contents, err = os.ReadFile(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil
		} else {
			// TODO find other errors to handle
			return []byte{}, err
		}

	}
	return contents, nil
}

func (m state) writeFile() error {
	fileContent := strings.Join(m.content, "\n")
	return os.WriteFile(m.fileName, []byte(fileContent), 0644)
}
