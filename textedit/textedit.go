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
			// TODO: break line
		case "backspace":
			// TODO implement removing previous character
		case "delete":
			//TODO implement removing current character
		case "esc":
			if len(m.message) == 0 {
				m.message = "No messages!"
			} else {
				m.message = ""
			}

		default:
			// TODO actually insert the character
		}

	}
	return m, nil
}

func (m state) View() string {
	var sb strings.Builder
	sb.WriteString(writeTitleLine(m.fileName, m.message))
	content, index := m.model.GetContent()

	for i := 0; i < len(content); i++ {
		b := content[i]
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
			default:
				sb.WriteByte(b)
			}
		}
	}
	if index == len(content) {
		sb.WriteString(ansiInvertRune(' '))
	}

	return sb.String()
}

func writeTitleLine(fileName string, message string) string {
	titleColor := "\033[42m"
	errorColor := "\033[41m"
	resumeColor := "\033[0m"
	fm := ""
	if len(message) > 0 {
		fm = fmt.Sprintf(" %s !! %s%s", errorColor, message, titleColor)
	}
	return fmt.Sprintf("%sFile: %s%s\n%s", titleColor, fileName, fm, resumeColor)
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
		m.model = Model{
			content: content,
		}
	} else {
		// TODO make sure the file doesn't exist already
		m.fileName = "untitled.txt"
		m.model = Model{
			content: []byte{},
		}
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
	fileContent, _ := m.model.GetContent()
	return os.WriteFile(m.fileName, fileContent, 0644)
}
