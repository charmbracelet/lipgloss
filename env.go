package lipgloss

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/lucasb-eyer/go-colorful"
)

// DetectColorProfile returns the color profile based on the terminal output,
// and environment variables.
//
// If the output is not a terminal, the color profile will be NoTTY unless
// CLICOLOR_FORCE=1 is set. This respects the NO_COLOR and
// CLICOLOR/CLICOLOR_FORCE environment variables.
//
// See https://no-color.org/ and https://bixense.com/clicolors/ for more information.
func DetectColorProfile(output io.Writer, environ []string) Profile {
	if environ == nil {
		environ = os.Environ()
	}

	env := environMap(environ)
	p := envColorProfile(env)
	out, ok := output.(term.File)
	isatty := ok && term.IsTerminal(out.Fd())

	if isTrue(env["NO_COLOR"]) {
		if isatty {
			return Ascii
		}
		return NoTTY
	} else if isTrue(env["CLICOLOR_FORCE"]) {
		return p
	}

	if !isatty {
		return NoTTY
	}

	return p
}

// QueryHasLightBackground returns true if the terminal has a light background.
func QueryHasLightBackground(in term.File, out term.File) bool {
	if !term.IsTerminal(out.Fd()) {
		return false
	}

	state, err := term.MakeRaw(in.Fd())
	if err != nil {
		return false
	}

	defer term.Restore(in.Fd(), state) // nolint:errcheck

	c, err := term.QueryBackgroundColor(in, out)
	if err != nil {
		return false
	}

	col, ok := colorful.MakeColor(c)
	if !ok {
		return false
	}

	_, _, l := col.Hsl()
	return l > 0.5
}

// EnvColorProfile returns the color profile based on environment variables.
//
// This respects the NO_COLOR and CLICOLOR/CLICOLOR_FORCE environment
// variables.
//
// See https://no-color.org/ and https://bixense.com/clicolors/ for more information.
func EnvColorProfile(environ []string) Profile {
	if environ == nil {
		environ = os.Environ()
	}

	env := environMap(environ)
	return envColorProfile(env)
}

// envColorProfile returns infers the color profile from the environment.
func envColorProfile(env map[string]string) (p Profile) {
	p = Ascii // Default to ASCII
	setProfile := func(profile Profile) {
		if profile < p {
			p = profile
		}
	}

	if isTrue(env["GOOGLE_CLOUD_SHELL"]) {
		setProfile(TrueColor)
	}

	term := env["TERM"]
	colorTerm := env["COLORTERM"]

	switch strings.ToLower(colorTerm) {
	case "24bit":
		fallthrough
	case "truecolor":
		if strings.HasPrefix(term, "screen") {
			// tmux supports TrueColor, screen only ANSI256
			if env["TERM_PROGRAM"] != "tmux" {
				setProfile(ANSI256)
			}
		}
		setProfile(TrueColor)
	case "yes":
		fallthrough
	case "true":
		setProfile(TrueColor)
	}

	switch term {
	case "xterm-kitty", "wezterm", "xterm-ghostty":
		setProfile(TrueColor)
	case "linux":
		setProfile(ANSI)
	}

	if strings.Contains(term, "256color") {
		setProfile(ANSI256)
	}
	if strings.Contains(term, "color") {
		setProfile(ANSI)
	}
	if strings.Contains(term, "ansi") {
		setProfile(ANSI)
	}

	return
}

// isTrue returns true if the string is a truthy value.
func isTrue(s string) bool {
	if s == "" {
		return false
	}
	v, _ := strconv.ParseBool(strings.ToLower(s))
	return v
}

// environMap converts an environment slice to a map.
func environMap(environ []string) map[string]string {
	m := make(map[string]string, len(environ))
	for _, e := range environ {
		parts := strings.SplitN(e, "=", 2)
		var value string
		if len(parts) == 2 {
			value = parts[1]
		}
		m[parts[0]] = value
	}
	return m
}
