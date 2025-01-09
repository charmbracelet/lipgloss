package tree_test

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
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

func ExampleSetValue() {
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
	glossier := t.Children().At(0)
	glossier.SetValue(tree.Root(glossier.Value()).Child(tree.Root("Apparel").Child("Pink Hoodie", "Baseball Cap")))
	fmt.Println(t.String())
	// Output:
	// hello
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
