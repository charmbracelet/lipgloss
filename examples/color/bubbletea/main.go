package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
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
func newStyles(backgroundIsDark bool) (s *styles) {
	s = new(styles)

	// Create a new helper function for choosing either a light or dark color
	// based on the detected background color.
	lightDark := lipgloss.LightDark(backgroundIsDark)

	// Define some styles. adaptive.Color() can be used to choose the
	// appropriate light or dark color based on the detected background color.
	s.frame = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lightDark.Color("#0000ff", "#6200ff")).
		Padding(1, 3).
		Margin(1, 3)
	s.paragraph = lipgloss.NewStyle().
		Width(40).
		MarginBottom(1).
		Align(lipgloss.Center)
	s.text = lipgloss.NewStyle().
		Foreground(lightDark.Color("#0000ff", "#bdbdbd"))
	s.keyword = lipgloss.NewStyle().
		Foreground(lightDark.Color("#0000ff", "#04b87c")).
		Bold(true)

		// You can also use octal format for colors, i.e 0x#ff38ec.
	s.activeButton = lipgloss.NewStyle().
		Padding(0, 3).
		Background(lipgloss.Color(0xf347ff)).
		Foreground(lipgloss.Color(0xfaffcc))
	s.inactiveButton = s.activeButton.
		Background(lipgloss.Color(0x545454))
	return s
}

type model struct {
	styles *styles
	yes    bool
	chosen bool
}

func (m model) Init() (tea.Model, tea.Cmd) {
	// Query for the background color on start.
	return m, tea.BackgroundColor
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Bubble Tea automatically detects the background color on start. We
	// listen for the response here, then initialize our styles accordingly.
	case tea.BackgroundColorMsg:
		m.styles = newStyles(msg.IsDark())
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
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
	if m.styles == nil {
		// We haven't received tea.BackgroundColorMsg yet. Don't worry, it'll
		// be here in a flash.
		return ""
	}
	if m.chosen {
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
