package lipgloss

import (
	"errors"
	"fmt"
	"image/color"
	"os"

	"github.com/charmbracelet/x/term"
)

// BackgroundColor queries the terminal's background color. Typically, you'll
// want to query against stdin and either stdout or stderr, depending on what
// you're writing to.
//
// This function is intended for standalone Lip Gloss use only. If you're using
// Bubble Tea, listen for tea.BackgroundColorMsg in your update function.
func BackgroundColor(in *os.File, out *os.File) (color.Color, error) {
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

// HasDarkBackground detects whether the terminal has a light or dark
// background. It's a convenience function that wraps [BackgroundColor] and
// [lipgloss.IsDarkColor].
//
// Typically, you'll want to query against stdin and either stdout or stderr
// depending on what you're writing to.
//
//	hasDarkBG, _ := HasDarkBackground(os.Stdin, os.Stdout)
//	adaptive := Adapt(hasDarkBG)
//	myHotColor := Color("#ff0000", "#0000ff")
//
// This is intedded for use in standalone Lip Gloss only. In Bubble Tea, listen
// for tea.BackgroundColorMsg in your update function.
//
//	case tea.BackgroundColorMsg:
//	    hasDarkBackground = msg.IsDark()
func HasDarkBackground(in *os.File, out *os.File) (bool, error) {
	bg, err := BackgroundColor(in, out)
	if err != nil {
		return true, fmt.Errorf("could not detect background color: %w", err)
	}
	if bg == nil {
		return true, errors.New("detected background color is nil")
	}

	return IsDarkColor(bg), nil
}
