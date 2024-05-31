package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

type file struct {
	name string
	tree.Node
}

type dir struct {
	name string
	tree.Node
}

func main() {
	t := tree.New().
		Enumerator(func(c tree.Children, i int) string {
			switch t := c.At(i).(type) {
			case file:
				if c.Length()-1 == i {
					return "╰─ "
				}
				return "├─ "
			default:
				if t.Hidden() {
					return "⏵ "
				}
				return "⏷ "
			}
		}).
		Root("~").
		Child(
			file{name: "README.md"},
			file{name: "Groceries"}.Hidden(),
			dir{name: "Emotes"},
			tree.New().Child(
				file{name: "KEK.png"},
				file{name: "LUL.png"},
				file{name: "OmegaLUL.png"},
				file{name: "Poggers.png"},
			),
			file{name: "Kittens"},
		)
	fmt.Println(t)
}
