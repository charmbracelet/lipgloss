package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/list"
)

func duckDuckGooseEnumerator(data list.Data, i int) string {
	if data.At(i).Name() == "Goose" {
		return "Honk →"
	}
	return ""
}

func main() {
	l := list.New("Duck", "Duck", "Duck", "Duck", "Goose", "Duck", "Duck")
	l.Enumerator(duckDuckGooseEnumerator)
	fmt.Println(l)
}
