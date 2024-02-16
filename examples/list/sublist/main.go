package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func main() {
	style1 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		MarginRight(1)
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		MarginRight(1)

	l := list.New().
		Item("Foo").
		Item(list.NewSublist("Bar", "foo2", "bar2")).
		Item(
			list.NewSublist("Qux", "aaa", "bbb").
				EnumeratorStyle(style1).
				Enumerator(list.Roman),
		).
		Item(
			list.NewSublist("Deep").
				EnumeratorStyle(style2).
				Enumerator(list.Alphabet).
				Item("foo").
				Item(
					list.NewSublist("Deeper").
						EnumeratorStyle(style1).
						Enumerator(list.Arabic).
						Item("a").
						Item("b").
						Item(
							list.NewSublist("Even Deeper, inherit parent renderer").
								Enumerator(list.Asterisk).
								EnumeratorStyle(style2).
								Item("sus").
								Item("d minor").
								Item("f#").
								Item(

									list.NewSublist("One ore level, with another renderer").
										EnumeratorStyle(style1).
										Enumerator(list.Dash).
										Item("a\nmultine\nstring").
										Item("hoccus poccus").
										Item("abra kadabra").
										Item(

											list.NewSublist("And finally, a tree within all this").
												EnumeratorStyle(style2).
												Item("another\nmultine\nstring").
												Item("something").
												Item(

													list.NewSublist("And finally, a tree within all this").
														EnumeratorStyle(style2).
														Item("another"),
												).
												Item("hallo").
												Item("wunderbar!"),
										).
										Item("this is a tree\nand other obvious statements"),
								),
						),
				).
				Item("bar"),
		).
		Item("Baz")
	fmt.Println(l)
}
