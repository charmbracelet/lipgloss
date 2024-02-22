package tree

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// Node is a node in a tree.
type Node interface {
	Name() string
	String() string
	Children() Data
	Hidden() bool
}

// StringNode is a node without children.
type StringNode string

// Children conforms with Node.
// StringNodes have no children.
func (StringNode) Children() Data { return nodeData(nil) }

// Name conforms with Node.
// Returns the value of the string itself.
func (s StringNode) Name() string { return string(s) }

// Hidden conforms with Node.
// Always returns false.
func (s StringNode) Hidden() bool { return false }

func (s StringNode) String() string { return s.Name() }

// Tree implements the Node interface with String data.
type Tree struct { //nolint:revive
	name         string
	renderer     *defaultRenderer
	rendererOnce sync.Once
	children     Data
	hide         bool
	offset       [2]int
}

// Returns true if this node is hidden.
func (n *Tree) Hidden() bool {
	return n.hide
}

// Hide sets whether or not to hide the tree.
// This is useful for collapsing / hiding sub-tree.
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

func (n *Tree) String() string {
	return n.ensureRenderer().Render(n, true, "")
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
// ├── foo
// ├── bar
// │   └── zaz
// └── qux
func (n *Tree) Item(item any) *Tree {
	switch item := item.(type) {
	case *Tree:
		newItem, rm := ensureParent(n.children, item)
		if rm >= 0 {
			n.children = n.children.Remove(rm)
		}
		n.children = n.children.Append(newItem)
	case Data:
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
	}
	return n
}

// Items add multiple items to the tree.
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
func ensureParent(nodes Data, item *Tree) (*Tree, int) {
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

func (n *Tree) ensureRenderer() *defaultRenderer {
	n.rendererOnce.Do(func() {
		n.renderer = newDefaultRenderer()
	})
	return n.renderer
}

// EnumeratorStyle sets the enumeration style.
func (n *Tree) EnumeratorStyle(style lipgloss.Style) *Tree {
	n.ensureRenderer().style.enumeratorFunc = func(Data, int) lipgloss.Style { return style }
	return n
}

// EnumeratorStyleFunc sets the enumeration style function.
func (n *Tree) EnumeratorStyleFunc(fn StyleFunc) *Tree {
	if fn == nil {
		fn = func(Data, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	n.ensureRenderer().style.enumeratorFunc = fn
	return n
}

// ItemStyle sets the item style.
func (n *Tree) ItemStyle(style lipgloss.Style) *Tree {
	n.ensureRenderer().style.itemFunc = func(Data, int) lipgloss.Style { return style }
	return n
}

// ItemStyleFunc sets the item style function.
func (n *Tree) ItemStyleFunc(fn StyleFunc) *Tree {
	if fn == nil {
		fn = func(Data, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	n.ensureRenderer().style.itemFunc = fn
	return n
}

// Enumerator sets the enumerator implementation.
func (n *Tree) Enumerator(enum Enumerator) *Tree {
	n.ensureRenderer().enumerator = enum
	return n
}

// Children returns the children of a string node.
func (n *Tree) Children() Data {
	var data []Node
	for i := n.offset[0]; i < n.children.Length()-n.offset[1]; i++ {
		data = append(data, n.children.At(i))
	}
	return nodeData(data)
}

// Root sets the tree node root name.
func (n *Tree) Root(root string) *Tree {
	n.name = root
	return n
}

// New returns a new tree.
func New() *Tree {
	return &Tree{
		children: nodeData(nil),
	}
}
