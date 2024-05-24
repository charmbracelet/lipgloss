package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	pink := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

	t := tree.New().
		Items(
			"Glossier",
			"Claire’s Boutique",
			tree.New().
				Root("Nyx").
				Items("Lip Gloss", "Foundation").
				EnumeratorStyle(pink),
			"Mac",
			"Milk",
		).
		EnumeratorStyle(purple)
	fmt.Println(t)
}
