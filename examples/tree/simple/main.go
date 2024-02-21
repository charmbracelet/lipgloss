package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	t := tree.New(
		"Root",
		"Item 1",
		"Item 2",
		tree.New(
			"Item 3",
			"Item 3.1",
		),
	)
	fmt.Println(t)
}
