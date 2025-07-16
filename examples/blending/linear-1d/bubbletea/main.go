// This example demonstrates how to use the colors.BlendLinear1D function to create
// beautiful color gradients in a Bubble Tea application.
package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/colors"
)

var gradients = []gradientData{
	{
		name: "Sunset",
		stops: []color.Color{
			lipgloss.Color("#FF6B6B"), // Coral
			lipgloss.Color("#FFB74D"), // Orange
			lipgloss.Color("#FFDFBA"), // Peach
		},
	},
	{
		name: "Ocean",
		stops: []color.Color{
			lipgloss.Color("#0077B6"), // Deep Blue
			lipgloss.Color("#48CAE4"), // Sky Blue
			lipgloss.Color("#ADE8F4"), // Light Blue
		},
	},
	{
		name: "Forest",
		stops: []color.Color{
			lipgloss.Color("#228B22"), // Forest Green
			lipgloss.Color("#90EE90"), // Light Green
			lipgloss.Color("#FFFFE0"), // Cream
		},
	},
	{
		name: "Purple Dream",
		stops: []color.Color{
			lipgloss.Color("#9370DB"), // Medium Purple
			lipgloss.Color("#DDA0DD"), // Plum
			lipgloss.Color("#FFB6C1"), // Light Pink
		},
	},
	{
		name: "Fire",
		stops: []color.Color{
			lipgloss.Color("#FF0000"), // Red
			lipgloss.Color("#FFA500"), // Orange
			lipgloss.Color("#FFFF00"), // Yellow
		},
	},
}

type gradientData struct {
	name  string
	stops []color.Color
}

// Style definitions.
type styles struct {
	// UI styles.
	title        lipgloss.Style
	gradientName lipgloss.Style
	info         lipgloss.Style
}

func newStyles(dark bool) (s *styles) {
	s = &styles{}

	lightDark := lipgloss.LightDark(dark)

	s.title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lightDark(lipgloss.Color("#2D3748"), lipgloss.Color("#E2E8F0"))).
		MarginBottom(1).
		Align(lipgloss.Center)

	s.gradientName = lipgloss.NewStyle().
		Bold(true).
		Foreground(lightDark(lipgloss.Color("#4A5568"), lipgloss.Color("#CBD5E0"))).
		PaddingRight(1)

	s.info = lipgloss.NewStyle().
		Foreground(lightDark(lipgloss.Color("#718096"), lipgloss.Color("#A0AEC0"))).
		Italic(true)

	return s
}

type model struct {
	width  int
	height int
	styles *styles
}

func (m model) Init() tea.Cmd {
	return tea.RequestBackgroundColor
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.BackgroundColorMsg:
		m.styles = newStyles(msg.IsDark())
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() (string, *tea.Cursor) {
	var maxTitleWidth int

	for _, gradient := range gradients {
		maxTitleWidth = max(maxTitleWidth, lipgloss.Width(m.styles.gradientName.Render(gradient.name)))
	}

	var content strings.Builder

	content.WriteString(m.styles.title.Render("Color Gradient Examples with BlendLinear1D"))
	content.WriteString("\n\n")

	var title string

	for _, gradient := range gradients {
		title = m.styles.gradientName.Width(maxTitleWidth).Render(gradient.name)
		content.WriteString(title)

		blendedColors := colors.BlendLinear1D(m.width-maxTitleWidth, gradient.stops...)

		for _, c := range blendedColors {
			content.WriteString(lipgloss.NewStyle().Background(c).Foreground(c).Render("█"))
		}

		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(m.styles.info.Render("Press Q to exit"))

	cursor := &tea.Cursor{}
	cursor.X = 0
	cursor.Y = 0

	return content.String(), cursor
}

func main() {
	_, err := tea.NewProgram(model{styles: newStyles(true)}, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Uh oh: %v", err)
		os.Exit(1)
	}
}
