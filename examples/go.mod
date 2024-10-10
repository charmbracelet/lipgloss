module examples

go 1.19

replace github.com/charmbracelet/lipgloss => ../

replace github.com/charmbracelet/lipgloss/list => ../list

replace github.com/charmbracelet/lipgloss/table => ../table

require (
	github.com/charmbracelet/bubbletea/v2 v2.0.0-alpha.1
	github.com/charmbracelet/colorprofile v0.1.2
	github.com/charmbracelet/lipgloss v0.13.1-0.20240822211938-b89f1a3db2a4
	github.com/charmbracelet/ssh v0.0.0-20240401141849-854cddfa2917
	github.com/charmbracelet/wish v1.4.0
	github.com/charmbracelet/x/term v0.2.0
	github.com/lucasb-eyer/go-colorful v1.2.0
)

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/bubbletea v0.25.0 // indirect
	github.com/charmbracelet/keygen v0.5.0 // indirect
	github.com/charmbracelet/log v0.4.0 // indirect
	github.com/charmbracelet/x/ansi v0.3.2 // indirect
	github.com/charmbracelet/x/errors v0.0.0-20240117030013-d31dba354651 // indirect
	github.com/charmbracelet/x/exp/term v0.0.0-20240328150354-ab9afc214dfd // indirect
	github.com/charmbracelet/x/input v0.2.0 // indirect
	github.com/charmbracelet/x/windows v0.2.0 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/creack/pty v1.1.21 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

// replace with log v2
replace github.com/charmbracelet/log => github.com/charmbracelet/log v0.4.1-0.20241010222913-47ce960d4847
