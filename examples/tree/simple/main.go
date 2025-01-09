package main

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/tree"
)

func main() {
	t := tree.Root(".").
		Child("macOS").
		Child(
			tree.New().
				Root("Linux").
				Child("NixOS").
				Child("Arch Linux (btw)").
				Child("Void Linux"),
		).
		Child(
			tree.New().
				Root("BSD").
				Child("FreeBSD").
				Child("OpenBSD"),
		)

	lipgloss.Println(t)
}
