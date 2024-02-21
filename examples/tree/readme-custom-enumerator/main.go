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
		tree.New("Nyx", "Foo", "Bar"),
		"Mac",
		"Milk",
	).Enumerator(func(atter tree.Data, i int, last bool) (indent string, prefix string) {
		return "->", "->"
	})
	fmt.Println(t)
}
