//go:build !windows

package lipgloss

import "io"

func cancelRead(io.Reader) bool {
	return false
}
