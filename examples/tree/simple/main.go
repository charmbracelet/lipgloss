package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	t := tree.New().
		Root(".").
		Item("Item 1").
		Item(
			tree.New().Root("Item 2").
				Item("Item 2.1").
				Item("Item 2.2").
				Item("Item 2.3"),
		).
		Item(
			tree.New().
				Root("Item 3").
				Item("Item 3.1").
				Item("Item 3.2"),
		)

	fmt.Println(t)
}
