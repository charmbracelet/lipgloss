package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func main() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).MarginRight(1)

	l := list.New(
		"Glossier",
		"Claire’s Boutique",
		"Nyx",
		"Mac",
		"Milk",
	).
		Enumerator(list.Roman).
		EnumeratorStyle(enumeratorStyle).
		ItemStyle(itemStyle)
	fmt.Println(l)
}
