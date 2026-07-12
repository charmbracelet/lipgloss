package lipgloss

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/charmbracelet/x/term"
)

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
	// Detecting the background color means writing an escape sequence and blocking
	// until the terminal answers it. Only a terminal ever answers, so both handles
	// have to be terminals.
	//
	// Windows used to be special-cased here: when a handle was not a terminal it
	// opened CONIN$/CONOUT$ and queried the console anyway. That made a program
	// whose stdio is redirected -- a service, a CI job -- query a console nobody is
	// watching, and wait for a reply that never comes.
	if !term.IsTerminal(in.Fd()) || !term.IsTerminal(out.Fd()) {
		return nil, fmt.Errorf("input/output is not a terminal")
	}

	if runtime.GOOS == "windows" {
		// Console input and output are separate handles on Windows, so query the
		// pair rather than treating either one as bidirectional.
		return backgroundColor(in, out)
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
