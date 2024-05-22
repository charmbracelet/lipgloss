package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	enumeratorStyle := lipgloss.NewStyle().Background(lipgloss.Color("240")).PaddingRight(1)
	itemStyle := lipgloss.NewStyle().Background(lipgloss.Color("99")).Bold(true).PaddingRight(1)

	t := tree.New().
		ItemStyle(itemStyle).
		EnumeratorStyle(enumeratorStyle).
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
