module charm.land/lipgloss/v2

retract v2.0.0-beta1 // We add a "." after the "beta" in the version number.

go 1.24.2

toolchain go1.24.4

replace github.com/charmbracelet/ultraviolet => ../ultraviolet/

require (
	github.com/aymanbagabas/go-udiff v0.3.1
	github.com/charmbracelet/colorprofile v0.3.2
	github.com/charmbracelet/ultraviolet v0.0.0-20251103104523-426d0f30df4f
	github.com/charmbracelet/x/ansi v0.10.3
	github.com/charmbracelet/x/exp/golden v0.0.0-20240806155701-69247e0abc2a
	github.com/charmbracelet/x/term v0.2.2
	github.com/clipperhouse/displaywidth v0.4.1
	github.com/lucasb-eyer/go-colorful v1.3.0
	github.com/rivo/uniseg v0.4.7
	golang.org/x/sys v0.37.0
)

require (
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/sync v0.17.0 // indirect
)
