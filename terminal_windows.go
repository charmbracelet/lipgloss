//go:build windows

package lipgloss

import (
	"io"

	"golang.org/x/sys/windows"
)

func cancelRead(r io.Reader) bool {
	f, ok := r.(interface{ Fd() uintptr })
	if !ok {
		return false
	}
	return windows.CancelIoEx(windows.Handle(f.Fd()), nil) == nil
}
