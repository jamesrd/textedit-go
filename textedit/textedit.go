package textedit

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type state struct {
	model     Model
	fileName  string
	height    int
	width     int
	pageStart int
	message   string
	content   []byte
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
				m.message = err.Error()
			}
			return m, nil

		// navigation
		case "up":
			m.model.MoveCursorY(-1)
		case "pgup":
			m.model.MoveCursorY(-m.height / 2)
		case "down":
			m.model.MoveCursorY(1)
		case "pgdown":
			m.model.MoveCursorY(m.height / 2)
		case "home":
			m.model.MoveCursorToLineStart()
		case "end":
			m.model.MoveCursorToLineEnd()

		case "left":
			m.model.MoveCursorX(-1)
		case "right":
			m.model.MoveCursorX(1)

		// text editing
		case "enter":
			m.model.Insert('\n')
		case "backspace":
			m.model.gapBuffer.RemoveLeft()
		case "delete":
			m.model.gapBuffer.RemoveRight()
		case "esc":
			if len(m.message) == 0 {
				m.message = "No messages!"
			} else {
				m.message = ""
			}

		default:
			m.model.Insert(key[0])
		}

	}
	m.model.index = m.model.gapBuffer.gapLeft
	return m, nil
}

func (m state) View() string {
	var sb strings.Builder
	sb.WriteString(m.writeTitleLine())

	index := m.model.index
	contentLen := m.model.gapBuffer.GetContentLen()

	for i := 0; i < contentLen; i++ {
		b := m.model.gapBuffer.GetByteAt(i)
		if i == index {
			switch b {
			case '\n':
				sb.WriteString(ansiInvertRune(' '))
				sb.WriteRune('\n')
			case '\t':
				sb.WriteString(ansiInvertRune(' '))
				sb.WriteString("       ")
			default:
				sb.WriteString(ansiInvertByte(b))
			}
		} else {
			switch b {
			case '\t':
				sb.WriteString("        ")
			case 0:
				sb.WriteByte('^')
			default:
				sb.WriteByte(b)
			}
		}
	}
	if index == contentLen {
		sb.WriteString(ansiInvertRune(' '))
	}

	return sb.String()
}

func (m *state) writeTitleLine() string {
	titleColor := "\033[45m"
	errorColor := "\033[41m"
	resumeColor := "\033[0m"
	fm := ""
	if len(m.message) > 0 {
		fm = fmt.Sprintf(" %s !! %s%s", errorColor, m.message, titleColor)
	}
	l := m.model.index
	r := m.model.gapBuffer.gapRight
	c := m.model.gapBuffer.GetContentLen()
	s := m.model.gapBuffer.size

	return fmt.Sprintf("%sFile: %s - %d,%d/%d/%d%s%s\n", titleColor, m.fileName, l, r, c, s, fm, resumeColor)
}

func ansiInvertRune(c rune) string {
	return fmt.Sprintf("\033[07m%c\033[27m", c)
}

func ansiInvertByte(c byte) string {
	if c == 0 {
		c = '^'
	}
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
		m.model = NewModel(content)
	} else {
		// TODO make sure the file doesn't exist already
		m.fileName = "untitled.txt"
		m.model = NewModel([]byte{})
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
	fileContent := m.model.gapBuffer.GetBytes()
	return os.WriteFile(m.fileName, fileContent, 0644)
}
