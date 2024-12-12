package tree_test

import (
	"fmt"

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
