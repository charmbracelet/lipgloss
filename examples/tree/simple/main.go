package main

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
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
