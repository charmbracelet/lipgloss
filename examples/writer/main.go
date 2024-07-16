package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/adaptive"
)

func main() {
	w := lipgloss.Writer{Forward: os.Stdout, Profile: lipgloss.Ascii}
	o := adaptive.Default()
	o.Writer = &w
	adaptive.SetDefault(o)
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF88CC"))

	adaptive.Println(style.Render("Hello, ANSI!"))
}
