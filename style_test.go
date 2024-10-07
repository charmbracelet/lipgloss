package lipgloss

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

func TestUnderline(t *testing.T) {
	r := NewRenderer(io.Discard)
	r.SetColorProfile(termenv.TrueColor)
	r.SetHasDarkBackground(true)
	t.Parallel()

	tt := []struct {
		style    Style
		expected string
	}{
		{
			r.NewStyle().Underline(true),
			"\x1b[4;4ma\x1b[0m\x1b[4;4mb\x1b[0m\x1b[4m \x1b[0m\x1b[4;4mc\x1b[0m",
		},
		{
			r.NewStyle().Underline(true).UnderlineSpaces(true),
			"\x1b[4;4ma\x1b[0m\x1b[4;4mb\x1b[0m\x1b[4m \x1b[0m\x1b[4;4mc\x1b[0m",
		},
		{
			r.NewStyle().Underline(true).UnderlineSpaces(false),
			"\x1b[4;4ma\x1b[0m\x1b[4;4mb\x1b[0m \x1b[4;4mc\x1b[0m",
		},
		{
			r.NewStyle().UnderlineSpaces(true),
			"ab\x1b[4m \x1b[0mc",
		},
	}

	for i, tc := range tt {
		s := tc.style.SetString("ab c")
		res := s.Render()
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}
}

func TestStrikethrough(t *testing.T) {
	r := NewRenderer(io.Discard)
	r.SetColorProfile(termenv.TrueColor)
	r.SetHasDarkBackground(true)
	t.Parallel()

	tt := []struct {
		style    Style
		expected string
	}{
		{
			r.NewStyle().Strikethrough(true),
			"\x1b[9ma\x1b[0m\x1b[9mb\x1b[0m\x1b[9m \x1b[0m\x1b[9mc\x1b[0m",
		},
		{
			r.NewStyle().Strikethrough(true).StrikethroughSpaces(true),
			"\x1b[9ma\x1b[0m\x1b[9mb\x1b[0m\x1b[9m \x1b[0m\x1b[9mc\x1b[0m",
		},
		{
			r.NewStyle().Strikethrough(true).StrikethroughSpaces(false),
			"\x1b[9ma\x1b[0m\x1b[9mb\x1b[0m \x1b[9mc\x1b[0m",
		},
		{
			r.NewStyle().StrikethroughSpaces(true),
			"ab\x1b[9m \x1b[0mc",
		},
	}

	for i, tc := range tt {
		s := tc.style.SetString("ab c")
		res := s.Render()
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}
}

