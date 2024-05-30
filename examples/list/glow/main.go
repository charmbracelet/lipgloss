package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

type Document struct {
	Name string
	Time string
}

var faint = lipgloss.NewStyle().Faint(true)

func (d Document) String() string {
	return d.Name + "\n" +
		faint.Render(d.Time)
}

var docs = []Document{
	{"README.md", "2 minutes ago"},
	{"Example.md", "1 hour ago"},
	{"secrets.md", "1 week ago"},
}

const selectedIndex = 1

func main() {
	baseStyle := lipgloss.NewStyle().
		MarginBottom(1).
		MarginLeft(1)
	dimColor := lipgloss.Color("250")
	hightlightColor := lipgloss.Color("#EE6FF8")

	l := list.New().
		Enumerator(func(_ list.Items, i int) string {
			if i == selectedIndex {
				return "│\n│"
			}
			return " "
		}).
		ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
			st := baseStyle
			if selectedIndex == i {
				return st.Foreground(hightlightColor)
			}
			return st.Foreground(dimColor)
		}).
		EnumeratorStyleFunc(func(_ list.Items, i int) lipgloss.Style {
			if selectedIndex == i {
				return lipgloss.NewStyle().Foreground(hightlightColor)
			}
			return lipgloss.NewStyle().Foreground(dimColor)
		})

	for _, d := range docs {
		l.Item(d.String())
	}

	fmt.Println()
	fmt.Println(l)
}
