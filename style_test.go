package lipgloss

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestUnderline(t *testing.T) {
	t.Parallel()

	tt := []struct {
		style    Style
		expected string
	}{
		{
			NewStyle().Underline(true),
			"\x1b[4;4ma\x1b[m\x1b[4;4mb\x1b[m\x1b[4m \x1b[m\x1b[4;4mc\x1b[m",
		},
		{
			NewStyle().Underline(true).UnderlineSpaces(true),
			"\x1b[4;4ma\x1b[m\x1b[4;4mb\x1b[m\x1b[4m \x1b[m\x1b[4;4mc\x1b[m",
		},
		{
			NewStyle().Underline(true).UnderlineSpaces(false),
			"\x1b[4;4ma\x1b[m\x1b[4;4mb\x1b[m \x1b[4;4mc\x1b[m",
		},
		{
			NewStyle().UnderlineSpaces(true),
			"ab\x1b[4m \x1b[mc",
		},
		{
			NewStyle().UnderlineStyle(UnderlineCurly),
			"\x1b[4;4:3ma\x1b[m\x1b[4;4:3mb\x1b[m\x1b[4m \x1b[m\x1b[4;4:3mc\x1b[m",
		},
		{
			NewStyle().UnderlineStyle(UnderlineCurly).UnderlineColor(Color("#FF0000")),
			"\x1b[4;58;2;255;0;0;4:3ma\x1b[m\x1b[4;58;2;255;0;0;4:3mb\x1b[m\x1b[58;2;255;0;0;4m \x1b[m\x1b[4;58;2;255;0;0;4:3mc\x1b[m",
		},
	}

	for i, tc := range tt {
		s := tc.style.SetString("ab c")
		res := s.Render()
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n`%q`\n\nActual output:\n`%q`\n\n",
				i, tc.expected,
				res)
		}
	}
}

func TestGetUnderlineColor(t *testing.T) {
	t.Parallel()

	red := Color("#FF0000")
	s := NewStyle().Underline(true).UnderlineColor(red)
	if s.GetUnderlineColor() != red {
		t.Errorf("GetUnderlineColor() = %v, want %v", s.GetUnderlineColor(), red)
	}
}

func TestStrikethrough(t *testing.T) {
	t.Parallel()

	tt := []struct {
		style    Style
		expected string
	}{
		{
			NewStyle().Strikethrough(true),
			"\x1b[9ma\x1b[m\x1b[9mb\x1b[m\x1b[9m \x1b[m\x1b[9mc\x1b[m",
		},
		{
			NewStyle().Strikethrough(true).StrikethroughSpaces(true),
			"\x1b[9ma\x1b[m\x1b[9mb\x1b[m\x1b[9m \x1b[m\x1b[9mc\x1b[m",
		},
		{
			NewStyle().Strikethrough(true).StrikethroughSpaces(false),
			"\x1b[9ma\x1b[m\x1b[9mb\x1b[m \x1b[9mc\x1b[m",
		},
		{
			NewStyle().StrikethroughSpaces(true),
			"ab\x1b[9m \x1b[mc",
		},
	}

	for i, tc := range tt {
		s := tc.style.SetString("ab c")
		res := s.Render()
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n`%q`\n\nActual output:\n`%q`\n\n",
				i, tc.expected,
				res)
		}
	}
}

func TestStyleRender(t *testing.T) {
	t.Parallel()

	tt := []struct {
		style    Style
		expected string
	}{
		{
			NewStyle().Foreground(Color("#5A56E0")),
			"\x1b[38;2;90;86;224mhello\x1b[m",
		},
		{
			NewStyle().Bold(true),
			"\x1b[1mhello\x1b[m",
		},
		{
			NewStyle().Italic(true),
			"\x1b[3mhello\x1b[m",
		},
		{
			NewStyle().Underline(true),
			"\x1b[4;4mh\x1b[m\x1b[4;4me\x1b[m\x1b[4;4ml\x1b[m\x1b[4;4ml\x1b[m\x1b[4;4mo\x1b[m",
		},
		{
			NewStyle().Blink(true),
			"\x1b[5mhello\x1b[m",
		},
		{
			NewStyle().Faint(true),
			"\x1b[2mhello\x1b[m",
		},
	}

	for i, tc := range tt {
		s := tc.style.SetString("hello")
		res := s.Render()
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n`%q`\n\nActual output:\n`%q`\n\n",
				i, tc.expected,
				res)
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
	s = NewStyle().Padding(1, 2, 3, 4).PaddingChar('x')
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

	requireEqual(t, 'x', s.GetPaddingChar())
	s = s.UnsetPaddingChar()
	requireEqual(t, ' ', s.GetPaddingChar())

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
			expected: "\x1b[1mbar foo\x1b[m",
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
			t.Errorf("Test %d, expected:\n`%q`\n\nActual output:\n`%q`\n\n",
				i, tc.expected,
				res)
		}
	}
}

