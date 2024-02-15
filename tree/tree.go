package tree

import "github.com/charmbracelet/lipgloss"

// Node is a node in a tree.
type Node interface {
	Name() string
	String() string
	Children() []Node
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
	renderer Renderer
	children []Node
}

// Name returns the root name of this node.
func (n *TreeNode) Name() string { return n.name }

func (n *TreeNode) String() string {
	var strs []string
	if name := n.Name(); name != "" {
		strs = append(strs, name)
	}
	if n.renderer == nil {
		n.renderer = DefaultRenderer()
	}
	strs = append(strs, n.renderer.Render(n.children, ""))
	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}

// Item appends an item to a list.
func (n *TreeNode) Item(item any) *TreeNode {
	switch item := item.(type) {
	case Node:
		n.children = append(n.children, item)
	case string:
		s := StringNode(item)
		n.children = append(n.children, &s)
	}
	return n
}

// Renderer sets the rendering function for a string node / tree.
func (n *TreeNode) Renderer(renderer Renderer) *TreeNode {
	n.renderer = renderer
	return n
}

// Children returns the children of a string node.
func (n *TreeNode) Children() []Node {
	return n.children
}

// New returns a new tree.
func New(name string, data ...any) *TreeNode {
	t := &TreeNode{name: name}
	for _, d := range data {
		t.Item(d)
	}
	return t
}
