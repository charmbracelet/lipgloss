package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const indentIncrement = 2

// Style is the styling applied to the list.
type Style struct {
	Base       lipgloss.Style
	Enumerator lipgloss.Style
	Item       lipgloss.Style
}

// List is a list of items.
type List struct {
	enumerator Enumerator
	hide       bool
	indent     int
	items      []any
	style      Style
}

// New returns a new list.
func New(items ...any) *List {
	return &List{
		enumerator: Bullet,
		indent:     0,
		items:      items,
		style: Style{
			Enumerator: lipgloss.NewStyle().MarginRight(1),
			Item:       lipgloss.NewStyle(),
			Base:       lipgloss.NewStyle(),
		},
	}
}

// Item appends an item to a list.
func (l *List) Item(item any) *List {
	l.items = append(l.items, item)
	return l
}

// Enumerator sets the enumeration type.
func (l *List) Enumerator(enumerator Enumerator) *List {
	l.enumerator = enumerator
	return l
}

// Indent sets the indent level.
func (l *List) Indent(indent int) *List {
	if indent < 0 {
		l.indent = 0
	} else {
		l.indent = indent
	}
	return l
}

// Hide sets whether or not to hide the list.
// This is useful for collapsing / hiding sub-lists.
func (l *List) Hide(hide bool) *List {
	l.hide = hide
	return l
}

// Render renders the list.
func (l *List) Render() string {
	return l.String()
}

// String returns the string representation of the list.
func (l *List) String() string {
	if l.hide {
		return ""
	}

	var s strings.Builder
	for i, item := range l.items {
		switch item := item.(type) {
		case *List:
			if item.indent <= 0 {
				item = item.Indent(l.indent + indentIncrement)
			}
			s.WriteString(item.String())
		default:
			s.WriteString(strings.Repeat(" ", l.indent))
			s.WriteString(l.style.Enumerator.Render(l.enumerator(l, i+1)))
			s.WriteString(l.style.Item.Render(fmt.Sprintf("%v", item)))
			s.WriteRune('\n')
		}
	}
	return s.String()
}
