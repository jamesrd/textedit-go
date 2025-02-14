package textedit

type Model struct {
	content []byte
	index   int
}

func (m *Model) MoveCursor(x int, y int) {
	if x != 0 {
		nindex := m.index + x
		nindex = max(0, nindex)
		m.index = min(nindex, len(m.content))
	}
	// TODO y will involve counting newlines
}
