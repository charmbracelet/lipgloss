package lipgloss

import (
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/exp/term"
	"github.com/xo/terminfo"
)

// DetectColorProfile returns the color profile based on the terminal output,
// and environment variables.
//
// If the output is not a terminal, the color profile will be NoTTY unless
// CLICOLOR_FORCE=1 is set. This respects the NO_COLOR and
// CLICOLOR/CLICOLOR_FORCE environment variables.
//
// See https://no-color.org/ and https://bixense.com/clicolors/ for more information.
func DetectColorProfile(stdout *os.File, environ []string) Profile {
	if environ == nil {
		environ = os.Environ()
	}

	env := environMap(environ)
	p := envColorProfile(env)
	if stdout == nil || !term.IsTerminal(stdout.Fd()) {
		p = NoTTY
	}

	if envNoColor(env) && p > Ascii {
		return Ascii
	}

	if cliColorForced(env) && p <= Ascii {
		p = ANSI
		if cp := envColorProfile(env); cp > p {
			p = cp
		}
	}

	return p
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
	p := envColorProfile(env)
	if envNoColor(env) && p > Ascii {
		return Ascii
	}

	if cliColorForced(env) && p <= Ascii {
		p = ANSI
		if cp := envColorProfile(env); cp > p {
			p = cp
		}
	}

	return p
}

// envNoColor returns true if the environment variables explicitly disable color output
// by setting NO_COLOR (https://no-color.org/)
// or CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If NO_COLOR is set, this will return true, ignoring CLICOLOR/CLICOLOR_FORCE
// If CLICOLOR is false, it will be true only if CLICOLOR_FORCE is also false or is unset.
func envNoColor(env map[string]string) bool {
	return isTrue(env["NO_COLOR"]) && !isTrue(env["CLICOLOR"]) && !cliColorForced(env)
}

func cliColorForced(env map[string]string) bool {
	if forced := env["CLICOLOR_FORCE"]; forced != "" {
		return isTrue(forced)
	}
	return false
}

// envColorProfile returns infers the color profile from the environment.
func envColorProfile(env map[string]string) (p Profile) {
	setProfile := func(profile Profile) {
		if profile > p {
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

	ti, _ := terminfo.Load(term)
	if ti != nil {
		extbools := ti.ExtBoolCapsShort()
		if _, ok := extbools["RGB"]; ok {
			setProfile(TrueColor)
		}

		if _, ok := extbools["Tc"]; ok {
			setProfile(TrueColor)
		}

		nums := ti.NumCapsShort()
		if colors, ok := nums["colors"]; ok {
			if colors >= 0x1000000 {
				setProfile(TrueColor)
			} else if colors >= 0x100 {
				setProfile(ANSI256)
			} else if colors >= 0x10 {
				setProfile(ANSI)
			}
		}
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
