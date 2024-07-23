package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	t := tree.New().
		Root("~").
		Child(
			"Documents",
			"Downloads",
			"Unfinished Projects")

	fmt.Println("A classic tree:\n" + t.String() + "\n")

	tr := tree.New().
		Root("~").
		Enumerator(tree.RoundedEnumerator).
		Child(
			"Documents",
			"Downloads",
			"Unfinished Projects")

	fmt.Println("A cool, rounded tree:\n" + tr.String() + "\n")

	ti := tree.New().
		Root("~").
		Child(
			list.New().
				Root("Documents").
				Items(
					"Important Documents",
					"Junk Drawer",
					"Books",
				),
			"Downloads",
			"Unfinished Projects")

	fmt.Println("A fancy, nested tree:\n" + ti.String() + "\n")

	documents := tree.New().
		Root("Documents").
		Child(
			"Important Documents",
			"Junk Drawer",
			"Books")

	treeAsRoot := tree.New().
		Root(documents).
		Child(
			"More Documents",
			tree.New().Root("Unfinished Projects").Child("Bubble Tea in Rust", "Zig Projects"))
	fmt.Println("A chaotic tree:\n" + treeAsRoot.String() + "\n")
}
