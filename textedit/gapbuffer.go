package textedit

import "fmt"

type GapBuffer struct {
	buffer   []byte
	size     int
	gapSize  int
	gapLeft  int
	gapRight int
}

// TODO add tests
func (g *GapBuffer) Left(n int) {
	for i := 0; i < n && g.gapLeft > 0; i++ {
		g.gapLeft--
		g.buffer[g.gapRight] = g.buffer[g.gapLeft]
		g.buffer[g.gapLeft] = 0
		g.gapRight--
	}
}

func (g *GapBuffer) Right(n int) {
	for i := 0; i < n && g.gapRight < g.size-1; i++ {
		g.gapRight++
		g.buffer[g.gapLeft] = g.buffer[g.gapRight]
		g.buffer[g.gapRight] = 0
		g.gapLeft++
	}
}

func (g *GapBuffer) Insert(c byte) {
	if g.gapLeft == g.gapRight {
		g.grow()
	}
	g.buffer[g.gapLeft] = c
	g.gapLeft++
}

func (g *GapBuffer) RemoveLeft() {
	if g.gapLeft > 0 {
		g.buffer[g.gapLeft] = 0
		g.gapLeft--
	}
}

func (g *GapBuffer) RemoveRight() {
	if g.gapRight < g.size-1 {
		g.buffer[g.gapRight] = 0
		g.gapRight++
	}
}

func (g *GapBuffer) grow() {
	newSize := g.size + g.gapSize - 1
	buffer := make([]byte, newSize)
	copy(buffer, g.buffer[0:g.gapLeft])
	newRight := g.gapSize + (g.gapLeft - 1)
	copy(buffer[newRight:], g.buffer[g.gapRight:])

	g.gapRight = newRight
	g.size = newSize
	g.buffer = buffer
}

func (g *GapBuffer) GetByteAt(pos int) byte {
	idx := pos
	if pos >= g.gapLeft {
		idx += g.gapRight - g.gapLeft + 1
	}
	if idx >= g.size {
		panic(fmt.Sprintf("Out of bounds. r %d c %d t %d", pos, idx, g.size))
	}
	return g.buffer[idx]
}

func (g *GapBuffer) GetContentLen() int {
	return g.size - (g.gapRight - g.gapLeft) - 1
}

func (g *GapBuffer) GetBytes() []byte {
	contentSize := g.GetContentLen()
	content := make([]byte, contentSize)
	copy(content[0:], g.buffer[0:g.gapLeft])
	copy(content[g.gapLeft:], g.buffer[g.gapRight+1:])
	return content
}

func NewGapBuffer(content []byte, gapSize int) GapBuffer {
	size := len(content) + gapSize
	buffer := make([]byte, size)
	copy(buffer[gapSize:], content)

	g := GapBuffer{
		buffer:   buffer,
		gapSize:  gapSize,
		gapLeft:  0,
		gapRight: gapSize - 1,
		size:     size,
	}

	return g
}
