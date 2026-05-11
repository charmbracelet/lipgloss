package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
)

func main() {
	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	pink := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

	t := tree.New().
		Child(
			"Glossier",
			"Claire’s Boutique",
			tree.Root("Nyx").
				Child("Lip Gloss", "Foundation").
				EnumeratorStyle(pink).
				IndenterStyle(purple),
			"Mac",
			"Milk",
		).
		EnumeratorStyle(purple).
		IndenterStyle(purple)
	fmt.Println(t)
}
