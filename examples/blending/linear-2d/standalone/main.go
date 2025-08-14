// This example demonstrates how to use the colors.Blend2D function to create
// beautiful 2D color gradients in a standalone Lip Gloss application.
package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

func main() {
	hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(hasDarkBG)

	gradients := []struct {
		name  string
		stops []color.Color
		angle float64
	}{
		{
			name: "Sunset Diagonal",
			stops: []color.Color{
				lipgloss.Color("#FF6B6B"), // Coral
				lipgloss.Color("#FFB74D"), // Orange
				lipgloss.Color("#FFDFBA"), // Peach
			},
			angle: 45,
		},
		{
			name: "Ocean Wave",
			stops: []color.Color{
				lipgloss.Color("#0077B6"), // Deep Blue
				lipgloss.Color("#48CAE4"), // Sky Blue
				lipgloss.Color("#ADE8F4"), // Light Blue
			},
			angle: 90,
		},
		{
			name: "Forest Mist",
			stops: []color.Color{
				lipgloss.Color("#228B22"), // Forest Green
				lipgloss.Color("#90EE90"), // Light Green
				lipgloss.Color("#FFFFE0"), // Cream
			},
			angle: 135,
		},
		{
			name: "Purple Dream",
			stops: []color.Color{
				lipgloss.Color("#9370DB"), // Medium Purple
				lipgloss.Color("#DDA0DD"), // Plum
				lipgloss.Color("#FFB6C1"), // Light Pink
			},
			angle: 180,
		},
		{
			name: "Fire Gradient",
			stops: []color.Color{
				lipgloss.Color("#FF0000"), // Red
				lipgloss.Color("#FFA500"), // Orange
				lipgloss.Color("#FFFF00"), // Yellow
			},
			angle: 225,
		},
	}

	// Create styles.
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lightDark(lipgloss.Color("#2D3748"), lipgloss.Color("#E2E8F0"))).
		MarginBottom(1).
		Align(lipgloss.Center)

	gradientStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lightDark(lipgloss.Color("#718096"), lipgloss.Color("#A0AEC0"))).
		MarginBottom(1)

	gradientNameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lightDark(lipgloss.Color("#4A5568"), lipgloss.Color("#CBD5E0"))).
		MarginBottom(1)

	var content strings.Builder

	content.WriteString(titleStyle.Render("2D Color Gradient Examples with Blend2D"))
	content.WriteString("\n\n")

	for _, gradient := range gradients {
		// Generate the gradient using Blend2D.
		width, height := 30, 12
		blendedColors := lipgloss.Blend2D(width, height, gradient.angle, gradient.stops...)

		// Create the gradient box using individual character styling.
		var gradientBox strings.Builder
		for y := range height { // Uses 1D row-major order.
			for x := range width {
				index := y*width + x
				gradientBox.WriteString(
					lipgloss.NewStyle().
						Foreground(blendedColors[index]).
						Render("█"),
				)
			}
			if y < height-1 { // End of row.
				gradientBox.WriteString("\n")
			}
		}

		content.WriteString(gradientNameStyle.Render(fmt.Sprintf("%s (Angle: %d°)", gradient.name, gradient.angle)))
		content.WriteString("\n")
		content.WriteString(gradientStyle.Render(gradientBox.String()))
		content.WriteString("\n")
	}

	lipgloss.Println(content.String())
}
