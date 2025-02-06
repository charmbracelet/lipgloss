package tree_test

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
	"github.com/charmbracelet/x/ansi"
)

// Leaf Examples

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

func ExampleLeaf_SetValue() {
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
	glossier := t.Children().At(0)
	glossier.SetValue("Il Makiage")
	fmt.Println(ansi.Strip(t.String()))
	// Output:
	//⁜ Makeup
	//├── Il Makiage
	//├── Fenty Beauty
	//│   ├── Gloss Bomb Universal Lip Luminizer
	//│   ╰── Hot Cheeks Velour Blushlighter
	//├── Nyx
	//├── Mac
	//╰── Milk
}

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
	//
}

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
