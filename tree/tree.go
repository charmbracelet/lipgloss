package tree

import (
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

// TreeNode implements the Node interface with String data.
type TreeNode struct { //nolint:revive
	name         string
	renderer     *defaultRenderer
	rendererOnce sync.Once
	children     Data
	hide         bool
}

// Returns true if this node is hidden.
func (n *TreeNode) Hidden() bool {
	return n.hide
}

// Hide sets whether or not to hide the tree.
// This is useful for collapsing / hiding sub-tree.
func (n *TreeNode) Hide(hide bool) *TreeNode {
	n.hide = hide
	return n
}

// Name returns the root name of this node.
func (n *TreeNode) Name() string { return n.name }

func (n *TreeNode) String() string {
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
//	New("foo", "bar", New("", "zaz"))
//	New("foo", New("bar", "zaz"))
//
// The resulting tree would be:
// - foo
// - bar
//   - zaz
func (n *TreeNode) Item(item any) *TreeNode {
	switch item := item.(type) {
	case *TreeNode:
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
	case string:
		s := StringNode(item)
		n.children = n.children.Append(&s)
	}
	return n
}

// Ensures the TreeItem being added is in good shape.
//
// If it has no name, and the current node list is empty, it will check the
// last item's of the list type:
// 1. IFF it's a TreeNode, it'll append item's children to it, and return it.
// 1. IFF it's a StringNode, it'll set its content as item's name, and remove it.
func ensureParent(nodes Data, item *TreeNode) (*TreeNode, int) {
	if item.Name() != "" || nodes.Length() == 0 {
		return item, -1
	}
	j := nodes.Length() - 1
	parent := nodes.At(j)
	switch parent := parent.(type) {
	case *TreeNode:
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

func (n *TreeNode) ensureRenderer() *defaultRenderer {
	n.rendererOnce.Do(func() {
		n.renderer = newDefaultRenderer()
	})
	return n.renderer
}

// EnumeratorStyle implements Renderer.
func (n *TreeNode) EnumeratorStyle(style lipgloss.Style) *TreeNode {
	n.ensureRenderer().style.enumeratorFunc = func(Data, int) lipgloss.Style { return style }
	return n
}

// EnumeratorStyleFunc implements Renderer.
func (n *TreeNode) EnumeratorStyleFunc(fn StyleFunc) *TreeNode {
	if fn == nil {
		fn = func(Data, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	n.ensureRenderer().style.enumeratorFunc = fn
	return n
}

// ItemStyle implements Renderer.
func (n *TreeNode) ItemStyle(style lipgloss.Style) *TreeNode {
	n.ensureRenderer().style.itemFunc = func(Data, int) lipgloss.Style { return style }
	return n
}

// ItemStyleFunc implements Renderer.
func (n *TreeNode) ItemStyleFunc(fn StyleFunc) *TreeNode {
	if fn == nil {
		fn = func(Data, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	n.ensureRenderer().style.itemFunc = fn
	return n
}

// Enumerator implements Renderer.
func (n *TreeNode) Enumerator(enum Enumerator) *TreeNode {
	n.ensureRenderer().enumerator = enum
	return n
}

// Children returns the children of a string node.
func (n *TreeNode) Children() Data {
	return n.children
}

// New returns a new tree.
func New(root string, data ...any) *TreeNode {
	t := &TreeNode{
		name:     root,
		children: nodeData(nil),
	}
	for _, d := range data {
		t = t.Item(d)
	}
	return t
}
