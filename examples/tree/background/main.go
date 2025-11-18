package main

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
)

func main() {
	darkBg := lipgloss.NewStyle().
		Background(lipgloss.Color("0")).
		Padding(0, 1)

	headerItemStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#ee6ff8")).
		Foreground(lipgloss.Color("#ecfe65")).
		Bold(true).
		Padding(0, 1)

	itemStyle := headerItemStyle.Background(lipgloss.Color("0"))

	t := tree.Root("# Table of Contents").
		RootStyle(itemStyle).
		ItemStyle(itemStyle).
		EnumeratorStyle(darkBg).
		IndenterStyle(darkBg).
		Child(
			tree.Root("## Chapter 1").
				Child("Chapter 1.1").
				Child("Chapter 1.2"),
		).
		Child(
			tree.Root("## Chapter 2").
				Child("Chapter 2.1").
				Child("Chapter 2.2"),
		)

	lipgloss.Println(t)
}
