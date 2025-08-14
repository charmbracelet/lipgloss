// This example demonstrates how to use the colors.Blend1D function to create
// beautiful color gradients in a standalone Lip Gloss application.
package main

import (
	"image/color"
	"os"
	"strings"

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
	hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(hasDarkBG)

	// Create styles.
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lightDark(lipgloss.Color("#2D3748"), lipgloss.Color("#E2E8F0"))).
		MarginBottom(1).
		Align(lipgloss.Center)

	gradientStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lightDark(lipgloss.Color("#718096"), lipgloss.Color("#A0AEC0")))

	var content strings.Builder

	content.WriteString(titleStyle.Render("Color Gradient Examples with Blend1D"))
	content.WriteString("\n")

	for _, gradient := range gradients {
		blendedColors := lipgloss.Blend1D(40, gradient...)

		var gradientBar strings.Builder
		for _, c := range blendedColors {
			blockStyle := lipgloss.NewStyle().Foreground(c)
			gradientBar.WriteString(blockStyle.Render("█"))
		}

		content.WriteString(gradientStyle.Render(gradientBar.String()))
		content.WriteString("\n")
	}

	lipgloss.Println(content.String())
}
