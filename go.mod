module github.com/charmbracelet/lipgloss

retract v0.7.0 // v0.7.0 introduces a bug that causes some apps to freeze.

go 1.18

require (
	github.com/mattn/go-runewidth v0.0.15
	github.com/muesli/reflow v0.3.0
	github.com/muesli/termenv v0.15.2
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
)
