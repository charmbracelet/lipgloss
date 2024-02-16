package tree

import (
	"github.com/charmbracelet/lipgloss"
)

// Node is a node in a tree.
type Node interface {
	Name() string
	String() string
	Children() []Node
}

// Atter returns a child node in a specified index.
type Atter interface {
	At(i int) Node
}

type atterImpl []Node

func (a atterImpl) At(i int) Node {
	if i >= 0 && i < len(a) {
		return a[i]
	}
	return nil
}

// StringNode is a node without children.
type StringNode string

// Children conforms with Node.
// StringNodes have no children.
func (StringNode) Children() []Node { return nil }

// Name conforms with Node.
// Returns the value of the string itself.
func (s StringNode) Name() string { return string(s) }

func (s StringNode) String() string { return s.Name() }

// TreeNode implements the Node interface with String data.
type TreeNode struct { //nolint:revive
	name     string
	renderer *defaultRenderer
	addItem  ItemAddFunc
	children []Node
}

// Name returns the root name of this node.
func (n *TreeNode) Name() string { return n.name }

func (n *TreeNode) String() string {
	return n.renderer.Render(n, true, "")
}

// ItemAddFunc is a function that, given the current list of nodes and a new
// item, adds the new item to the list and returns it.
// This is to be used internally only.
type ItemAddFunc func(nodes []Node, item any) []Node

func treeItemAddFn(children []Node, item any) []Node {
	switch item := item.(type) {
	case Node:
		children = append(children, item)
	case string:
		s := StringNode(item)
		children = append(children, &s)
	}
	return children
}

// ItemAddFunc changes the way items are added to a tree.
func (n *TreeNode) ItemAddFunc(fn ItemAddFunc) *TreeNode {
	n.addItem = fn
	return n
}

// Item appends an item to a list.
func (n *TreeNode) Item(item any) *TreeNode {
	n.children = n.addItem(n.children, item)
	return n
}

// EnumeratorStyle implements Renderer.
func (n *TreeNode) EnumeratorStyle(style lipgloss.Style) *TreeNode {
	n.renderer.style.enumeratorFunc = func(Atter, int) lipgloss.Style { return style }
	return n
}

// EnumeratorStyleFunc implements Renderer.
func (n *TreeNode) EnumeratorStyleFunc(fn StyleFunc) *TreeNode {
	if fn == nil {
		fn = func(Atter, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	n.renderer.style.enumeratorFunc = fn
	return n
}

// ItemStyle implements Renderer.
func (n *TreeNode) ItemStyle(style lipgloss.Style) *TreeNode {
	n.renderer.style.itemFunc = func(Atter, int) lipgloss.Style { return style }
	return n
}

// ItemStyleFunc implements Renderer.
func (n *TreeNode) ItemStyleFunc(fn StyleFunc) *TreeNode {
	if fn == nil {
		fn = func(Atter, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	n.renderer.style.enumeratorFunc = fn
	return n
}

// Enumerator implements Renderer.
func (n *TreeNode) Enumerator(enum Enumerator) *TreeNode {
	n.renderer.enumerator = enum
	return n
}

// Children returns the children of a string node.
func (n *TreeNode) Children() []Node {
	return n.children
}

// New returns a new tree.
func New(root string, data ...any) *TreeNode {
	t := &TreeNode{
		name:     root,
		renderer: newDefaultRenderer(),
		addItem:  treeItemAddFn,
	}
	for _, d := range data {
		t = t.Item(d)
	}
	return t
}
