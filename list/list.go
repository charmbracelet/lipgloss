package list

import "github.com/charmbracelet/lipgloss/tree"

// New returns a new list.
func New(items ...any) *tree.TreeNode {
	return tree.New("", items...).
		Enumerator(Bullet)
}
