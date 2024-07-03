package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func main() {
	enumeratorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("0")).
		Padding(0, 1)

	headerItemStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#ee6ff8")).
		Foreground(lipgloss.Color("#ecfe65")).
		Bold(true).
		Padding(0, 1)

	itemStyle := headerItemStyle.Background(lipgloss.Color("0"))

	t := list.New(
		"# Table of Contents",
		list.New(
			"## Chapter 1",
			list.New(
				"Chapter 1.1",
				"Chapter 1.2",
			).Enumerator(list.Tree).Indenter(list.TreeIndenter),
			"## Chapter 2",
			list.New(
				"Chapter 2.1",
				"Chapter 2.2",
			).Enumerator(list.Tree).Indenter(list.TreeIndenter),
		).Enumerator(list.Tree).Indenter(list.TreeIndenter),
	).
		ItemStyle(itemStyle).
		Enumerator(list.Tree).
		Indenter(list.TreeIndenter).
		EnumeratorStyle(enumeratorStyle)

	fmt.Println(t)
}
