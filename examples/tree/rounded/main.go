package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	itemStyle := lipgloss.NewStyle().MarginRight(1)
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginRight(1)

	t := tree.New().
		Root("Groceries").
		Items(
			tree.New().
				Root("Fruits").
				Items(
					"Blood Orange",
					"Papaya",
					"Dragonfruit",
					"Yuzu",
				),
			tree.New().
				Root("Items").
				Items(
					"Cat Food",
					"Nutella",
					"Powdered Sugar",
				),
			tree.New().
				Root("Veggies").
				Items(
					"Leek",
					"Artichoke",
				),
		).ItemStyle(itemStyle).EnumeratorStyle(enumeratorStyle).Enumerator(tree.RoundedEnumerator)

	fmt.Println(t)
}
