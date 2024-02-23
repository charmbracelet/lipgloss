package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
	humanize "github.com/dustin/go-humanize"
)

type Document struct {
	Name string
	Date time.Time
}

var faint = lipgloss.NewStyle().Faint(true)

func (d Document) String() string {
	return d.Name + "\n" +
		faint.Render(humanize.Time(d.Date))
}

var docs = []Document{
	{"README.md", time.Now().Add(-time.Minute * 2)},
	{"Example.md", time.Now().Add(-time.Hour)},
	{"secrets.md", time.Now().Add(-time.Hour * 24 * 7)},
}

const selectedIndex = 1

func main() {
	baseStyle := lipgloss.NewStyle().
		MarginBottom(1).
		MarginLeft(1)
	dimColor := lipgloss.Color("250")
	hightlightColor := lipgloss.Color("#EE6FF8")

	l := list.New().
		Enumerator(func(_ list.Data, i int) string {
			if i == selectedIndex {
				return "│\n│"
			}
			return " "
		}).
		ItemStyleFunc(func(_ tree.Data, i int) lipgloss.Style {
			st := baseStyle.Copy()
			if selectedIndex == i {
				return st.Foreground(hightlightColor)
			}
			return st.Foreground(dimColor)
		}).
		EnumeratorStyleFunc(func(_ tree.Data, i int) lipgloss.Style {
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
