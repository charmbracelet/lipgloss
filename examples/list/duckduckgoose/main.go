package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func DuckDuckGooseEnumerator(l *list.List, i int) string {
	if l.At(i) == "Goose" {
		return "Honk →"
	}
	return " "
}

func main() {
	enumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00d787")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	baseStyle := lipgloss.NewStyle().Padding(1)

	l := list.New("Duck", "Duck", "Duck", "Goose", "Duck")
	l.Enumerator(DuckDuckGooseEnumerator)
	l.EnumeratorStyle(enumStyle)
	l.ItemStyle(itemStyle)
	l.BaseStyle(baseStyle)
	fmt.Println(l)
}
