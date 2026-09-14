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

func TestAlignTextHorizontal(t *testing.T) {
	for _, tt := range []struct {
		name, input, want string
		pos               Position
		width             int
	}{
		{"left", "one\nlonger", "one   \nlonger", Left, 0},
		{"right", "one\nlonger", "     one\n  longer", Right, 8},
		{"center", "one\nlonger", "  one   \n longer ", Center, 8},
		{"narrow", "one\nlonger", "one   \nlonger", Left, 2},
		{"empty", "", "   ", Left, 3},
		{"trailing newline", "a\n", "a\n ", Left, 0},
		{"normalization", "a\tb\r\nx", "a    b\nx     ", Left, 0},
		{"unicode", "日本語\n🐳e\u0301", "日本語\n🐳e\u0301   ", Left, 0},
		{"ansi", "\x1b[31mx\x1b[m\nabc", "\x1b[31mx\x1b[m  \nabc", Left, 0},
		{"link", "\x1b]8;;https://example.com\ax\x1b]8;;\a\nabc", "\x1b]8;;https://example.com\ax\x1b]8;;\a  \nabc", Left, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := alignTextHorizontal(tt.input, tt.pos, tt.width, nil); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
	style := ansi.NewStyle().BackgroundColor(ansi.Red)
	if got, want := alignTextHorizontal("a\nabc", Left, 0, &style), "a"+style.Styled("  ")+"\nabc"; got != want {
		t.Fatalf("styled padding: got %q, want %q", got, want)
	}
}

func BenchmarkAlignTextHorizontal(b *testing.B) {
	for _, tt := range []struct{ name, input string }{
		{"plain", "short\na longer line"},
		{"ansi", strings.Repeat("\x1b[38;2;120;180;240m日本語 🐳 styled content\x1b[m\nshort\n", 20)},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = alignTextHorizontal(tt.input, Left, 80, nil)
			}
		})
	}
}
