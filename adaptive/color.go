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

// Output represents an adaptive color output based on the terminal color
// profile and background color capability.
type Output struct {
	BackgroundColor lipgloss.TerminalColor
	Profile         lipgloss.Profile
}

// NewOutput returns a new adaptive color output that can be used to create
// adaptive colors.
func NewOutput(input io.Reader, output io.Writer, environ []string) *Output {
	o := &Output{
		BackgroundColor: color.Black, // Default to black background
		Profile:         lipgloss.DetectColorProfile(output, environ),
	}
	if c, ok := queryBackgroundColor(input, output); ok && c != nil {
		o.BackgroundColor = c
	}
	return o
}

// QueryBackgroundColor queries and returns the terminal background color.
func (o *Output) QueryBackgroundColor(input io.Reader, output io.Writer) (color.Color, bool) {
	c, ok := queryBackgroundColor(input, output)
	if ok && c != nil {
		o.BackgroundColor = c
	}
	return c, ok
}

// HasLightBackground returns whether the output has a light background.
func (o Output) HasLightBackground() bool {
	if o.BackgroundColor == nil {
		return false
	}
	if col, ok := colorful.MakeColor(o.BackgroundColor); ok {
		_, _, l := col.Hsl()
		return l > 0.5
	}
	return false
}

// AdaptiveColor returns a color appropriate for the terminal background color.
//
// Example usage:
//
//	color := adaptive.AdaptiveColor("#0000ff", "#000099")
func (o Output) AdaptiveColor(light, dark string) lipgloss.TerminalColor {
	if o.HasLightBackground() {
		return lipgloss.Color(light)
	}
	return lipgloss.Color(dark)
}

// CompleteColor specifies exact values for TrueColor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed, instead the
// appropriate color will be returned at based on the terminal color profile.
func (o Output) CompleteColor(trueColor, ansi256, ansi string) lipgloss.TerminalColor {
	switch o.Profile {
	case lipgloss.TrueColor:
		return lipgloss.Color(trueColor)
	case lipgloss.ANSI256:
		return lipgloss.Color(ansi256)
	case lipgloss.ANSI:
		return lipgloss.Color(ansi)
	default:
		return lipgloss.NoColor{}
	}
}

// AdaptiveColor returns a color appropriate for the terminal background color.
//
// This function is a convenience function that uses the default adaptive color
// output which is based on the standard input and output and environment
// variables.
//
// Example usage:
//
//	color := adaptive.AdaptiveColor("#0000ff", "#000099")
func AdaptiveColor(light, dark string) lipgloss.TerminalColor {
	return Default().AdaptiveColor(light, dark)
}

// CompleteColor specifies exact values for TrueColor, ANSI256, and ANSI color
// profiles. Automatic color degradation will not be performed, instead the
// appropriate color will be returned at based on the terminal color profile.
//
// This function is a convenience function that uses the default adaptive color
// output which is based on the standard input and output and environment
// variables.
func CompleteColor(trueColor, ansi256, ansi string) lipgloss.TerminalColor {
	return Default().CompleteColor(trueColor, ansi256, ansi)
}

// HasLightBackground returns whether the output has a light background.
//
// This function is a convenience function that uses the default adaptive color
// output which is based on the standard input and output and environment
// variables.
func HasLightBackground() bool {
	return Default().HasLightBackground()
}

// QueryBackgroundColor queries the terminal for the background color.
//
// This function is a convenience function that uses the default adaptive color
// output which is based on the standard input and output and environment
// variables.
func QueryBackgroundColor(in io.Reader, out io.Writer) (c color.Color, ok bool) {
	return Default().QueryBackgroundColor(in, out)
}
