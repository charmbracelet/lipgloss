package main

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/tree"
)

func main() {
	backgroundColor := lipgloss.Color("#ee6ff8")
	foregroundColor := lipgloss.Color("#ecfe65")

	darkBg := lipgloss.NewStyle().
		Background(backgroundColor).
		Padding(0, 1)

	headerItemStyle := lipgloss.NewStyle().
		Foreground(foregroundColor).
		Bold(true).
		Padding(0, 1)

	itemStyle := headerItemStyle.Background(backgroundColor)

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
