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

	t := tree.New(
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
			tree.New("Something", "foo", "bar"),
		).
			Renderer(
				tree.NewDefaultRenderer().
					EnumeratorStyle(style2),
			),
		"Foobar",
	).
		Renderer(
			tree.NewDefaultRenderer().
				EnumeratorStyle(style1),
		)

	fmt.Println(t)
}
