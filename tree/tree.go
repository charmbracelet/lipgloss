package tree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Node is a node in a tree.
type Node interface {
	Name() string
	String() string
	Children() []Node
}

// StyleFunc allows the list to be styled per item.
type StyleFunc func(i int) lipgloss.Style

// Style is the styling applied to the list.
type Style struct {
	PrefixFunc StyleFunc
	ItemFunc   StyleFunc
}

// Renderer is the function that allow customization of the indentation of
// the tree.
type Renderer interface {
	Render(children []Node, prefix string) string
}

// Enumerator returns the branch and tree prefixes of a given item.
type Enumerator func(last bool) (branch string, prefix string)

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
	var sb strings.Builder
	if n.Name() != "" {
		sb.WriteString(n.Name() + "\n")
	}
	sb.WriteString(n.renderer.Render(n.children, ""))
	return sb.String()
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

// DefaultStyles is the default tree styles.
func DefaultStyles() Style {
	return Style{
		PrefixFunc: func(int) lipgloss.Style {
			return lipgloss.NewStyle()
		},
		ItemFunc: func(i int) lipgloss.Style {
			return lipgloss.NewStyle().MarginLeft(1)
		},
	}
}

type defaultRenderer struct {
	style      Style
	enumerator Enumerator
}

// NewDefaultRenderer returns the default renderer with the given style and
// enumerator.
func NewDefaultRenderer(style Style, enumerator Enumerator) Renderer {
	return &defaultRenderer{
		style:      style,
		enumerator: enumerator,
	}
}

// DefaultEnumerator enumerates items.
func DefaultEnumerator(last bool) (branch, prefix string) {
	if last {
		return "   ", "└──"
	}
	return "│  ", "├──"
}

func (r *defaultRenderer) Render(children []Node, prefix string) string {
	var sb strings.Builder
	for i, child := range children {
		last := i == len(children)-1
		branch, treePrefix := r.enumerator(last)

		for i, line := range strings.Split(child.Name(), "\n") {
			if i == 0 {
				sb.WriteString(
					r.style.PrefixFunc(i).Render(prefix+treePrefix) +
						r.style.ItemFunc(i).Render(line) +
						"\n",
				)
				continue
			}
			sb.WriteString(
				r.style.PrefixFunc(i).Render(prefix+branch) +
					r.style.ItemFunc(i).Render(line) +
					"\n",
			)
		}

		if len(child.Children()) > 0 {
			sb.WriteString(r.Render(child.Children(), prefix+branch))
		}
	}
	return sb.String()
}

// New returns a new tree.
func New(name string, data ...any) *TreeNode {
	var children []Node
	for _, d := range data {
		switch d := d.(type) {
		case Node:
			children = append(children, d)
		case string:
			s := StringNode(d)
			children = append(children, &s)
		}
	}
	return &TreeNode{
		name:     name,
		children: children,
		renderer: NewDefaultRenderer(DefaultStyles(), DefaultEnumerator),
	}
}
