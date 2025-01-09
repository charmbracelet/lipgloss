package main

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/tree"
)

func main() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	t := tree.
		Root("⁜ Makeup").
		Child(
			"Glossier",
			"Fenty Beauty",
			tree.New().Child(
				"Gloss Bomb Universal Lip Luminizer",
				"Hot Cheeks Velour Blushlighter",
			),
			"Nyx",
			"Mac",
			"Milk",
		).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyle(itemStyle)

	lipgloss.Println(t)
}
