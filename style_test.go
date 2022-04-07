package lipgloss

import (
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/require"
)

func TestStyleRender(t *testing.T) {
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
		SetColorProfile(termenv.TrueColor)
		res := tc.style.Render("hello")
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

	require.Equal(t, s.GetBold(), i.GetBold())
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

	require.Equal(t, s.GetBold(), i.GetBold())
	require.Equal(t, s.GetItalic(), i.GetItalic())
	require.Equal(t, s.GetUnderline(), i.GetUnderline())
	require.Equal(t, s.GetStrikethrough(), i.GetStrikethrough())
	require.Equal(t, s.GetBlink(), i.GetBlink())
	require.Equal(t, s.GetFaint(), i.GetFaint())
	require.Equal(t, s.GetForeground(), i.GetForeground())
	require.Equal(t, s.GetBackground(), i.GetBackground())

	require.NotEqual(t, s.GetMarginLeft(), i.GetMarginLeft())
	require.NotEqual(t, s.GetMarginRight(), i.GetMarginRight())
	require.NotEqual(t, s.GetMarginTop(), i.GetMarginTop())
	require.NotEqual(t, s.GetMarginBottom(), i.GetMarginBottom())
	require.NotEqual(t, s.GetPaddingLeft(), i.GetPaddingLeft())
	require.NotEqual(t, s.GetPaddingRight(), i.GetPaddingRight())
	require.NotEqual(t, s.GetPaddingTop(), i.GetPaddingTop())
	require.NotEqual(t, s.GetPaddingBottom(), i.GetPaddingBottom())
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

	require.Equal(t, s.GetBold(), i.GetBold())
	require.Equal(t, s.GetItalic(), i.GetItalic())
	require.Equal(t, s.GetUnderline(), i.GetUnderline())
	require.Equal(t, s.GetStrikethrough(), i.GetStrikethrough())
	require.Equal(t, s.GetBlink(), i.GetBlink())
	require.Equal(t, s.GetFaint(), i.GetFaint())
	require.Equal(t, s.GetForeground(), i.GetForeground())
	require.Equal(t, s.GetBackground(), i.GetBackground())

	require.Equal(t, s.GetMarginLeft(), i.GetMarginLeft())
	require.Equal(t, s.GetMarginRight(), i.GetMarginRight())
	require.Equal(t, s.GetMarginTop(), i.GetMarginTop())
	require.Equal(t, s.GetMarginBottom(), i.GetMarginBottom())
	require.Equal(t, s.GetPaddingLeft(), i.GetPaddingLeft())
	require.Equal(t, s.GetPaddingRight(), i.GetPaddingRight())
	require.Equal(t, s.GetPaddingTop(), i.GetPaddingTop())
	require.Equal(t, s.GetPaddingBottom(), i.GetPaddingBottom())
}

func TestStyleUnset(t *testing.T) {
	t.Parallel()

	s := NewStyle().Bold(true)
	require.True(t, s.GetBold())
	s.UnsetBold()
	require.False(t, s.GetBold())

	s = NewStyle().Italic(true)
	require.True(t, s.GetItalic())
	s.UnsetItalic()
	require.False(t, s.GetItalic())

	s = NewStyle().Underline(true)
	require.True(t, s.GetUnderline())
	s.UnsetUnderline()
	require.False(t, s.GetUnderline())

	s = NewStyle().Strikethrough(true)
	require.True(t, s.GetStrikethrough())
	s.UnsetStrikethrough()
	require.False(t, s.GetStrikethrough())

	s = NewStyle().Reverse(true)
	require.True(t, s.GetReverse())
	s.UnsetReverse()
	require.False(t, s.GetReverse())

	s = NewStyle().Blink(true)
	require.True(t, s.GetBlink())
	s.UnsetBlink()
	require.False(t, s.GetBlink())

	s = NewStyle().Faint(true)
	require.True(t, s.GetFaint())
	s.UnsetFaint()
	require.False(t, s.GetFaint())

	s = NewStyle().Inline(true)
	require.True(t, s.GetInline())
	s.UnsetInline()
	require.False(t, s.GetInline())

	// colors
	col := Color("#ffffff")
	s = NewStyle().Foreground(col)
	require.Equal(t, col, s.GetForeground())
	s.UnsetForeground()
	require.NotEqual(t, col, s.GetForeground())

	s = NewStyle().Background(col)
	require.Equal(t, col, s.GetBackground())
	s.UnsetBackground()
	require.NotEqual(t, col, s.GetBackground())

	// margins
	s = NewStyle().Margin(1, 2, 3, 4)
	require.Equal(t, 1, s.GetMarginTop())
	s.UnsetMarginTop()
	require.Equal(t, 0, s.GetMarginTop())

	require.Equal(t, 2, s.GetMarginRight())
	s.UnsetMarginRight()
	require.Equal(t, 0, s.GetMarginRight())

	require.Equal(t, 3, s.GetMarginBottom())
	s.UnsetMarginBottom()
	require.Equal(t, 0, s.GetMarginBottom())

	require.Equal(t, 4, s.GetMarginLeft())
	s.UnsetMarginLeft()
	require.Equal(t, 0, s.GetMarginLeft())

	// padding
	s = NewStyle().Padding(1, 2, 3, 4)
	require.Equal(t, 1, s.GetPaddingTop())
	s.UnsetPaddingTop()
	require.Equal(t, 0, s.GetPaddingTop())

	require.Equal(t, 2, s.GetPaddingRight())
	s.UnsetPaddingRight()
	require.Equal(t, 0, s.GetPaddingRight())

	require.Equal(t, 3, s.GetPaddingBottom())
	s.UnsetPaddingBottom()
	require.Equal(t, 0, s.GetPaddingBottom())

	require.Equal(t, 4, s.GetPaddingLeft())
	s.UnsetPaddingLeft()
	require.Equal(t, 0, s.GetPaddingLeft())

	// border
	s = NewStyle().Border(normalBorder, true, true, true, true)
	require.True(t, s.GetBorderTop())
	s.UnsetBorderTop()
	require.False(t, s.GetBorderTop())

	require.True(t, s.GetBorderRight())
	s.UnsetBorderRight()
	require.False(t, s.GetBorderRight())

	require.True(t, s.GetBorderBottom())
	s.UnsetBorderBottom()
	require.False(t, s.GetBorderBottom())

	require.True(t, s.GetBorderLeft())
	s.UnsetBorderLeft()
	require.False(t, s.GetBorderLeft())
}

func BenchmarkStyleRender(b *testing.B) {
	s := NewStyle().
		Bold(true).
		Foreground(Color("#ffffff"))

	for i := 0; i < b.N; i++ {
		s.Render("Hello world")
	}
}
