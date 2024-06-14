package main

import (
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// w := lipgloss.Writer{Forward: os.Stdout, Profile: lipgloss.ANSI}
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF88CC"))

	lipgloss.Println(style.Render("Hello, ANSI!"))
}
