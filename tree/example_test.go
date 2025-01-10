package tree_test

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
	"github.com/charmbracelet/x/ansi"
)

func ExampleLeaf_Insert() {
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
		Enumerator(tree.RoundedEnumerator)
	// Adds a new Tree Node to a Leaf (Mac).
	t.Replace(3, t.Children().At(3).Insert(0, "Glow Play Cushion Blush"))
	fmt.Println(ansi.Strip(t.String()))
	// Output:
	//⁜ Makeup
	//├── Glossier
	//├── Fenty Beauty
	//│   ├── Gloss Bomb Universal Lip Luminizer
	//│   ╰── Hot Cheeks Velour Blushlighter
	//├── Nyx
	//├── Mac
	//│   ╰── Glow Play Cushion Blush
	//╰── Milk
}

func ExampleLeaf_Replace() {
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
		Enumerator(tree.RoundedEnumerator)
	// Add Glow Play Cushion Blush to Mac Leaf.
	t.Replace(3, t.Children().At(3).Replace(0, "Glow Play Cushion Blush"))
	fmt.Println(ansi.Strip(t.String()))
	// Output:
	//⁜ Makeup
	//├── Glossier
	//├── Fenty Beauty
	//│   ├── Gloss Bomb Universal Lip Luminizer
	//│   ╰── Hot Cheeks Velour Blushlighter
	//├── Nyx
	//├── Mac
	//│   ╰── Glow Play Cushion Blush
	//╰── Milk
}

// Tree Examples

func ExampleTree_Replace() {
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
		Enumerator(tree.RoundedEnumerator)
	// Add a Tree as a Child of "Glossier". At this stage "Glossier" is a Leaf,
	// so we re-assign the value of "Glossier" in the "Makeup" Tree to its new
	// Tree value returned from Child().
	t.Replace(0, t.Children().At(0).Child(
		tree.Root("Apparel").Child("Pink Hoodie", "Baseball Cap"),
	))

	// Add a Leaf as a Child of "Glossier". At this stage "Glossier" is a Tree,
	// so we don't need to use [Tree.Replace] on the parent tree.
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
		Enumerator(tree.RoundedEnumerator)
	// Adds a new Tree Node after Fenty Beauty.
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
