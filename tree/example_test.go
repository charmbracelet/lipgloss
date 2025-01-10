package tree_test

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/charmbracelet/x/ansi"
)

// Leaf Examples

func ExampleLeaf_SetHidden() {
	tr := tree.New().
		Child(
			"Foo",
			tree.Root("Bar").
				Child(
					"Qux",
					tree.Root("Quux").
						Child("Hello!"),
					"Quuux",
				),
			"Baz",
		)

	tr.Children().At(1).Children().At(2).SetHidden(true)
	fmt.Println(tr.String())
	// Output:
	//
	// ├── Foo
	// ├── Bar
	// │   ├── Qux
	// │   └── Quux
	// │       └── Hello!
	// └── Baz
	//
}

func ExampleNewLeaf() {
	tr := tree.New().
		Child(
			"Foo",
			tree.Root("Bar").
				Child(
					"Qux",
					tree.Root("Quux").
						Child(
							tree.NewLeaf("This should be hidden", true),
							tree.NewLeaf(
								tree.Root("I am groot").Child("leaves"), false),
						),
					"Quuux",
				),
			"Baz",
		)

	fmt.Println(tr.String())
	// Output:
	// ├── Foo
	// ├── Bar
	// │   ├── Qux
	// │   ├── Quux
	// │   │   └── I am groot
	// │   │       └── leaves
	// │   └── Quuux
	// └── Baz
	//
}

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

// Tree Examples

func ExampleTree_Hide() {
	tr := tree.New().
		Child(
			"Foo",
			tree.Root("Bar").
				Child(
					"Qux",
					tree.Root("Quux").
						Child("Foo", "Bar").
						Hide(true),
					"Quuux",
				),
			"Baz",
		)

	fmt.Println(tr.String())
	// Output:
	// ├── Foo
	// ├── Bar
	// │   ├── Qux
	// │   └── Quuux
	// └── Baz
}

func ExampleTree_SetHidden() {
	tr := tree.New().
		Child(
			"Foo",
			tree.Root("Bar").
				Child(
					"Qux",
					tree.Root("Quux").
						Child("Foo", "Bar"),
					"Quuux",
				),
			"Baz",
		)

	// Hide a tree after its creation. We'll hide Quux.
	tr.Children().At(1).Children().At(1).SetHidden(true)
	// Output:
	// ├── Foo
	// ├── Bar
	// │   ├── Qux
	// │   └── Quuux
	// └── Baz
	//
	fmt.Println(tr.String())
}
