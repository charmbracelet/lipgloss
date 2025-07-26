// This example demonstrates how to use the colors.Lighten and colors.Darken functions
// to create progressive brightness variations in a standalone Lip Gloss application.
package main

import (
	"image/color"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/colors"
)

func main() {
	hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(hasDarkBG)

	// Base colors to demonstrate lightening and darkening.
	baseColors := map[string]color.Color{
		"Red":   lipgloss.Color("#FF0000"),
		"Blue":  lipgloss.Color("#0066FF"),
		"Green": lipgloss.Color("#00FF00"),
		"Gray":  lipgloss.Color("#808080"),
	}

	// Percentage to lighten/darken by.
	percentage := 0.05 // 5%

	// Number of steps to generate.
	steps := 20

	colorNameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lightDark(lipgloss.Color("#2D3748"), lipgloss.Color("#E2E8F0")))

	var content strings.Builder

	for name, baseColor := range baseColors {
		content.WriteString(colorNameStyle.Render(name))
		content.WriteString("\n")

		// Create lightened variations.
		var lightenedBox strings.Builder
		lightenedBox.WriteString("Lightened: ")
		for i := range steps {
			lightenedBox.WriteString(
				lipgloss.NewStyle().
					Foreground(colors.Lighten(baseColor, percentage*(float64(i)+1))).
					Render("██"),
			)
		}
		content.WriteString(lightenedBox.String())
		content.WriteString("\n")

		// Create darkened variations.
		var darkenedBox strings.Builder
		darkenedBox.WriteString("Darkened:  ")
		for i := range steps {
			darkenedBox.WriteString(
				lipgloss.NewStyle().
					Foreground(colors.Darken(baseColor, percentage*(float64(i)+1))).
					Render("██"),
			)
		}
		content.WriteString(darkenedBox.String())
		content.WriteString("\n\n")
	}

	lipgloss.Println(content.String())
}
