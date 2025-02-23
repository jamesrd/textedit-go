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
	tabstop     int
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
	lines, y, x := m.model.GetPageLines(m.height - 2)
	sb.WriteString(m.writeTitleLine(y+1, x+1))

	linesWritten := len(lines)
	maxLength := m.width - 1

	for cY, line := range lines {
		if cY > 0 {
			sb.WriteRune('\n')
		}
		wLine := line
		if cY == y {
			lLen := len(line)
			aX := x
			if lLen > maxLength {
				sl := x - (maxLength / 2)
				el := x + (maxLength / 2)
				if sl < 0 {
					el += -sl
					sl = 0
				}
				if el > lLen {
					sl -= (el - lLen)
					el = lLen
				}
				wLine = line[sl:el]

				aX = x - sl
			}
			if aX < len(wLine) {
				wLine = wLine[:aX] + "\033[07m" + string(wLine[aX]) + "\033[27m" + wLine[aX+1:]
			}
		}
		sb.WriteString(wLine)
		if y == cY && x == len(line) {
			sb.WriteString(ansiInvertRune(' '))
		}
	}
	for linesWritten < m.height-1 {
		sb.WriteString("\n ~")
		linesWritten++
	}

	return sb.String()
}

func (m *state) writeTitleLine(line int, col int) string {
	titleColor := "\033[45m\033[30m"
	errorColor := "\033[101m\033[97m"
	resumeColor := "\033[0m"

	mt := m.model.GetStatus()
	titleBase := fmt.Sprintf(" %s - y: %d x: %d %s ", m.fileName, line, col, mt)
	messgeBase := ""
	if len(m.message) > 0 {
		messgeBase = fmt.Sprintf("<! %s !> ", m.message)
	}

	padLen := m.width - (len(titleBase) + len(messgeBase))

	var sb strings.Builder
	sb.WriteString(titleColor)
	sb.WriteString(titleBase)
	sb.WriteString(errorColor)
	sb.WriteString(messgeBase)
	sb.WriteString(titleColor)
	sb.WriteString(fmt.Sprintf("%*c", padLen-1, ' '))
	sb.WriteString(resumeColor)
	sb.WriteRune('\n')

	return sb.String()

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

func InitModelWithFile(fileName string, gap int) state {
	m := state{
		tabstop: 4,
	}
	if len(fileName) > 0 {
		m.fileName = fileName
		content, err := readFile(fileName)
		if err != nil {
			panic(err)
		}
		m.model = NewModel(content, gap)
	} else {
		m.fileName = getUnusedFilename("untitled", "txt")
		m.model = NewModel([]byte{}, gap)
	}
	return m
}

func readFile(name string) ([]byte, error) {
	var contents, err = os.ReadFile(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil
		} else {
			return []byte{}, err
		}

	}
	return contents, nil
}

func getUnusedFilename(base string, ext string) string {
	name := fmt.Sprintf("%s.%s", base, ext)
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		return name
	}
	found := false
	i := 0
	for !found {
		i++
		name = fmt.Sprintf("%s%d.%s", base, i, ext)
		_, err := os.Stat(name)
		found = errors.Is(err, os.ErrNotExist)
	}
	return name
}

func (m state) writeFile() error {
	fileContent, _ := m.model.GetContent()
	return os.WriteFile(m.fileName, fileContent, 0644)
}
