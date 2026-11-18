// Package compat is a compatibility layer for Lip Gloss that provides a way to
// deal with the hassle of setting up a writer. It's impure because it uses
// global variables, is not thread-safe, and only works with the default
// standard I/O streams.
//
// In case you want [os.Stderr] to be used as the default writer, you can set
// both [Profile] and [HasDarkBackground] to use [os.Stderr] with
// the following code:
//
//	import (
//		"os"
//
//		"charm.land/lipgloss/v2"
//		"charm.land/lipgloss/v2/compat"
//		"github.com/charmbracelet/colorprofile"
//	)
//
//	func init() {
//		compat.Profile = colorprofile.Detect(os.Stderr, os.Environ())
//		compat.HasDarkBackground = lipgloss.HasDarkBackground(os.Stdin, os.Stderr)
//	}
package compat
