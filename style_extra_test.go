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
