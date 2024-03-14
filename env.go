package lipgloss

import (
	"strconv"
	"strings"
)

// envNoColor returns true if the environment variables explicitly disable color output
// by setting NO_COLOR (https://no-color.org/)
// or CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If NO_COLOR is set, this will return true, ignoring CLICOLOR/CLICOLOR_FORCE
// If CLICOLOR=="0", it will be true only if CLICOLOR_FORCE is also "0" or is unset.
func (o *Renderer) envNoColor() bool {
	return o.environ["NO_COLOR"] != "" || (o.environ["CLICOLOR"] == "0" && !o.cliColorForced())
}

// envColorProfile returns the color profile based on environment variables set
// Supports NO_COLOR (https://no-color.org/)
// and CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If none of these environment variables are set, this behaves the same as ColorProfile()
// It will return the Ascii color profile if EnvNoColor() returns true
// If the terminal does not support any colors, but CLICOLOR_FORCE is set and not "0"
// then the ANSI color profile will be returned.
func (o *Renderer) envColorProfile() Profile {
	if o.envNoColor() {
		return Ascii
	}
	p := o.detectColorProfile()
	if o.cliColorForced() && p == Ascii {
		return ANSI
	}
	return p
}

func (o *Renderer) cliColorForced() bool {
	if forced := o.environ["CLICOLOR_FORCE"]; forced != "" {
		return !isTrue(forced)
	}
	return false
}

// detectColorProfile returns the supported color profile:
// Ascii, ANSI, ANSI256, or TrueColor.
func (o *Renderer) detectColorProfile() (p Profile) {
	if !o.isatty {
		return Ascii
	}

	setProfile := func(profile Profile) {
		if profile > p {
			p = profile
		}
	}

	if isTrue(o.environ["GOOGLE_CLOUD_SHELL"]) {
		setProfile(TrueColor)
	}

	term := o.environ["TERM"]
	colorTerm := o.environ["COLORTERM"]

	switch strings.ToLower(colorTerm) {
	case "24bit":
		fallthrough
	case "truecolor":
		if strings.HasPrefix(term, "screen") {
			// tmux supports TrueColor, screen only ANSI256
			if o.environ["TERM_PROGRAM"] != "tmux" {
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
