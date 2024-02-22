package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	t := tree.New().
		Root("Root").
		Items(
			"Item 1",
			"Item 2",
			tree.New().
				Root("Item 3").
				Item("Item 3.1"),
		)
	fmt.Println(t)
}
