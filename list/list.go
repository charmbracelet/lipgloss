// Package list provides an API to create printable list structures that can be
// included in any command line application. This package includes full support
// for nested lists. Here's how you do it:
//
//	l := list.New(
//		"A",
//		"B",
//		"C",
//		list.New(
//			"D",
//			"E",
//			"F",
//		).Enumerator(list.Roman),
//	)
//	fmt.Println(l)
//
// The list package provides built-in enumerator styles to help glamourize your
// lists. This package wraps the tree package with list-specific styling. Lists
// are fully customizable, so let your creativity flow.
package list

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

// Enumerator is the type of enumeration to use for the list styling.
// It is the prefix for the list.
type Enumerator func(items Items, i int) string

// List is a list of items.
type List struct {
	tree *tree.Tree
}

// Items is the interface that wraps the basic methods of a list model.
type Items tree.Children

// StyleFunc allows the list to be styled per item.
type StyleFunc tree.StyleFunc

// Returns true if this node is hidden.
func (n *List) Hidden() bool {
	return n.tree.Hidden()
}

// Hide sets whether or not to hide the tree.
// This is useful for collapsing or hiding sub-lists.
func (n *List) Hide(hide bool) *List {
	n.tree.Hide(hide)
	return n
}

// Offset sets the tree children offsets.
func (n *List) OffsetStart(offset int) *List {
	n.tree.OffsetStart(offset)
	return n
}

// Offset sets the tree children offsets.
func (n *List) OffsetEnd(offset int) *List {
	n.tree.OffsetEnd(offset)
	return n
}

// Name returns the root name of this node.
func (n *List) Name() string { return n.tree.Value() }

func (n *List) String() string {
	return n.tree.String()
}

// EnumeratorStyle sets the enumeration style.
// Margins and paddings should usually be set only in ItemStyle or ItemStyleFunc.
func (n *List) EnumeratorStyle(style lipgloss.Style) *List {
	n.tree.EnumeratorStyle(style)
	return n
}

// EnumeratorStyleFunc sets the enumeration style function. Use this function
// for conditional styling.
//
// Margins and paddings should usually be set only in ItemStyle/ItemStyleFunc.
//
//	l := list.New().
//		EnumeratorStyleFunc(func(_ tree.Data, i int) lipgloss.Style {
//		    if i == 1 {
//		        return lipgloss.NewStyle().Foreground(hightlightColor)
//		    }
//		    return lipgloss.NewStyle().Foreground(dimColor)
//		})
func (n *List) EnumeratorStyleFunc(fn StyleFunc) *List {
	n.tree.EnumeratorStyleFunc(tree.StyleFunc(fn))
	return n
}

// ItemStyle sets the item style.
//
//	l := tree.New("Duck", "Duck", "Duck", "Goose", "Duck").
//		ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(255)))
func (n *List) ItemStyle(style lipgloss.Style) *List {
	n.tree.ItemStyle(style)
	return n
}

// ItemStyleFunc sets the item style function. Use this for conditional styling.
// For example:
//
//	l := list.New().
//		ItemStyleFunc(func(_ tree.Data, i int) lipgloss.Style {
//			st := baseStyle.Copy()
//			if selectedIndex == i {
//				return st.Foreground(hightlightColor)
//			}
//			return st.Foreground(dimColor)
//		})
func (n *List) ItemStyleFunc(fn StyleFunc) *List {
	n.tree.ItemStyleFunc(tree.StyleFunc(fn))
	return n
}

// Item appends an item to a list. Lists support nesting.
//
//	l := list.New().
//	Item("Item 1").
//	Item(list.New("Item 1.1", "Item 1.2"))
func (n *List) Item(item any) *List {
	switch item := item.(type) {
	case *List:
		n.tree.Item(item.tree)
	default:
		n.tree.Item(item)
	}
	return n
}

// Items add multiple items to the tree.
func (n *List) Items(items ...any) *List {
	for _, item := range items {
		n.Item(item)
	}
	return n
}

// Enumerator sets the enumerator implementation. Lipgloss includes
// predefined enumerators including bullets, roman numerals, and more. For
// example, you can have a numbered list:
//
//	list.New().
//		Enumerator(Arabic)
func (n *List) Enumerator(enum Enumerator) *List {
	n.tree.Enumerator(func(data tree.Children, i int) (string, string) {
		return " ", enum(data, i)
	})
	return n
}

// New returns a new list with a bullet enumerator.
func New(items ...any) *List {
	l := &List{
		tree: tree.New(),
	}
	return l.Items(items...).Enumerator(Bullet)
}
