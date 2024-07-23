package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/list"
)

func main() {
	t := list.New().
		Root("~").
		Items(
			"Documents",
			"Downloads",
			"Unfinished Projects").
		Enumerator(list.Tree)

	fmt.Println("A classic tree:\n" + t.String() + "\n")

	tr := list.New().
		Root("~").
		Items(
			"Documents",
			"Downloads",
			"Unfinished Projects").
		Enumerator(list.TreeRounded)

	fmt.Println("A cool, rounded tree:\n" + tr.String() + "\n")

	ti := list.New().
		Root("~").
		Items(
			list.New().
				Root("Documents").
				Items(
					"Important Documents",
					"Junk Drawer",
					"Books",
				).Enumerator(list.Tree),
			"Downloads",
			"Unfinished Projects").
		Enumerator(list.Tree).
		Indenter(list.TreeIndenter)

	fmt.Println("A fancy, nested tree:\n" + ti.String() + "\n")

	documents := list.New().
		Root("Documents").
		Items(
			"Important Documents",
			"Junk Drawer",
			"Books").
		Enumerator(list.Tree)

	treeAsRoot := list.New().
		Root(documents).
		Items(
			"More Documents",
			"Unfinished Projects",
			list.New("Bubble Tea in Rust", "Zig Projects").Enumerator(list.Tree)).
		Enumerator(list.Tree).
		Indenter(list.TreeIndenter)

	fmt.Println("A chaotic tree:\n" + treeAsRoot.String() + "\n")
}
