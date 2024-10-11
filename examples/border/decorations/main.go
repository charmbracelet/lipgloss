package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func main() {

	m := struct {
		index  int
		count  int
		status string
	}{
		index:  5,
		count:  10,
		status: ":)",
	}

	reverseStyle := lipgloss.NewStyle().Reverse(true)

	t := lipgloss.NewStyle().
		Width(40).
		Height(10).
		Border(lipgloss.NormalBorder()).
		BorderDecoration(lipgloss.NewBorderDecoration(
			lipgloss.BorderTop,
			lipgloss.Center,
			reverseStyle.Padding(0, 1).Render("BIG TITLE"),
		)).
		BorderDecoration(lipgloss.NewBorderDecoration(
			lipgloss.BorderBottom,
			lipgloss.Right,
			func(width int, middle string) string {
				return reverseStyle.Render(fmt.Sprintf("[%d/%d]", m.index+1, m.count)) + middle
			},
		)).
		BorderDecoration(lipgloss.NewBorderDecoration(
			lipgloss.BorderBottom,
			lipgloss.Left,
			reverseStyle.Padding(0, 1).SetString(fmt.Sprintf("Status: %s", m.status)).String,
		)).
		BorderDecoration(lipgloss.NewBorderDecoration(
			lipgloss.BorderLeft,
			lipgloss.Center,
			reverseStyle.SetString("VERTICAL").String,
		)).
		BorderDecoration(lipgloss.NewBorderDecoration(
			lipgloss.BorderRight,
			lipgloss.Top,
			func(row int, width int, middle string) string {
				if row == 6 {
					return "▶"
				}
				return middle
			},
		))

	fmt.Println()
	fmt.Println(t)
	fmt.Println()
}
