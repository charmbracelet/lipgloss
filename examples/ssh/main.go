package main

// This example demonstrates how to use a custom Lip Gloss renderer with Wish,
// a package for building custom SSH servers.
//
// For details on wish see: https://github.com/charmbracelet/wish/

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/wish"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/muesli/termenv"
)

type sshOutput struct {
	ssh.Session
	tty *os.File
}

func (s *sshOutput) Write(p []byte) (int, error) {
	return s.Session.Write(p)
}

func (s *sshOutput) Fd() uintptr {
	return s.tty.Fd()
}

type sshEnviron struct {
	environ []string
}

func (s *sshEnviron) Getenv(key string) string {
	for _, v := range s.environ {
		if strings.HasPrefix(v, key+"=") {
			return v[len(key)+1:]
		}
	}
	return ""
}

func (s *sshEnviron) Environ() []string {
	return s.environ
}

func outputFromSession(s ssh.Session) *termenv.Output {
	sshPty, _, _ := s.Pty()
	_, tty, err := pty.Open()
	if err != nil {
		panic(err)
	}
	o := &sshOutput{
		Session: s,
		tty:     tty,
	}
	environ := s.Environ()
	environ = append(environ, fmt.Sprintf("TERM=%s", sshPty.Term))
	e := &sshEnviron{
		environ: environ,
	}
	return termenv.NewOutput(o, termenv.WithEnvironment(e))
}

func main() {
	addr := ":3456"
	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath("ssh_example"),
		wish.WithMiddleware(
			func(sh ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					output := outputFromSession(s)
					pty, _, active := s.Pty()
					if !active {
						sh(s)
						return
					}
					w, _ := pty.Window.Width, pty.Window.Height

					renderer := lipgloss.NewRenderer(lipgloss.WithTermenvOutput(output),
						lipgloss.WithColorProfile(termenv.TrueColor))
					str := strings.Builder{}
					fmt.Fprintf(&str, "\n%s %s %s %s %s",
						renderer.NewStyle().SetString("bold").Bold(true),
						renderer.NewStyle().SetString("faint").Faint(true),
						renderer.NewStyle().SetString("italic").Italic(true),
						renderer.NewStyle().SetString("underline").Underline(true),
						renderer.NewStyle().SetString("crossout").Strikethrough(true),
					)

					fmt.Fprintf(&str, "\n%s %s %s %s %s %s %s",
						renderer.NewStyle().SetString("red").Foreground(lipgloss.Color("#E88388")),
						renderer.NewStyle().SetString("green").Foreground(lipgloss.Color("#A8CC8C")),
						renderer.NewStyle().SetString("yellow").Foreground(lipgloss.Color("#DBAB79")),
						renderer.NewStyle().SetString("blue").Foreground(lipgloss.Color("#71BEF2")),
						renderer.NewStyle().SetString("magenta").Foreground(lipgloss.Color("#D290E4")),
						renderer.NewStyle().SetString("cyan").Foreground(lipgloss.Color("#66C2CD")),
						renderer.NewStyle().SetString("gray").Foreground(lipgloss.Color("#B9BFCA")),
					)

					fmt.Fprintf(&str, "\n%s %s %s %s %s %s %s\n\n",
						renderer.NewStyle().SetString("red").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#E88388")),
						renderer.NewStyle().SetString("green").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#A8CC8C")),
						renderer.NewStyle().SetString("yellow").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#DBAB79")),
						renderer.NewStyle().SetString("blue").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#71BEF2")),
						renderer.NewStyle().SetString("magenta").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#D290E4")),
						renderer.NewStyle().SetString("cyan").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#66C2CD")),
						renderer.NewStyle().SetString("gray").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#B9BFCA")),
					)

					fmt.Fprintf(&str, "%s %t\n", renderer.NewStyle().SetString("Has dark background?").Bold(true), renderer.HasDarkBackground())
					fmt.Fprintln(&str)

					wish.WriteString(s, renderer.Place(w, lipgloss.Height(str.String()), lipgloss.Center, lipgloss.Center, str.String()))

					sh(s)
				}
			},
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on %s", addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
