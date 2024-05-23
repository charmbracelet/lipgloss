package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func directoryEnumerator(data tree.Data, i int) (indent, prefix string) {
	node := data.At(i)
	if node.Children().Length() == 0 {
		return " ", "▶"
	}
	return "  ", "▼"
}

func main() {
	t := tree.New().
		Enumerator(directoryEnumerator).
		Items("Pics").
		Items(
			tree.New().
				Root("Friends").
				Items(
					"README.md",
					"Ayman",
					tree.New().
						Enumerator(tree.RoundedEnumerator).
						Root("Carlos").Items(
						"Dabatoulle.png",
						"KEK.png",
						"Chefkiss.png",
						"Thinkies.png",
					),
					"Bash",
					"Carlos",
					"Maas",
					"Muesli",
				),
		)

	fmt.Println(t)
}
