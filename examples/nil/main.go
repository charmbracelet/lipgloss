package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	style := lipgloss.NewStyle().Render("Oops")
	fmt.Println(style)
}
