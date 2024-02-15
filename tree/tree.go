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
	Enumerator(enum Enumerator) Renderer
	Styles(style Style) Renderer
}

// Enumerator returns the branch and tree prefixes of a given item.
type Enumerator func(i int, last bool) (indent string, prefix string)

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

// Enumerator implements Renderer.
func (r *defaultRenderer) Enumerator(enum Enumerator) Renderer {
	r.enumerator = enum
	return r
}

// Styles implements Renderer.
func (r *defaultRenderer) Styles(style Style) Renderer {
	r.style = style
	return r
}

// DefaultRenderer returns the default renderer with the given style and
// enumerator.
func DefaultRenderer() Renderer {
	return &defaultRenderer{
		style:      DefaultStyles(),
		enumerator: DefaultEnumerator,
	}
}

// DefaultEnumerator enumerates items.
func DefaultEnumerator(_ int, last bool) (indent, prefix string) {
	if last {
		return "   ", "└──"
	}
	return "│  ", "├──"
}

func (r *defaultRenderer) Render(children []Node, prefix string) string {
	var strs []string
	for i, child := range children {
		last := i == len(children)-1
		branch, treePrefix := r.enumerator(i, last)

		for i, line := range strings.Split(child.Name(), "\n") {
			if i == 0 {
				strs = append(strs, lipgloss.JoinHorizontal(
					lipgloss.Left,
					r.style.PrefixFunc(i).Render(prefix+treePrefix),
					r.style.ItemFunc(i).Render(line),
				))
				continue
			}
			strs = append(strs, lipgloss.JoinHorizontal(
				lipgloss.Left,
				r.style.PrefixFunc(i).Render(prefix+branch),
				r.style.ItemFunc(i).Render(line),
			))
		}

		if len(child.Children()) > 0 {
			switch child := child.(type) {
			case *TreeNode:
				strs = append(strs, child.renderer.Render(child.Children(), prefix+branch))
			default:
				strs = append(strs, r.Render(child.Children(), prefix+branch))
			}
		}
	}
	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}

// New returns a new tree.
func New(name string, data ...any) *TreeNode {
	t := &TreeNode{
		name:     name,
		renderer: DefaultRenderer(),
	}
	for _, d := range data {
		t = t.Item(d)
	}

	return t
}
