package list

import "github.com/charmbracelet/lipgloss/tree"

// New returns a new list.
func New(items ...any) *tree.TreeNode {
	return tree.New("", items...).
		Enumerator(Bullet)
}

// NewSublist returns a new node with the given name and subitems.
func NewSublist(parent string, items ...any) *tree.TreeNode {
	return tree.New(parent, items...)
}
