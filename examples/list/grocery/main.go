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
	return "•"
}

var dimEnumStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	MarginRight(1)

var highlightedEnumStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10")).
	MarginRight(1)

func EnumStyleFunc(l *list.List, i int) lipgloss.Style {
	for _, p := range purchased {
		if l.At(i) == p {
			return highlightedEnumStyle
		}
	}
	return dimEnumStyle
}

func main() {
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	l := list.New(
		"Artichoke",
		"Baking Flour", "Bananas", "Barley", "Bean Sprouts",
		"Cashew Apple", "Cashews", "Coconut Milk", "Curry Paste", "Currywurst",
		"Dill", "Dragonfruit", "Dried Shrimp",
		"Eggs",
		"Fish Cake", "Furikake",
		"Jicama",
		"Kohlrabi",
		"Leeks", "Lentils", "Licorice Root",
	).
		ItemStyle(itemStyle).
		Enumerator(GroceryEnumerator).
		EnumeratorStyleFunc(EnumStyleFunc)

	fmt.Println(l)
}
