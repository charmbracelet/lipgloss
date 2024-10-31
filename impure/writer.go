package impure

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/colorprofile"
)

// Println to stdout, automatically downsampling colors when necessary, ending
// with a trailing newline.
//
// Example:
//
//	str := NewStyle().
//	    Foreground(lipgloss.Color("#6a00ff")).
//	    Render("breakfast")
//
//	Println("Time for a", str, "sandwich!")
func Println(v ...interface{}) (int, error) {
	return fmt.Fprintln(Writer, v...) //nolint:wrapcheck
}

// Print formatted text to stdout, automatically downsampling colors when
// necessary.
//
// Example:
//
//	str := NewStyle().
//	  Foreground(lipgloss.Color("#6a00ff")).
//	  Render("knuckle")
//
//	Printf("Time for a %s sandwich!\n", str)
func Printf(format string, v ...interface{}) (int, error) {
	return fmt.Fprintf(Writer, format, v...) //nolint:wrapcheck
}

// Print to stdout, automatically downsampling colors when necessary.
//
// Example:
//
//	str := NewStyle().
//	    Foreground(lipgloss.Color("#6a00ff")).
//	    Render("Who wants marmalade?\n")
//
//	Print(str)
func Print(v ...interface{}) (int, error) {
	return fmt.Fprint(Writer, v...) //nolint:wrapcheck
}

// Fprint pritnts to the given writer, automatically downsampling colors when
// necessary.
//
// Example:
//
//	str := NewStyle().
//	    Foreground(lipgloss.Color("#6a00ff")).
//	    Render("guzzle")
//
//	Fprint(os.Stderr, "I %s horchata pretty much all the time.\n", str)
func Fprint(w io.Writer, v ...interface{}) (int, error) {
	return fmt.Fprint(colorprofile.NewWriter(w, os.Environ()), v...) //nolint:wrapcheck
}

// Fprint pritnts to the given writer, automatically downsampling colors when
// necessary, and ending with a trailing newline.
//
// Example:
//
//	str := NewStyle().
//	    Foreground(lipgloss.Color("#6a00ff")).
//	    Render("Sandwich time!")
//
//	Fprintln(os.Stderr, str)
func Fprintln(w io.Writer, v ...interface{}) (int, error) {
	return fmt.Fprintln(colorprofile.NewWriter(w, os.Environ()), v...) //nolint:wrapcheck
}

// Fprintf prints text to a writer, against the given format, automatically
// downsampling colors when necessary.
//
// Example:
//
//	str := NewStyle().
//	    Foreground(lipgloss.Color("#6a00ff")).
//	    Render("artichokes")
//
//	Fprintf(os.Stderr, "I really love %s!\n", food)
func Fprintf(w io.Writer, format string, v ...interface{}) (int, error) {
	return fmt.Fprintf(colorprofile.NewWriter(w, os.Environ()), format, v...) //nolint:wrapcheck
}
