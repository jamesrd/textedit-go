package textedit

import (
	"reflect"
	"testing"
)

var basicTestContent = []byte("0123456789")

func checkGapPositions(t *testing.T, action string, eLeft int, eRight int, gapLeft int, gapRight int) {
	t.Helper()
	if gapLeft != eLeft || gapRight != eRight {
		t.Errorf(`Move %s: Expected %d %d - got %d %d`, action, eLeft, eRight, gapLeft, gapRight)
	}
}

func TestRightLeftSteps(t *testing.T) {
	const gapSize int = 2
	g := NewGapBuffer(basicTestContent, gapSize)
	assert(t, g.size == len(basicTestContent)+gapSize, `GapBuffer size not correct %v`, g.size)

	eLeft := 0
	eRight := gapSize - 1
	for eLeft < g.GetContentLen() {
		eLeft++
		eRight++
		g.Right(1)
		checkGapPositions(t, "Right", eLeft, eRight, g.gapLeft, g.gapRight)
	}

	// moving one more won't change the position
	g.Right(1)
	if g.gapLeft != eLeft || g.gapRight != eRight {
		t.Errorf(`Move Right: Expected %v %v - got %v %v`, eLeft, eRight, g.gapLeft, g.gapRight)
	}

	for eLeft > 0 {
		eLeft--
		eRight--
		g.Left(1)
		if g.gapLeft != eLeft || g.gapRight != eRight {
			t.Errorf(`Move Left: Expected %v %v - got %v %v`, eLeft, eRight, g.gapLeft, g.gapRight)
		}
	}

	// moving one more won't change the position
	g.Left(1)
	if g.gapLeft != eLeft || g.gapRight != eRight {
		t.Errorf(`Move Left: Expected %v %v - got %v %v`, eLeft, eRight, g.gapLeft, g.gapRight)
	}

}

func TestRightLeftJump(t *testing.T) {
	gapSize := 10
	g := NewGapBuffer(basicTestContent, gapSize)
	jump := 100
	eLeft := len(basicTestContent)
	eRight := eLeft + gapSize - 1
	g.Right(jump)
	if g.gapLeft != eLeft || g.gapRight != eRight {
		t.Errorf(`Move Right: Expected %v %v - got %v %v`, eLeft, eRight, g.gapLeft, g.gapRight)
	}

	eLeft = 0
	eRight = eLeft + gapSize - 1
	g.Left(jump)
	if g.gapLeft != eLeft || g.gapRight != eRight {
		t.Errorf(`Move Left: Expected %v %v - got %v %v`, eLeft, eRight, g.gapLeft, g.gapRight)
	}

}

// test helpers
// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()
	if !condition {
		tb.Fatalf(msg, v...)
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Errorf("unexpected error: %s", err.Error())
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(exp, act) {
		tb.Errorf("exp: %#v\n\n\tgot: %#v", exp, act)
	}
}
