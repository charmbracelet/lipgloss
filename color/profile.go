package color

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Profile is a color profile.
type Profile termenv.Profile

// NewProfile returns a new color profile.
// If env is nil, it will use OsEnviron.
func NewProfile(env Environ) Profile {
	if env == nil {
		env = OsEnviron{}
	}

	if envNoColor(env) {
		return Profile(termenv.Ascii)
	}
	p := colorProfile(env)
	if cliColorForced(env) && p == termenv.Ascii {
		return Profile(termenv.ANSI)
	}
	return Profile(p)
}

// Color returns a TerminalColor based on the color profile.
func (p Profile) Color(c string) lipgloss.TerminalColor {
	return termenv.Profile(p).Color(c)
}

func colorProfile(env Environ) termenv.Profile {
	if env.Getenv("GOOGLE_CLOUD_SHELL") == "true" {
		return termenv.TrueColor
	}

	term := env.Getenv("TERM")
	colorTerm := env.Getenv("COLORTERM")

	switch strings.ToLower(colorTerm) {
	case "24bit":
		fallthrough
	case "truecolor":
		if strings.HasPrefix(term, "screen") {
			// tmux supports TrueColor, screen only ANSI256
			if env.Getenv("TERM_PROGRAM") != "tmux" {
				return termenv.ANSI256
			}
		}
		return termenv.TrueColor
	case "yes":
		fallthrough
	case "true":
		return termenv.ANSI256
	}

	switch term {
	case "xterm-kitty", "wezterm", "xterm-ghostty":
		return termenv.TrueColor
	case "linux":
		return termenv.ANSI
	}

	if strings.Contains(term, "256color") {
		return termenv.ANSI256
	}
	if strings.Contains(term, "color") {
		return termenv.ANSI
	}
	if strings.Contains(term, "ansi") {
		return termenv.ANSI
	}

	return termenv.Ascii
}

func envNoColor(env Environ) bool {
	return env.Getenv("NO_COLOR") != "" || (env.Getenv("CLICOLOR") == "0" && !cliColorForced(env))
}

func cliColorForced(env Environ) bool {
	if forced := env.Getenv("CLICOLOR_FORCE"); forced != "" {
		return forced != "0"
	}
	return false
}
