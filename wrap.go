package lipgloss

import (
	"strings"

	"github.com/charmbracelet/x/exp/term/ansi"
)

// wrap wraps a string to a given width. It will break the string at spaces
// and hyphens, and will remove any leading or trailing whitespace from the
// wrapped lines.
func wrap(str string, width int) string {
	wrapped := ansi.Wordwrap(str, width, "")
	lines := strings.Split(wrapped, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		linew := ansi.StringWidth(line)
		if linew <= width {
			continue
		}

		wline := ansi.Wordwrap(line, width, "")
		wlines := strings.Split(wline, "\n")
		for j := 0; j < len(wlines); j++ {
			if ansi.StringWidth(wlines[j]) > width {
				wline = ansi.Wrap(line, width, false)
				wlines = strings.Split(wline, "\n")
				break
			}
		}

		if len(wlines) > 0 {
			lines[i] = wlines[0]
		}

		if len(wlines) > 1 && i+1 < len(lines) {
			endsWithHyphen := strings.HasSuffix(wlines[1], "-")
			if endsWithHyphen {
				lines[i+1] = wlines[1] + lines[i+1]
			} else {
				lines[i+1] = wlines[1] + " " + lines[i+1]
			}
		}
	}
	return strings.Join(lines, "\n")
}
