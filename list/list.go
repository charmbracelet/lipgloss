package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StyleFunc allows the list to be styled per item.
type StyleFunc func(_ *List, i int) lipgloss.Style

// Style is the styling applied to the list.
type Style struct {
	Base           lipgloss.Style
	EnumeratorFunc StyleFunc
	ItemFunc       StyleFunc
}

// List is a list of items.
type List struct {
	enumerator Enumerator
	hide       bool
	items      []any
	style      Style
}

// New returns a new list.
func New(items ...any) *List {
	return &List{
		items: items,

		enumerator: Bullet,
		style: Style{
			EnumeratorFunc: func(_ *List, i int) lipgloss.Style {
				return lipgloss.NewStyle().MarginRight(1)
			},
			ItemFunc: func(_ *List, i int) lipgloss.Style {
				return lipgloss.NewStyle()
			},
			Base: lipgloss.NewStyle(),
		},
	}
}

// Item appends an item to a list.
func (l *List) Item(item any) *List {
	l.items = append(l.items, item)
	return l
}

// At returns the item at index.
func (l *List) At(i int) any {
	if i < 0 || i >= len(l.items) {
		return ""
	}

	return l.items[i]
}

// Enumerator sets the enumeration type.
func (l *List) Enumerator(enumerator Enumerator) *List {
	l.enumerator = enumerator
	return l
}

// EnumeratorStyle sets the enumerator style.
func (l *List) EnumeratorStyle(style lipgloss.Style) *List {
	l.style.EnumeratorFunc = func(_ *List, _ int) lipgloss.Style {
		return style
	}
	return l
}

// ItemStyle sets the item style.
func (l *List) ItemStyle(style lipgloss.Style) *List {
	l.style.ItemFunc = func(_ *List, _ int) lipgloss.Style {
		return style
	}
	return l
}

// EnumeratorStyleFunc sets the enumerator style function.
//
// This option is mutually exclusive with EnumeratorStyle.
func (l *List) EnumeratorStyleFunc(style StyleFunc) *List {
	if style == nil {
		l.EnumeratorStyle(lipgloss.NewStyle())
	}
	l.style.EnumeratorFunc = style
	return l
}

// ItemStyle sets the item style style function.
//
// This option is mutually exclusive with ItemStyle.
func (l *List) ItemStyleFunc(style StyleFunc) *List {
	if style == nil {
		l.ItemStyle(lipgloss.NewStyle())
	}
	l.style.ItemFunc = style
	return l
}

// BaseStyle sets the base style.
func (l *List) BaseStyle(style lipgloss.Style) *List {
	l.style.Base = style
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

	// find the longest enumerator value of this list.
	var maxLen int
	for i := 0; i < len(l.items); i++ {
		enum := l.style.EnumeratorFunc(l, i).Render(l.enumerator(l, i))
		maxLen = max(lipgloss.Width(enum), maxLen)
	}

	var s strings.Builder
	for i, item := range l.items {
		enum := l.style.EnumeratorFunc(l, i).Render(l.enumerator(l, i))
		enumLen := lipgloss.Width(enum)
		enum = strings.Repeat(" ", maxLen-enumLen) + enum
		item := l.style.ItemFunc(l, i).Render(fmt.Sprintf("%v", item))
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, enum, item))
		if i != len(l.items)-1 {
			s.WriteRune('\n')
		}
	}
	return l.style.Base.Render(s.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
