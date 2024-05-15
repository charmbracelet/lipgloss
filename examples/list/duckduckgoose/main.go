package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func DuckDuckGooseEnumerator(l list.List) list.Enumerator {
	return func(i int) string {
		if l.At(i) == "Goose" {
			return "Honk →"
		}
		return ""
	}
}

func main() {
	l := list.New("Duck", "Duck", "Duck", "Duck", "Goose", "Duck", "Duck")

	l = l.Enumerator(DuckDuckGooseEnumerator(l)).EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("48")).MarginRight(1))

	fmt.Println(l)
}
