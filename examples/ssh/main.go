package main

// This example demonstrates how to use a custom Lip Gloss renderer with Wish,
// a package for building custom SSH servers.
//
// The big advantage to using custom renderers here is that we can accurately
// detect the background color and color profile for each client and render
// against that accordingly.
//
// For details on wish see: https://github.com/charmbracelet/wish/

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/adaptive"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/lucasb-eyer/go-colorful"
)

// Available styles.
type styles struct {
	bold          lipgloss.Style
	faint         lipgloss.Style
	italic        lipgloss.Style
	underline     lipgloss.Style
	strikethrough lipgloss.Style
	red           lipgloss.Style
	green         lipgloss.Style
	yellow        lipgloss.Style
	blue          lipgloss.Style
	magenta       lipgloss.Style
	cyan          lipgloss.Style
	gray          lipgloss.Style
}

// Create new styles against a given renderer.
func makeStyles() styles {
	return styles{
		bold:          lipgloss.NewStyle().SetString("bold").Bold(true),
		faint:         lipgloss.NewStyle().SetString("faint").Faint(true),
		italic:        lipgloss.NewStyle().SetString("italic").Italic(true),
		underline:     lipgloss.NewStyle().SetString("underline").Underline(true),
		strikethrough: lipgloss.NewStyle().SetString("strikethrough").Strikethrough(true),
		red:           lipgloss.NewStyle().SetString("red").Foreground(lipgloss.Color("#E88388")),
		green:         lipgloss.NewStyle().SetString("green").Foreground(lipgloss.Color("#A8CC8C")),
		yellow:        lipgloss.NewStyle().SetString("yellow").Foreground(lipgloss.Color("#DBAB79")),
		blue:          lipgloss.NewStyle().SetString("blue").Foreground(lipgloss.Color("#71BEF2")),
		magenta:       lipgloss.NewStyle().SetString("magenta").Foreground(lipgloss.Color("#D290E4")),
		cyan:          lipgloss.NewStyle().SetString("cyan").Foreground(lipgloss.Color("#66C2CD")),
		gray:          lipgloss.NewStyle().SetString("gray").Foreground(lipgloss.Color("#B9BFCA")),
	}
}

// Handle SSH requests.
func handler(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		pty, _, active := sess.Pty()
		if !active {
			next(sess)
			return
		}

		environ := sess.Environ()
		environ = append(environ, fmt.Sprintf("TERM=%s", pty.Term))
		output := adaptive.NewOutput(pty.Slave, pty.Slave, environ)
		width := pty.Window.Width

		// Initialize new styles against the renderer.
		styles := makeStyles()

		str := strings.Builder{}

		fmt.Fprintf(&str, "\n\nProfile: %s\n%s %s %s %s %s",
			output.ColorProfile().String(),
			styles.bold,
			styles.faint,
			styles.italic,
			styles.underline,
			styles.strikethrough,
		)

		fmt.Fprintf(&str, "\n%s %s %s %s %s %s %s",
			styles.red,
			styles.green,
			styles.yellow,
			styles.blue,
			styles.magenta,
			styles.cyan,
			styles.gray,
		)

		fmt.Fprintf(&str, "\n%s %s %s %s %s %s %s\n\n",
			styles.red,
			styles.green,
			styles.yellow,
			styles.blue,
			styles.magenta,
			styles.cyan,
			styles.gray,
		)

		col, _ := colorful.MakeColor(output.BackgroundColor)
		fmt.Fprintf(&str, "%s %t %s\n\n", styles.bold.UnsetString().Render("Has dark background?"),
			output.HasDarkBackground(),
			col.Hex())

		block := lipgloss.Place(width,
			lipgloss.Height(str.String()), lipgloss.Center, lipgloss.Center, str.String(),
			lipgloss.WithWhitespaceChars("/"),
			lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Foreground(output.AdaptiveColor(lipgloss.Color(250), lipgloss.Color(236)))),
		)

		// Render to client.
		output.Println(block)

		next(sess)
	}
}

func main() {
	port := 3456
	s, err := wish.NewServer(
		ssh.AllocatePty(),
		wish.WithAddress(fmt.Sprintf(":%d", port)),
		wish.WithHostKeyPath("ssh_example"),
		wish.WithMiddleware(handler),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("SSH server listening on port %d", port)
	log.Printf("To connect from your local machine run: ssh localhost -p %d", port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
