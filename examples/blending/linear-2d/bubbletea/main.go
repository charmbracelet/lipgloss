// This example demonstrates how to use the colors.BlendLinear2D function to create
// beautiful 2D color gradients in a Bubble Tea application.
package main

import (
	"cmp"
	"fmt"
	"image/color"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/colors"
)

var gradients = [][]color.Color{
	{
		lipgloss.Color("#FF6B6B"), // Coral
		lipgloss.Color("#FFB74D"), // Orange
		lipgloss.Color("#FFDFBA"), // Peach
	},
	{
		lipgloss.Color("#0077B6"), // Deep Blue
		lipgloss.Color("#48CAE4"), // Sky Blue
		lipgloss.Color("#ADE8F4"), // Light Blue
	},
	{
		lipgloss.Color("#228B22"), // Forest Green
		lipgloss.Color("#90EE90"), // Light Green
		lipgloss.Color("#FFFFE0"), // Cream
	},
	{
		lipgloss.Color("#9370DB"), // Medium Purple
		lipgloss.Color("#DDA0DD"), // Plum
		lipgloss.Color("#FFB6C1"), // Light Pink
	},
	{
		lipgloss.Color("#9900FF"), // Purple
		lipgloss.Color("#00FA68"), // Lime
		lipgloss.Color("#ED5353"), // Red
	},
}

func main() {
	m := model{
		boxWidth:         20,
		boxHeight:        10,
		angle:            45,
		selectedGradient: 0,

		infoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginTop(1),
		controlsStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginTop(1),
		gradientBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")),
	}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	// UI state.
	windowWidth      int
	windowHeight     int
	boxWidth         int
	boxHeight        int
	angle            int
	selectedGradient int

	// UI styles.
	infoStyle        lipgloss.Style
	controlsStyle    lipgloss.Style
	gradientBoxStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.calcBoxSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "a":
			m.angle = (m.angle + 15) % 360
		case "d":
			m.angle = (m.angle - 15 + 360) % 360
		case "left":
			m.boxWidth -= 2
			m.calcBoxSize()
		case "right":
			m.boxWidth += 2
			m.calcBoxSize()
		case "up":
			m.boxHeight--
			m.calcBoxSize()
		case "down":
			m.boxHeight++
			m.calcBoxSize()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.selectedGradient = max(0, min(int(msg.String()[0]-'1'), len(gradients)-1))
		}
	case tea.MouseClickMsg:
		switch msg.Mouse().Button {
		case tea.MouseLeft:
			m.boxWidth = msg.Mouse().X
			m.boxHeight = msg.Mouse().Y
			m.calcBoxSize()
		}
	}
	return m, nil
}

func (m *model) calcBoxSize() {
	m.boxWidth = clamp(m.boxWidth, 5, m.windowWidth-m.gradientBoxStyle.GetHorizontalFrameSize())
	m.boxHeight = clamp(m.boxHeight, 3, m.windowHeight-m.gradientBoxStyle.GetVerticalFrameSize()-m.infoStyle.GetVerticalFrameSize()-m.controlsStyle.GetVerticalFrameSize()-2)
}

func (m model) View() string {
	gradientColors := colors.BlendLinear2D(m.boxWidth, m.boxHeight, m.angle, gradients[m.selectedGradient]...)

	// Build the gradient content.
	gradientContent := strings.Builder{}
	for y := range m.boxHeight {
		for x := range m.boxWidth {
			index := y*m.boxWidth + x
			gradientContent.WriteString(
				lipgloss.NewStyle().
					Background(gradientColors[index]).
					Render(" "),
			)
		}
		if y < m.boxHeight-1 { // End of row.
			gradientContent.WriteString("\n")
		}
	}

	gradient := m.gradientBoxStyle.Render(gradientContent.String())

	info := m.infoStyle.Width(m.windowWidth).Render(fmt.Sprintf(
		"Size: %dx%d | Angle: %d° | Colors: %d",
		m.boxWidth,
		m.boxHeight,
		m.angle,
		len(gradients[m.selectedGradient]),
	))

	controls := m.controlsStyle.Width(m.windowWidth).Render(fmt.Sprintf(
		"Controls: a/d (angle) | ←→ (width) | ↑↓ (height) | 1-%d (color scheme) | mouse click",
		len(gradients),
	))

	return lipgloss.NewStyle().
		Width(m.windowWidth).
		Height(m.windowHeight).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().
				Width(m.windowWidth).
				Height(m.windowHeight-lipgloss.Height(info)-lipgloss.Height(controls)).
				Render(gradient),
			info,
			controls,
		))
}

func clamp[T cmp.Ordered](v, low, high T) T {
	return min(high, max(low, v))
}
