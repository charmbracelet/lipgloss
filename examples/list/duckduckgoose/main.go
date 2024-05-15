package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func main() {
	l := list.New("Duck", "Duck", "Duck", "Duck", "Goose", "Duck", "Duck")

	var DuckDuckGooseEnumerator = func(i int) string {
		if l.At(i) == "Goose" {
			return "Honk →"
		}
		return ""
	}

	l = l.Enumerator(DuckDuckGooseEnumerator).EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("48")).MarginRight(1))

	fmt.Println(l)
}
