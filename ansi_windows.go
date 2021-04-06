// +build windows

package lipgloss

import (
	"os"
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
		stdout := windows.Handle(os.Stdout.Fd())
		var originalMode uint32

		windows.GetConsoleMode(stdout, &originalMode)
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	})
}
