package tree

import (
	"strings"
	"testing"
)

func TestTree(t *testing.T) {
	tree := New(
		"",
		"Foo",
		New(
			"Bar",
			"Qux",
			New(
				"Quux",
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

func TestTreeRoot(t *testing.T) {
	tree := New(
		"The Root",
		"Foo",
		New(
			"Bar",
			"Qux",
			"Quuux",
		),
		"Baz",
	)

	expected := strings.TrimPrefix(`
The Root
├── Foo
├── Bar
│  ├── Qux
│  └── Quuux
└── Baz
`, "\n")

	if tree.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, tree)
	}
}

func TestTreeStartsWithSubtree(t *testing.T) {
	tree := New(
		"",
		New(
			"Bar",
			"Qux",
			"Quuux",
		),
		"Baz",
	)

	expected := strings.TrimPrefix(`
├── Bar
│  ├── Qux
│  └── Quuux
└── Baz
`, "\n")

	if tree.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, tree)
	}
}

func TestTreeLastNodeIsSubTree(t *testing.T) {
	tree := New(
		"",
		"Foo",
		New(
			"Bar",
			"Qux",
			New(
				"Quux",
				"Foo",
				"Bar",
			),
			"Quuux",
		),
	)

	expected := strings.TrimPrefix(`
├── Foo
└── Bar
   ├── Qux
   ├── Quux
   │  ├── Foo
   │  └── Bar
   └── Quuux
`, "\n")

	if tree.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, tree)
	}
}

func TestTreeNil(t *testing.T) {
	tree := New(
		"",
		nil,
		New(
			"Bar",
			"Qux",
			New(
				"Quux",
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

func arrowIndent(children []Node, prefix string) string {
	var sb strings.Builder
	if prefix == "" {
		prefix = "-> "
	}
	for _, node := range children {
		sb.WriteString(prefix + node.Name() + "\n")
		if len(node.Children()) > 0 {
			sb.WriteString(arrowIndent(node.Children(), prefix+"-> "))
		}
	}
	return sb.String()
}

func TestTreeCustom(t *testing.T) {
	quuux := StringNode("Quuux")
	tree := New(
		"",
		"Foo",
		New(
			"Bar",
			StringNode("Qux"),
			New(
				"Quux",
				"Foo",
				"Bar",
			),
			&quuux,
		),
		"Baz",
	).Indent(arrowIndent)

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

func TestTreeMultilineNode(t *testing.T) {
	tree := New(
		"Multiline\nRoot\nNode",
		"Foo",
		New(
			"Bar",
			"Qux\nLine 2\nLine 3\nLine 4",
			New(
				"Quux",
				"Foo",
				"Bar",
			),
			"Quuux",
		),
		"Baz\nLine 2",
	)

	expected := strings.TrimPrefix(`
Multiline
Root
Node
├── Foo
├── Bar
│  ├── Qux
│  │   Line 2
│  │   Line 3
│  │   Line 4
│  ├── Quux
│  │  ├── Foo
│  │  └── Bar
│  └── Quuux
└── Baz
    Line 2
`, "\n")

	if tree.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, tree)
	}
}
