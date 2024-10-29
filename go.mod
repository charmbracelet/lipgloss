module github.com/charmbracelet/lipgloss

retract v0.7.0 // v0.7.0 introduces a bug that causes some apps to freeze.

retract v0.11.1 // v0.11.1 uses a broken version of x/ansi StringWidth that causes some lines to wrap incorrectly.

go 1.18

require (
	github.com/aymanbagabas/go-udiff v0.2.0
	github.com/charmbracelet/x/ansi v0.4.2
	github.com/charmbracelet/x/exp/golden v0.0.0-20240806155701-69247e0abc2a
	github.com/muesli/termenv v0.15.2
	github.com/rivo/uniseg v0.4.7
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	golang.org/x/sys v0.19.0 // indirect
)
