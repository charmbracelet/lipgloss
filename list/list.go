package list

import "github.com/charmbracelet/lipgloss/tree"

// NewDefaultRenderer returns the default list renderer.
func NewDefaultRenderer() *tree.DefaultRenderer {
	return tree.NewDefaultRenderer().Enumerator(Bullet)
}

// New returns a new list.
func New(items ...any) *tree.TreeNode {
	return tree.New("", items...).
		Renderer(NewDefaultRenderer())
}

// NewSublist returns a new node with the given name and subitems.
func NewSublist(parent string, items ...any) *tree.TreeNode {
	return tree.New(parent, items...)
}
