package list

import "github.com/charmbracelet/lipgloss/tree"

// New returns a new list.
func New(items ...any) *tree.TreeNode {
	t := tree.New("").
		Enumerator(Bullet)
	for _, item := range items {
		t = t.Item(item)
	}
	return t
}
