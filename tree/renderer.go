package tree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StyleFunc allows the list to be styled per item.
type StyleFunc func(data Data, i int) lipgloss.Style

// Style is the styling applied to the list.
type Style struct {
	enumeratorFunc StyleFunc
	itemFunc       StyleFunc
}

// NewDefaultRenderer returns the default renderer with the given style and
// enumerator.
func newDefaultRenderer() *defaultRenderer {
	return &defaultRenderer{
		style: Style{
			enumeratorFunc: func(Data, int) lipgloss.Style {
				return lipgloss.NewStyle().MarginRight(1)
			},
			itemFunc: func(Data, int) lipgloss.Style {
				return lipgloss.NewStyle()
			},
		},
		enumerator: DefaultEnumerator,
	}
}

// defaultRenderer is the default renderer used by the tree.
type defaultRenderer struct {
	style      Style
	enumerator Enumerator
}

// Render conforms with the Renderer interface.
func (r *defaultRenderer) Render(node Node, root bool, prefix string) string {
	if node.Hidden() {
		return ""
	}
	var strs []string
	var maxLen int
	children := node.Children()
	enumerator := r.enumerator

	// print the root node name if its not empty.
	if name := node.Name(); name != "" && root {
		strs = append(strs, r.style.itemFunc(children, -1).Render(name))
	}

	for i := 0; i < children.Length(); i++ {
		_, prefix := enumerator(children, i)
		prefix = r.style.enumeratorFunc(children, i).Render(prefix)
		maxLen = max(lipgloss.Width(prefix), maxLen)
	}

	for i := 0; i < children.Length(); i++ {
		child := children.At(i)
		if child.Hidden() {
			continue
		}
		indent, nodePrefix := enumerator(children, i)
		enumStyle := r.style.enumeratorFunc(children, i)
		itemStyle := r.style.itemFunc(children, i)

		nodePrefix = enumStyle.Render(nodePrefix)
		if l := maxLen - lipgloss.Width(nodePrefix); l > 0 {
			nodePrefix = strings.Repeat(" ", l) + nodePrefix
		}

		for j, line := range strings.Split(child.Name(), "\n") {
			if j == 0 {
				strs = append(strs, lipgloss.JoinHorizontal(
					lipgloss.Left,
					prefix+nodePrefix,
					itemStyle.Render(line),
				))
				continue
			}
			strs = append(strs, lipgloss.JoinHorizontal(
				lipgloss.Left,
				prefix+enumStyle.Render(indent),
				itemStyle.Render(line),
			))
		}

		if children.Length() > 0 {
			// here we see if the child has a custom renderer, which means the
			// user set a custom enumerator, style, etc.
			// if it has one, we'll use it to render itself.
			// otherwise, we keep using the current renderer.
			renderer := r
			switch child := child.(type) {
			case *TreeNode:
				if child.renderer != nil {
					renderer = child.renderer
				}
			}
			if s := renderer.Render(
				child,
				false,
				prefix+enumStyle.Render(indent),
			); s != "" {
				strs = append(strs, s)
			}
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
