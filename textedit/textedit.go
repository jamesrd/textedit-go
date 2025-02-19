package textedit

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type state struct {
	model       Model
	fileName    string
	height      int
	width       int
	pageUpLines int
	message     string
}

func (m state) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.pageUpLines = m.height / 2

	case tea.KeyMsg:
		return m.processKey(msg)

	}
	return m, nil
}

func (m state) processKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key := msg.String(); key {
	// application control
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit

	case "ctrl+s":
		err := m.writeFile()
		if err != nil {
			m.message = err.Error()
		}

	// navigation
	case "up":
		m.model.MoveCursorY(-1)
	case "pgup":
		m.model.MoveCursorY(-m.pageUpLines)
	case "down":
		m.model.MoveCursorY(1)
	case "pgdown":
		m.model.MoveCursorY(m.pageUpLines)

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
		m.model.Backspace()
	case "delete":
		m.model.Delete()
	case "esc":
		if len(m.message) == 0 {
			m.message = "No messages!"
		} else {
			m.message = ""
		}
	case "tab":
		m.model.Insert('\t')

	default:
		if len(key) > 1 {
			m.message = fmt.Sprintf("Unhandled [%s]", key)
		} else {
			m.model.Insert(key[0])
		}
	}
	return m, nil
}

func (m state) View() string {
	var sb strings.Builder
	sb.WriteString(m.writeTitleLine())

	sIdx, eIdx, index := m.model.GetPageByLines(m.height - 2)

	linesWritten := 0

	for i := sIdx; i < eIdx; i++ {
		b := m.model.gapBuffer.GetByteAt(i)
		if i == index {
			switch b {
			case '\n':
				sb.WriteString(ansiInvertRune(' '))
				if linesWritten < m.height-1 {
					sb.WriteRune('\n')
				}
			case '\t':
				// TODO handle variable tabstop
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
			case '\n':
				if linesWritten < m.height-1 {
					sb.WriteByte(b)
				}
			default:
				sb.WriteByte(b)
			}
		}
		if b == '\n' {
			linesWritten++
		}

	}
	if index == eIdx {
		sb.WriteString(ansiInvertRune(' '))
	}
	for linesWritten < m.height-2 {
		sb.WriteString("\n ~")
		linesWritten++
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
	mt := m.model.GetStatus()

	return fmt.Sprintf("%sFile: %s - %s%s%s\n", titleColor, m.fileName, mt, fm, resumeColor)
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
	fileContent, _ := m.model.GetContent()
	return os.WriteFile(m.fileName, fileContent, 0644)
}
