package tree

import (
	"strings"

	"github.com/charmbracelet/lipgloss/list"
)

// Node is a node in a tree.
type Node interface {
	Name() string
	String() string
	Children() []Node
}

// Enumerator is the function that allow customization of the indentation of
// the tree.
type Enumerator func(children []Node, prefix string) string

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
	name       string
	enumerator Enumerator
	children   []Node
}

// Name returns the root name of this node.
func (n TreeNode) Name() string { return n.name }

func (n TreeNode) String() string {
	var sb strings.Builder
	if n.Name() != "" {
		sb.WriteString(n.Name() + "\n")
	}
	sb.WriteString(n.enumerator(n.children, ""))
	return sb.String()
}

// Enumerator sets the indentation function for a string node / tree.
func (n *TreeNode) Enumerator(indentFunc Enumerator) *TreeNode {
	n.enumerator = indentFunc
	return n
}

// Children returns the children of a string node.
func (n TreeNode) Children() []Node {
	return n.children
}

func ListEnumerator(children []Node, prefix string) string {
	var sb strings.Builder
	for _, child := range children {
		for _, line := range strings.Split(child.Name(), "\n") {
			sb.WriteString(prefix + line + "\n")
		}

		if len(child.Children()) > 0 {
			sb.WriteString(DefaultEnumerator(child.Children(), "  "+prefix))
		}
	}
	return sb.String()
}

func DefaultEnumerator(children []Node, prefix string) string {
	var sb strings.Builder
	for i, child := range children {
		var treePrefix string
		var branch string

		if i == len(children)-1 {
			treePrefix = "└──"
			branch = "   "
		} else {
			treePrefix = "├──"
			branch = "│  "
		}

		for i, line := range strings.Split(child.Name(), "\n") {
			if i == 0 {
				sb.WriteString(prefix + treePrefix + " " + line + "\n")
				continue
			}
			sb.WriteString(prefix + branch + " " + line + "\n")
		}

		if len(child.Children()) > 0 {
			sb.WriteString(DefaultEnumerator(child.Children(), prefix+branch))
		}
	}
	return sb.String()
}

// New returns a new tree.
func New(name string, data ...any) *TreeNode {
	var children []Node
	for _, d := range data {
		switch d := d.(type) {
		case *TreeNode:
			children = append(children, d)
		case StringNode:
			children = append(children, d)
		case *StringNode:
			children = append(children, d)
		case *list.List:
			children = append(children, StringNode(d.Render()))
		case string:
			s := StringNode(d)
			children = append(children, &s)
		}
	}
	return &TreeNode{
		name:       name,
		children:   children,
		enumerator: DefaultEnumerator,
	}
}
