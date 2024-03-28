module github.com/charmbracelet/lipgloss

retract v0.7.0 // v0.7.0 introduces a bug that causes some apps to freeze.

go 1.18

require (
	github.com/charmbracelet/x/exp/term v0.0.0-20240328150354-ab9afc214dfd
	github.com/muesli/termenv v0.15.2
	github.com/rivo/uniseg v0.4.7
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	golang.org/x/sys v0.18.0 // indirect
)