func TestCustomPaddingChar(t *testing.T) {
	s := NewStyle().Padding(0, 3).PaddingChar('x')
	requireEqual(t, "xxxTESTxxx", s.Render("TEST"))
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
				for i := range n / 2 {
					rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
				}
				return string(rune)
			},
			"犬 yzal eht revo depmuj 狐 nworb kciuq ehT",
		},
	} {
		res := NewStyle().Bold(true).Transform(tc.fn).Render(tc.input)
		expected := "\x1b[1m" + tc.expected + "\x1b[m"
		if res != expected {
			t.Errorf("Test #%d:\nExpected: %q\nGot:      %q", i+1, expected, res)
		}
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

func requireEqual(tb testing.TB, a, b any) {
	tb.Helper()
	if !reflect.DeepEqual(a, b) {
		tb.Errorf("%v != %v", a, b)
		tb.FailNow()
	}
}

func requireNotEqual(tb testing.TB, a, b any) {
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

func TestWidth(t *testing.T) {
	tests := []struct {
		name  string
		style Style
	}{
		{"width with borders", NewStyle().Padding(0, 2).Border(NormalBorder(), true)},
		{"width no borders", NewStyle().Padding(0, 2)},
		{"width unset borders", NewStyle().Padding(0, 2).Border(NormalBorder(), true).BorderLeft(false).BorderRight(false)},
		{"width single-sided border", NewStyle().Padding(0, 2).Border(NormalBorder(), true).UnsetBorderBottom().UnsetBorderTop().UnsetBorderRight()},
	}
	{
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				content := "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
				contentWidth := 80 - tc.style.GetHorizontalFrameSize()
				rendered := tc.style.Width(contentWidth).Render(content)
				if Width(rendered) != contentWidth {
					t.Log("\n" + rendered)
					t.Fatalf("got: %d\n, want: %d", Width(rendered), contentWidth)
				}
			})
		}
	}
}

func TestHeight(t *testing.T) {
	tests := []struct {
		name  string
		style Style
	}{
		{"height with borders", NewStyle().Width(80).Padding(0, 2).Border(NormalBorder(), true)},
		{"height no borders", NewStyle().Width(80).Padding(0, 2)},
		{"height unset borders", NewStyle().Width(80).Padding(0, 2).Border(NormalBorder(), true).BorderBottom(false).BorderTop(false)},
		{"height single-sided border", NewStyle().Width(80).Padding(0, 2).Border(NormalBorder(), true).UnsetBorderLeft().UnsetBorderBottom().UnsetBorderRight()},
	}
	{
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				content := "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
				contentHeight := 20 - tc.style.GetVerticalFrameSize()
				rendered := tc.style.Height(contentHeight).Render(content)
				if Height(rendered) != contentHeight {
					t.Log("\n" + rendered)
					t.Fatalf("got: %d\n, want: %d", Height(rendered), contentHeight)
				}
			})
		}
	}
}

func TestHyperlink(t *testing.T) {
	tests := []struct {
		name     string
		style    Style
		expected string
	}{
		{
			name:     "hyperlink",
			style:    NewStyle().Hyperlink("https://example.com").SetString("https://example.com"),
			expected: "\x1b]8;;https://example.com\x07https://example.com\x1b]8;;\x07",
		},
		{
			name:     "hyperlink with text",
			style:    NewStyle().Hyperlink("https://example.com", "id=123").SetString("example"),
			expected: "\x1b]8;id=123;https://example.com\x07example\x1b]8;;\x07",
		},
		{
			name: "hyperlink with text and style",
			style: NewStyle().Hyperlink("https://example.com", "id=123").SetString("example").
				Bold(true).Foreground(Color("234")),
			expected: "\x1b]8;id=123;https://example.com\x07\x1b[1;38;5;234mexample\x1b[m\x1b]8;;\x07",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.style.String() != tc.expected {
				t.Fatalf("got: %q, want: %q", tc.style.String(), tc.expected)
			}
		})
	}
}

