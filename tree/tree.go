// Package tree allows you to build trees, as simple or complicated as you need.
//
// Define a tree with a root node, and children, set rendering properties (such
// as style, enumerators, etc...), and print it.
//
//	t := tree.New().
//		Items(
//			".git",
//			tree.New().
//				Root("examples/").
//				Items(
//					tree.New().
//						Root("list/").
//						Items("main.go").
//					tree.New().
//						Root("table/").
//						Items("main.go").
//				).
//			tree.New().
//				Root("list/").
//				Items("list.go", "list_test.go").
//			tree.New().
//				Root("table/").
//				Items("table.go", "table_test.go").
//			"align.go",
//			"align_test.go",
//			"join.go",
//			"join_test.go",
//		)
package tree

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// Node defines a node in a tree.
type Node interface {
	fmt.Stringer
	Value() string
	Children() Children
	Hidden() bool
}

// Leaf is a node without children.
type Leaf struct {
	value  string
	hidden bool
}

// Children of a Leaf node are always empty.
func (Leaf) Children() Children {
	return NodeChildren(nil)
}

// Value of a leaf node returns its value.
func (s Leaf) Value() string {
	return s.value
}

// Hidden returns whether a Leaf node is hidden.
func (s Leaf) Hidden() bool {
	return s.hidden
}

// String returns the string representation of a Leaf node.
func (s Leaf) String() string {
	return s.Value()
}

// Tree implements a Node.
type Tree struct { //nolint:revive
	value    string
	hidden   bool
	offset   [2]int
	children Children

	r     *renderer
	ronce sync.Once
}

// Hidden returns whether this node is hidden.
func (t *Tree) Hidden() bool {
	return t.hidden
}

// Hide sets whether to hide the tree node.
func (t *Tree) Hide(hide bool) *Tree {
	t.hidden = hide
	return t
}

// Offset sets the tree children offsets.
func (t *Tree) Offset(start, end int) *Tree {
	if start > end {
		_start := start
		start = end
		end = _start
	}

	if start < 0 {
		start = 0
	}
	if end < 0 || end > t.children.Length() {
		end = t.children.Length()
	}

	t.offset[0] = start
	t.offset[1] = end
	return t
}

// Value returns the root name of this node.
func (t *Tree) Value() string {
	return t.value
}

// String returns the string representation of the tree node.
func (t *Tree) String() string {
	return t.ensureRenderer().render(t, true, "")
}

// Child appends an item to a list.
//
// If the tree being added is a new TreeNode without a name, we add its
// children to the previous string node.
//
// This is mostly syntactic sugar for adding items to lists.
//
// Both of these should result in the same thing:
//
//	New().Root("foo").Items("bar", New().Child("zaz"), "qux")
//	New().Root("foo").Items(New().Root("bar").Child("zaz"), "qux")
//
// The resulting tree would be:
//
//	├── foo
//	├── bar
//	│   └── zaz
//	└── qux
//
// You may also change the tree style using Enumerator.
func (t *Tree) Child(children ...any) *Tree {
	for _, child := range children {
		switch item := child.(type) {
		case *Tree:
			newItem, rm := ensureParent(t.children, item)
			if rm >= 0 {
				t.children = t.children.(NodeChildren).Remove(rm)
			}
			t.children = t.children.(NodeChildren).Append(newItem)
		case Children:
			for i := 0; i < item.Length(); i++ {
				t.children = t.children.(NodeChildren).Append(item.At(i))
			}
		case Node:
			t.children = t.children.(NodeChildren).Append(item)
		case fmt.Stringer:
			s := Leaf{value: item.String()}
			t.children = t.children.(NodeChildren).Append(s)
		case string:
			s := Leaf{value: item}
			t.children = t.children.(NodeChildren).Append(&s)
		// XXX: passing []any and []string would be the most common errors it
		// seems, this fixes it, but doesn't handle other types of slices...
		// Maybe it's best to not handle any of these?
		case []any:
			return t.Child(item...)
		case []string:
			ss := make([]any, 0, len(item))
			for _, s := range item {
				ss = append(ss, s)
			}
			return t.Child(ss...)
		case nil:
			continue
		default:
			// optimistically try to convert to a string...
			return t.Child(fmt.Sprintf("%v", item))
		}
	}
	return t
}

