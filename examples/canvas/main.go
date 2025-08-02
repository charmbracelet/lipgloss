package main

import (
	"image/color"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// newField fills a rectangular area with a given character in a given color.
func newField(rows, cols int, color color.Color) string {
	fieldSetyle := lipgloss.NewStyle().Foreground(color)
	fieldBuilder := strings.Builder{}
	for i := range rows {
		for range cols {
			fieldBuilder.WriteString("/")
		}
		if i < rows-1 {
			fieldBuilder.WriteString("\n")
		}
	}
	return fieldSetyle.Render(fieldBuilder.String())
}

// newCard creates a little card with rounded borders and a text label.
func newCard(darkMode bool, text string) string {
	lightDark := lipgloss.LightDark(darkMode)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Charple).
		Foreground(lightDark(charmtone.Iron, charmtone.Butter)).
		Height(9).
		Width(16).
		PaddingTop(3).
		Align(lipgloss.Center).
		Render(text)
}

func main() {
	darkMode := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(darkMode)

	// A few text blocks.
	lighterField := newField(17, 43, lightDark(charmtone.Smoke, charmtone.Charcoal))
	darkerField := newField(17, 43, lightDark(charmtone.Squid, charmtone.Pepper))

	// A few layers. Layers are created from strings (or blocks of text).
	pickles := lipgloss.NewLayer(newCard(darkMode, "Pickles"))
	melon := lipgloss.NewLayer(newCard(darkMode, "Bitter Melon"))
	sriracha := lipgloss.NewLayer(newCard(darkMode, "Sriracha"))

	// A canvas is simply a collection of layers.
	canvas := lipgloss.NewCanvas(
		// Layers can have X, Y, and Z offsets. By default, X, Y, and
		// Z are all 0.
		lipgloss.NewLayer(lighterField).X(5).Y(2),

		// Layers can be nested.
		lipgloss.NewLayer(darkerField).AddLayers(
			pickles.X(4).Y(2).Z(1), // the Z index places this layer above the others
			melon.X(22).Y(1),
			sriracha.X(11).Y(7),
		),
	)

	lipgloss.Println(canvas.Render())
}