func TestUnsetHyperlink(t *testing.T) {
	tests := []struct {
		name     string
		style    Style
		expected string
	}{
		{
			name:     "unset hyperlink",
			style:    NewStyle().Hyperlink("https://example.com").SetString("https://example.com").UnsetHyperlink(),
			expected: "https://example.com",
		},
		{
			name:     "unset hyperlink with text",
			style:    NewStyle().Hyperlink("https://example.com", "id=123").SetString("example").UnsetHyperlink(),
			expected: "example",
		},
		{
			name: "unset hyperlink with text and style",
			style: NewStyle().Hyperlink("https://example.com", "id=123").SetString("example").
				Bold(true).Foreground(Color("234")).UnsetHyperlink(),
			expected: "\x1b[1;38;5;234mexample\x1b[m",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.style.String() != tc.expected {
				t.Fatalf("got: %q, want: %q", tc.style.String(), tc.expected)
			}
		})
	}
}

func BenchmarkPad(b *testing.B) {
	tests := []struct {
		name string
		str  string
		n    int
	}{
		{name: "pad-10", str: "foo bar", n: 10},
		{name: "pad-100", str: "foo bar", n: 100},
		{name: "pad-negative-10", str: "foo bar", n: -10},
		{name: "pad-negative-100", str: "foo bar", n: -100},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				pad(tt.str, tt.n, nil, ' ')
			}
		})
	}
}

func BenchmarkStyleRender(b *testing.B) {
	tests := []struct {
		name  string
		style Style
		input string
	}{
		{
			name: "simple-1-line",
			style: NewStyle().
				Bold(true).
				Foreground(Color("#ffffff")),
			input: "Hello world",
		},
		{
			name: "simple-5-lines",
			style: NewStyle().
				Bold(true).
				Foreground(Color("#ffffff")),
			input: strings.Repeat("Hello world\n", 5),
		},
		{
			name: "simple-5-lines-inline",
			style: NewStyle().
				Bold(true).
				Foreground(Color("#ffffff")).
				Inline(true),
			input: strings.Repeat("Hello world\n", 5),
		},
		{
			name: "simple-10-lines-5-height-40-width",
			style: NewStyle().
				Bold(true).
				Foreground(Color("#ffffff")).
				Height(5).
				Width(40),
			input: strings.Repeat("Hello world\n", 10),
		},
		{
			name: "simple-10-lines-width-maxwidth",
			style: NewStyle().
				Bold(true).
				Foreground(Color("#ffffff")).
				Width(40).
				MaxWidth(40),
			input: strings.Repeat("Hello world\n", 10),
		},
		{
			name: "simple-10-lines-width-maxwidth-borders",
			style: NewStyle().
				Bold(true).
				Foreground(Color("#ffffff")).
				Width(40).
				MaxWidth(40).
				Border(RoundedBorder(), true),
			input: strings.Repeat("Hello world\n", 10),
		},
		{
			name: "simple-10-lines-width-maxwidth-borders-padding-margins",
			style: NewStyle().
				Bold(true).
				Foreground(Color("#ffffff")).
				Width(40).
				MaxWidth(40).
				Border(RoundedBorder(), true).
				Padding(1, 1).
				Margin(1, 1),
			input: strings.Repeat("Hello world\n", 10),
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				tt.style.Render(tt.input)
			}
		})
	}
}

