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
}

func NewModel(content []byte, gap int) Model {
	model := Model{
		gapBuffer: NewGapBuffer(content, gap),
		index:     0,
		virtualX:  0,
	}

	return model
}

func (m *Model) MoveCursorX(d int) {
	if d < 0 {
		m.gapBuffer.Left(-d)
	} else {
		m.gapBuffer.Right(d)
	}
	m.virtualX = m.gapBuffer.gapLeft - m.findLineStart(m.gapBuffer.gapLeft)
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
	return fmt.Sprintf("%d/%d [%d] col: %d", m.gapBuffer.gapLeft, size, gap, m.virtualX)
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

	sIdx := m.findLineStart(m.gapBuffer.gapLeft)
	m.virtualX = m.gapBuffer.gapLeft - sIdx
}

func (m *Model) findLineEnd(cIdx int) int {
	if cIdx < m.gapBuffer.GetContentLen() && m.gapBuffer.GetByteAt(cIdx) != '\n' {
		nIdx := m.scanNewLine(cIdx, forward)
		return nIdx
	}
	return cIdx
}

func (m *Model) MoveCursorY(d int) {
	tVirtualX := m.virtualX
	if d == 0 {
		return
	} else if d < 0 {
		m.MoveCursorToLineStart()
		for i := 0; i > d && m.gapBuffer.gapLeft > 0; i-- {
			m.MoveCursorX(-1)
			m.MoveCursorToLineStart()
		}
	} else if d > 0 {
		for i := 0; i < d && m.gapBuffer.gapLeft < m.gapBuffer.GetContentLen(); i++ {
			m.MoveCursorToLineEnd()
			m.MoveCursorX(1)
		}
		m.MoveCursorToLineStart()
	}

	dVx := m.findLineEnd(m.gapBuffer.gapLeft) - m.gapBuffer.gapLeft

	if tVirtualX > dVx {
		m.MoveCursorX(dVx)
		m.virtualX = tVirtualX
	} else {
		m.MoveCursorX(tVirtualX)
	}
}

func (m *Model) GetPageByLines(totalLines int) (int, int, int) {
	sIdx := m.findLineStart(m.gapBuffer.gapLeft)
	eIdx := m.findLineEnd(m.gapBuffer.gapLeft)
	linesFound := 0
	keepScanning := true
	contentLen := m.gapBuffer.GetContentLen()
	for linesFound < totalLines && keepScanning {
		if sIdx > 0 {
			sIdx = m.findLineStart(sIdx - 1)
			linesFound++
		}

		if linesFound < totalLines && eIdx < contentLen {
			eIdx = m.findLineEnd(eIdx + 1)
			linesFound++
		}

		keepScanning = sIdx > 0 || eIdx < contentLen
	}

	return sIdx, eIdx, m.gapBuffer.gapLeft
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
