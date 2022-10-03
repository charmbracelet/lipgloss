package lipgloss

import (
	"reflect"
	"testing"

	"github.com/muesli/termenv"
)

func TestStyleRender(t *testing.T) {
	renderer.SetColorProfile(termenv.TrueColor)
	t.Parallel()

	tt := []struct {
		style    Style
		expected string
	}{
		{
			NewStyle().Foreground(Color("#5A56E0")),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
		{
			NewStyle().Bold(true),
			"\x1b[1mhello\x1b[0m",
		},
		{
			NewStyle().Italic(true),
			"\x1b[3mhello\x1b[0m",
		},
		{
			NewStyle().Underline(true),
			"\x1b[4;4mh\x1b[0m\x1b[4;4me\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4mo\x1b[0m",
		},
		{
			NewStyle().Blink(true),
			"\x1b[5mhello\x1b[0m",
		},
		{
			NewStyle().Faint(true),
			"\x1b[2mhello\x1b[0m",
		},
	}

	for i, tc := range tt {
		s := tc.style.Copy().SetString("hello")
		res := s.Render()
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}
}

func TestValueCopy(t *testing.T) {
	t.Parallel()

	s := NewStyle().
		Bold(true)

	i := s
	i.Bold(false)

	requireEqual(t, s.GetBold(), i.GetBold())
}

func TestStyleInherit(t *testing.T) {
	t.Parallel()

	s := NewStyle().
		Bold(true).
		Italic(true).
		Underline(true).
		Strikethrough(true).
		Blink(true).
		Faint(true).
		Foreground(Color("#ffffff")).
		Background(Color("#111111")).
		Margin(1, 1, 1, 1).
		Padding(1, 1, 1, 1)

	i := NewStyle().Inherit(s)

	requireEqual(t, s.GetBold(), i.GetBold())
	requireEqual(t, s.GetItalic(), i.GetItalic())
	requireEqual(t, s.GetUnderline(), i.GetUnderline())
	requireEqual(t, s.GetStrikethrough(), i.GetStrikethrough())
	requireEqual(t, s.GetBlink(), i.GetBlink())
	requireEqual(t, s.GetFaint(), i.GetFaint())
	requireEqual(t, s.GetForeground(), i.GetForeground())
	requireEqual(t, s.GetBackground(), i.GetBackground())

	requireNotEqual(t, s.GetMarginLeft(), i.GetMarginLeft())
	requireNotEqual(t, s.GetMarginRight(), i.GetMarginRight())
	requireNotEqual(t, s.GetMarginTop(), i.GetMarginTop())
	requireNotEqual(t, s.GetMarginBottom(), i.GetMarginBottom())
	requireNotEqual(t, s.GetPaddingLeft(), i.GetPaddingLeft())
	requireNotEqual(t, s.GetPaddingRight(), i.GetPaddingRight())
	requireNotEqual(t, s.GetPaddingTop(), i.GetPaddingTop())
	requireNotEqual(t, s.GetPaddingBottom(), i.GetPaddingBottom())
}

func TestStyleCopy(t *testing.T) {
	t.Parallel()

	s := NewStyle().
		Bold(true).
		Italic(true).
		Underline(true).
		Strikethrough(true).
		Blink(true).
		Faint(true).
		Foreground(Color("#ffffff")).
		Background(Color("#111111")).
		Margin(1, 1, 1, 1).
		Padding(1, 1, 1, 1)

	i := s.Copy()

	requireEqual(t, s.GetBold(), i.GetBold())
	requireEqual(t, s.GetItalic(), i.GetItalic())
	requireEqual(t, s.GetUnderline(), i.GetUnderline())
	requireEqual(t, s.GetStrikethrough(), i.GetStrikethrough())
	requireEqual(t, s.GetBlink(), i.GetBlink())
	requireEqual(t, s.GetFaint(), i.GetFaint())
	requireEqual(t, s.GetForeground(), i.GetForeground())
	requireEqual(t, s.GetBackground(), i.GetBackground())

	requireEqual(t, s.GetMarginLeft(), i.GetMarginLeft())
	requireEqual(t, s.GetMarginRight(), i.GetMarginRight())
	requireEqual(t, s.GetMarginTop(), i.GetMarginTop())
	requireEqual(t, s.GetMarginBottom(), i.GetMarginBottom())
	requireEqual(t, s.GetPaddingLeft(), i.GetPaddingLeft())
	requireEqual(t, s.GetPaddingRight(), i.GetPaddingRight())
	requireEqual(t, s.GetPaddingTop(), i.GetPaddingTop())
	requireEqual(t, s.GetPaddingBottom(), i.GetPaddingBottom())
}

func TestStyleUnset(t *testing.T) {
	t.Parallel()

	s := NewStyle().Bold(true)
	requireTrue(t, s.GetBold())
	s.UnsetBold()
	requireFalse(t, s.GetBold())

	s = NewStyle().Italic(true)
	requireTrue(t, s.GetItalic())
	s.UnsetItalic()
	requireFalse(t, s.GetItalic())

	s = NewStyle().Underline(true)
	requireTrue(t, s.GetUnderline())
	s.UnsetUnderline()
	requireFalse(t, s.GetUnderline())

	s = NewStyle().Strikethrough(true)
	requireTrue(t, s.GetStrikethrough())
	s.UnsetStrikethrough()
	requireFalse(t, s.GetStrikethrough())

	s = NewStyle().Reverse(true)
	requireTrue(t, s.GetReverse())
	s.UnsetReverse()
	requireFalse(t, s.GetReverse())

	s = NewStyle().Blink(true)
	requireTrue(t, s.GetBlink())
	s.UnsetBlink()
	requireFalse(t, s.GetBlink())

	s = NewStyle().Faint(true)
	requireTrue(t, s.GetFaint())
	s.UnsetFaint()
	requireFalse(t, s.GetFaint())

	s = NewStyle().Inline(true)
	requireTrue(t, s.GetInline())
	s.UnsetInline()
	requireFalse(t, s.GetInline())

	// colors
	col := Color("#ffffff")
	s = NewStyle().Foreground(col)
	requireEqual(t, col, s.GetForeground())
	s.UnsetForeground()
	requireNotEqual(t, col, s.GetForeground())

	s = NewStyle().Background(col)
	requireEqual(t, col, s.GetBackground())
	s.UnsetBackground()
	requireNotEqual(t, col, s.GetBackground())

	// margins
	s = NewStyle().Margin(1, 2, 3, 4)
	requireEqual(t, 1, s.GetMarginTop())
	s.UnsetMarginTop()
	requireEqual(t, 0, s.GetMarginTop())

	requireEqual(t, 2, s.GetMarginRight())
	s.UnsetMarginRight()
	requireEqual(t, 0, s.GetMarginRight())

	requireEqual(t, 3, s.GetMarginBottom())
	s.UnsetMarginBottom()
	requireEqual(t, 0, s.GetMarginBottom())

	requireEqual(t, 4, s.GetMarginLeft())
	s.UnsetMarginLeft()
	requireEqual(t, 0, s.GetMarginLeft())

	// padding
	s = NewStyle().Padding(1, 2, 3, 4)
	requireEqual(t, 1, s.GetPaddingTop())
	s.UnsetPaddingTop()
	requireEqual(t, 0, s.GetPaddingTop())

	requireEqual(t, 2, s.GetPaddingRight())
	s.UnsetPaddingRight()
	requireEqual(t, 0, s.GetPaddingRight())

	requireEqual(t, 3, s.GetPaddingBottom())
	s.UnsetPaddingBottom()
	requireEqual(t, 0, s.GetPaddingBottom())

	requireEqual(t, 4, s.GetPaddingLeft())
	s.UnsetPaddingLeft()
	requireEqual(t, 0, s.GetPaddingLeft())

	// border
	s = NewStyle().Border(normalBorder, true, true, true, true)
	requireTrue(t, s.GetBorderTop())
	s.UnsetBorderTop()
	requireFalse(t, s.GetBorderTop())

	requireTrue(t, s.GetBorderRight())
	s.UnsetBorderRight()
	requireFalse(t, s.GetBorderRight())

	requireTrue(t, s.GetBorderBottom())
	s.UnsetBorderBottom()
	requireFalse(t, s.GetBorderBottom())

	requireTrue(t, s.GetBorderLeft())
	s.UnsetBorderLeft()
	requireFalse(t, s.GetBorderLeft())
}

func TestStyleValue(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		style    Style
		expected string
	}{
		{
			name:     "empty",
			style:    NewStyle(),
			expected: "foo",
		},
		{
			name:     "set string",
			style:    NewStyle().SetString("bar"),
			expected: "bar foo",
		},
		{
			name:     "set string with bold",
			style:    NewStyle().SetString("bar").Bold(true),
			expected: "\x1b[1mbar foo\x1b[0m",
		},
		{
			name:     "new style with string",
			style:    NewStyle("bar", "foobar"),
			expected: "bar foobar foo",
		},
	}

	for i, tc := range tt {
		res := tc.style.Render("foo")
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}

}

func BenchmarkStyleRender(b *testing.B) {
	s := NewStyle().
		Bold(true).
		Foreground(Color("#ffffff"))

	for i := 0; i < b.N; i++ {
		s.Render("Hello world")
	}
}

func requireTrue(tb testing.TB, b bool) {
	requireEqual(tb, true, b)
}

func requireFalse(tb testing.TB, b bool) {
	requireEqual(tb, false, b)
}

func requireEqual(tb testing.TB, a, b interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(a, b) {
		tb.Errorf("%v != %v", a, b)
		tb.FailNow()
	}
}

func requireNotEqual(tb testing.TB, a, b interface{}) {
	tb.Helper()
	if reflect.DeepEqual(a, b) {
		tb.Errorf("%v == %v", a, b)
		tb.FailNow()
	}
}
