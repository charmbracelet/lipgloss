package adaptive

import (
	"image/color"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	defaultOutput     atomic.Pointer[Output]
	defaultOutputOnce sync.Once
)

// Default returns the default adaptive color output.
func Default() *Output {
	d := defaultOutput.Load()
	if d == nil {
		defaultOutputOnce.Do(func() {
			defaultOutput.Store(NewOutput(os.Stdin, os.Stdout, os.Environ()))
		})
		d = defaultOutput.Load()
	}
	return d
}

// SetDefault sets the default adaptive color output.
func SetDefault(o *Output) {
	defaultOutput.Store(o)
}

// Output represents an adaptive color output based on the terminal color
// profile and background color capability.
type Output struct {
	in  io.Reader
	out io.Writer

	Writer          *lipgloss.Writer
	BackgroundColor lipgloss.TerminalColor
}

// NewOutput returns a new adaptive color output that can be used to create
// adaptive colors.
func NewOutput(input io.Reader, output io.Writer, environ []string) *Output {
	o := &Output{
		in:              input,
		out:             output,
		Writer:          lipgloss.NewWriter(output, environ),
		BackgroundColor: color.Black, // Default to black background
	}
	if c, ok := queryBackgroundColor(input, output); ok && c != nil {
		o.BackgroundColor = c
	}
	return o
}

// Print writes the given text to the underlying writer.
func (o Output) Print(text string) (n int, err error) {
	return o.Writer.Print(text)
}

// Println writes the given text to the underlying writer followed by a newline.
func (o Output) Println(text string) (n int, err error) {
	return o.Writer.Println(text)
}

// Printf writes the given text to the underlying writer with the given format.
func (o Output) Printf(format string, a ...interface{}) (n int, err error) {
	return o.Writer.Printf(format, a...)
}

// ColorProfile returns the color profile of the output.
//
// This is a convenience function that returns the color profile of the
// underlying Lip Gloss writer.
func (o Output) ColorProfile() lipgloss.Profile {
	return o.Writer.Profile
}

// QueryBackgroundColor queries, caches, and returns the terminal background
// color.
//
// It returns false if the background color could not be determined.
func (o *Output) QueryBackgroundColor() (color.Color, bool) {
	c, ok := queryBackgroundColor(o.in, o.out)
	if ok && c != nil {
		o.BackgroundColor = c
	}
	return c, ok
}

// HasDarkBackground returns whether the output has a dark background.
//
// If the background color is not set, it will default to true.
func (o Output) HasDarkBackground() bool {
	if o.BackgroundColor == nil {
		return true
	}
	if col, ok := colorful.MakeColor(o.BackgroundColor); ok {
		_, _, l := col.Hsl()
		return l < 0.5
	}
	return true
}

// AdaptiveColor returns a color appropriate for the terminal background color.
//
// Example usage:
//
//	light := lipgloss.Color("#0000ff")
//	dark := lipgloss.Color("#000099")
//	color := adaptive.AdaptiveColor(light, dark)
func (o Output) AdaptiveColor(light, dark lipgloss.TerminalColor) lipgloss.TerminalColor {
	if o.HasDarkBackground() {
		return dark
	}
	return light
}

// CompleteColor specifies exact values for TrueColor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed, instead the
// appropriate color will be returned at based on the terminal color profile.
func (o Output) CompleteColor(trueColor, ansi256, ansi lipgloss.TerminalColor) lipgloss.TerminalColor {
	switch o.ColorProfile() {
	case lipgloss.TrueColor:
		return trueColor
	case lipgloss.ANSI256:
		return ansi256
	case lipgloss.ANSI:
		return ansi
	default:
		return lipgloss.NoColor{}
	}
}
