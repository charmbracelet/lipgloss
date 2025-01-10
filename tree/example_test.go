package tree_test

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/charmbracelet/x/ansi"
)

// Leaf Examples

func ExampleNodeChildren_Replace() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	t := tree.
		Root("⁜ Makeup").
		Child(
			"Glossier",
			"Fenty Beauty",
			tree.New().Child(
				"Gloss Bomb Universal Lip Luminizer",
				"Hot Cheeks Velour Blushlighter",
			),
			"Nyx",
			"Mac",
			"Milk",
		).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyle(itemStyle)
	// Add a Tree as a Child of "Glossier"
	// This rewrites all of the tree data to do one replace. Not the most
	// efficient. Maybe we can improve on this.
	//
	// That is how we're handling any Child manipulation in the Child() func as
	// well. Because the children are an interface it's a bit trickier. We need
	// to do an assignment, can't just manipulate the children directly.
	t.SetChildren(t.Children().(tree.NodeChildren).
		Replace(0, t.Children().At(0).Child(
			tree.Root("Apparel").Child("Pink Hoodie", "Baseball Cap"),
		)))

	// Add a Leaf as a Child of "Glossier"
	t.Children().At(0).Child("Makeup")
	fmt.Println(ansi.Strip(t.String()))

	// Output:
	// ⁜ Makeup
	// ├── Glossier
	// │   ├── Apparel
	// │   │   ├── Pink Hoodie
	// │   │   ╰── Baseball Cap
	// │   ╰── Makeup
	// ├── Fenty Beauty
	// │   ├── Gloss Bomb Universal Lip Luminizer
	// │   ╰── Hot Cheeks Velour Blushlighter
	// ├── Nyx
	// ├── Mac
	// ╰── Milk
	//
}
