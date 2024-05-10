module github.com/charmbracelet/lipgloss

retract v0.7.0 // v0.7.0 introduces a bug that causes some apps to freeze.

go 1.18

require (
	github.com/charmbracelet/x/exp/term v0.0.0-20240425164147-ba2a9512b05f
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/rivo/uniseg v0.4.7
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e
	golang.org/x/sys v0.19.0
)

require (
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
)
