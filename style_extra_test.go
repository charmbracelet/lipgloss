package lipgloss_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/cockroachdb/datadriven"
	lipglossc "github.com/knz/lipgloss-convert"
	"github.com/muesli/termenv"
)

func Example_string() {
	s := lipgloss.NewStyle().
		Width(10).
		Align(lipgloss.Right).
		SetString("hello")

	fmt.Println(s.Value())
	fmt.Println(s.String())
	fmt.Println(s.UnsetString().Value())

	// Output:
	// hello
	//      hello
	//
}

// Example_padding exercises the computed padding getters.
func Example_padding() {
	s := lipgloss.NewStyle().Padding(10001, 10010, 10100, 11000)

	fmt.Println(s.GetPadding())
	fmt.Println(s.GetHorizontalPadding())
	fmt.Println(s.GetVerticalPadding())

	// Output:
	// 10001 10010 10100 11000
	// 21010
	// 20101
}

// Example_margin exercises the computed margin getters.
func Example_margin() {
	s := lipgloss.NewStyle().Margin(10001, 10010, 10100, 11000)

	fmt.Println(s.GetMargin())
	fmt.Println(s.GetHorizontalMargins())
	fmt.Println(s.GetVerticalMargins())

	// Output:
	// 10001 10010 10100 11000
	// 21010
	// 20101
}

// Example_border exercises the computed border getters.
func Example_border() {
	fmt.Println(lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).GetBorderStyle())
	fmt.Println(lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).GetBorderStyle())
	fmt.Println(lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder()).GetBorderStyle())
	fmt.Println(lipgloss.NewStyle().BorderStyle(lipgloss.DoubleBorder()).GetBorderStyle())
	fmt.Println(lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder()).GetBorderStyle())

	for _, b := range []bool{false, true} {
		fmt.Println("border-enabed:", b)
		s := lipgloss.NewStyle().Border(
			lipgloss.Border{"x", "xx", "xxx", "xxxx", "xxxxx", "xxxxxx", "a", "a", "a", "a", "a", "a", "a"},
			b,
		)

		fmt.Println(s.GetBorder())

		// Note: the border size computations seem to be wrong.
		// The top/bottom borders should be size 1 (just 1 line)
		// and the left/right borders should _add_ the rune sizes,
		// not compute the maximum.
		// See: https://github.com/charmbracelet/lipgloss/issues/112
		fmt.Println(s.GetBorderTopSize(), s.GetBorderBottomSize())
		fmt.Println(s.GetBorderLeftSize(), s.GetBorderRightSize())

		fmt.Println(s.GetHorizontalBorderSize())
		fmt.Println(s.GetVerticalBorderSize())
	}

	// Output:
	// {─ ─ │ │ ┌ ┐ ┘ └}
	// {─ ─ │ │ ╭ ╮ ╯ ╰}
	// {━ ━ ┃ ┃ ┏ ┓ ┛ ┗}
	// {═ ═ ║ ║ ╔ ╗ ╝ ╚}
	// {               }
	// border-enabed: false
	// {x xx xxx xxxx a a a a} false false false false
	// 0 0
	// 0 0
	// 2
	// 2
	// border-enabed: true
	// {x xx xxx xxxx a a a a} true true true true
	// 1 1
	// 1 1
	// 2
	// 2
}

// Example_frame exercises the computed frame size getters.
func Example_frame() {
	s := lipgloss.NewStyle().
		Margin(10001, 10010, 10100, 11000).
		Padding(10001, 10010, 10100, 11000).
		Border(lipgloss.NormalBorder(), true)

	fmt.Println(s.GetFrameSize())
	fmt.Println(s.GetHorizontalFrameSize())
	fmt.Println(s.GetVerticalFrameSize())

	// Output:
	// 42022 40204
	// 42022
	// 40204
}

type S = lipgloss.Style

