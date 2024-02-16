package tree

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/internal/golden"
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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
	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestTreeSubTreeWithCustomRenderer(t *testing.T) {
	tree := New(
		"The Root Node(tm)",
		New("Parent", "child 1", "child 2").
			ItemStyleFunc(func(atter Atter, i int) lipgloss.Style {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("240"))
			}).
			EnumeratorStyleFunc(func(_ Atter, i int) lipgloss.Style {
				color := "212"
				if i%2 == 0 {
					color = "99"
				}
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color(color)).
					MarginRight(1)
			}),
		"Baz",
	)

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestTreeStyleNilFuncs(t *testing.T) {
	tree := New(
		"Multiline",
		"Foo",
		"Baz",
	).ItemStyleFunc(nil).
		EnumeratorStyleFunc(nil)

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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
