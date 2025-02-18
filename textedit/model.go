package textedit

import "fmt"

type direction int

const (
	backward direction = iota
	forward
)

type Model struct {
	gapBuffer GapBuffer
	index     int
	virtualX  int
	tabstop   int
}

func NewModel(content []byte) Model {
	model := Model{
		index:    0,
		virtualX: 0,
		tabstop:  8,
	}
	// TODO: probably want a larger gap size
	model.gapBuffer = NewGapBuffer(content, 100)

	return model
}

func (m *Model) MoveCursorX(d int) {
	if d < 0 {
		m.gapBuffer.Left(-d)
	} else {
		m.gapBuffer.Right(d)
	}
}

func (m *Model) Insert(b byte) {
	m.gapBuffer.Insert(b)
}

func (m *Model) Backspace() {
	m.gapBuffer.RemoveLeft()
}

func (m *Model) Delete() {
	m.gapBuffer.RemoveRight()
}

func (m *Model) GetContent() ([]byte, int) {
	return m.gapBuffer.GetBytes(), m.gapBuffer.gapLeft
}

func (m *Model) GetStatus() string {
	gap := m.gapBuffer.gapRight - m.gapBuffer.gapLeft + 1
	size := m.gapBuffer.size - gap
	return fmt.Sprintf("%d/%d [%d]", m.gapBuffer.gapLeft, size, gap)
}

func (m *Model) MoveCursorToLineStart() {
	target := m.findLineStart(m.gapBuffer.gapLeft)
	m.gapBuffer.Left(m.gapBuffer.gapLeft - target)
	m.virtualX = 0
}

func (m *Model) findLineStart(cIdx int) int {
	rIdx := m.scanNewLine(cIdx, backward)
	if rIdx > 0 {
		rIdx++
	}
	return rIdx
}

func (m *Model) MoveCursorToLineEnd() {
	target := m.findLineEnd(m.gapBuffer.gapLeft)
	m.gapBuffer.Right(target - m.gapBuffer.gapLeft)

	m.virtualX = 0
}

func (m *Model) findLineEnd(cIdx int) int {
	if cIdx < m.gapBuffer.GetContentLen() && m.gapBuffer.GetByteAt(cIdx) != '\n' {
		nIdx := m.scanNewLine(cIdx, forward)
		return nIdx
	}
	return cIdx
}

// TODO make following functions work with gapBuffer

func (m *Model) MoveCursorY(d int) {
	var newLineStart int
	if d == 0 {
		return
	} else if d < 0 {
		newLineStart = m.findPreviousLineStart(m.index)
		for i := -1; i > d; i-- {
			newLineStart = m.findPreviousLineStart(newLineStart)
		}
	} else if d > 0 {
		newLineStart = m.findNextLineStart(m.index)
		for i := 1; i < d; i++ {
			newLineStart = m.findNextLineStart(newLineStart)
		}
	}

	newLineEnd := m.findLineEnd(newLineStart)
	if newLineStart+m.virtualX < newLineEnd {
		m.index = newLineStart + m.virtualX
	} else {
		m.index = newLineStart
	}
}

func (m *Model) findPreviousLineStart(idx int) int {
	cls := m.findLineStart(idx)
	if cls > 0 {
		cls = m.findLineStart(cls - 1)
	}
	return cls
}

func (m *Model) findNextLineStart(idx int) int {
	// if idx >= len(m.content) {
	// 	return len(m.content)
	// }
	// if m.content[idx] == '\n' {
	// 	return idx + 1
	// }
	// cls := m.scanNewLine(idx, forward) + 1
	// if cls > len(m.content) {
	// 	cls = m.findLineStart(cls - 1)
	// }
	// return cls
	return idx
}

func (m *Model) scanNewLine(c int, d direction) int {
	inc := -1
	if d == forward {
		inc = 1
	}

	idx := c
	found := false

	for !found && idx+inc >= 0 && idx+inc < m.gapBuffer.GetContentLen() {
		idx += inc
		found = '\n' == m.gapBuffer.GetByteAt(idx)
	}
	return idx
}
