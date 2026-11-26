package lipgloss

import (
	"fmt"
	"image/color"
	"os"
	"runtime"

	"github.com/charmbracelet/x/term"
)

func openConsole(name string) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_RDWR, 0o644) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", name, err)
	}
	if !term.IsTerminal(f.Fd()) {
		_ = f.Close()
		return nil, fmt.Errorf("%s is not a terminal", name)
	}
	return f, nil
}

func backgroundColor(in term.File, out term.File) (color.Color, error) {
	state, err := term.MakeRaw(in.Fd())
	if err != nil {
		return nil, fmt.Errorf("error setting raw state to detect background color: %w", err)
	}

	defer term.Restore(in.Fd(), state) //nolint:errcheck

	bg, err := queryBackgroundColor(in, out)
	if err != nil {
		return nil, err
	}

	return bg, nil
}

// BackgroundColor queries the terminal's background color. Typically, you'll
// want to query against stdin and either stdout or stderr, depending on what
// you're writing to.
//
// This function is intended for standalone Lip Gloss use only. If you're using
// Bubble Tea, listen for tea.BackgroundColorMsg in your update function.
func BackgroundColor(in term.File, out term.File) (bg color.Color, err error) {
	inIsTerm := term.IsTerminal(in.Fd())
	outIsTerm := term.IsTerminal(out.Fd())

	if runtime.GOOS == "windows" { //nolint:nestif
		if !inIsTerm && !outIsTerm {
			return nil, fmt.Errorf("input/output is not a terminal")
		}

		// On Windows, when input or output is redirected or piped, we need to
		// open the console explicitly.
		// See https://learn.microsoft.com/en-us/windows/console/getstdhandle#remarks
		if !inIsTerm {
			f, err := openConsole("CONIN$")
			if err != nil {
				return nil, err
			}
			defer f.Close() //nolint:errcheck
			in = f
		}
		if !outIsTerm {
			f, err := openConsole("CONOUT$")
			if err != nil {
				return nil, err
			}
			defer f.Close() //nolint:errcheck
			out = f
		}
		return backgroundColor(in, out)
	}

	// NOTE: On Unix, both input and output must be a tty.
	if !inIsTerm || !outIsTerm {
		return nil, fmt.Errorf("input/output is not a terminal")
	}

	for _, f := range []term.File{in, out} {
		if bg, err = backgroundColor(f, f); err == nil {
			return bg, nil
		}
	}

	return bg, err
}

// HasDarkBackground detects whether the terminal has a light or dark
// background.
//
// Typically, you'll want to query against stdin and either stdout or stderr
// depending on what you're writing to.
//
//	hasDarkBG := HasDarkBackground(os.Stdin, os.Stdout)
//	lightDark := LightDark(hasDarkBG)
//	myHotColor := lightDark("#ff0000", "#0000ff")
//
// This is intended for use in standalone Lip Gloss only. In Bubble Tea, listen
// for tea.BackgroundColorMsg in your Update function.
//
//	case tea.BackgroundColorMsg:
//	    hasDarkBackground = msg.IsDark()
//
// By default, this function will return true if it encounters an error.
func HasDarkBackground(in term.File, out term.File) bool {
	bg, err := BackgroundColor(in, out)
	if err != nil || bg == nil {
		return true
	}
	return isDarkColor(bg)
}
