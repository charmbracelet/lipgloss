package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
)

// This example abuses both lists and trees.
// Its a list that goes deep, and items varies from other lists, trees, tables,
// etc.

func main() {
	style1 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		MarginRight(1)
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		MarginRight(1)

	baseStyle := lipgloss.NewStyle().
		MarginBottom(1).
		MarginLeft(1)
	dimColor := lipgloss.Color("250")
	hightlightColor := lipgloss.Color("#EE6FF8")

	l := list.New().
		Item("Item 1").
		Item("Item 2").
		Item(list.New("Item 2.1", "Item 2.2")).
		Item("Item 3").
		Item(
			list.New("Item 3.1", "Item 3.2").
				EnumeratorStyle(style1).
				Enumerator(list.Roman),
		).
		Item("Item 4").
		Item(
			list.New().
				EnumeratorStyle(style2).
				Enumerator(list.Alphabet).
				Item("Item 4.1").
				Item("Item 4.2").
				Item(
					list.New().
						EnumeratorStyle(style1).
						Enumerator(list.Arabic).
						Item("Item 4.2.1").
						Item("Item 4.2.2").
						Item("Item 4.2.3").
						Item(
							list.New().
								Enumerator(list.Asterisk).
								EnumeratorStyle(style2).
								Item("Item 4.2.3.1").
								Item("Item 4.2.3.3").
								Item("Item 4.2.3.4").
								Item("Item 4.2.3.5").
								Item(
									list.New().
										EnumeratorStyle(style1).
										Enumerator(list.Dash).
										Item("Item 4.2.3.3.1\nis a multiline string").
										Item("Item 4.2.3.3.2").
										Item("Item 4.2.3.3.3").
										Item("Item 4.2.3.3.4").
										Item(

											tree.New().
												EnumeratorStyle(style2).
												Item("Item 4.2.3.3.4.1\nThis is a tree within a list.\nLists are actually syntactic sugar for trees!").
												Item("Item 4.2.3.3.4.2").
												Item("Item 4.2.3.3.4.3").
												Item(
													tree.New().
														EnumeratorStyle(style2).
														Item("Item 4.2.3.3.4.3.1\nanother tree").
														Item("Item 4.2.3.3.4.3.2").
														Item(
															lipgloss.NewStyle().
																Bold(true).
																Foreground(lipgloss.Color("#FAFAFA")).
																Background(lipgloss.Color("#7D56F4")).
																AlignHorizontal(lipgloss.Center).
																AlignVertical(lipgloss.Center).
																Padding(1, 3).
																Width(40).
																Render("Item 4.2.3.3.4.3.3\n\nItems can be any string, including tables!"),
														).
														Item(
															list.New("A list within a tree", "a", "b", "c").
																Enumerator(list.Roman),
														).
														Item(
															table.New().
																Width(40).
																BorderStyle(style1.MarginRight(0)).
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
														Item("Item 4.2.3.3.4.3.3").
														Item(
															list.New().
																Enumerator(func(_ list.Items, i int) string {
																	if i == 1 {
																		return "│\n│"
																	}
																	return " "
																}).
																ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
																	st := baseStyle
																	if i == 1 {
																		return st.Foreground(hightlightColor)
																	}
																	return st.Foreground(dimColor)
																}).
																EnumeratorStyleFunc(func(_ list.Items, i int) lipgloss.Style {
																	if i == 1 {
																		return lipgloss.NewStyle().Foreground(hightlightColor)
																	}
																	return lipgloss.NewStyle().Foreground(dimColor)
																}).
																Item("Item a\n1 day ago").
																Item("Item b\n2 days ago").
																Item("Item c\n10 minutes ago").
																Item("Item d\n1 month ago"),
														).
														Item("Item 4.2.3.3.4.3.4"),
												).
												Item("Item 4.2.3.3.4.4").
												Item("Item 4.2.3.3.4.5"),
										).
										Item("Item 4.2.3.3.5"),
								),
						),
				).
				Item("Item 4.3"),
		).
		Item("item 5")
	fmt.Println(l)
}
