package tree_test

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/charmbracelet/x/ansi"
)

func ExampleTree_Replace() {
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
	t.Replace(0, t.Children().At(0).Child(
		tree.Root("Apparel").Child("Pink Hoodie", "Baseball Cap"),
	))

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

func ExampleTree_Insert() {
	// Styles are here in case we want to test that styles are properly inherited...
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
	// Adds a new Tree Node after Fenty Beauty
	t.Insert(2, tree.Root("Lancôme").Child("Juicy Tubes Lip Gloss", "Lash Idôle", "Teint Idôle Highlighter"))

	// Adds a new Tree Node in Fenty Beauty
	t.Replace(1, t.Children().At(1).Insert(0, "Blurring Skin Tint"))

	// Adds a new Tree Node to a Leaf (Mac)
	t.Replace(4, t.Children().At(4).Insert(0, "Glow Play Cushion Blush"))
	fmt.Println(ansi.Strip(t.String()))
	// Output:
	//⁜ Makeup
	//├── Glossier
	//├── Fenty Beauty
	//│   ├── Blurring Skin Tint
	//│   ├── Gloss Bomb Universal Lip Luminizer
	//│   ╰── Hot Cheeks Velour Blushlighter
	//├── Lancôme
	//│   ├── Juicy Tubes Lip Gloss
	//│   ├── Lash Idôle
	//│   ╰── Teint Idôle Highlighter
	//├── Nyx
	//├── Mac
	//│   ╰── Glow Play Cushion Blush
	//╰── Milk
}
