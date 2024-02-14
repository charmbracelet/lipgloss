package tree

import (
	"strings"
)

// Node is a node in a tree.
type Node interface {
	Name() string
	String() string
	Children() []Node
}

// StringNode is a node without children.
type StringNode string

func (StringNode) Children() []Node { return nil }
func (s StringNode) Name() string   { return string(s) }
func (s StringNode) String() string { return s.Name() }

// IndentFunc is the function that allow customization of the indentation of
// the tree.
type IndentFunc func(children []Node, prefix string) string

// TreeNode implements the Node interface with String data.
type TreeNode struct {
	name       string
	indentFunc IndentFunc
	children   []Node
}

// Name returns the root name of this node.
func (n TreeNode) Name() string {
	return n.name
}

func (n TreeNode) String() string {
	var sb strings.Builder
	if n.Name() != "" {
		sb.WriteString(n.Name() + "\n")
	}
	sb.WriteString(n.indentFunc(n.children, ""))
	return sb.String()
}

// Indent sets the indentation function for a string node / tree.
func (n *TreeNode) Indent(indentFunc IndentFunc) *TreeNode {
	n.indentFunc = indentFunc
	return n
}

// Children returns the children of a string node.
func (n TreeNode) Children() []Node {
	return n.children
}

func defaultIndentFunc(children []Node, prefix string) string {
	var sb strings.Builder
	for i, info := range children {
		var treePrefix string
		var branch string

		if i == len(children)-1 {
			treePrefix = "└──"
			branch = "   "
		} else {
			treePrefix = "├──"
			branch = "│  "
		}

		sb.WriteString(prefix + treePrefix + " " + info.Name() + "\n")
		if len(info.Children()) > 0 {
			sb.WriteString(defaultIndentFunc(info.Children(), prefix+branch))
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
		case string:
			s := StringNode(d)
			children = append(children, &s)
		}
	}
	return &TreeNode{
		name:       name,
		children:   children,
		indentFunc: defaultIndentFunc,
	}
}
