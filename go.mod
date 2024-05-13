module github.com/charmbracelet/lipgloss

retract v0.7.0 // v0.7.0 introduces a bug that causes some apps to freeze.

go 1.18

require (
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/mattn/go-runewidth v0.0.13
	github.com/muesli/reflow v0.2.1-0.20210115123740-9e1d0d53df68
	github.com/muesli/termenv v0.11.1-0.20220204035834-5ac8409525e0
	github.com/charmbracelet/x/exp/term v0.0.0-20240408110044-525ba71bb562
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
