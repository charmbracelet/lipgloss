package tree

import "strings"

// Node is a node in a tree.
type Node interface {
	String() string
	Children() []Node
}

// IndentFunc is the function that allow customization of the indentation of
// the tree.
type IndentFunc func(n Node, level, index int, last bool) string

// StringNode implements the Node interface with String data.
type StringNode struct {
	indentFunc IndentFunc
	data       *string
	children   []Node
}

func (n StringNode) string(level, index int, last bool) string {
	var s strings.Builder

	if n.data != nil {
		if n.indentFunc != nil {
			s.WriteString(n.indentFunc(n, level, index, last))
		}
		s.WriteString(*n.data + "\n")
	}

	for i, child := range n.children {
		c := child.(*StringNode)
		c.indentFunc = n.indentFunc
		s.WriteString(c.string(level+1, i, i == len(n.children)-1))
	}

	return s.String()
}

func (n StringNode) String() string {
	return n.string(-1, 0, false)
}

// Indent sets the indentation function for a string node / tree.
func (n *StringNode) Indent(indentFunc IndentFunc) *StringNode {
	n.indentFunc = indentFunc
	return n
}

// Children returns the children of a string node.
func (n StringNode) Children() []Node {
	return n.children
}

func defaultIndentFunc(_ Node, level, _ int, last bool) string {
	var s strings.Builder

	if level >= 0 {
		s.WriteString(strings.Repeat("│  ", level))
	}

	if last {
		s.WriteString("└── ")
	} else {
		s.WriteString("├── ")
	}

	return s.String()
}

// New returns a new tree.
func New(data ...any) *StringNode {
	var children []Node

	for _, d := range data {
		switch d := d.(type) {
		case *StringNode:
			children = append(children, d)
		case string:
			children = append(children, &StringNode{data: &d, indentFunc: defaultIndentFunc})
		}
	}

	return &StringNode{
		children:   children,
		indentFunc: defaultIndentFunc,
	}
}
