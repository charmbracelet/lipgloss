package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	t := tree.Root(".").
		Child("Item 1").
		Child(
			tree.New().
				Root("Item 2").
				Child("Item 2.1").
				Child("Item 2.2").
				Child("Item 2.3"),
		).
		Child(
			tree.New().
				Root("Item 3").
				Child("Item 3.1").
				Child("Item 3.2"),
		)

	fmt.Println(t)
}
