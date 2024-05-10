package lipgloss

import (
	"os"

	"github.com/charmbracelet/x/term"
)

var (
	// ColorProfile is the color profile used by lipgloss.
	// This is the default color profile used to create new styles.
	// By default, it allows for 24-bit color (TrueColor), decorations, and
	// doesn't do color conversion.
	ColorProfile Profile

	// HasLightBackground is true if the terminal has a light background.
	// This is the default value used to create new styles.
	HasLightBackground bool
)

// UseDefault will set the default color profile and background color detection
// from the given terminal file descriptors and environment variables.
func UseDefault(in term.File, out term.File, env []string) {
	ColorProfile = DetectColorProfile(out, env)
	HasLightBackground = QueryHasLightBackground(in, out)
}

// UseStdDefaults will set the default color profile and background color
// detection from the standard input, output, and OS environment variables.
func UseStdDefaults() {
	UseDefault(os.Stdin, os.Stdout, os.Environ())
}
