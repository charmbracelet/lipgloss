package tree

import (
	"strings"
	"testing"
)

func TestTree(t *testing.T) {
	tree := New(
		"Foo",
		"Bar",
		New(
			"Qux",
			"Quux",
			New(
				"Foo",
				"Bar",
			),
			"Quuux",
		),
		"Baz",
	)

	expected := strings.TrimPrefix(`
├── Foo
├── Bar
│  ├── Qux
│  ├── Quux
│  │  ├── Foo
│  │  └── Bar
│  └── Quuux
└── Baz
`, "\n")

	if tree.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, tree)
	}
}

func TestTreeNil(t *testing.T) {
	tree := New(
		nil,
		"Bar",
		New(
			"Qux",
			"Quux",
			New(
				nil,
				"Bar",
			),
			"Quuux",
		),
		"Baz",
	)

	expected := strings.TrimPrefix(`
├── Bar
│  ├── Qux
│  ├── Quux
│  │  └── Bar
│  └── Quuux
└── Baz
`, "\n")

	if tree.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, tree)
	}
}

func TestTreeCustom(t *testing.T) {
	tree := New(
		"Foo",
		"Bar",
		New(
			"Qux",
			"Quux",
			New(
				"Foo",
				"Bar",
			),
			"Quuux",
		),
		"Baz",
	).Indent(func(t Node, level, index int, last bool) string {
		return strings.Repeat("-> ", level+1)
	})

	expected := strings.TrimPrefix(`
-> Foo
-> Bar
-> -> Qux
-> -> Quux
-> -> -> Foo
-> -> -> Bar
-> -> Quuux
-> Baz
`, "\n")

	if tree.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, tree)
	}
}