func TestStyleRender(t *testing.T) {
	r := NewRenderer(io.Discard)
	r.SetColorProfile(termenv.TrueColor)
	r.SetHasDarkBackground(true)
	t.Parallel()

	tt := []struct {
		style    Style
		expected string
	}{
		{
			r.NewStyle().Foreground(Color("#5A56E0")),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
		{
			r.NewStyle().Foreground(AdaptiveColor{Light: "#fffe12", Dark: "#5A56E0"}),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
		{
			r.NewStyle().Bold(true),
			"\x1b[1mhello\x1b[0m",
		},
		{
			r.NewStyle().Italic(true),
			"\x1b[3mhello\x1b[0m",
		},
		{
			r.NewStyle().Underline(true),
			"\x1b[4;4mh\x1b[0m\x1b[4;4me\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4mo\x1b[0m",
		},
		{
			r.NewStyle().Blink(true),
			"\x1b[5mhello\x1b[0m",
		},
		{
			r.NewStyle().Faint(true),
			"\x1b[2mhello\x1b[0m",
		},
	}

	for i, tc := range tt {
		s := tc.style.SetString("hello")
		res := s.Render()
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}
}

func TestStyleCustomRender(t *testing.T) {
	r := NewRenderer(io.Discard)
	r.SetHasDarkBackground(false)
	r.SetColorProfile(termenv.TrueColor)
	tt := []struct {
		style    Style
		expected string
	}{
		{
			r.NewStyle().Foreground(Color("#5A56E0")),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
		{
			r.NewStyle().Foreground(AdaptiveColor{Light: "#fffe12", Dark: "#5A56E0"}),
			"\x1b[38;2;255;254;18mhello\x1b[0m",
		},
		{
			r.NewStyle().Bold(true),
			"\x1b[1mhello\x1b[0m",
		},
		{
			r.NewStyle().Italic(true),
			"\x1b[3mhello\x1b[0m",
		},
		{
			r.NewStyle().Underline(true),
			"\x1b[4;4mh\x1b[0m\x1b[4;4me\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4ml\x1b[0m\x1b[4;4mo\x1b[0m",
		},
		{
			r.NewStyle().Blink(true),
			"\x1b[5mhello\x1b[0m",
		},
		{
			r.NewStyle().Faint(true),
			"\x1b[2mhello\x1b[0m",
		},
		{
			NewStyle().Faint(true).Renderer(r),
			"\x1b[2mhello\x1b[0m",
		},
	}

	for i, tc := range tt {
		s := tc.style.SetString("hello")
		res := s.Render()
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}
}

func TestStyleRenderer(t *testing.T) {
	r := NewRenderer(io.Discard)
	s1 := NewStyle().Bold(true)
	s2 := s1.Renderer(r)
	if s1.r == s2.r {
		t.Fatalf("expected different renderers")
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
	requireEqual(t, s.GetUnderlineSpaces(), i.GetUnderlineSpaces())
	requireEqual(t, s.GetStrikethrough(), i.GetStrikethrough())
	requireEqual(t, s.GetStrikethroughSpaces(), i.GetStrikethroughSpaces())
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
		Padding(1, 1, 1, 1).
		TabWidth(2)

	i := s // copy

	requireEqual(t, s.GetBold(), i.GetBold())
	requireEqual(t, s.GetItalic(), i.GetItalic())
	requireEqual(t, s.GetUnderline(), i.GetUnderline())
	requireEqual(t, s.GetUnderlineSpaces(), i.GetUnderlineSpaces())
	requireEqual(t, s.GetStrikethrough(), i.GetStrikethrough())
	requireEqual(t, s.GetStrikethroughSpaces(), i.GetStrikethroughSpaces())
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
	requireEqual(t, s.GetTabWidth(), i.GetTabWidth())
}

func TestStyleUnset(t *testing.T) {
	t.Parallel()

	s := NewStyle().Bold(true)
	requireTrue(t, s.GetBold())
	s = s.UnsetBold()
	requireFalse(t, s.GetBold())

	s = NewStyle().Italic(true)
	requireTrue(t, s.GetItalic())
	s = s.UnsetItalic()
	requireFalse(t, s.GetItalic())

	s = NewStyle().Underline(true)
	requireTrue(t, s.GetUnderline())
	s = s.UnsetUnderline()
	requireFalse(t, s.GetUnderline())

	s = NewStyle().UnderlineSpaces(true)
	requireTrue(t, s.GetUnderlineSpaces())
	s = s.UnsetUnderlineSpaces()
	requireFalse(t, s.GetUnderlineSpaces())

	s = NewStyle().Strikethrough(true)
	requireTrue(t, s.GetStrikethrough())
	s = s.UnsetStrikethrough()
	requireFalse(t, s.GetStrikethrough())

	s = NewStyle().StrikethroughSpaces(true)
	requireTrue(t, s.GetStrikethroughSpaces())
	s = s.UnsetStrikethroughSpaces()
	requireFalse(t, s.GetStrikethroughSpaces())

	s = NewStyle().Reverse(true)
	requireTrue(t, s.GetReverse())
	s = s.UnsetReverse()
	requireFalse(t, s.GetReverse())

	s = NewStyle().Blink(true)
	requireTrue(t, s.GetBlink())
	s = s.UnsetBlink()
	requireFalse(t, s.GetBlink())

	s = NewStyle().Faint(true)
	requireTrue(t, s.GetFaint())
	s = s.UnsetFaint()
	requireFalse(t, s.GetFaint())

	s = NewStyle().Inline(true)
	requireTrue(t, s.GetInline())
	s = s.UnsetInline()
	requireFalse(t, s.GetInline())

	// colors
	col := Color("#ffffff")
	s = NewStyle().Foreground(col)
	requireEqual(t, col, s.GetForeground())
	s = s.UnsetForeground()
	requireNotEqual(t, col, s.GetForeground())

	s = NewStyle().Background(col)
	requireEqual(t, col, s.GetBackground())
	s = s.UnsetBackground()
	requireNotEqual(t, col, s.GetBackground())

	// margins
	s = NewStyle().Margin(1, 2, 3, 4)
	requireEqual(t, 1, s.GetMarginTop())
	s = s.UnsetMarginTop()
	requireEqual(t, 0, s.GetMarginTop())

	requireEqual(t, 2, s.GetMarginRight())
	s = s.UnsetMarginRight()
	requireEqual(t, 0, s.GetMarginRight())

	requireEqual(t, 3, s.GetMarginBottom())
	s = s.UnsetMarginBottom()
	requireEqual(t, 0, s.GetMarginBottom())

	requireEqual(t, 4, s.GetMarginLeft())
	s = s.UnsetMarginLeft()
	requireEqual(t, 0, s.GetMarginLeft())

	// padding
	s = NewStyle().Padding(1, 2, 3, 4)
	requireEqual(t, 1, s.GetPaddingTop())
	s = s.UnsetPaddingTop()
	requireEqual(t, 0, s.GetPaddingTop())

	requireEqual(t, 2, s.GetPaddingRight())
	s = s.UnsetPaddingRight()
	requireEqual(t, 0, s.GetPaddingRight())

	requireEqual(t, 3, s.GetPaddingBottom())
	s = s.UnsetPaddingBottom()
	requireEqual(t, 0, s.GetPaddingBottom())

	requireEqual(t, 4, s.GetPaddingLeft())
	s = s.UnsetPaddingLeft()
	requireEqual(t, 0, s.GetPaddingLeft())

	// border
	s = NewStyle().Border(normalBorder, true, true, true, true)
	requireTrue(t, s.GetBorderTop())
	s = s.UnsetBorderTop()
	requireFalse(t, s.GetBorderTop())

	requireTrue(t, s.GetBorderRight())
	s = s.UnsetBorderRight()
	requireFalse(t, s.GetBorderRight())

	requireTrue(t, s.GetBorderBottom())
	s = s.UnsetBorderBottom()
	requireFalse(t, s.GetBorderBottom())

	requireTrue(t, s.GetBorderLeft())
	s = s.UnsetBorderLeft()
	requireFalse(t, s.GetBorderLeft())

	// tab width
	s = NewStyle().TabWidth(2)
	requireEqual(t, s.GetTabWidth(), 2)
	s = s.UnsetTabWidth()
	requireNotEqual(t, s.GetTabWidth(), 4)
}

func TestStyleValue(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		text     string
		style    Style
		expected string
	}{
		{
			name:     "empty",
			text:     "foo",
			style:    NewStyle(),
			expected: "foo",
		},
		{
			name:     "set string",
			text:     "foo",
			style:    NewStyle().SetString("bar"),
			expected: "bar foo",
		},
		{
			name:     "set string with bold",
			text:     "foo",
			style:    NewStyle().SetString("bar").Bold(true),
			expected: "\x1b[1mbar foo\x1b[0m",
		},
		{
			name:     "new style with string",
			text:     "foo",
			style:    NewStyle().SetString("bar", "foobar"),
			expected: "bar foobar foo",
		},
		{
			name:     "margin right",
			text:     "foo",
			style:    NewStyle().MarginRight(1),
			expected: "foo ",
		},
		{
			name:     "margin left",
			text:     "foo",
			style:    NewStyle().MarginLeft(1),
			expected: " foo",
		},
		{
			name:     "empty text margin right",
			text:     "",
			style:    NewStyle().MarginRight(1),
			expected: " ",
		},
		{
			name:     "empty text margin left",
			text:     "",
			style:    NewStyle().MarginLeft(1),
			expected: " ",
		},
	}

	for i, tc := range tt {
		res := tc.style.Render(tc.text)
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}
}

func TestTabConversion(t *testing.T) {
	s := NewStyle()
	requireEqual(t, "[    ]", s.Render("[\t]"))
	s = NewStyle().TabWidth(2)
	requireEqual(t, "[  ]", s.Render("[\t]"))
	s = NewStyle().TabWidth(0)
	requireEqual(t, "[]", s.Render("[\t]"))
	s = NewStyle().TabWidth(-1)
	requireEqual(t, "[\t]", s.Render("[\t]"))
}

func TestStringTransform(t *testing.T) {
	for i, tc := range []struct {
		input    string
		fn       func(string) string
		expected string
	}{
		// No-op.
		{
			"hello",
			func(s string) string { return s },
			"hello",
		},
		// Uppercase.
		{
			"raow",
			strings.ToUpper,
			"RAOW",
		},
		// English and Chinese.
		{
			"The quick brown 狐 jumped over the lazy 犬",
			func(s string) string {
				n := 0
				rune := make([]rune, len(s))
				for _, r := range s {
					rune[n] = r
					n++
				}
				rune = rune[0:n]
				for i := 0; i < n/2; i++ {
					rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
				}
				return string(rune)
			},
			"犬 yzal eht revo depmuj 狐 nworb kciuq ehT",
		},
	} {
		res := NewStyle().Bold(true).Transform(tc.fn).Render(tc.input)
		expected := "\x1b[1m" + tc.expected + "\x1b[0m"
		if res != expected {
			t.Errorf("Test #%d:\nExpected: %q\nGot:      %q", i+1, expected, res)
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
	tb.Helper()
	requireEqual(tb, true, b)
}

func requireFalse(tb testing.TB, b bool) {
	tb.Helper()
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

func TestCarriageReturnInRender(t *testing.T) {
	out := fmt.Sprintf("%s\r\n%s\r\n", "Super duper california oranges", "Hello world")
	testStyle := NewStyle().
		MarginLeft(1)
	got := testStyle.Render(string(out))
	want := testStyle.Render(fmt.Sprintf("%s\n%s\n", "Super duper california oranges", "Hello world"))

	if got != want {
		t.Logf("got(detailed):\n%q\nwant(detailed):\n%q", got, want)
		t.Fatalf("got(string):\n%s\nwant(string):\n%s", got, want)
	}
}
