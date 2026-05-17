// Package compat is a compatibility layer for Lip Gloss that provides a way to
// deal with the hassle of setting up a writer. It's impure because it uses
// global variables, is not thread-safe, and only works with the default
// standard I/O streams.
//
// Background and color profile detection use [os.Stderr] by default (matching
// log and most Bubble Tea apps). Override detection with [SetHasDarkBackground]
// or by assigning [HasDarkBackground] before rendering any adaptive colors.
//
// In Bubble Tea, sync from tea.BackgroundColorMsg:
//
//	case tea.BackgroundColorMsg:
//	    compat.SetHasDarkBackground(msg.IsDark())
package compat
