// Package list allows you to build lists, as simple or complicated as you need.
//
// Simply, define a list with some items and set it's rendering properties, like
// enumerator and styling:
//
//	groceries := list.New(
//		"Bananas",
//		"Barley",
//		"Cashews",
//		"Milk",
//		list.New(
//			"Almond Milk"
//			"Coconut Milk"
//			"Full Fat Milk"
//		)
//		"Eggs",
//		"Fish Cake",
//		"Leeks",
//		"Papaya",
//	)
//
//	fmt.Println(groceries)
package list

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

// List represents a list of items that can be displayed. Lists can contain
// lists as items, they will be rendered as nested (sub)lists.
//
// In fact, lists can contain anything as items, like lipgloss.Table or lipgloss.Tree.
type List struct{ tree *tree.Tree }

// New returns a new list with the given items.
//
//	alphabet := list.New(
//		"A",
//		"B",
//		"C",
//		"D",
//		"E",
//		"F",
//		...
//	)
//
// Items can be other lists, trees, tables, rendered markdown;
// anything you want, really.
func New(items ...any) *List {
	l := &List{tree: tree.New()}
	return l.Items(items...).
		Enumerator(Bullet).
		Indenter(func(Items, int) string { return " " })
}

// Items represents the list items.
type Items tree.Children

// StyleFunc is the style function that determines the style of an item.
//
// It takes the list items and index of the list and determines the lipgloss
// Style to use for that index.
//
// Example:
//
//	l := list.New().
//		Item("Red").
//		Item("Green").
//		Item("Blue").
//		ItemStyleFunc(func(items list.Items, i int) lipgloss.Style {
//			switch {
//			case i == 0:
//				return RedStyle
//			case i == 1:
//				return GreenStyle
//			case i == 2:
//				return BlueStyle
//			}
//	})
type StyleFunc func(items Items, index int) lipgloss.Style

// Hidden returns whether this list is hidden.
func (l *List) Hidden() bool {
	return l.tree.Hidden()
}

// Hide hides this list.
// If this list is hidden, it will not be shown when rendered.
func (l *List) Hide(hide bool) *List {
	l.tree.Hide(hide)
	return l
}

// Offset sets the start and end offset for the list.
//
//	Example:
//		l := list.New("A", "B", "C", "D").
//			Offset(1, -1)
//
//	fmt.Println(l)
//
//	 • B
//	 • C
//	 • D
func (l *List) Offset(start, end int) *List {
	l.tree.Offset(start, end)
	return l
}

// Value returns the value of this node.
func (l *List) Value() string {
	return l.tree.Value()
}

func (l *List) String() string {
	return l.tree.String()
}

// EnumeratorStyle sets the enumerator style for all enumerators.
//
// To set the enumerator style conditionally based on the item value or index,
// use [EnumeratorStyleFunc].
func (l *List) EnumeratorStyle(style lipgloss.Style) *List {
	l.tree.EnumeratorStyle(style)
	return l
}

// EnumeratorStyleFunc sets the enumerator style function for the list items.
//
// Use this to conditionally set different styles based on the current items,
// sibling items, or index values (i.e. even or odd).
//
// Example:
//
//	l := list.New().
//		EnumeratorStyleFunc(func(_ list.Items, i int) lipgloss.Style {
//			if selected == i {
//				return lipgloss.NewStyle().Foreground(brightPink)
//			}
//			return lipgloss.NewStyle()
//		})
func (l *List) EnumeratorStyleFunc(f StyleFunc) *List {
	l.tree.EnumeratorStyleFunc(func(children tree.Children, index int) lipgloss.Style {
		return f(children, index)
	})
	return l
}

// Indenter sets the indenter implementation. This is used to change the way
// the tree is indented. The default indentor places a border connecting sibling
// elements and no border for the last child.
//
//	└── Foo
//	    └── Bar
//	        └── Baz
//	            └── Qux
//	                └── Quux
//
// You can define your own indenter.
//
//	func ArrowIndenter(children tree.Children, index int) string {
//		return "→ "
//	}
//
//	→ Foo
//	→ → Bar
//	→ → → Baz
//	→ → → → Qux
//	→ → → → → Quux
func (l *List) Indenter(indenter Indenter) *List {
	l.tree.Indenter(
		func(children tree.Children, index int) string {
			return indenter(children, index)
		},
	)
	return l
}

// ItemStyle sets the item style for all items.
//
// To set the item style conditionally based on the item value or index,
// use [ItemStyleFunc].
func (l *List) ItemStyle(style lipgloss.Style) *List {
	l.tree.ItemStyle(style)
	return l
}

// ItemStyleFunc sets the item style function for the list items.
//
// Use this to conditionally set different styles based on the current items,
// sibling items, or index values.
//
// Example:
//
//	l := list.New().
//		ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
//			if selected == i {
//				return lipgloss.NewStyle().Foreground(brightPink)
//			}
//			return lipgloss.NewStyle()
//		})
func (l *List) ItemStyleFunc(f StyleFunc) *List {
	l.tree.ItemStyleFunc(func(children tree.Children, index int) lipgloss.Style {
		return f(children, index)
	})
	return l
}

// Item appends an item to the list.
//
//	l := list.New().
//		Item("Foo").
//		Item("Bar").
//		Item("Baz")
func (l *List) Item(item any) *List {
	switch item := item.(type) {
	case *List:
		l.tree.Child(item.tree)
	default:
		l.tree.Child(item)
	}
	return l
}

// Items appends multiple items to the list.
//
//	l := list.New().
//		Items("Foo", "Bar", "Baz"),
func (l *List) Items(items ...any) *List {
	for _, item := range items {
		l.Item(item)
	}
	return l
}

// Enumerator sets the list enumerator.
//
// There are several predefined enumerators:
// • Alphabet
// • Arabic
// • Bullet
// • Dash
// • Roman
//
// Or, define your own.
//
//	func echoEnumerator(items list.Items, i int) string {
//		return items.At(i).Value() + ". "
//	}
//
//	l := list.New("Foo", "Bar", "Baz").Enumerator(echoEnumerator)
//	fmt.Println(l)
//
//	 Foo. Foo
//	 Bar. Bar
//	 Baz. Baz
func (l *List) Enumerator(enumerator Enumerator) *List {
	l.tree.Enumerator(func(c tree.Children, i int) string { return enumerator(c, i) })
	return l
}