// TestStyleMethods validates most of the Get, Set and Unset methods.
func TestStyleMethod(t *testing.T) {
	g := lipgloss.Color("#0f0")
	r := lipgloss.Color("#f00")
	b := lipgloss.Color("#00f")
	w := lipgloss.Color("#fff")

	td := []struct {
		changeStyle func(S) S
		repr        string
	}{
		{func(s S) S { return s.AlignHorizontal(0.5) }, `align-horizontal: 0.5;`},
		{func(s S) S { return s.Background(g) }, `background: #0f0;`},
		{func(s S) S { return s.Blink(true) }, `blink: true;`},
		{func(s S) S { return s.Bold(true) }, `bold: true;`},
		{func(s S) S { return s.BorderBottom(true) }, `border-bottom: true;`},
		{func(s S) S { return s.BorderBottomBackground(g) }, `border-bottom-background: #0f0;`},
		{func(s S) S { return s.BorderBottomForeground(g) }, `border-bottom-foreground: #0f0;`},
		{func(s S) S { return s.BorderLeft(true) }, `border-left: true;`},
		{func(s S) S { return s.BorderLeftBackground(g) }, `border-left-background: #0f0;`},
		{func(s S) S { return s.BorderLeftForeground(g) }, `border-left-foreground: #0f0;`},
		{func(s S) S { return s.BorderRight(true) }, `border-right: true;`},
		{func(s S) S { return s.BorderRightBackground(g) }, `border-right-background: #0f0;`},
		{func(s S) S { return s.BorderRightForeground(g) }, `border-right-foreground: #0f0;`},
		{func(s S) S {
			return s.BorderStyle(lipgloss.Border{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"})
		}, `border-style: border("a","b","c","d","e","f","g","h");`},
		{func(s S) S { return s.BorderTop(true) }, `border-top: true;`},
		{func(s S) S { return s.BorderTopBackground(g) }, `border-top-background: #0f0;`},
		{func(s S) S { return s.BorderTopForeground(g) }, `border-top-foreground: #0f0;`},
		{func(s S) S {
			return s.Border(lipgloss.Border{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"}, true)
		}, `border-bottom: true; border-left: true; border-right: true; ` +
			`border-style: border("a","b","c","d","e","f","g","h"); ` +
			`border-top: true;`},
		{func(s S) S { return s.BorderBackground(g) }, `border-bottom-background: #0f0; border-left-background: #0f0; border-right-background: #0f0; border-top-background: #0f0;`},
		{func(s S) S { return s.BorderForeground(g) }, `border-bottom-foreground: #0f0; border-left-foreground: #0f0; border-right-foreground: #0f0; border-top-foreground: #0f0;`},

		{func(s S) S { return s.ColorWhitespace(true) }, `color-whitespace: true;`},
		{func(s S) S { return s.Faint(true) }, `faint: true;`},
		{func(s S) S { return s.Foreground(g) }, `foreground: #0f0;`},
		{func(s S) S { return s.Height(3) }, `height: 3;`},
		{func(s S) S { return s.Inline(true) }, `inline: true;`},
		{func(s S) S { return s.Italic(true) }, `italic: true;`},
		{func(s S) S { return s.Margin(1, 2, 3, 4) }, `margin-bottom: 3; margin-left: 4; margin-right: 2; margin-top: 1;`},
		{func(s S) S { return s.MarginBottom(3) }, `margin-bottom: 3;`},
		{func(s S) S { return s.MarginLeft(3) }, `margin-left: 3;`},
		{func(s S) S { return s.MarginRight(3) }, `margin-right: 3;`},
		{func(s S) S { return s.MarginTop(3) }, `margin-top: 3;`},
		{func(s S) S { return s.MarginBackground(g) }, `margin-background: #0f0;`},
		{func(s S) S { return s.MaxHeight(3) }, `max-height: 3;`},
		{func(s S) S { return s.MaxWidth(3) }, `max-width: 3;`},
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
		// Variable size setters.
		{func(s S) S { return s.Padding(1, 2, 3, 4) }, `padding-bottom: 3; padding-left: 4; padding-right: 2; padding-top: 1;`},
		{func(s S) S { return s.Padding(1, 2, 3) }, `padding-bottom: 3; padding-left: 2; padding-right: 2; padding-top: 1;`},
		{func(s S) S { return s.Padding(1, 2) }, `padding-bottom: 1; padding-left: 2; padding-right: 2; padding-top: 1;`},
		{func(s S) S { return s.Padding(1) }, `padding-bottom: 1; padding-left: 1; padding-right: 1; padding-top: 1;`},
		{func(s S) S { return s.Padding() }, ``},
		{func(s S) S { return s.Padding(1, 2, 3, 4, 5) }, ``},
		{func(s S) S { return s.Margin(1, 2, 3, 4) }, `margin-bottom: 3; margin-left: 4; margin-right: 2; margin-top: 1;`},
		{func(s S) S { return s.Margin(1, 2, 3) }, `margin-bottom: 3; margin-left: 2; margin-right: 2; margin-top: 1;`},
		{func(s S) S { return s.Margin(1, 2) }, `margin-bottom: 1; margin-left: 2; margin-right: 2; margin-top: 1;`},
		{func(s S) S { return s.Margin(1) }, `margin-bottom: 1; margin-left: 1; margin-right: 1; margin-top: 1;`},
		{func(s S) S { return s.Margin() }, ``},
		{func(s S) S { return s.Margin(1, 2, 3, 4, 5) }, ``},
		{func(s S) S {
			return s.Border(lipgloss.Border{}, true, true, true, true)
		}, `border-bottom: true; border-left: true; border-right: true; ` +
			`border-top: true;`},
		{func(s S) S {
			return s.Border(lipgloss.Border{}, true, true, true)
		}, `border-bottom: true; border-left: true; border-right: true; ` +
			`border-top: true;`},
		{func(s S) S {
			return s.Border(lipgloss.Border{}, true, true)
		}, `border-bottom: true; border-left: true; border-right: true; ` +
			`border-top: true;`},
		{func(s S) S {
			return s.Border(lipgloss.Border{})
		}, `border-bottom: true; border-left: true; border-right: true; ` +
			`border-top: true;`},
		{func(s S) S {
			return s.Border(lipgloss.Border{}, true, true, true, true, true)
		}, `border-bottom: true; border-left: true; border-right: true; ` +
			`border-top: true;`},
		{func(s S) S { return s.BorderBackground(g, r, b, w) }, `border-bottom-background: #00f; border-left-background: #fff; border-right-background: #f00; border-top-background: #0f0;`},
		{func(s S) S { return s.BorderBackground(g, r, b) }, `border-bottom-background: #00f; border-left-background: #f00; border-right-background: #f00; border-top-background: #0f0;`},
		{func(s S) S { return s.BorderBackground(g, r) }, `border-bottom-background: #0f0; border-left-background: #f00; border-right-background: #f00; border-top-background: #0f0;`},
		{func(s S) S { return s.BorderBackground() }, ``},
		{func(s S) S { return s.BorderBackground(g, r, b, w, g) }, ``},
		{func(s S) S { return s.BorderForeground(g, r, b, w) }, `border-bottom-foreground: #00f; border-left-foreground: #fff; border-right-foreground: #f00; border-top-foreground: #0f0;`},
		{func(s S) S { return s.BorderForeground(g, r, b) }, `border-bottom-foreground: #00f; border-left-foreground: #f00; border-right-foreground: #f00; border-top-foreground: #0f0;`},
		{func(s S) S { return s.BorderForeground(g, r) }, `border-bottom-foreground: #0f0; border-left-foreground: #f00; border-right-foreground: #f00; border-top-foreground: #0f0;`},
		{func(s S) S { return s.BorderForeground() }, ``},
		{func(s S) S { return s.BorderForeground(g, r, b, w, g) }, ``},
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

		if tc.repr == `` {
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

func TestRender(t *testing.T) {
	curProfile := lipgloss.ColorProfile()
	defer lipgloss.SetColorProfile(curProfile)

	lipgloss.SetColorProfile(termenv.TrueColor)

	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		d := driver{
			s:    lipgloss.NewStyle(),
			text: "hello!",
		}
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			return d.renderTest(t, td)
		})
	})
}

type driver struct {
	s      lipgloss.Style
	text   string
	spaces bool
}

func (d *driver) renderTest(t *testing.T, td *datadriven.TestData) string {
	switch td.Cmd {
	case "text":
		d.text = td.Input
		return "ok"

	case "spvis":
		d.spaces = !d.spaces
		return "ok"

	case "inherit":
		// Make the current style inherit from the specified parent style.
		parentStyle, err := lipglossc.Import(lipgloss.NewStyle(), td.Input)
		if err != nil {
			t.Fatalf("%s: invalid style: %v", td.Pos, err)
		}
		d.s = d.s.Inherit(parentStyle)
		return lipglossc.Export(d.s, lipglossc.WithSeparator("\n"))

	case "set":
		newStyle, err := lipglossc.Import(d.s, td.Input)
		if err != nil {
			t.Fatalf("%s: invalid style: %v", td.Pos, err)
		}
		d.s = newStyle

		o := d.s.Render(d.text)
		o = strings.ReplaceAll(o, "\n", "␤\n")
		if !d.spaces {
			o = strings.ReplaceAll(o, " ", "·")
		}
		// Add a "no newline at end" marker if there was no newline at the end.
		if len(o) == 0 || o[len(o)-1] != '\n' {
			o += "🛇"
		}
		return o

	default:
		t.Fatalf("%s: unknown command: %q", td.Pos, td.Cmd)
		return "" // unreachable
	}
}
