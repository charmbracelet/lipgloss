package list

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type (
	// Enumerator is the type of enumeration to use for the list styling.
	// It is the prefix for the list.
	Enumerator func(data Data, i int) string

	// List is a list of items.
	List struct {
		inner *tree.TreeNode
	}

	// Data is the interface that wraps the basic methods of a list model.
	Data tree.Data

	// StyleFunc allows the list to be styled per item.
	StyleFunc tree.StyleFunc
)

// Returns true if this node is hidden.
func (n *List) Hidden() bool {
	return n.inner.Hidden()
}

// Hide sets whether or not to hide the tree.
// This is useful for collapsing / hiding sub-tree.
func (n *List) Hide(hide bool) *List {
	n.inner.Hide(hide)
	return n
}

// Offset sets the tree children offsets.
func (n *List) OffsetStart(offset int) *List {
	n.inner.OffsetStart(offset)
	return n
}

// Offset sets the tree children offsets.
func (n *List) OffsetEnd(offset int) *List {
	n.inner.OffsetEnd(offset)
	return n
}

// Name returns the root name of this node.
func (n *List) Name() string { return n.inner.Name() }

func (n *List) String() string {
	return n.inner.String()
}

// EnumeratorStyle implements Renderer.
func (n *List) EnumeratorStyle(style lipgloss.Style) *List {
	n.inner.EnumeratorStyle(style)
	return n
}

// EnumeratorStyleFunc implements Renderer.
func (n *List) EnumeratorStyleFunc(fn StyleFunc) *List {
	n.inner.EnumeratorStyleFunc(tree.StyleFunc(fn))
	return n
}

// ItemStyle implements Renderer.
func (n *List) ItemStyle(style lipgloss.Style) *List {
	n.inner.ItemStyle(style)
	return n
}

// ItemStyleFunc implements Renderer.
func (n *List) ItemStyleFunc(fn StyleFunc) *List {
	n.inner.ItemStyleFunc(tree.StyleFunc(fn))
	return n
}

// Item appends an item to a list.
func (n *List) Item(item any) *List {
	switch item := item.(type) {
	case *List:
		n.inner.Item(item.inner)
	default:
		n.inner.Item(item)
	}
	return n
}

// Items add multiple items to the tree.
func (n *List) Items(items ...any) *List {
	n.inner.Items(items)
	return n
}

// Enumerator implements Renderer.
func (n *List) Enumerator(enum Enumerator) *List {
	n.inner.Enumerator(func(data tree.Data, i int) (string, string) {
		return " ", enum(data, i)
	})
	return n
}

// New returns a new list.
func New(items ...any) *List {
	l := &List{
		inner: tree.New().Items(items...),
	}
	return l.Enumerator(Bullet)
}
