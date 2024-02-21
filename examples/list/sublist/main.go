package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
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
		Item("Bar").
		Item(list.New("foo2", "bar2")).
		Item("Qux").
		Item(
			list.New("aaa", "bbb").
				EnumeratorStyle(style1).
				Enumerator(list.Roman),
		).
		Item("Deep").
		Item(
			list.New().
				EnumeratorStyle(style2).
				Enumerator(list.Alphabet).
				Item("foo").
				Item("Deeper").
				Item(
					list.New().
						EnumeratorStyle(style1).
						Enumerator(list.Arabic).
						Item("a").
						Item("b").
						Item("Even Deeper, inherit parent renderer").
						Item(
							list.New().
								Enumerator(list.Asterisk).
								EnumeratorStyle(style2).
								Item("sus").
								Item("d minor").
								Item("f#").
								Item("One ore level, with another renderer").
								Item(
									list.New().
										EnumeratorStyle(style1).
										Enumerator(list.Dash).
										Item("a\nmultine\nstring").
										Item("hoccus poccus").
										Item("abra kadabra").
										Item("And finally, a tree within all this").
										Item(

											tree.New("").
												EnumeratorStyle(style2).
												Item("another\nmultine\nstring").
												Item("something").
												Item("a subtree").
												Item(
													tree.New("").
														EnumeratorStyle(style2).
														Item("yup").
														Item("many itens").
														Item(
															lipgloss.NewStyle().
																Bold(true).
																Foreground(lipgloss.Color("#FAFAFA")).
																Background(lipgloss.Color("#7D56F4")).
																AlignHorizontal(lipgloss.Center).
																AlignVertical(lipgloss.Center).
																Padding(1, 3).
																Width(22).
																Render("charming"),
														).
														Item(
															table.New().
																Width(40).
																BorderStyle(style1.Copy().MarginRight(0)).
																StyleFunc(func(row, col int) lipgloss.Style {
																	style := lipgloss.NewStyle()
																	if col == 1 {
																		style = style.Align(lipgloss.Center)
																	}
																	if row == 0 {
																		return style.Bold(true)
																	}
																	return style.Faint(true)
																}).
																Headers("ITEM", "QTY").
																Row("Banana", "10").
																Row("Orange", "2").
																Row("Apple", "6").
																Row("Strawberry", "12"),
														).
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
