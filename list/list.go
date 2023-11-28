package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
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
	Items      []any
	style      Style
}

// New returns a new list.
func New(items ...any) *List {
	return &List{
		Items: items,

		enumerator: Bullet,
		indent:     0,
		style: Style{
			Enumerator: lipgloss.NewStyle().MarginRight(1),
			Item:       lipgloss.NewStyle(),
			Base:       lipgloss.NewStyle(),
		},
	}
}

// Item appends an item to a list.
func (l *List) Item(item any) *List {
	l.Items = append(l.Items, item)
	return l
}

// Enumerator sets the enumeration type.
func (l *List) Enumerator(enumerator Enumerator) *List {
	l.enumerator = enumerator
	return l
}

// Style sets the list style.
func (l *List) Style(style Style) *List {
	l.style = style
	return l
}

// EnumeratorStyle sets the enumerator style.
func (l *List) EnumeratorStyle(style lipgloss.Style) *List {
	l.style.Enumerator = style
	return l
}

// ItemStyle sets the item style.
func (l *List) ItemStyle(style lipgloss.Style) *List {
	l.style.Item = style
	return l
}

// BaseStyle sets the base style.
func (l *List) BaseStyle(style lipgloss.Style) *List {
	l.style.Base = style
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

	// find the longest enumerator value of this list.
	var maxLen int
	for i := 0; i < len(l.Items); i++ {
		enum := l.style.Enumerator.Render(l.enumerator(l, i+1))
		maxLen = max(runewidth.StringWidth(enum), maxLen)
	}

	var s strings.Builder
	for i, item := range l.Items {
		switch item := item.(type) {
		case *List:
			if item.indent <= 0 {
				item = item.Indent(l.indent + indentIncrement)
			}
			s.WriteString(item.String())
		default:
			enum := l.style.Enumerator.Render(l.enumerator(l, i+1))
			enumLen := runewidth.StringWidth(enum)
			s.WriteString(strings.Repeat(" ", l.indent))
			s.WriteString(strings.Repeat(" ", maxLen-enumLen))
			s.WriteString(enum)
			s.WriteString(l.style.Item.Render(fmt.Sprintf("%v", item)))
			s.WriteRune('\n')
		}
	}
	return s.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
