// Package defaults is a helper package that sets the default Lip Gloss profile
// and background color detection from the standard input, output, and OS
// environment variables.
//
// This is useful for standalone Lip Gloss applications that only use the
// standard input and output i.e. only run locally in a terminal.
// You can simply import this package to set the default profile and background
// color detection for Lip Gloss styles.
//
//	package main
//
//	import (
//		"github.com/charmbracelet/lipgloss"
//		_ "github.com/charmbracelet/lipgloss/defaults" // use std profile defaults
//	)
//
//	func main() {
//		// Your code here
//	}
package defaults

import "github.com/charmbracelet/lipgloss"

func init() {
	lipgloss.UseStdDefaults()
}
