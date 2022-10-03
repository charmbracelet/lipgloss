package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
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
			bm.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				output := outputFromSession(s)
				renderer := lipgloss.NewRenderer(lipgloss.WithOutput(output))
				return model{
					detectedProfile: renderer.ColorProfile(),
					renderer:        renderer,
				}, nil
			}),
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
