package textedit

type direction int

const (
	backward direction = iota
	forward
)

type Model struct {
	content  []byte
	index    int
	virtualX int
	tabstop  int
}

func (m *Model) MoveCursorX(d int) {
	nindex := m.index + d
	nindex = max(0, nindex)
	m.index = min(nindex, len(m.content))
	m.virtualX = m.index - m.findLineStart(m.index)
}

func (m *Model) MoveCursorToLineStart() {
	m.index = m.findLineStart(m.index)
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
	m.index = m.findLineEnd(m.index)
	m.virtualX = m.index - m.findLineStart(m.index)
}

func (m *Model) findLineEnd(cIdx int) int {
	if cIdx < len(m.content) && m.content[cIdx] != '\n' {
		return m.scanNewLine(m.index, forward)
	}
	return cIdx
}

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
	if idx >= len(m.content) {
		return len(m.content)
	}
	if m.content[idx] == '\n' {
		return idx + 1
	}
	cls := m.scanNewLine(idx, forward) + 1
	if cls > len(m.content) {
		cls = m.findLineStart(cls - 1)
	}
	return cls
}

func (m *Model) GetContent() ([]byte, int) {
	return m.content, m.index
}

func (m *Model) GetLines(idx int, lines int) ([]byte, int) {
	startIdx, _ := m.countLinesBetween(idx, m.index, lines)
	endIdx := m.findNthNewline(startIdx, lines)
	return m.content[startIdx:endIdx], startIdx
}

func (m *Model) countLinesBetween(startIdx int, endIdx int, maxLines int) (int, int) {
	lines := 0
	if endIdx < startIdx {
		startIdx = m.findLineStart(endIdx)
	}
	idx := min(endIdx, len(m.content)-1)
	for lines < maxLines && idx >= startIdx && idx >= 0 {
		if m.content[idx] == '\n' {
			lines++
		}
		idx--
	}
	if lines == maxLines {
		return idx + 1, lines
	}
	return startIdx, lines
}

func (m *Model) findNthNewline(start int, n int) int {
	idx := start
	length := len(m.content)

	lines := 0

	for lines < n && idx < length {
		if m.content[idx] == '\n' {
			lines++
		}
		idx++
	}
	return idx
}

func (m *Model) scanNewLine(c int, d direction) int {
	inc := -1
	if d == forward {
		inc = 1
	}

	idx := c
	found := false

	for !found && idx+inc >= 0 && idx+inc < len(m.content) {
		idx += inc
		found = '\n' == m.content[idx]
	}
	return idx
}
