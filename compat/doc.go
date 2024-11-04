// Package impure is a compatibility layer for Lip Gloss that provides a way to
// deal with the hassle of setting up a writer. It's impure because it uses
// global variables, is not thread-safe, and only works with the default
// standard I/O streams.
//
// In case you want [os.Stderr] to be used as the default writer, you can set
// both [Writer] and [HasDarkBackground] to use [os.Stderr] with
// the following code:
//
//	import (
//		"os"
//
//		"github.com/charmbracelet/colorprofile"
//		"github.com/charmbracelet/lipgloss/v2/impure"
//	)
//
//	func init() {
//		impure.Writer = colorprofile.NewWriter(os.Stderr, os.Environ())
//		impure.HasDarkBackground, _ = lipgloss.HasDarkBackground(os.Stdin, os.Stderr)
//	}
package compat
