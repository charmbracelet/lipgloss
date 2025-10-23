module github.com/charmbracelet/lipgloss/v2

retract v2.0.0-beta1 // We add a "." after the "beta" in the version number.

go 1.24.2

toolchain go1.24.4

require (
	github.com/aymanbagabas/go-udiff v0.2.0
	github.com/charmbracelet/colorprofile v0.3.2
	github.com/charmbracelet/ultraviolet v0.0.0-20250915111650-81d4262876ef
	github.com/charmbracelet/x/ansi v0.10.1
	github.com/charmbracelet/x/cellbuf v0.0.13
	github.com/charmbracelet/x/exp/golden v0.0.0-20240806155701-69247e0abc2a
	github.com/charmbracelet/x/term v0.2.1
	github.com/lucasb-eyer/go-colorful v1.3.0
	github.com/rivo/uniseg v0.4.7
	golang.org/x/sys v0.36.0
)

require (
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/sync v0.17.0 // indirect
)
