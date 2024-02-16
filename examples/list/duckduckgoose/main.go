package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
)

func duckDuckGooseEnumerator(atter tree.Atter, i int, _ bool) (string, string) {
	if atter.At(i).Name() == "Goose" {
		return "", "Honk →"
	}
	return "", " "
}

func main() {
	enumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00d787")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	l := list.New("Duck", "Duck", "Duck", "Goose", "Duck").
		ItemStyle(itemStyle).
		EnumeratorStyle(enumStyle).
		Enumerator(duckDuckGooseEnumerator)
	fmt.Println(l)
}
