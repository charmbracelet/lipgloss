package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/wish"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/lucasb-eyer/go-colorful"
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

					renderer := lipgloss.NewRenderer(lipgloss.WithOutput(output))
					str := strings.Builder{}
					fmt.Fprintf(&str, "\n%s %s %s %s %s",
						renderer.NewStyle("bold").Bold(true),
						renderer.NewStyle("faint").Faint(true),
						renderer.NewStyle("italic").Italic(true),
						renderer.NewStyle("underline").Underline(true),
						renderer.NewStyle("crossout").Strikethrough(true),
					)

					fmt.Fprintf(&str, "\n%s %s %s %s %s %s %s",
						renderer.NewStyle("red").Foreground(lipgloss.Color("#E88388")),
						renderer.NewStyle("green").Foreground(lipgloss.Color("#A8CC8C")),
						renderer.NewStyle("yellow").Foreground(lipgloss.Color("#DBAB79")),
						renderer.NewStyle("blue").Foreground(lipgloss.Color("#71BEF2")),
						renderer.NewStyle("magenta").Foreground(lipgloss.Color("#D290E4")),
						renderer.NewStyle("cyan").Foreground(lipgloss.Color("#66C2CD")),
						renderer.NewStyle("gray").Foreground(lipgloss.Color("#B9BFCA")),
					)

					fmt.Fprintf(&str, "\n%s %s %s %s %s %s %s\n\n",
						renderer.NewStyle("red").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#E88388")),
						renderer.NewStyle("green").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#A8CC8C")),
						renderer.NewStyle("yellow").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#DBAB79")),
						renderer.NewStyle("blue").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#71BEF2")),
						renderer.NewStyle("magenta").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#D290E4")),
						renderer.NewStyle("cyan").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#66C2CD")),
						renderer.NewStyle("gray").Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#B9BFCA")),
					)

					fmt.Fprintf(&str, "%s %t\n", renderer.NewStyle("Has dark background?").Bold(true), renderer.HasDarkBackground())
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

type model struct {
	detectedProfile termenv.Profile
	renderer        *lipgloss.Renderer
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	m.renderer.SetColorProfile(m.detectedProfile)
	fmt.Fprintf(&s, "Detected Color Profile: %s\n", colorProfile(m.renderer.ColorProfile()))

	// Color grid
	fmt.Fprintf(&s, "%s\n", colors(m.renderer))

	// Enforced color profile
	m.renderer.SetColorProfile(termenv.TrueColor)
	fmt.Fprintf(&s, "Enforce True Color\n")

	// Color grid
	fmt.Fprintf(&s, "%s\n", colors(m.renderer))

	// Detect dark background
	fmt.Fprintf(&s, "Has Dark Background: %v\n", m.renderer.HasDarkBackground())

	return s.String()
}

func colors(r *lipgloss.Renderer) string {
	colors := colorGrid(14, 8)

	b := strings.Builder{}
	for _, x := range colors {
		for _, y := range x {
			s := r.NewStyle("  ").Background(lipgloss.Color(y))
			b.WriteString(s.String())
		}
		b.WriteRune('\n')
	}

	return b.String()
}

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}

func colorProfile(p termenv.Profile) string {
	switch p {
	case termenv.TrueColor:
		return "True Color"
	case termenv.ANSI256:
		return "ANSI 256"
	case termenv.ANSI:
		return "ANSI"
	case termenv.Ascii:
		return "No Color"
	default:
		return "Unknown"
	}
}
