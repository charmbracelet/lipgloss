package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/list"
)

func main() {
	t := tree().Items(
		".",
		"Item 1",
		"Item 2",
		"Item 2.1",
		"Item 2.2",
		"Item 2.3",
		tree().Items(
			"Item 3",
			"Item 3.1",
			"Item 3.2",
		),
	)

	fmt.Printf(".\n%s\n", t)
}

func tree() *list.List {
	return list.New().
		Indenter(list.TreeIndenter).
		Enumerator(list.Tree)
}
