package main

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/list"
)

func main() {
	l := list.New(
		"A",
		"B",
		"C",
		list.New(
			"D",
			"E",
			"F",
		).Enumerator(list.Roman),
		"G",
	)
	lipgloss.Println(l)
}
