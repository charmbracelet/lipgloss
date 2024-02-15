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
	st := Style{
		ItemFunc: func(i int) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).MarginLeft(1)
		},
		PrefixFunc: func(i int) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		},
	}
	r := DefaultRenderer().
		Styles(st).Enumerator(func(i int, last bool) (branch string, prefix string) {
		return "-> ", "->"
	})
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
	).Renderer(r)

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
