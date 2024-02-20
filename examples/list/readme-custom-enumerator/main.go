package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	var duckDuckGooseEnumerator tree.Enumerator = func(atter tree.Atter, i int, last bool) (string, string) {
		if atter.At(i).Name() == "Goose" {
			return "", "Honk →"
		}
		return "", ""
	}
	l := list.New("Duck", "Duck", "Duck", "Duck", "Goose", "Duck", "Duck")
	l.Enumerator(duckDuckGooseEnumerator)
	fmt.Println(l)
}
