package tree

// Node is a node in a tree.
type Node interface {
	Name() string
	String() string
	Children() []Node
	Atter
}

// Atter returns a child node in a specified index.
type Atter interface {
	At(i int) Node
}

type atterImpl []Node

func (a atterImpl) At(i int) Node { return a[i] }

// StringNode is a node without children.
type StringNode string

// Children conforms with Node.
// StringNodes have no children.
func (StringNode) Children() []Node { return nil }

// Name conforms with Node.
// Returns the value of the string itself.
func (s StringNode) Name() string { return string(s) }

// At conforms with Atter.
// Always returns nil.
func (s StringNode) At(int) Node { return nil }

func (s StringNode) String() string { return s.Name() }

// TreeNode implements the Node interface with String data.
type TreeNode struct { //nolint:revive
	name     string
	renderer Renderer
	children []Node
}

// Name returns the root name of this node.
func (n *TreeNode) Name() string { return n.name }

// At returns the child at the specific index, if the index is within bounds.
// It'll return nil otherwise.
func (n *TreeNode) At(i int) Node {
	if i >= 0 && i < len(n.children) {
		return n.children[i]
	}
	return nil
}

func (n *TreeNode) String() string {
	if n.renderer == nil {
		n.renderer = NewDefaultRenderer()
	}
	return n.renderer.Render(n, true, "")
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
func New(root string, data ...any) *TreeNode {
	t := &TreeNode{name: root}
	for _, d := range data {
		t = t.Item(d)
	}
	return t
}
