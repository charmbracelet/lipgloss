package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

type dir struct {
	name   string
	files  []file
	hidden bool
}

type file struct {
	name   string
	hidden bool
}

func (d dir) At(i int) tree.Node      { return d.files[i] }
func (d dir) Children() tree.Children { return d }
func (d dir) Hidden() bool            { return d.hidden }
func (d dir) Length() int             { return len(d.files) }

func (f file) Children() tree.Children { return nil }
func (f file) Hidden() bool            { return f.hidden }
func (f file) String() string          { return f.name }
func (f file) Value() string           { return f.name }

func main() {
	t := tree.New().
		Root("~").
		Child(
			file{name: "README.md"},
			file{name: "Groceries", hidden: false},
			tree.Root("Emotes").
				Child(
					file{name: "KEK.png"},
					file{name: "LUL.png"},
					file{name: "OmegaLUL.png"},
					file{name: "Poggers.png"},
				),
			file{name: "Kittens"},
		)
	fmt.Println(t)
}
