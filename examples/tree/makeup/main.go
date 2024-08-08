package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

	t := tree.
		Root("Makeup").
		Child(
			"Glossier",
			"Claire’s Boutique",
			"Nyx",
			"Mac",
			"Milk",
		).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		ItemStyle(itemStyle).
		RootStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")))

	fmt.Println(t)
}
