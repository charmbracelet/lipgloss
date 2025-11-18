package tree

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// StyleFunc allows the tree to be styled per item.
type StyleFunc func(children Children, i int) lipgloss.Style

// Style is the styling applied to the tree.
type Style struct {
	enumeratorFunc StyleFunc
	indenterFunc   StyleFunc
	itemFunc       StyleFunc
	root           lipgloss.Style
}

// newRenderer returns the renderer used to render a tree.
func newRenderer() *renderer {
	return &renderer{
		style: Style{
			enumeratorFunc: func(Children, int) lipgloss.Style {
				return lipgloss.NewStyle().PaddingRight(1)
			},
			indenterFunc: func(Children, int) lipgloss.Style {
				return lipgloss.NewStyle().PaddingRight(1)
			},
			itemFunc: func(Children, int) lipgloss.Style {
				return lipgloss.NewStyle()
			},
		},
		enumerator: DefaultEnumerator,
		indenter:   DefaultIndenter,
	}
}

type renderer struct {
	style      Style
	enumerator Enumerator
	indenter   Indenter
	width      int
}

// render is responsible for actually rendering the tree.
func (r *renderer) render(node Node, root bool, prefix string) string {
	if node.Hidden() {
		return ""
	}

	var maxLen int
	children := node.Children()
	enumerator := r.enumerator
	indenter := r.indenter

	strs := make([]string, 0, children.Length())

	// print the root node name if its not empty.
	if name := node.Value(); name != "" && root {
		strs = append(strs, r.style.root.Render(name))
	}

	for i := range children.Length() {
		if i < children.Length()-1 {
			if child := children.At(i + 1); child.Hidden() {
				// Don't count the last child if its hidden. This renders the
				// last visible element with the right prefix
				//
				// The only type of Children is NodeChildren.
				children = children.(NodeChildren).Remove(i + 1)
			}
		}
		prefix := enumerator(children, i)
		prefix = r.style.enumeratorFunc(children, i).Render(prefix)
		maxLen = max(lipgloss.Width(prefix), maxLen)
	}

	for i := range children.Length() {
		child := children.At(i)
		if child.Hidden() {
			continue
		}
		indentStyle := r.style.indenterFunc(children, i)
		enumStyle := r.style.enumeratorFunc(children, i)

		itemStyle := r.style.itemFunc(children, i)

		indent := indentStyle.Render(indenter(children, i))
		nodePrefix := enumStyle.Render(enumerator(children, i))

		// Preserve the background color of the enumerator when adding the padding
		enumBgStyle := lipgloss.NewStyle().Background(enumStyle.GetBackground())

		// Add padding to the left of the node to align it with the longest prefix of its siblings
		if l := maxLen - lipgloss.Width(nodePrefix); l > 0 {
			nodePrefix = enumBgStyle.Render(strings.Repeat(" ", l)) + nodePrefix
		}

		item := itemStyle.Render(child.Value())
		multineLinePrefix := enumBgStyle.Render(prefix)

		// This dance below is to account for multiline prefixes, e.g. "|\n|".
		// In that case, we need to make sure that both the parent prefix and
		// the current node's prefix have the same height.
		for lipgloss.Height(item) > lipgloss.Height(nodePrefix) {
			nodePrefix = lipgloss.JoinVertical(
				lipgloss.Left,
				nodePrefix,
				indent,
			)
		}
		for lipgloss.Height(nodePrefix) > lipgloss.Height(multineLinePrefix) {
			multineLinePrefix = lipgloss.JoinVertical(
				lipgloss.Left,
				multineLinePrefix,
				prefix,
			)
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Top,
			multineLinePrefix,
			nodePrefix,
			item,
		)

		// If the line is shorter than the desired width, we pad it with spaces.
		if pad := r.width - lipgloss.Width(line); pad > 0 {
			line = line + itemStyle.Render(strings.Repeat(" ", pad))
		}
		strs = append(
			strs,
			line,
		)

		if children.Length() > 0 {
			// Here we see if the child has a custom renderer, which means the
			// user set a custom enumerator/indenter/item style, etc.
			// If it has one, we'll use it to render itself.
			// otherwise, we keep using the current renderer.
			// Note that the renderer doesn't inherit its parent's styles.
			renderer := r
			switch child := child.(type) {
			case *Tree:
				if child.r != nil {
					renderer = child.r
				}
			}
			if s := renderer.render(
				child,
				false,
				prefix+indent,
			); s != "" {
				strs = append(strs, s)
			}
		}
	}
	return strings.Join(strs, "\n")
}
