package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

type dir string
type file string

func (d dir) String() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Height(2).Render(string(d) + "/")
}

func (f file) String() string {
	return lipgloss.NewStyle().Faint(true).Height(2).Render(string(f))
}

func main() {
	l := tree(
		dir("home"),
		tree(
			file("README.md"),
			file("Groceries"),
			dir("Emotes"),
			tree(
				file("KEK.png"),
				file("LUL.png"),
				file("OmegaLUL.png"),
				file("Poggers.png"),
			),
		),
	)
	fmt.Println(l)
}

func tree(items ...any) *list.List {
	return list.New(items...).Enumerator(list.Tree).Indenter(list.TreeIndent)
}
