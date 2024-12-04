package main

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
)

func main() {
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render

	t := table.New()
	t.Row("Bubble Tea", s("Milky"))
	t.Row("Milk Tea", s("Also milky"))
	t.Row("Actual milk", s("Milky as well"))
	lipgloss.Println(t.Render())
}
