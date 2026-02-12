package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	tbl := table.New().
		Headers("Name", "Age", "City").
		Row("Alice", "25", "New York").
		Row("Bob", "30", "Los Angeles").
		Row("Charlie", "35", "Chicago").
		Footers("Total", "3", "3 cities").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4A90E2"))
			}
			if row == table.FooterRow {
				return lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#FFFFFF")).
					Background(lipgloss.Color("#34495E"))
			}
			return lipgloss.NewStyle()
		}).
		Border(lipgloss.RoundedBorder()).
		BorderHeader(true).
		BorderFooter(true).
		BorderColumn(true).
		BorderRow(false)

	fmt.Println(tbl)
}
