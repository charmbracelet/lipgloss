package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	t := tree.New().
		Items(
			"Glossier",
			"Claire’s Boutique",
			tree.New().
				Root("Nyx").
				Items(
					"Qux",
					"Quux",
				),
			"Mac",
			"Milk",
		).
		Enumerator(func(tree.Data, int) (string, string) {
			return "->", "->"
		})
	fmt.Println(t)
}
