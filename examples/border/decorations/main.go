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
			lipgloss.Top,
			lipgloss.Center,
			reverseStyle.Padding(0, 1).Render("BIG TITLE"),
		)).
		BorderDecoration(lipgloss.NewBorderDecoration(
			lipgloss.Bottom,
			lipgloss.Right,
			func(width int, middle string) string {
				return reverseStyle.Render(fmt.Sprintf("[%d/%d]", m.index+1, m.count)) + middle
			},
		)).
		BorderDecoration(lipgloss.NewBorderDecoration(
			lipgloss.Bottom,
			lipgloss.Left,
			reverseStyle.Padding(0, 1).SetString(fmt.Sprintf("Status: %s", m.status)).String,
		))

	fmt.Println()
	fmt.Println(t)
	fmt.Println()

}
