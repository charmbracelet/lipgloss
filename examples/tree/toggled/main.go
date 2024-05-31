package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

type file struct {
	name   string
	hidden bool

	files []file
}

func (d file) At(i int) tree.Node {
	if i < 0 || i >= len(d.files) {
		return nil
	}
	return d.files[i]
}

func (d file) Length() int {
	return len(d.files)
}

func (f file) Children() tree.Children {
	if len(f.files) == 0 {
		return nil
	}
	return f
}

func (f file) Hidden() bool {
	return f.hidden
}

func (f file) String() string {
	return f.name
}

func (f file) Value() string {
	return f.name
}

func main() {
	t := tree.Root("~").Child(
		file{
			files: []file{
				{name: "README.md"},
				{name: "Groceries", hidden: false},
				{name: "Emotes",
					files: []file{
						{name: "KEK.png"},
						{name: "LUL.png"},
						{name: "OmegaLUL.png"},
						{name: "Poggers.png"},
					}},
				{name: "GIFs",
					files: []file{
						{name: "KEK.gif", hidden: true},
						{name: "LUL.gif", hidden: true},
						{name: "OmegaLUL.gif", hidden: true},
						{name: "Poggers.gif", hidden: true},
					}},
				{name: "Kittens"},
			}})
	fmt.Println(t)
}
