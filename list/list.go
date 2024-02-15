package list

import (
	"github.com/charmbracelet/lipgloss/tree"
)

// List is a list of items.
type List struct {
	hide   bool
	data   Data
	indent int
	offset int
	height int
}

func DefaultRenderer() tree.Renderer {
	return tree.DefaultRenderer().Enumerator(Bullet)
}

// New returns a new list.
func New(items ...any) *tree.TreeNode {
	return tree.New("", items).Renderer(DefaultRenderer())
}

func NewSublist(parent string, items ...any) *tree.TreeNode {
	return tree.New(parent, items...).Renderer(DefaultRenderer())
}

// At returns the item at index.
func (l *List) At(i int) any {
	return l.data.At(i)
}

// Data sets the list data.
func (l *List) Data(data Data) *List {
	l.data = data
	return l
}

// Indent sets the list indent.
func (l *List) Indent(indent int) *List {
	l.indent = indent
	return l
}

// Offset sets the list offset.
func (l *List) Offset(offset int) *List {
	l.offset = offset
	return l
}

// Height sets the list height.
func (l *List) Height(height int) *List {
	l.height = height
	return l
}

// Hide sets whether or not to hide the list.
// This is useful for collapsing / hiding sub-lists.
func (l *List) Hide(hide bool) *List {
	l.hide = hide
	return l
}
