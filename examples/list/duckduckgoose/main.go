package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/list"
)

func DuckDuckGooseEnumerator(l *list.List, i int) string {
	if l.At(i) == "Goose" {
		return "→ "
	}
	return "  "
}

func main() {
	l := list.New("Duck", "Duck", "Duck", "Duck", "Goose", "Duck", "Duck")
	l.Enumerator(DuckDuckGooseEnumerator)
	fmt.Println(l)
}
