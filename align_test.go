package lipgloss

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestAlignTextVertical(t *testing.T) {
	tests := []struct {
		str    string
		pos    Position
		height int
		want   string
	}{
		{str: "Foo", pos: Top, height: 2, want: "Foo\n"},
		{str: "Foo", pos: Center, height: 5, want: "\n\nFoo\n\n"},
		{str: "Foo", pos: Bottom, height: 5, want: "\n\n\n\nFoo"},

		{str: "Foo\nBar", pos: Bottom, height: 5, want: "\n\n\nFoo\nBar"},
		{str: "Foo\nBar", pos: Center, height: 5, want: "\nFoo\nBar\n\n"},
		{str: "Foo\nBar", pos: Top, height: 5, want: "Foo\nBar\n\n\n"},

		{str: "Foo\nBar\nBaz", pos: Bottom, height: 5, want: "\n\nFoo\nBar\nBaz"},
		{str: "Foo\nBar\nBaz", pos: Center, height: 5, want: "\nFoo\nBar\nBaz\n"},

		{str: "Foo\nBar\nBaz", pos: Bottom, height: 3, want: "Foo\nBar\nBaz"},
		{str: "Foo\nBar\nBaz", pos: Center, height: 3, want: "Foo\nBar\nBaz"},
		{str: "Foo\nBar\nBaz", pos: Top, height: 3, want: "Foo\nBar\nBaz"},

		{str: "Foo\n\n\n\nBar", pos: Bottom, height: 5, want: "Foo\n\n\n\nBar"},
		{str: "Foo\n\n\n\nBar", pos: Center, height: 5, want: "Foo\n\n\n\nBar"},
		{str: "Foo\n\n\n\nBar", pos: Top, height: 5, want: "Foo\n\n\n\nBar"},

		{str: "Foo\nBar\nBaz", pos: Center, height: 9, want: "\n\n\nFoo\nBar\nBaz\n\n\n"},
		{str: "Foo\nBar\nBaz", pos: Center, height: 10, want: "\n\n\nFoo\nBar\nBaz\n\n\n\n"},
	}

	for _, test := range tests {
		got := alignTextVertical(test.str, test.pos, test.height, nil)
		if got != test.want {
			t.Errorf("alignTextVertical(%q, %v, %d) = %q, want %q", test.str, test.pos, test.height, got, test.want)
		}
	}
}

func TestAlignTextHorizontalStyledPadding(t *testing.T) {
	tests := []struct {
		name string
		pos  Position
		want string
	}{
		{name: "left", pos: Left, want: "x    "},
		{name: "center", pos: Center, want: "  x  "},
		{name: "right", pos: Right, want: "    x"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rendered := NewStyle().
				Width(5).
				Background(Color("#ff0000")).
				Align(test.pos).
				Render("x")

			if got := ansi.Strip(rendered); got != test.want {
				t.Errorf("Align(%v) = %q, want %q", test.pos, got, test.want)
			}
			if got := ansi.StringWidth(rendered); got != 5 { //nolint:mnd
				t.Errorf("Align(%v) width = %d, want 5", test.pos, got)
			}
		})
	}
}

func TestAlignTextHorizontalPadsToWidestLine(t *testing.T) {
	rendered := NewStyle().
		Width(6).
		Background(Color("#ff0000")).
		Align(Right).
		Render("a\nbbb")

	lines := strings.Split(ansi.Strip(rendered), "\n")
	want := []string{"     a", "   bbb"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d", len(lines), len(want))
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestAlignTextVerticalShorterThanContent(t *testing.T) {
	got := alignTextVertical("Foo\nBar\nBaz", Center, 2, nil)
	want := "Foo\nBar\nBaz"
	if got != want {
		t.Errorf("alignTextVertical() = %q, want %q", got, want)
	}
}
