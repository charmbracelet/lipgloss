// This example demonstrates how to use the colors.Blend2D function to create
// beautiful 2D color gradients in a Bubble Tea application.
package main

import (
	"cmp"
	"fmt"
	"image/color"
	"math"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
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
	angle            float64
	selectedGradient int

	// UI styles.
	infoStyle        lipgloss.Style
	controlsStyle    lipgloss.Style
	gradientBoxStyle lipgloss.Style
	gradients        []color.Color
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.updateGradient()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "a":
			m.angle = math.Mod(m.angle+15, 360)
			m.updateGradient()
		case "d":
			m.angle = math.Mod(m.angle-15+360, 360)
			m.updateGradient()
		case "left":
			m.boxWidth -= 2
			m.updateGradient()
		case "right":
			m.boxWidth += 2
			m.updateGradient()
		case "up":
			m.boxHeight--
			m.updateGradient()
		case "down":
			m.boxHeight++
			m.updateGradient()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.selectedGradient = max(0, min(int(msg.String()[0]-'1'), len(gradients)-1))
			m.updateGradient()
		}
	case tea.MouseClickMsg:
		switch msg.Mouse().Button {
		case tea.MouseLeft:
			m.boxWidth = msg.Mouse().X
			m.boxHeight = msg.Mouse().Y
			m.updateGradient()
		}
	}
	return m, nil
}

func (m *model) updateGradient() {
	m.boxWidth = clamp(m.boxWidth, 5, m.windowWidth-m.gradientBoxStyle.GetHorizontalFrameSize())
	m.boxHeight = clamp(m.boxHeight, 3, m.windowHeight-m.gradientBoxStyle.GetVerticalFrameSize()-m.infoStyle.GetVerticalFrameSize()-m.controlsStyle.GetVerticalFrameSize()-2)

	// Since gradients that might be large can take up more memory, only generate gradients when
	// the box size (potentially) changes. If you have much smaller gradients, this is less of
	// an issue.
	m.gradients = lipgloss.Blend2D(m.boxWidth, m.boxHeight, m.angle, gradients[m.selectedGradient]...)
}

func (m model) View() string {
	if len(m.gradients) == 0 || m.windowWidth == 0 || m.windowHeight == 0 {
		return "" // Wait until we generate the initial gradient/get window size.
	}

	// Build the gradient content.
	gradientContent := strings.Builder{}
	for y := range m.boxHeight { // Uses 1D row-major order.
		for x := range m.boxWidth {
			index := y*m.boxWidth + x
			gradientContent.WriteString(
				lipgloss.NewStyle().
					Background(m.gradients[index]).
					Render(" "),
			)
		}
		if y < m.boxHeight-1 { // End of row.
			gradientContent.WriteString("\n")
		}
	}

	gradient := m.gradientBoxStyle.Render(gradientContent.String())

	info := m.infoStyle.Width(m.windowWidth).Render(fmt.Sprintf(
		"Size: %dx%d | Angle: %.1f° | Colors: %d",
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
