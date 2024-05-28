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
	Name() string
	Children() Children
	Hidden() bool
}

// StringNode is a node without children and with a string describing it.
type StringNode string

// Children conforms with Node.
//
// Always returns no children.
func (StringNode) Children() Children { return NodeData(nil) }

// Name conforms with Node.
//
// Returns the value of the string itself.
func (s StringNode) Name() string { return string(s) }

// Hidden conforms with Node.
//
// Always returns false.
func (s StringNode) Hidden() bool { return false }

// String returns conforms with Stringer.
func (s StringNode) String() string { return s.Name() }

// Tree implements the Node interface.
// It has a name and, optionally, children.
type Tree struct { //nolint:revive
	name         string
	renderer     *renderer
	rendererOnce sync.Once
	children     Children
	hide         bool
	offset       [2]int
}

// Returns true if this node is hidden.
func (n *Tree) Hidden() bool {
	return n.hide
}

// Hide sets whether or not to hide the tree.
// This is useful for collapsing or hiding sub-trees.
func (n *Tree) Hide(hide bool) *Tree {
	n.hide = hide
	return n
}

// Offset sets the tree children offsets.
func (n *Tree) OffsetStart(offset int) *Tree {
	n.offset[0] = offset
	return n
}

// Offset sets the tree children offsets.
func (n *Tree) OffsetEnd(offset int) *Tree {
	n.offset[1] = offset
	return n
}

// Name returns the root name of this node.
func (n *Tree) Name() string { return n.name }

// String conforms with Stringer.
func (n *Tree) String() string {
	return n.ensureRenderer().render(n, true, "")
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
func (n *Tree) Item(item any) *Tree {
	switch item := item.(type) {
	case *Tree:
		newItem, rm := ensureParent(n.children, item)
		if rm >= 0 {
			n.children = n.children.Remove(rm)
		}
		n.children = n.children.Append(newItem)
	case Children:
		for i := 0; i < item.Length(); i++ {
			n.children = n.children.Append(item.At(i))
		}
	case Node:
		n.children = n.children.Append(item)
	case fmt.Stringer:
		s := StringNode(item.String())
		n.children = n.children.Append(&s)
	case string:
		s := StringNode(item)
		n.children = n.children.Append(&s)
	// XXX: passing []any and []string would be the most common errors it
	// seems, this fixes it, but doesn't handle other types of slices...
	// Maybe it's best to not handle any of these?
	case []any:
		return n.Items(item...)
	case []string:
		ss := make([]any, 0, len(item))
		for _, s := range item {
			ss = append(ss, s)
		}
		return n.Items(ss...)
	case nil:
		return n
	default:
		// optimistically try to convert to a string...
		return n.Item(fmt.Sprintf("%v", item))
	}
	return n
}

// Items add multiple items to the tree.
//
//	t := tree.New().
//		Root("Nyx").
//		Items("Qux", "Quux").
func (n *Tree) Items(items ...any) *Tree {
	for _, item := range items {
		n.Item(item)
	}
	return n
}

// Ensures the TreeItem being added is in good shape.
//
// If it has no name, and the current node list is empty, it will check the
// last item's of the list type:
// 1. IFF it's a TreeNode, it'll append item's children to it, and return it.
// 1. IFF it's a StringNode, it'll set its content as item's name, and remove it.
func ensureParent(nodes Children, item *Tree) (*Tree, int) {
	if item.Name() != "" || nodes.Length() == 0 {
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
		item.name = parent.Name()
		return item, j
	case *StringNode:
		item.name = parent.Name()
		return item, j
	}
	return item, -1
}

// Ensure the tree node has a renderer.
func (n *Tree) ensureRenderer() *renderer {
	n.rendererOnce.Do(func() {
		n.renderer = newRenderer()
	})
	return n.renderer
}

// EnumeratorStyle sets the enumeration style.
// Margins and paddings should usually be set only in ItemStyle or ItemStyleFunc.
//
//	t := tree.New("Duck", "Duck", "Duck", "Goose", "Duck").
//		EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#00d787")))
func (n *Tree) EnumeratorStyle(style lipgloss.Style) *Tree {
	n.ensureRenderer().style.enumeratorFunc = func(Children, int) lipgloss.Style { return style }
	return n
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
func (n *Tree) EnumeratorStyleFunc(fn StyleFunc) *Tree {
	if fn == nil {
		fn = func(Children, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	n.ensureRenderer().style.enumeratorFunc = fn
	return n
}

// ItemStyle sets the item style.
//
//	t := tree.New("Duck", "Duck", "Duck", "Goose", "Duck").
//		ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(255)))
func (n *Tree) ItemStyle(style lipgloss.Style) *Tree {
	n.ensureRenderer().style.itemFunc = func(Children, int) lipgloss.Style { return style }
	return n
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
func (n *Tree) ItemStyleFunc(fn StyleFunc) *Tree {
	if fn == nil {
		fn = func(Children, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	n.ensureRenderer().style.itemFunc = fn
	return n
}

// Enumerator sets the enumerator implementation. This can be used to change the way the branches indicators look.
// Lipgloss includes predefined enumerators including bullets, roman numerals, and more. For
// example, you can have a numbered list:
//
//	tree.New().
//		Enumerator(Arabic)
func (n *Tree) Enumerator(enum Enumerator) *Tree {
	n.ensureRenderer().enumerator = enum
	return n
}

// Children returns the children of a node.
func (n *Tree) Children() Children {
	var data []Node
	for i := n.offset[0]; i < n.children.Length()-n.offset[1]; i++ {
		data = append(data, n.children.At(i))
	}
	return NodeData(data)
}

// Root sets the tree node root name.
func (n *Tree) Root(root string) *Tree {
	n.name = root
	return n
}

// New returns a new tree.
func New() *Tree {
	return &Tree{
		children: NodeData(nil),
	}
}
