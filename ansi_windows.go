//go:build windows
// +build windows

package lipgloss

import (
	"sync"

	"golang.org/x/sys/windows"
)

var enableANSI sync.Once

// enableANSIColors enables support for ANSI color sequences in the Windows
// default console (cmd.exe and the PowerShell application). Note that this
// only works with Windows 10. Also note that Windows Terminal supports colors
// by default.
func enableLegacyWindowsANSI() {
	enableANSI.Do(func() {
		handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
		if err != nil {
			return
		}

		var mode uint32
		err = windows.GetConsoleMode(handle, &mode)
		if err != nil {
			return
		}

		// See https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
		if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING {
			vtpmode := mode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
			if err := windows.SetConsoleMode(handle, vtpmode); err != nil {
				return
			}
		}
	})
}
