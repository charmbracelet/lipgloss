// Package tree provides an API to create printable tree-like structures that
// can be included in any command line application. It goes something like:
//
//	t := tree.New().
//		Root(".").
//		Item("Item 1").
//		Item(
//		tree.New().Root("Item 2").
//			Item("Item 2.1").
//			Item("Item 2.2").
//			Item("Item 2.3"),
//		).
//		Item(
//		tree.New().
//			Root("Item 3").
//			Item("Item 3.1").
//			Item("Item 3.2"),
//		)
//
//	fmt.Println(t)
//
// If you're looking to create a list, you can use the list package which wraps
// the tree package with bulleted enumerations. Trees are fully customizable, so
// don't be shy, give 'em the 'ol razzle dazzle.
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

// StringNode is a node without children and with a string describing it.
type StringNode string

// Children conforms with Node.
//
// Always returns no children.
func (StringNode) Children() Children { return NodeChildren(nil) }

// Value conforms with Node.
//
// Returns the value of the string itself.
func (s StringNode) Value() string { return string(s) }

// Hidden conforms with Node.
//
// Always returns false.
func (s StringNode) Hidden() bool { return false }

// String returns conforms with Stringer.
func (s StringNode) String() string { return s.Value() }

// Tree implements the Node interface.
// It has a name and, optionally, children.
type Tree struct { //nolint:revive
	value        string
	renderer     *renderer
	rendererOnce sync.Once
	children     Children
	hide         bool
	offset       [2]int
}

// Returns true if this node is hidden.
func (t *Tree) Hidden() bool {
	return t.hide
}

// Hide sets whether or not to hide the tree.
// This is useful for collapsing or hiding sub-trees.
func (t *Tree) Hide(hide bool) *Tree {
	t.hide = hide
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

// String conforms with Stringer.
func (t *Tree) String() string {
	return t.ensureRenderer().render(t, true, "")
}

// Item appends an item to a list.
//
// If the tree being added is a new TreeNode without a name, we add its
// children to the previous string node.
//
// This is mostly syntactic sugar for adding items to lists.
//
// Both of these should result in the same thing:
//
//	New().Root("foo").Items("bar", New().Item("zaz"), "qux")
//	New().Root("foo").Items(New().Root("bar").Item("zaz"), "qux")
//
// The resulting tree would be:
//
//	├── foo
//	├── bar
//	│   └── zaz
//	└── qux
//
// You may also change the tree style using Enumerator.
func (t *Tree) Item(item any) *Tree {
	switch item := item.(type) {
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
		s := StringNode(item.String())
		t.children = t.children.(NodeChildren).Append(s)
	case string:
		s := StringNode(item)
		t.children = t.children.(NodeChildren).Append(&s)
	// XXX: passing []any and []string would be the most common errors it
	// seems, this fixes it, but doesn't handle other types of slices...
	// Maybe it's best to not handle any of these?
	case []any:
		return t.Items(item...)
	case []string:
		ss := make([]any, 0, len(item))
		for _, s := range item {
			ss = append(ss, s)
		}
		return t.Items(ss...)
	case nil:
		return t
	default:
		// optimistically try to convert to a string...
		return t.Item(fmt.Sprintf("%v", item))
	}
	return t
}

// Items add multiple items to the tree.
//
//	t := tree.New().
//		Root("Nyx").
//		Items("Qux", "Quux").
func (t *Tree) Items(items ...any) *Tree {
	for _, item := range items {
		t.Item(item)
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
			parent.Item(item.children.At(i))
		}
		return parent, j
	case StringNode:
		item.value = parent.Value()
		return item, j
	case *StringNode:
		item.value = parent.Value()
		return item, j
	}
	return item, -1
}

// Ensure the tree node has a renderer.
func (t *Tree) ensureRenderer() *renderer {
	t.rendererOnce.Do(func() {
		t.renderer = newRenderer()
	})
	return t.renderer
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
//		ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(255)))
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

// Children returns the children of a node.
func (t *Tree) Children() Children {
	var data []Node
	for i := t.offset[0]; i < t.children.Length()-t.offset[1]; i++ {
		data = append(data, t.children.At(i))
	}
	return NodeChildren(data)
}

// Root sets the tree node root name.
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
