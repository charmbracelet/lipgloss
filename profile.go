package lipgloss

import (
	"fmt"
	"strconv"
	"strings"

	ansi "github.com/charmbracelet/x/exp/term/ansi/style"
	"github.com/lucasb-eyer/go-colorful"
)

// Profile is a color profile: Ascii, ANSI, ANSI256, or TrueColor.
type Profile int

const (
	// Ascii, uncolored profile
	Ascii = Profile(iota) //nolint:revive
	// ANSI, 4-bit color profile
	ANSI
	// ANSI256, 8-bit color profile
	ANSI256
	// TrueColor, 24-bit color profile
	TrueColor
)

// Convert transforms a given Color to a Color supported within the Profile.
func (p Profile) Convert(c ansi.Color) ansi.Color {
	if p == Ascii {
		return nil
	}

	switch v := c.(type) {
	case ansi.BasicColor:
		return v

	case ansi.ExtendedColor:
		if p == ANSI {
			return ansi256ToANSIColor(v)
		}
		return v

	case ansi.TrueColor:
		// TODO: improve this
		r, g, b, _ := v.RGBA()
		h, err := colorful.Hex(fmt.Sprintf("#%02x%02x%02x", r, g, b))
		if err != nil {
			return nil
		}
		if p != TrueColor {
			ac := hexToANSI256Color(h)
			if p == ANSI {
				return ansi256ToANSIColor(ac)
			}
			return ac
		}
		return v
	}

	return c
}

// Color creates a Color from a string. Valid inputs are hex colors, as well as
// ANSI color codes (0-15, 16-255).
func (p Profile) Color(s string) ansi.Color {
	if len(s) == 0 {
		return ansi.Black
	}

	var c ansi.Color
	if strings.HasPrefix(s, "#") {
		col, err := colorful.Hex(s)
		if err != nil {
			return nil
		}
		c = ansi.TrueColor(rgbToHex(
			uint32(col.R*255.0),
			uint32(col.G*255.0),
			uint32(col.B*255.0),
		))
	} else {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}

		if i < 16 {
			c = ansi.BasicColor(i)
		} else {
			c = ansi.ExtendedColor(i)
		}
	}

	return p.Convert(c)
}

func envColorProfile(isatty bool, env Environ) Profile {
	if envNoColor(env) {
		return Ascii
	}
	p := colorProfile(isatty, env)
	if cliColorForced(env) && p == Ascii {
		return ANSI
	}
	return p
}

// colorProfile returns the color profile for the given environment.
func colorProfile(isatty bool, env Environ) Profile {
	if !isatty {
		return Ascii
	}

	if env.Getenv("GOOGLE_CLOUD_SHELL") == "true" {
		return TrueColor
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
				return ANSI256
			}
		}
		return TrueColor
	case "yes":
		fallthrough
	case "true":
		return ANSI256
	}

	switch term {
	case "xterm-kitty", "wezterm", "xterm-ghostty":
		return TrueColor
	case "linux":
		return ANSI
	}

	if strings.Contains(term, "256color") {
		return ANSI256
	}
	if strings.Contains(term, "color") {
		return ANSI
	}
	if strings.Contains(term, "ansi") {
		return ANSI
	}

	return Ascii
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

func (p Profile) string() style {
	return style{p: p}
}
