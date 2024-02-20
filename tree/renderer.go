package tree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StyleFunc allows the list to be styled per item.
type StyleFunc func(atter Atter, i int) lipgloss.Style

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
			enumeratorFunc: func(Atter, int) lipgloss.Style {
				return lipgloss.NewStyle().MarginRight(1)
			},
			itemFunc: func(Atter, int) lipgloss.Style {
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
	var strs []string
	var maxLen int
	children := node.Children()
	atter := atterImpl(children)
	enumerator := r.enumerator

	// print the root node name if its not empty.
	if name := node.Name(); name != "" && root {
		strs = append(strs, r.style.itemFunc(atter, -1).Render(name))
	}

	for i := range children {
		_, prefix := enumerator(atter, i, i == len(children)-1)
		prefix = r.style.enumeratorFunc(atter, i).Render(prefix)
		maxLen = max(lipgloss.Width(prefix), maxLen)
	}

	for i, child := range children {
		last := i == len(children)-1
		indent, nodePrefix := enumerator(atter, i, last)
		enumStyle := r.style.enumeratorFunc(atter, i)
		itemStyle := r.style.itemFunc(atter, i)

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

		if len(child.Children()) > 0 {
			renderer := r
			switch child := child.(type) {
			case *TreeNode:
				if child.renderer != nil {
					renderer = child.renderer
				}
			}
			strs = append(
				strs,
				renderer.Render(
					child,
					false,
					prefix+enumStyle.Render(indent),
				),
			)
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
