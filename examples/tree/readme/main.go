package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).MarginRight(1)

	t := tree.New(
		"",
		"Glossier",
		"Claire’s Boutique",
		tree.New("Nyx", "Foo", "Bar"),
		"Mac",
		"Milk",
	).
		Enumerator(func(tree.Atter, int, bool) (indent string, prefix string) {
			return "->", "->"
		}).
		EnumeratorStyle(enumeratorStyle).
		ItemStyle(itemStyle)
	fmt.Println(t)
}