func TestGetHorizontalFrameSize(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		style Style
		want  int
	}{
		{
			name:  "default zero",
			style: NewStyle(),
			want:  0,
		},
		{
			name:  "padding only",
			style: NewStyle().Padding(2, 4),
			want:  8, // left(4) + right(4)
		},
		{
			name:  "padding asymmetric",
			style: NewStyle().Padding(0, 1, 0, 3),
			want:  4, // left(3) + right(1)
		},
		{
			name:  "margins only",
			style: NewStyle().Margin(1, 2),
			want:  4, // left(2) + right(2)
		},
		{
			name:  "normal border all sides",
			style: NewStyle().Border(NormalBorder(), true),
			want:  2, // left(1) + right(1)
		},
		{
			name:  "rounded border all sides",
			style: NewStyle().Border(RoundedBorder(), true),
			want:  2,
		},
		{
			name:  "double border all sides",
			style: NewStyle().Border(DoubleBorder(), true),
			want:  2,
		},
		{
			name:  "padding plus border",
			style: NewStyle().Padding(0, 3).Border(NormalBorder(), true),
			want:  8, // left(1) + right(1) + left-pad(3) + right-pad(3)
		},
		{
			name:  "padding plus border plus margins",
			style: NewStyle().Padding(0, 2).Margin(0, 1).Border(NormalBorder(), true),
			want:  8, // left-margin(1) + left-border(1) + left-pad(2) + right-pad(2) + right-border(1) + right-margin(1)
		},
		{
			name:  "border without sides defaults to all sides",
			style: NewStyle().Border(NormalBorder()),
			want:  2,
		},
		{
			name:  "single-sided border left only",
			style: NewStyle().Border(NormalBorder(), true).BorderLeft(true).UnsetBorderRight().UnsetBorderTop().UnsetBorderBottom(),
			want:  1, // left(1) only
		},
		{
			name:  "no border when all sides unset",
			style: NewStyle().Border(NormalBorder()).BorderTop(false).BorderRight(false).BorderBottom(false).BorderLeft(false),
			want:  0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.style.GetHorizontalFrameSize()
			if got != tc.want {
				t.Errorf("GetHorizontalFrameSize() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestGetVerticalFrameSize(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name  string
		style Style
		want  int
	}{
		{
			name:  "default zero",
			style: NewStyle(),
			want:  0,
		},
		{
			name:  "padding only",
			style: NewStyle().Padding(2, 0),
			want:  4, // top(2) + bottom(2)
		},
		{
			name:  "padding asymmetric",
			style: NewStyle().Padding(1, 0, 3, 0),
			want:  4, // top(1) + bottom(3)
		},
		{
			name:  "margins only",
			style: NewStyle().Margin(2, 0),
			want:  4, // top(2) + bottom(2)
		},
		{
			name:  "normal border all sides",
			style: NewStyle().Border(NormalBorder(), true),
			want:  2, // top(1) + bottom(1)
		},
		{
			name:  "thick border all sides",
			style: NewStyle().Border(ThickBorder(), true),
			want:  2,
		},
		{
			name:  "padding plus border",
			style: NewStyle().Padding(2, 0).Border(NormalBorder(), true),
			want:  6, // top(1) + top-pad(2) + bottom-pad(2) + bottom(1)
		},
		{
			name:  "padding plus border plus margins",
			style: NewStyle().Padding(1, 0).Margin(2, 0).Border(NormalBorder(), true),
			want:  8, // top-margin(2) + top-border(1) + top-pad(1) + bottom-pad(1) + bottom-border(1) + bottom-margin(2)
		},
		{
			name:  "single-sided border top only",
			style: NewStyle().Border(NormalBorder()).BorderTop(true).UnsetBorderRight().UnsetBorderBottom().UnsetBorderLeft(),
			want:  1,
		},
		{
			name:  "hidden border",
			style: NewStyle().Border(HiddenBorder(), true),
			want:  2, // hidden border still has size 1 per side
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.style.GetVerticalFrameSize()
			if got != tc.want {
				t.Errorf("GetVerticalFrameSize() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestGetFrameSize(t *testing.T) {
	t.Parallel()

	style := NewStyle().Padding(1, 2).Margin(0, 1).Border(NormalBorder(), true)
	x, y := style.GetFrameSize()

	wantX := 8 // left-margin(1) + left-border(1) + left-pad(2) + right-pad(2) + right-border(1) + right-margin(1)
	wantY := 4 // top-margin(0) + top-border(1) + top-pad(1) + bottom-pad(1) + bottom-border(1) + bottom-margin(0)

	if x != wantX {
		t.Errorf("GetFrameSize() x = %d, want %d", x, wantX)
	}
	if y != wantY {
		t.Errorf("GetFrameSize() y = %d, want %d", y, wantY)
	}
}