// Ensures the TreeItem being added is in good shape.
//
// If it has no name, and the current node list is empty, it will check the
// last item's of the list type:
// 1. IFF it's a TreeNode, it'll append item's children to it, and return it.
// 1. IFF it's a StringNode, it'll set its content as item's name, and remove it.
func ensureParent(nodes Children, item *Tree) (*Tree, int) {
	if item.Value() != "" || nodes.Length() == 0 {
		return item, -1
	}
	j := nodes.Length() - 1
	parent := nodes.At(j)
	switch parent := parent.(type) {
	case *Tree:
		for i := 0; i < item.Children().Length(); i++ {
			parent.Child(item.children.At(i))
		}
		return parent, j
	case Leaf:
		item.value = parent.Value()
		return item, j
	case *Leaf:
		item.value = parent.Value()
		return item, j
	}
	return item, -1
}

// Ensure the tree node has a renderer.
func (t *Tree) ensureRenderer() *renderer {
	t.ronce.Do(func() {
		t.r = newRenderer()
	})
	return t.r
}

// EnumeratorStyle sets the enumeration style.
// Margins and paddings should usually be set only in ItemStyle or ItemStyleFunc.
//
//	t := tree.New("Duck", "Duck", "Duck", "Goose", "Duck").
//		EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00d787")))
func (t *Tree) EnumeratorStyle(style lipgloss.Style) *Tree {
	t.ensureRenderer().style.enumeratorFunc = func(Children, int) lipgloss.Style {
		return style
	}
	return t
}

// EnumeratorStyleFunc sets the enumeration style function. Use this function
// for conditional styling.
//
// Margins and paddings should usually be set only in ItemStyle/ItemStyleFunc.
//
//	t := tree.New().
//		EnumeratorStyleFunc(func(_ tree.Data, i int) lipgloss.Style {
//		    if i == 1 {
//		        return lipgloss.NewStyle().Foreground(hightlightColor)
//		    }
//		    return lipgloss.NewStyle().Foreground(dimColor)
//		})
func (t *Tree) EnumeratorStyleFunc(fn StyleFunc) *Tree {
	if fn == nil {
		fn = func(Children, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	t.ensureRenderer().style.enumeratorFunc = fn
	return t
}

// ItemStyle sets the item style.
//
//	t := tree.New("Duck", "Duck", "Duck", "Goose", "Duck").
//		ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("255")))
func (t *Tree) ItemStyle(style lipgloss.Style) *Tree {
	t.ensureRenderer().style.itemFunc = func(Children, int) lipgloss.Style { return style }
	return t
}

// ItemStyleFunc sets the item style function. Use this for conditional styling.
// For example:
//
//	t := tree.New().
//		ItemStyleFunc(func(_ tree.Data, i int) lipgloss.Style {
//			st := baseStyle.Copy()
//			if selectedIndex == i {
//				return st.Foreground(hightlightColor)
//			}
//			return st.Foreground(dimColor)
//		})
func (t *Tree) ItemStyleFunc(fn StyleFunc) *Tree {
	if fn == nil {
		fn = func(Children, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	t.ensureRenderer().style.itemFunc = fn
	return t
}

// Enumerator sets the enumerator implementation. This can be used to change the way the branches indicators look.
// Lipgloss includes predefined enumerators including bullets, roman numerals, and more. For
// example, you can have a numbered list:
//
//	tree.New().
//		Enumerator(Arabic)
func (t *Tree) Enumerator(enum Enumerator) *Tree {
	t.ensureRenderer().enumerator = enum
	return t
}

// Indenter sets the indenter implementation.
func (t *Tree) Indenter(indenter Indenter) *Tree {
	t.ensureRenderer().indenter = indenter
	return t
}

// Children returns the children of a node.
func (t *Tree) Children() Children {
	var data []Node
	for i := t.offset[0]; i < t.children.Length()-t.offset[1]; i++ {
		data = append(data, t.children.At(i))
	}
	return NodeChildren(data)
}

// Root returns a new tree with the root set.
//
//	tree.Root(root)
//
// It is a shorthand for:
//
//	tree.New().Root(root)
func Root(root string) *Tree {
	return New().Root(root)
}

// Root sets the root value of this tree.
func (t *Tree) Root(root string) *Tree {
	t.value = root
	return t
}

// New returns a new tree.
func New() *Tree {
	return &Tree{
		children: NodeChildren(nil),
	}
}
