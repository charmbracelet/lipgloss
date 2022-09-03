package lipgloss_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	lipglossc "github.com/knz/lipgloss-convert"
)

type S = lipgloss.Style

func TestStyle(t *testing.T) {
	td := []struct {
		changeStyle func(S) S
		repr        string
	}{
		{func(s S) S { return s.AlignHorizontal(0.5) }, `align-horizontal: 0.5;`},
		{func(s S) S { return s.Background(lipgloss.Color("#0f0")) }, `background: #0f0;`},
		{func(s S) S { return s.Blink(true) }, `blink: true;`},
		{func(s S) S { return s.Bold(true) }, `bold: true;`},
		{func(s S) S { return s.BorderBottom(true) }, `border-bottom: true;`},
		{func(s S) S { return s.BorderBottomBackground(lipgloss.Color("#0f0")) }, `border-bottom-background: #0f0;`},
		{func(s S) S { return s.BorderBottomForeground(lipgloss.Color("#0f0")) }, `border-bottom-foreground: #0f0;`},
		{func(s S) S { return s.BorderLeft(true) }, `border-left: true;`},
		{func(s S) S { return s.BorderLeftBackground(lipgloss.Color("#0f0")) }, `border-left-background: #0f0;`},
		{func(s S) S { return s.BorderLeftForeground(lipgloss.Color("#0f0")) }, `border-left-foreground: #0f0;`},
		{func(s S) S { return s.BorderRight(true) }, `border-right: true;`},
		{func(s S) S { return s.BorderRightBackground(lipgloss.Color("#0f0")) }, `border-right-background: #0f0;`},
		{func(s S) S { return s.BorderRightForeground(lipgloss.Color("#0f0")) }, `border-right-foreground: #0f0;`},
		{func(s S) S {
			return s.BorderStyle(lipgloss.Border{"a", "b", "c", "d", "e", "f", "g", "h"})
		}, `border-style: border("a","b","c","d","e","f","g","h");`},
		{func(s S) S { return s.BorderTop(true) }, `border-top: true;`},
		{func(s S) S { return s.BorderTopBackground(lipgloss.Color("#0f0")) }, `border-top-background: #0f0;`},
		{func(s S) S { return s.BorderTopForeground(lipgloss.Color("#0f0")) }, `border-top-foreground: #0f0;`},
		{func(s S) S {
			return s.Border(lipgloss.Border{"a", "b", "c", "d", "e", "f", "g", "h"}, true)
		}, `border-bottom: true; border-left: true; border-right: true; ` +
			`border-style: border("a","b","c","d","e","f","g","h"); ` +
			`border-top: true;`},
		{func(s S) S { return s.BorderBackground(lipgloss.Color("#0f0")) }, `border-bottom-background: #0f0; border-left-background: #0f0; border-right-background: #0f0; border-top-background: #0f0;`},
		{func(s S) S { return s.BorderForeground(lipgloss.Color("#0f0")) }, `border-bottom-foreground: #0f0; border-left-foreground: #0f0; border-right-foreground: #0f0; border-top-foreground: #0f0;`},

		{func(s S) S { return s.ColorWhitespace(true) }, `color-whitespace: true;`},
		{func(s S) S { return s.Faint(true) }, `faint: true;`},
		{func(s S) S { return s.Foreground(lipgloss.Color("#0f0")) }, `foreground: #0f0;`},
		{func(s S) S { return s.Height(3) }, `height: 3;`},
		{func(s S) S { return s.Inline(true) }, `inline: true;`},
		{func(s S) S { return s.Italic(true) }, `italic: true;`},
		{func(s S) S { return s.Margin(1, 2, 3, 4) }, `margin-bottom: 3; margin-left: 4; margin-right: 2; margin-top: 1;`},
		{func(s S) S { return s.MarginBottom(3) }, `margin-bottom: 3;`},
		{func(s S) S { return s.MarginLeft(3) }, `margin-left: 3;`},
		{func(s S) S { return s.MarginRight(3) }, `margin-right: 3;`},
		{func(s S) S { return s.MarginTop(3) }, `margin-top: 3;`},
		{func(s S) S { return s.MarginBackground(lipgloss.Color("#0f0")) }, `margin-background: #0f0;`},
		{func(s S) S { return s.MaxHeight(3) }, `max-height: 3;`},
		{func(s S) S { return s.MaxWidth(3) }, `max-width: 3;`},
		{func(s S) S { return s.Padding(1, 2, 3, 4) }, `padding-bottom: 3; padding-left: 4; padding-right: 2; padding-top: 1;`},
		{func(s S) S { return s.PaddingBottom(3) }, `padding-bottom: 3;`},
		{func(s S) S { return s.PaddingLeft(3) }, `padding-left: 3;`},
		{func(s S) S { return s.PaddingRight(3) }, `padding-right: 3;`},
		{func(s S) S { return s.PaddingTop(3) }, `padding-top: 3;`},
		{func(s S) S { return s.Reverse(true) }, `reverse: true;`},
		{func(s S) S { return s.Strikethrough(true) }, `strikethrough: true;`},
		{func(s S) S { return s.StrikethroughSpaces(true) }, `strikethrough-spaces: true;`},
		{func(s S) S { return s.Underline(true) }, `underline: true;`},
		{func(s S) S { return s.UnderlineSpaces(true) }, `underline-spaces: true;`},
		{func(s S) S { return s.Width(3) }, `width: 3;`},
	}

	for _, tc := range td {
		// Apply the style change and compare to the reference.
		s := lipgloss.NewStyle()
		s = tc.changeStyle(s)
		repr := lipglossc.Export(s)
		if repr != tc.repr {
			t.Errorf("expected %q, got %q", tc.repr, repr)
			continue
		}

		// Apply the unset function and assert the resulting style is
		// empty.
		r := strings.Split(tc.repr, ":")
		if len(r) > 2 {
			// Special case: border-background / border-foreground.
			if strings.HasPrefix(r[0], "border-") && strings.HasSuffix(r[0], "ground") {
				r[0] = "border-" + r[0][strings.LastIndexByte(r[0], '-')+1:]
			} else {
				// Special case: padding, margin, border etc.
				r[0] = r[0][:strings.IndexByte(r[0], '-')]
			}
		}
		unset := r[0] + ": unset"
		u, err := lipglossc.Import(s, unset)
		if err != nil {
			t.Error(err)
			continue
		}
		repr = lipglossc.Export(u)
		if repr != "" {
			t.Errorf("expected empty style, got %q", repr)
		}
	}
}
