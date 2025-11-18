module examples

go 1.24.3

toolchain go1.24.4

replace charm.land/lipgloss/v2 => ../

require (
	charm.land/bubbletea/v2 v2.0.0-rc.1.0.20251106192006-06c0cda318b3
	charm.land/lipgloss/v2 v2.0.0-beta.3.0.20251106192539-4b304240aab7
	github.com/charmbracelet/colorprofile v0.3.3
	github.com/charmbracelet/ssh v0.0.0-20241211182756-4fe22b0f1b7c
	github.com/charmbracelet/wish/v2 v2.0.0-20251106193208-3cd15da8229f
	github.com/charmbracelet/x/exp/charmtone v0.0.0-20250627134340-c144409e381c
	github.com/charmbracelet/x/term v0.2.2
	github.com/rivo/uniseg v0.4.7
)

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/charmbracelet/keygen v0.5.1 // indirect
	github.com/charmbracelet/log/v2 v2.0.0-20251106192421-eb64aaa963a0 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20251106190538-99ea45596692 // indirect
	github.com/charmbracelet/x/ansi v0.11.1 // indirect
	github.com/charmbracelet/x/conpty v0.1.0 // indirect
	github.com/charmbracelet/x/errors v0.0.0-20240508181413-e8d8b6e2de86 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.5.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/creack/pty v1.1.21 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.3.0 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
)

// replace with log v2
replace github.com/charmbracelet/log => github.com/charmbracelet/log v0.4.1-0.20241010222913-47ce960d4847
