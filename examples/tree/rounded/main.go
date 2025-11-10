package main

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
)

func main() {
	itemStyle := lipgloss.NewStyle().MarginRight(1)
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginRight(1)

	t := tree.Root("Groceries").
		Child(
			tree.Root("Fruits").
				Child(
					"Blood Orange",
					"Papaya",
					"Dragonfruit",
					"Yuzu",
				),
			tree.Root("Items").
				Child(
					"Cat Food",
					"Nutella",
					"Powdered Sugar",
				),
			tree.Root("Veggies").
				Child(
					"Leek",
					"Artichoke",
				),
		).ItemStyle(itemStyle).EnumeratorStyle(enumeratorStyle).Enumerator(tree.RoundedEnumerator)

	lipgloss.Println(t)
}
