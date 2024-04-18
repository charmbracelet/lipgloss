package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func enumerator(data tree.Data, i int) (indent, prefix string) {
	if data.At(i).Hidden() {
		return "   ", "▶ "
	}
	return "  ", "▼ "
}

func normalEnumerator(data tree.Data, _ int) (indent, prefix string) {
	return "  ", "• "
}

func main() {
	pink := lipgloss.NewStyle().MarginRight(1)
	gray := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginRight(1)

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
				).EnumeratorStyle(gray).Enumerator(normalEnumerator),
			tree.New().
				Root("Items").
				Items(
					"Cat Food",
					"Nutella",
					"Powdered Sugar",
				).EnumeratorStyle(gray).Enumerator(normalEnumerator).Hide(true),
			tree.New().
				Root("Veggies").
				Items(
					"Leek",
					"Artichoke",
				).EnumeratorStyle(gray).Enumerator(normalEnumerator),
		).
		ItemStyle(pink).
		EnumeratorStyle(gray).
		Enumerator(enumerator)

	fmt.Println(t)
}
