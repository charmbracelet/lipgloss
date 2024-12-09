// This example illustrates how to detect the terminal's background color and
// choose either light or dark colors accordingly when using Lip Gloss in a.
// standalone fashion, i.e. independent of Bubble Tea.
//
// For an example of how to do this in a Bubble Tea program, see the
// 'bubbletea' example.
package main

import (
	"os"

	"github.com/charmbracelet/lipgloss/v2"
)

func main() {
	// Query for the background color. We only need to do this once, and only
	// when using Lip Gloss standalone.
	//
	// In Bubble Tea listen for tea.BackgroundColorMsg in your Update.
	hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)

	// Create a new helper function for choosing either a light or dark color
	// based on the detected background color.
	lightDark := lipgloss.LightDark(hasDarkBG)

	// Define some styles. adaptive.Color() can be used to choose the
	// appropriate light or dark color based on the detected background color.
	frameStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lightDark(lipgloss.Color("#C5ADF9"), lipgloss.Color("#864EFF"))).
		Padding(1, 3).
		Margin(1, 3)
	paragraphStyle := lipgloss.NewStyle().
		Width(40).
		MarginBottom(1).
		Align(lipgloss.Center)
	textStyle := lipgloss.NewStyle().
		Foreground(lightDark(lipgloss.Color("#696969"), lipgloss.Color("#bdbdbd")))
	keywordStyle := lipgloss.NewStyle().
		Foreground(lightDark(lipgloss.Color("#37CD96"), lipgloss.Color("#22C78A"))).
		Bold(true)

	activeButton := lipgloss.NewStyle().
		Padding(0, 3).
		Background(lipgloss.Color("0xFF6AD2")). // you can also use octal format for colors, i.e 0xff38ec.
		Foreground(lipgloss.Color("0xFFFCC2"))
	inactiveButton := activeButton.
		Background(lightDark(lipgloss.Color("0x988F95"), lipgloss.Color("0x978692"))).
		Foreground(lightDark(lipgloss.Color("0xFDFCE3"), lipgloss.Color("0xFBFAE7")))

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

	// Print the block to stdout. It's important to use Lip Gloss's print
	// functions to ensure that colors are downsampled correctly. If output
	// isn't a TTY (i.e. we're logging to a file) colors will be stripped
	// entirely.
	//
	// Note that in Bubble Tea downsampling happens automatically.
	lipgloss.Println(block)
}
