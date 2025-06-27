package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

var (
	frameColor      = compat.AdaptiveColor{Light: lipgloss.Color("#C5ADF9"), Dark: lipgloss.Color("#864EFF")}
	textColor       = compat.AdaptiveColor{Light: lipgloss.Color("#696969"), Dark: lipgloss.Color("#bdbdbd")}
	keywordColor    = compat.AdaptiveColor{Light: lipgloss.Color("#37CD96"), Dark: lipgloss.Color("#22C78A")}
	inactiveBgColor = compat.AdaptiveColor{Light: lipgloss.Color("#988F95"), Dark: lipgloss.Color("#978692")}
	inactiveFgColor = compat.AdaptiveColor{Light: lipgloss.Color("#FDFCE3"), Dark: lipgloss.Color("#FBFAE7")}
)

// Style definitions.
type styles struct {
	frame,
	paragraph,
	text,
	keyword,
	activeButton,
	inactiveButton lipgloss.Style
}

// Styles are initialized based on the background color of the terminal.
func newStyles() (s styles) {
	// Define some styles. adaptive.Color() can be used to choose the
	// appropriate light or dark color based on the detected background color.
	s.frame = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(frameColor).
		Padding(1, 3).
		Margin(1, 3)
	s.paragraph = lipgloss.NewStyle().
		Width(40).
		MarginBottom(1).
		Align(lipgloss.Center)
	s.text = lipgloss.NewStyle().
		Foreground(textColor)
	s.keyword = lipgloss.NewStyle().
		Foreground(keywordColor).
		Bold(true)

	s.activeButton = lipgloss.NewStyle().
		Padding(0, 3).
		Background(lipgloss.Color("#FF6AD2")).
		Foreground(lipgloss.Color("#FFFCC2"))
	s.inactiveButton = s.activeButton.
		Background(inactiveBgColor).
		Foreground(inactiveFgColor)
	return s
}

type model struct {
	styles  styles
	yes     bool
	chosen  bool
	aborted bool
}

func (m model) Init() tea.Cmd {
	m.yes = true
	m.styles = newStyles()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.aborted = true
			return m, tea.Quit
		case "enter":
			m.chosen = true
			return m, tea.Quit
		case "left", "right", "h", "l":
			m.yes = !m.yes
		case "y":
			m.yes = true
			m.chosen = true
			return m, tea.Quit
		case "n":
			m.yes = false
			m.chosen = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.chosen || m.aborted {
		// We're about to exit, so wipe the UI.
		return ""
	}

	var (
		s = m.styles
		y = "Yes"
		n = "No"
	)

	if m.yes {
		y = s.activeButton.Render(y)
		n = s.inactiveButton.Render(n)
	} else {
		y = s.inactiveButton.Render(y)
		n = s.activeButton.Render(n)
	}

	return s.frame.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			s.paragraph.Render(
				s.text.Render("Are you sure you want to eat that ")+
					s.keyword.Render("moderatly ripe")+
					s.text.Render(" banana?"),
			),
			y+"  "+n,
		),
	)
}

func main() {
	m, err := tea.NewProgram(model{}).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Uh oh: %v", err)
		os.Exit(1)
	}

	if m := m.(model); m.chosen {
		if m.yes {
			fmt.Println("Are you sure? It's not ripe yet.")
		} else {
			fmt.Println("Well, alright. It was probably good, though.")
		}
	}
}
