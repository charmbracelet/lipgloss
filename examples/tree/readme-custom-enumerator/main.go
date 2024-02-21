package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	t := tree.New(
		"",
		"Glossier",
		"Claire’s Boutique",
		tree.New(
			"Nyx",
			"Qux",
			"Quux",
		),
		"Mac",
		"Milk",
	).Enumerator(func(data tree.Data, i int, last bool) (indent string, prefix string) {
		return "->", "->"
	})
	fmt.Println(t)
}
