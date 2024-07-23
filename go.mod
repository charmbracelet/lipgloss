module github.com/charmbracelet/lipgloss

retract v0.7.0 // v0.7.0 introduces a bug that causes some apps to freeze.

retract v0.11.1 // v0.11.1 uses a broken version of x/ansi StringWidth that causes some lines to wrap incorrectly.

go 1.18

require (
	github.com/aymanbagabas/go-udiff v0.2.0
	github.com/charmbracelet/x/ansi v0.1.4
	github.com/charmbracelet/x/term v0.1.1
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/rivo/uniseg v0.4.7
	golang.org/x/sys v0.20.0
)

require (
	github.com/charmbracelet/x/input v0.1.0 // indirect
	github.com/charmbracelet/x/windows v0.1.0 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
)
