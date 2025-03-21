package lipgloss_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

type blackhole struct{}

func (bh blackhole) Write(b []byte) (int, error) { return len(b), nil }

func TestPlaceHorizontal(t *testing.T) {
	testCases := []struct {
		w   int
		s   string
		pos lipgloss.Position
		exp string
	}{
		// odd spacing
		{w: 10, s: "Hello", pos: lipgloss.Left, exp: "Hello     "},
		{w: 10, s: "Hello", pos: 0, exp: "Hello     "},
		{w: 10, s: "Hello", pos: 0.000000001, exp: "Hello     "},
		{w: 10, s: "Hello", pos: lipgloss.Right, exp: "     Hello"},
		{w: 10, s: "Hello", pos: 1, exp: "     Hello"},
		{w: 10, s: "Hello", pos: 0.999999999, exp: "     Hello"},
		{w: 10, s: "Hello", pos: 0.49, exp: "  Hello   "},
		{w: 10, s: "Hello", pos: lipgloss.Center, exp: "  Hello   "},
		{w: 10, s: "Hello", pos: 0.51, exp: "   Hello  "},
	}

	for i, testCase := range testCases {
		r := lipgloss.NewRenderer(blackhole{})

		act := r.PlaceHorizontal(testCase.w, testCase.pos, testCase.s)
		exp := testCase.exp

		if exp != act {
			t.Errorf("Test %d: expected %q, got %q", i, exp, act)
		}
	}
}

func TestPlaceVertical(t *testing.T) {
	testCases := []struct {
		height   int
		content  string
		position lipgloss.Position
		expected string
	}{
		{height: 3, content: "Hello", position: lipgloss.Top, expected: "Hello\n     \n     "},
		{height: 3, content: "Hello", position: 0, expected: "Hello\n     \n     "},
		{height: 3, content: "Hello", position: 0.000000001, expected: "Hello\n     \n     "},
		{height: 3, content: "Hello", position: lipgloss.Bottom, expected: "     \n     \nHello"},
		{height: 3, content: "Hello", position: 1, expected: "     \n     \nHello"},
		{height: 3, content: "Hello", position: 0.999999999, expected: "     \n     \nHello"},
		{height: 4, content: "Hello", position: 0.49, expected: "     \nHello\n     \n     "},
		{height: 4, content: "Hello", position: lipgloss.Center, expected: "     \nHello\n     \n     "},
		{height: 4, content: "Hello", position: 0.51, expected: "     \n     \nHello\n     "},
	}

	for i, testCase := range testCases {
		r := lipgloss.NewRenderer(blackhole{})

		act := r.PlaceVertical(testCase.height, testCase.position, testCase.content)
		exp := testCase.expected

		if exp != act {
			t.Errorf("Test %d: expected %q, got %q", i, exp, act)
		}
	}
}
