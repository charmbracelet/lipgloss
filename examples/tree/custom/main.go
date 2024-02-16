package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

	t := tree.New(
		"Home",
		"Documents",
		"Downloads",
		"Music",
		"Photos",
		"Movies",
	).
		EnumeratorStyle(enumeratorStyle).
		ItemStyle(itemStyle)
	fmt.Println(t)
}
