package textedit

import (
	"testing"
)

func TestLeft(t *testing.T) {
	g := NewGapBuffer([]byte{'a'}, 1)
	g.Left(100)
	// also need to check gapright ... etc...
	if g.gapLeft != 0 {
		t.Fatalf(`Expected 0; got %v`, g.gapLeft)
	}
}

func TestRight(t *testing.T) {
	g := NewGapBuffer([]byte{'a'}, 1)
	g.Right(100)
	if g.gapLeft != g.GetContentLen() {
		t.Fatalf(`Expected 0; got %v`, g.gapLeft)
	}
}
