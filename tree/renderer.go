package tree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Renderer is the function that allow customization of the indentation of
// the tree.
type Renderer interface {
	Render(node Node, root bool, prefix string) string
}

// StyleFunc allows the list to be styled per item.
type StyleFunc func(atter Atter, i int) lipgloss.Style

// Style is the styling applied to the list.
type Style struct {
	EnumeratorFunc StyleFunc
	ItemFunc       StyleFunc
}

// NewDefaultRenderer returns the default renderer with the given style and
// enumerator.
func NewDefaultRenderer() *DefaultRenderer {
	return &DefaultRenderer{
		style: Style{
			EnumeratorFunc: func(Atter, int) lipgloss.Style {
				return lipgloss.NewStyle().MarginRight(1)
			},
			ItemFunc: func(Atter, int) lipgloss.Style {
				return lipgloss.NewStyle()
			},
		},
		enumerator: DefaultEnumerator,
	}
}

// DefaultRenderer is the default renderer used by the tree.
type DefaultRenderer struct {
	style      Style
	enumerator Enumerator
}

// EnumeratorStyle implements Renderer.
func (r *DefaultRenderer) EnumeratorStyle(style lipgloss.Style) *DefaultRenderer {
	return r.EnumeratorStyleFunc(func(Atter, int) lipgloss.Style { return style })
}

// EnumeratorStyleFunc implements Renderer.
func (r *DefaultRenderer) EnumeratorStyleFunc(fn StyleFunc) *DefaultRenderer {
	if fn == nil {
		fn = func(Atter, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	r.style.EnumeratorFunc = fn
	return r
}

// ItemStyle implements Renderer.
func (r *DefaultRenderer) ItemStyle(style lipgloss.Style) *DefaultRenderer {
	return r.ItemStyleFunc(func(Atter, int) lipgloss.Style { return style })
}

// ItemStyleFunc implements Renderer.
func (r *DefaultRenderer) ItemStyleFunc(fn StyleFunc) *DefaultRenderer {
	if fn == nil {
		fn = func(Atter, int) lipgloss.Style { return lipgloss.NewStyle() }
	}
	r.style.EnumeratorFunc = fn
	return r
}

// Enumerator implements Renderer.
func (r *DefaultRenderer) Enumerator(enum Enumerator) *DefaultRenderer {
	r.enumerator = enum
	return r
}

// Render conforms with the Renderer interface.
func (r *DefaultRenderer) Render(node Node, root bool, prefix string) string {
	var strs []string
	var maxLen int
	children := node.Children()
	atter := atterImpl(children)

	// print the root node name if its not empty.
	if name := node.Name(); name != "" && root {
		strs = append(strs, r.style.ItemFunc(atter, -1).Render(name))
	}

	for i := range children {
		_, prefix := r.enumerator(atter, i, i == len(children)-1)
		prefix = r.style.EnumeratorFunc(atter, i).Render(prefix)
		maxLen = max(lipgloss.Width(prefix), maxLen)
	}

	for i, child := range children {
		last := i == len(children)-1
		indent, nodePrefix := r.enumerator(atter, i, last)

		nodePrefix = r.style.EnumeratorFunc(atter, i).Render(nodePrefix)
		if l := maxLen - lipgloss.Width(nodePrefix); l > 0 {
			nodePrefix = strings.Repeat(" ", l) + nodePrefix
		}

		for j, line := range strings.Split(child.Name(), "\n") {
			if j == 0 {
				strs = append(strs, lipgloss.JoinHorizontal(
					lipgloss.Left,
					prefix+nodePrefix,
					r.style.ItemFunc(atter, i).Render(line),
				))
				continue
			}
			strs = append(strs, lipgloss.JoinHorizontal(
				lipgloss.Left,
				prefix+indent,
				r.style.ItemFunc(atter, i).Render(line),
			))
		}

		if len(child.Children()) > 0 {
			var renderer Renderer = r
			switch child := child.(type) {
			case *TreeNode:
				if child.renderer != nil {
					renderer = child.renderer
				}
			}
			strs = append(strs, renderer.Render(child, false, prefix+indent))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
