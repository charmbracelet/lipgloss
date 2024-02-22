package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	style1 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		MarginRight(1)
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		MarginRight(1)

	t := tree.New().
		Root("Items").
		Items(
			tree.New().
				Root("Vegetables").
				Items(
					"Lettuce",
					"Cabbage",
				),
			tree.New().
				Root("Drinks").
				Items(
					"Beer",
					"Wine",
					"Whiskey",
				),
			"Bread",
			tree.New().
				Root("Meats").
				Items(
					"Beef",
					"Pork",
					tree.New().
						Root("Birds").
						Items("Chicken", "Duck"),
				).EnumeratorStyle(style2),
			"Fruit",
		).EnumeratorStyle(style1)

	fmt.Println(t)
}
