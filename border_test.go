package lipgloss

import (
	"strings"
	"testing"
)

var secretBorderFantasy = Border{
	Top:         "._.:*:",
	Bottom:      "._.:*:",
	Left:        "|*",
	Right:       "|*",
	TopLeft:     "*",
	TopRight:    "*",
	BottomLeft:  "*",
	BottomRight: "*",
}

// NB: The superfluous TrimSpace calls and newlines here are for readability.

func TestBorder(t *testing.T) {
	for i, tc := range []struct {
		name     string
		style    Style
		expected string
	}{
		{
			name: "default border via border style",
			style: NewStyle().
				BorderStyle(NormalBorder()).
				SetString("Hello"),
			expected: strings.TrimSpace(`
┌─────┐
│Hello│
└─────┘`),
		},

		{
			name: "default border via all-in-one, implicit",
			style: NewStyle().
				Border(RoundedBorder()).
				SetString("Hello"),
			expected: strings.TrimSpace(`
╭─────╮
│Hello│
╰─────╯`),
		},

		{
			name: "default border via all-in-one, explicit",
			style: NewStyle().
				Border(RoundedBorder(), true, true, true, true).
				SetString("Hello"),
			expected: strings.TrimSpace(`
╭─────╮
│Hello│
╰─────╯`),
		},

		{
			name: "rounded border via all-in-one, no bottom",
			style: NewStyle().
				Border(RoundedBorder(), true, true, false, true).
				SetString("Hello"),
			expected: strings.TrimSpace(`
╭─────╮
│Hello│`),
		},

		{
			name: "rounded border method-by-method, no right",
			style: NewStyle().
				BorderStyle(RoundedBorder()).
				BorderTop(true).
				BorderBottom(true).
				BorderLeft(true).
				SetString("Hello"),
			expected: strings.TrimSpace(`
╭─────
│Hello
╰─────`),
		},

		{
			name: "rounded border method-by-method, no left",
			style: NewStyle().
				BorderStyle(RoundedBorder()).
				BorderTop(true).
				BorderBottom(true).
				BorderLeft(false).
				BorderRight(true).
				SetString("Hello"),
			expected: strings.TrimSpace(`
─────╮
Hello│
─────╯`),
		},

		{
			name: "border via methods, no actual border set",
			style: NewStyle().
				BorderTop(true).
				BorderRight(true).
				BorderBottom(true).
				BorderLeft(true).
				SetString("Hello"),
			expected: "Hello",
		},

		{
			name: "custom border",
			style: NewStyle().
				BorderStyle(Border{
					Left:        "|",
					Right:       "|",
					Top:         ">",
					Bottom:      "<",
					TopLeft:     "+",
					TopRight:    ">",
					BottomLeft:  "<",
					BottomRight: "+",
				}).
				Padding(0, 1).
				SetString("Hello"),
			expected: strings.TrimSpace(`
+>>>>>>>>
| Hello |
<<<<<<<<+`),
		},

		{
			name: "corners only",
			style: NewStyle().
				BorderStyle(Border{
					TopLeft:     "+",
					TopRight:    "+",
					BottomLeft:  "+",
					BottomRight: "+",
				}).
				SetString("Hello"),
			expected: strings.TrimSpace("\n" +
				`+     +` + "\n" +
				` Hello ` + "\n" +
				`+     +`),
		},

		{
			name:     "set top via method",
			style:    NewStyle().BorderTop(true).SetString("Hello"),
			expected: "Hello",
		},
		{
			name:     "set bottom via method",
			style:    NewStyle().BorderTop(true).SetString("Hello"),
			expected: "Hello",
		},

		{
			name:     "inline, set right via method",
			style:    NewStyle().BorderRight(true).SetString("Hello"),
			expected: `Hello`,
		},

		{
			name: "set right via border style",
			style: NewStyle().
				BorderStyle(Border{
					Right: "|",
				}).
				SetString("Hello"),
			expected: "Hello|",
		},

		{
			name: "left via border style only",
			style: NewStyle().
				BorderStyle(Border{
					Left: "|",
				}).
				SetString("Hello"),
			expected: "|Hello",
		},

		{
			name: "left and right via border style only",
			style: NewStyle().
				BorderStyle(Border{
					Left:  "(",
					Right: ")",
				}).
				SetString("你好"),
			expected: `(你好)`,
		},

		{
			name: "inline, left and right via border style only",
			style: NewStyle().
				BorderStyle(Border{
					Left:  "「",
					Right: "」",
				}).
				Inline(true).
				SetString("你好"),
			expected: `「你好」`,
		},

		{
			name: "left and right, two cells high",
			style: NewStyle().
				BorderStyle(Border{
					Left:  "(",
					Right: ")",
				}).
				Padding(0, 1).
				SetString("你\n好"),
			expected: strings.TrimSpace(`
( 你 )
( 好 )`),
		},

		{
			name: "left and right with vertical padding",
			style: NewStyle().
				BorderStyle(Border{
					Left:  "(",
					Right: ")",
				}).
				Padding(1).
				SetString("你\n好"),
			expected: strings.TrimSpace(`
(    )
( 你 )
( 好 )
(    )`),
		},

		{
			name: "right only by deduction, two cells high",
			style: NewStyle().
				BorderStyle(Border{
					Left:  "(",
					Right: ")",
				}).
				BorderLeft(false).
				BorderRight(true).
				Padding(0, 1).
				SetString("你\n好"),
			expected: ` 你 )
 好 )`,
		},

		{
			name: "all but left via shorthand, two cells high",
			style: NewStyle().
				Border(DoubleBorder(), true, true, true, false).
				SetString("你\n好"),
			expected: strings.TrimSpace(`
══╗
你║
好║
══╝`),
		},

		{
			name: "outrageous border",
			style: NewStyle().
				BorderStyle(secretBorderFantasy).
				Padding(1, 2).
				SetString("Kitty\nCat").
				Align(Center),
			expected: strings.TrimSpace(`
*._.:*:._.*
|         |
*  Kitty  *
|   Cat   |
*         *
*._.:*:._.*`),
		},
	} {
		res := tc.style.String()
		if res != tc.expected {
			t.Errorf(
				"Test #%d (%s):\nExpected:\n%s\nGot:     \n%s",
				i+1,
				tc.name,
				showHiddenChars(tc.expected),
				showHiddenChars(res),
			)
		}
	}
}

func showHiddenChars(s string) string {
	s = strings.ReplaceAll(s, " ", "•")
	s = strings.ReplaceAll(s, "\t", "→")
	return strings.ReplaceAll(s, "\n", "¶\n")
}
