package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
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

func groceryEnumerator(atter tree.Atter, i int, _ bool) (string, string) {
	for _, p := range purchased {
		if atter.At(i).Name() == p {
			return "", "✓"
		}
	}
	return "", "•"
}

var dimEnumStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	MarginRight(1)

var highlightedEnumStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10")).
	MarginRight(1)

func enumStyleFunc(atter tree.Atter, i int) lipgloss.Style {
	for _, p := range purchased {
		if atter.At(i).Name() == p {
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
	).Renderer(
		list.DefaultRenderer().
			Enumerator(groceryEnumerator).
			EnumeratorStyleFunc(enumStyleFunc).
			ItemStyle(itemStyle),
	)

	fmt.Println(l)
}
