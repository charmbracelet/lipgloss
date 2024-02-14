package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StyleFunc allows the list to be styled per item.
type StyleFunc func(i int) lipgloss.Style

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
	data       Data
	style      Style
	indent     int
	offset     int
	height     int
}

// New returns a new list.
func New(items ...string) *List {
	return &List{
		data:   NewStringData(items...),
		height: -1,

		enumerator: Bullet,
		style: Style{
			EnumeratorFunc: func(i int) lipgloss.Style {
				return lipgloss.NewStyle().MarginRight(1)
			},
			ItemFunc: func(i int) lipgloss.Style {
				return lipgloss.NewStyle()
			},
			Base: lipgloss.NewStyle(),
		},
	}
}

// Item appends an item to a list.
func (l *List) Item(item string) *List {
	switch l.data.(type) {
	case *StringData:
		l.data.(*StringData).Append(item)
	}
	return l
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

// Enumerator sets the enumeration type.
func (l *List) Enumerator(enumerator Enumerator) *List {
	l.enumerator = enumerator
	return l
}

// EnumeratorStyle sets the enumerator style.
func (l *List) EnumeratorStyle(style lipgloss.Style) *List {
	l.style.EnumeratorFunc = func(_ int) lipgloss.Style {
		return style
	}
	return l
}

// ItemStyle sets the item style.
func (l *List) ItemStyle(style lipgloss.Style) *List {
	l.style.ItemFunc = func(_ int) lipgloss.Style {
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

	last := l.data.Length() - l.offset
	if l.height > 0 {
		last = l.height + l.offset
	}
	if last <= 0 || last > l.data.Length() {
		return l.style.Base.Render("")
	}

	// find the longest enumerator value of this list.
	var maxLen int
	for i := l.offset; i < last; i++ {
		enum := l.style.EnumeratorFunc(i).Render(l.enumerator(i))
		maxLen = max(lipgloss.Width(enum), maxLen)
	}

	var s strings.Builder

	for i := l.offset; i < last; i++ {
		enum := l.style.EnumeratorFunc(i).Render(l.enumerator(i))
		enumLen := lipgloss.Width(enum)
		prefix := strings.Repeat(" ", l.indent+maxLen-enumLen) + enum
		item := l.style.ItemFunc(i).Render(fmt.Sprintf("%v", l.data.At(i)))
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, prefix, item))
		if i != last-1 {
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
