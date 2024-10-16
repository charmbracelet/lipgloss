// This example illustrates how to detect the terminal's background color and
// choose either light or dark colors accordingly when using Lip Gloss in a.
// standalone fashion, i.e. independent of Bubble Tea.
//
// For an example of how to do this in a Bubble Tea program, see the
// 'bubbletea' example.
package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/standalone"
)

func main() {
	// Query for the background color. We only need to do this once, and only
	// when using Lip Gloss standalone.
	//
	// In Bubble Tea listen for tea.BackgroundColorMsg in your Update.
	hasDarkBG, err := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not detect background color: %v\n", err)
		os.Exit(1)
	}

	// Create a new helper function for choosing either a light or dark color
	// based on the detected background color.
	adaptive := lipgloss.Adapt(hasDarkBG)

	// Define some styles. adaptive.Color() can be used to choose the
	// appropriate light or dark color based on the detected background color.
	frameStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(adaptive.Color("#0000ff", "#6200ff")).
		Padding(1, 3).
		Margin(1, 3)
	paragraphStyle := lipgloss.NewStyle().
		Width(40).
		MarginBottom(1).
		Align(lipgloss.Center)
	textStyle := lipgloss.NewStyle().
		Foreground(adaptive.Color("#0000ff", "#bdbdbd"))
	keywordStyle := lipgloss.NewStyle().
		Foreground(adaptive.Color("#0000ff", "#04b87c")).
		Bold(true)

		// You can also use octal format for colors, i.e 0x#ff38ec.
	activeButton := lipgloss.NewStyle().
		Padding(0, 3).
		Background(lipgloss.Color(0xf347ff)).
		Foreground(lipgloss.Color(0xfaffcc))
	inactiveButton := activeButton.
		Background(lipgloss.Color(0x545454))

	// Build layout.
	text := paragraphStyle.Render(
		textStyle.Render("Are you sure you want to eat that ") +
			keywordStyle.Render("moderatly ripe") +
			textStyle.Render(" banana?"),
	)
	buttons := activeButton.Render("Yes") + "  " + inactiveButton.Render("No")
	block := frameStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, text, buttons),
	)

	// Print the block to stdout. It's important to use the standalone package
	// to ensure that colors are downsampled correctly. If output isn't a
	// TTY (i.e. we're logging to a file) colors will be stripped entirely.
	standalone.Println(block)
}
