package tree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Renderer is the function that allow customization of the indentation of
// the tree.
type Renderer interface {
	Render(children []Node, prefix string) string
	Enumerator(enum Enumerator) Renderer
	Styles(style Style) Renderer
}

// DefaultRenderer returns the default renderer with the given style and
// enumerator.
func DefaultRenderer() Renderer {
	return &defaultRenderer{
		style:      DefaultStyles(),
		enumerator: DefaultEnumerator,
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r *defaultRenderer) Render(children []Node, prefix string) string {
	var strs []string
	var maxLen int

	for i := range children {
		_, prefix := r.enumerator(i, i == len(children)-1)
		prefix = r.style.PrefixFunc(i).Render(prefix)
		maxLen = max(lipgloss.Width(prefix), maxLen)
	}

	for i, child := range children {
		last := i == len(children)-1
		indent, nodePrefix := r.enumerator(i, last)

		nodePrefix = r.style.PrefixFunc(i).Render(nodePrefix)
		if l := maxLen - lipgloss.Width(nodePrefix); l > 0 {
			nodePrefix = strings.Repeat(" ", l) + nodePrefix
		}

		for i, line := range strings.Split(child.Name(), "\n") {
			if i == 0 {
				strs = append(strs, lipgloss.JoinHorizontal(
					lipgloss.Left,
					prefix+nodePrefix,
					r.style.ItemFunc(i).Render(line),
				))
				continue
			}
			strs = append(strs, lipgloss.JoinHorizontal(
				lipgloss.Left,
				prefix+indent,
				r.style.ItemFunc(i).Render(line),
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
			strs = append(strs, renderer.Render(child.Children(), prefix+indent))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
