package tree

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/internal/require"
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

	expected := `
├── Foo
├── Bar
│   ├── Qux
│   ├── Quux
│   │   ├── Foo
│   │   └── Bar
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
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
	expected := `
The Root
├── Foo
├── Bar
│   ├── Qux
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
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

	expected := `
├── Bar
│   ├── Qux
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeAddTwoSubTreesWithoutName(t *testing.T) {
	tree := New(
		"",
		"bar",
		"foo",
		New(
			"",
			"Bar 11",
			"Bar 12",
			"Bar 13",
			"Bar 14",
			"Bar 15",
		),
		New(
			"",
			"Bar 21",
			"Bar 22",
			"Bar 23",
			"Bar 24",
			"Bar 25",
		),
		"Baz",
	)

	expected := `
├── bar
├── foo
│   ├── Bar 11
│   ├── Bar 12
│   ├── Bar 13
│   ├── Bar 14
│   ├── Bar 15
│   ├── Bar 21
│   ├── Bar 22
│   ├── Bar 23
│   ├── Bar 24
│   └── Bar 25
└── Baz
	`
	require.Equal(t, expected, tree.String())
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

	expected := `
├── Foo
└── Bar
    ├── Qux
    ├── Quux
    │   ├── Foo
    │   └── Bar
    └── Quuux
	`
	require.Equal(t, expected, tree.String())
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

	expected := `
├── Bar
│   ├── Qux
│   ├── Quux
│   │   └── Bar
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
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
	).
		ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("9"))).
		EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).MarginRight(1)).
		Enumerator(func(Atter, int, bool) (indent string, prefix string) {
			return "->", "->"
		})

	expected := `
-> Foo
-> Bar
-> -> Qux
-> -> Quux
-> -> -> Foo
-> -> -> Bar
-> -> Quuux
-> Baz
	`
	require.Equal(t, expected, tree.String())
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

	expected := `
Multiline
Root
Node
├── Foo
├── Bar
│   ├── Qux
│   │   Line 2
│   │   Line 3
│   │   Line 4
│   ├── Quux
│   │   ├── Foo
│   │   └── Bar
│   └── Quuux
└── Baz
    Line 2
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeSubTreeWithCustomRenderer(t *testing.T) {
	tree := New(
		"The Root Node(tm)",
		New("Parent", "child 1", "child 2").
			ItemStyleFunc(func(Atter, int) lipgloss.Style {
				return lipgloss.NewStyle().
					SetString("*")
			}).
			EnumeratorStyleFunc(func(_ Atter, i int) lipgloss.Style {
				return lipgloss.NewStyle().
					SetString("+").
					MarginRight(1)
			}),
		"Baz",
	)

	expected := `
The Root Node(tm)
├── Parent
│   + ├── * child 1
│   + └── * child 2
└── Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeMixedEnumeratorSize(t *testing.T) {
	tree := New(
		"The Root Node(tm)",
		"child 1",
		"child 2",
		"child 3",
		"child 4",
		"child 5",
	).Enumerator(func(_ Atter, i int, _ bool) (indent string, prefix string) {
		romans := map[int]string{
			1: "I",
			2: "II",
			3: "III",
			4: "IV",
			5: "V",
			6: "VI",
		}
		return "", romans[i+1]
	})

	expected := `
The Root Node(tm)
  I child 1
 II child 2
III child 3
 IV child 4
  V child 5
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeStyleNilFuncs(t *testing.T) {
	tree := New(
		"Multiline",
		"Foo",
		"Baz",
	).ItemStyleFunc(nil).
		EnumeratorStyleFunc(nil)

	expected := `
Multiline
├──Foo
└──Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeStyleAt(t *testing.T) {
	tree := New(
		"Multiline",
		"Foo",
		"Baz",
	).Enumerator(func(atter Atter, i int, _ bool) (indent string, prefix string) {
		if atter.At(i).Name() == "Foo" {
			return "", ">"
		}
		return "", "-"
	})

	expected := `
Multiline
> Foo
- Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestAtter(t *testing.T) {
	nodes := atterImpl([]Node{
		StringNode("foo"),
		StringNode("bar"),
	})
	t.Run("0", func(t *testing.T) {
		if s := nodes.At(0).String(); s != "foo" {
			t.Errorf("expected 'foo', got '%s'", s)
		}
	})
	t.Run("10", func(t *testing.T) {
		if n := nodes.At(10); n != nil {
			t.Errorf("expected nil, got '%s'", n)
		}
	})
	t.Run("-1", func(t *testing.T) {
		if n := nodes.At(10); n != nil {
			t.Errorf("expected nil, got '%s'", n)
		}
	})
}
