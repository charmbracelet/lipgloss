package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

var purchased = []string{
	"Bananas",
	"Barley",
	"Cashews",
	"Coconut Milk",
	"Dill",
	"Eggs",
	"Fish Cake",
	"Leeks",
	"Papaya",
}

func GroceryEnumerator(l *list.List, i int) string {
	for _, p := range purchased {
		if l.At(i) == p {
			return "✓"
		}
	}
	return " "
}

func newList(items ...any) *list.List {
	enumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		MarginRight(1)

	return list.New(items...).
		Enumerator(GroceryEnumerator).
		EnumeratorStyle(enumStyle)
}

func main() {
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	enumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).MarginRight(1)

	l := list.New(
		"A", newList("Artichoke"),
		"B", newList("Baking Flour", "Bananas", "Barley", "Bean Sprouts"),
		"C", newList("Cashew Apple", "Cashews", "Coconut Milk", "Curry Paste", "Currywurst"),
		"D", newList("Dill", "Dragonfruit", "Dried Shrimp"),
		"E", newList("Eggs"),
		"F", newList("Fish Cake", "Furikake"),
		"J", newList("Jicama"),
		"K", newList("Kohlrabi"),
		"L", newList("Leeks", "Lentils", "Licorice Root"),
	).
		ItemStyle(itemStyle).
		EnumeratorStyle(enumStyle)

	fmt.Println(l)
}
