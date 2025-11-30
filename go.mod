module github.com/charmbracelet/lipgloss

retract v0.7.0 // v0.7.0 introduces a bug that causes some apps to freeze.

retract v0.11.1 // v0.11.1 uses a broken version of x/ansi StringWidth that causes some lines to wrap incorrectly.

go 1.24.0

require (
	github.com/aymanbagabas/go-udiff v0.3.1
	github.com/charmbracelet/x/ansi v0.10.2
	github.com/charmbracelet/x/cellbuf v0.0.13
	github.com/charmbracelet/x/exp/golden v0.0.0-20250609102027-b60490452b30
	github.com/clipperhouse/displaywidth v0.6.1
	github.com/muesli/termenv v0.16.0
	github.com/rivo/uniseg v0.4.7
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/colorprofile v0.2.3-0.20250311203215-f60798e515dc // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.17 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/sys v0.30.0 // indirect
)
