package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	fmt.Println(
		tree.New(
			"Items",
			tree.New(
				"Vegetables",
				"Lettuce",
				"Cabbage",
			),
			tree.New(
				"Drinks",
				"Beer",
				"Wine",
				"Whiskey",
			),
			"Bread",
			tree.New(
				"Meats",
				"Beef",
				"Pork",
				"Chicken",
			),
		),
	)
}
