package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render

	t := table.New()
	t.Row("Bubble Tea", s("Milky"))
	t.Row("Milk Tea", s("Also milky"))
	t.Row("Actual milk", s("Milky as well"))
	fmt.Println(t.Render())
}
