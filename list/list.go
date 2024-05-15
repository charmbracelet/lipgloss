// Package list defines an API to build lists.
//
// A list is a enumerated collection of items. Items are rendered vertically
// stacked on top of each other. Like the following:
//
//	list.New("A", "B", "C").
//	  Enumerator(list.Bullet).
//	  String()
//
// • A
// • B
// • C
package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Item is a list item. It allows the list to print it's elements, a list can
// contain nested lists.
type Item interface {
	any | List
}

// List is the representation of a lipgloss List.
//
// It can be printed to display a human-readable list.
//
//	fmt.Println(list.New("A", "B", "C"))
//
//	• A
//	• B
//	• C
type List struct {
	indent              int
	separator           lipgloss.Style
	items               []Item
	itemStyleFunc       StyleFunc
	enumerator          Enumerator
	enumeratorStyleFunc StyleFunc
}

// StyleFunc defines a list style function that returns the correct style for
// the list element at the given index.
type StyleFunc func(i int) lipgloss.Style

// New returns a new list given the list items. List items may be of any
// primitive type and other Lists.
//
//	list.New(
//	  "Foo",
//	  "Bar",
//	  NewBaz{},
//	  list.New(
//	    "Qux",
//	    "Quux",
//	  ),
//	)
func New(items ...Item) List {
	return List{
		items:      items,
		separator:  lipgloss.NewStyle().SetString("\n"),
		enumerator: func(_ int) string { return "•" },
		enumeratorStyleFunc: func(_ int) lipgloss.Style {
			return lipgloss.NewStyle().MarginRight(1)
		},
		itemStyleFunc: func(_ int) lipgloss.Style {
			return lipgloss.NewStyle()
		},
	}
}

// Item appends the given item to the list.
// Alias for AddItem for a more fluent API.
//
//	New().
//	  Item("A").
//	  Item("B").
//	  Item("C").
//	  Item("D").
func (l List) Item(item Item) List {
	return l.AddItem(item)
}

// AddItem appends the given item to the list.
func (l List) AddItem(item Item) List {
	l.items = append(l.items, item)
	return l
}

// Separator sets the list separator style, which separates list items.
// The default separator is a newline.
func (l List) Separator(s lipgloss.Style) List {
	l.separator = s
	return l
}

// Prefix sets a static enumerator.
//
// list.New(...).Prefix(s) is equivalent to:
//
//	list.New(...).Enumerator(func(_ int) string { return s })
func (l List) Prefix(s string) List {
	l.enumerator = func(_ int) string { return s }
	return l
}

// Enumerator sets list enumeration function. The function will be called with
// the index of the item in the list.
func (l List) Enumerator(e Enumerator) List {
	l.enumerator = e
	return l
}

// EnumeratorStyleFunc sets the style function for the list enumeration. The
// function will be called with the index of the item in the list.
func (l List) EnumeratorStyleFunc(f StyleFunc) List {
	l.enumeratorStyleFunc = f
	return l
}

// ItemStyleFunc sets the style function for the list enumeration. The
// function will be called with the index of the item in the list.
func (l List) ItemStyleFunc(f StyleFunc) List {
	l.itemStyleFunc = f
	return l
}

// EnumeratorStyle sets the style for the list enumeration. The function will be
// called with the index of the item in the list.
func (l List) EnumeratorStyle(style lipgloss.Style) List {
	l.enumeratorStyleFunc = func(_ int) lipgloss.Style { return style }
	return l
}

// ItemStyle sets the style function for the list enumeration. The
// function will be called with the index of the item in the list.
func (l List) ItemStyle(style lipgloss.Style) List {
	l.itemStyleFunc = func(_ int) lipgloss.Style { return style }
	return l
}

// At returns the item at index i.
func (l List) At(i int) Item {
	if i >= 0 && i < len(l.items) {
		return l.items[i]
	}
	return nil
}

const indent = 2

// String returns a string representation of the list.
func (l List) String() string {
	var s strings.Builder

	// find the longest enumerator value of this list.
	last := len(l.items) - 1
	var maxLen int
	for i := 0; i <= last; i++ {
		enum := l.enumeratorStyleFunc(i).Render(l.enumerator(i))
		maxLen = max(lipgloss.Width(enum), maxLen)
	}

	for i, item := range l.items {
		switch item := item.(type) {
		case List:
			item.indent = l.indent + indent
			s.WriteString(item.String())
			if i != last {
				s.WriteString(l.separator.String())
			}
		default:
			indent := strings.Repeat(" ", l.indent)
			enumerator := l.enumeratorStyleFunc(i).
				Width(maxLen - 1).
				Align(lipgloss.Right).
				Render(l.enumerator(i))
			listItem := lipgloss.JoinHorizontal(
				lipgloss.Top,
				indent,
				enumerator,
				l.itemStyleFunc(i).Render(fmt.Sprintf("%v", item)),
			)
			s.WriteString(listItem)
			if i != last {
				s.WriteString(l.separator.String())
			}
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
