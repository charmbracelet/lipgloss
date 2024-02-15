package list

import "github.com/charmbracelet/lipgloss/tree"

// DefaultRenderer returns the default list renderer.
func DefaultRenderer() tree.Renderer {
	return tree.DefaultRenderer().Enumerator(Bullet)
}

// New returns a new list.
func New(items ...any) *tree.TreeNode {
	return tree.New("", items).
		Renderer(DefaultRenderer())
}

// NewSublist returns a new node with the given name and subitems.
func NewSublist(parent string, items ...any) *tree.TreeNode {
	return tree.New(parent, items...)
}
