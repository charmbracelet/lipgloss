package main

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	borderRotationFPS   = 15
	borderRotationSteps = 5
)

type borderRotationTickMsg struct {
	Value int
}

func borderRotationTick(current int) tea.Cmd {
	return tea.Tick(time.Second/time.Duration(borderRotationFPS), func(_ time.Time) tea.Msg {
		return borderRotationTickMsg{Value: current + borderRotationSteps}
	})
}

type model struct {
	borderRotation int
}

func (m model) Init() tea.Cmd {
	return borderRotationTick(0)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case borderRotationTickMsg:
		m.borderRotation = msg.Value
		return m, borderRotationTick(msg.Value)
	}

	return m, nil
}

func (m model) View() tea.View {
	v := tea.NewView(lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForegroundBlend(
			lipgloss.Color("#00FA68"),
			lipgloss.Color("#9900FF"),
			lipgloss.Color("#ED5353"),
			lipgloss.Color("#9900FF"),
			lipgloss.Color("#00FA68"),
		).
		BorderForegroundBlendOffset(m.borderRotation).
		Width(60).
		Height(15).
		Render("Hello, world!"))
	v.AltScreen = true
	return v
}

func main() {
	_, err := tea.NewProgram(model{}).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Uh oh: %v", err)
		os.Exit(1)
	}
}
