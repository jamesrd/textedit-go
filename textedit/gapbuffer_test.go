package textedit

import (
	"reflect"
	"slices"
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

func TestInsert(t *testing.T) {
	gapSize := 10
	g := NewGapBuffer(basicTestContent, gapSize)

	var foo byte = 'a'
	g.Right(1)
	pos := g.gapLeft
	g.Insert(foo)
	assert(t, pos+1 == g.gapLeft, `Cursor did not move on insert %v - %v`, pos, g.gapLeft)
	equals(t, foo, g.GetByteAt(pos))
}

func TestInsertGrow(t *testing.T) {
	gapSize := 10
	eGap := gapSize - 1
	g := NewGapBuffer(basicTestContent, gapSize)

	i := 0
	cGap := g.gapRight - g.gapLeft
	equals(t, eGap, cGap)
	// insert until gap is closed
	for i < gapSize-1 {
		i++
		g.Insert('a')
		cGap = g.gapRight - g.gapLeft
		equals(t, eGap-i, cGap)
	}
	assert(t, cGap == 0, `Gap should be 0 but it's %v`, cGap)
	g.Insert('a')
	cGap = g.gapRight - g.gapLeft
	equals(t, eGap-1, cGap)

}

func TestGetByteAt(t *testing.T) {
	gapSize := 10
	g := NewGapBuffer(basicTestContent, gapSize)

	// loop with the gap buffer at all positions
	// check the byte at each

	for i := 0; i < len(basicTestContent); i++ {
		for x, c := range basicTestContent {
			gC := g.GetByteAt(x)
			assert(t, c == gC, `Expected %v got %v with gapLeft = %v`, c, gC, g.gapLeft)
		}
		g.Right(1)
	}

	// test range errors
	assertPanicGetByteAt(t, g.GetByteAt, g.size)
	assertPanicGetByteAt(t, g.GetByteAt, -1)

}

func TestRemoveLeft(t *testing.T) {
	gapSize := 10
	g := NewGapBuffer(basicTestContent, gapSize)

	// When all the way left, remove left doesn't remove anything
	g.RemoveLeft()
	equals(t, len(basicTestContent), g.GetContentLen())
	equals(t, basicTestContent[0], g.GetByteAt(0))

	// Move right one, then remove left. gap[0] == source[1]
	g.Right(1)
	g.RemoveLeft()
	equals(t, len(basicTestContent)-1, g.GetContentLen())
	equals(t, basicTestContent[1], g.GetByteAt(0))
}

func TestRemoveRight(t *testing.T) {
	gapSize := 10
	g := NewGapBuffer(basicTestContent, gapSize)

	// When all the way left, after remove right: gap[1] = source[2]
	g.RemoveRight()
	equals(t, len(basicTestContent)-1, g.GetContentLen())
	equals(t, basicTestContent[2], g.GetByteAt(1))

	// Move all the way to the right, remove right won't remove anything
	g.Right(g.GetContentLen())
	g.RemoveRight()
	equals(t, len(basicTestContent)-1, g.GetContentLen())
}

func TestGetBytes(t *testing.T) {
	gapSize := 10
	g := NewGapBuffer(basicTestContent, gapSize)

	equals(t, basicTestContent, g.GetBytes())

	var foo byte = 'a'
	g.Insert(foo)
	equals(t, slices.Insert(basicTestContent, 0, foo), g.GetBytes())

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

// TODO try to make this generic
func assertPanicGetByteAt(t testing.TB, f func(v int) byte, v int) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f(v)
}
