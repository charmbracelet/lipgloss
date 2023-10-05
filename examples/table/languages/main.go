package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")
)

func main() {
	re := lipgloss.NewRenderer(os.Stdout)

	var (
		// HeaderStyle is the lipgloss style used for the table headers.
		HeaderStyle = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		// CellStyle is the base lipgloss style used for the table rows.
		CellStyle = re.NewStyle().Padding(0, 1).Width(14)
		// OddRowStyle is the lipgloss style used for odd-numbered table rows.
		OddRowStyle = CellStyle.Copy().Foreground(gray)
		// EvenRowStyle is the lipgloss style used for even-numbered table rows.
		EvenRowStyle = CellStyle.Copy().Foreground(lightGray)
		// BorderStyle is the lipgloss style used for the table border.
		BorderStyle = lipgloss.NewStyle().Foreground(purple)
	)

	rows := [][]string{
		{"Chinese", "您好", "你好"},
		{"Japanese", "こんにちは", "やあ"},
		{"Arabic", "أهلين", "أهلا"},
		{"Russian", "Здравствуйте", "Привет"},
		{"Spanish", "Hola", "¿Qué tal?"},
		{"English", "You look absolutely fabulous.", "How's it going?"},
	}

	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				style = EvenRowStyle
			default:
				style = OddRowStyle
			}

			// Make the second column a little wider.
			if col == 1 {
				style = style.Copy().Width(22)
			}

			// Arabic is a right-to-left language, so right align the text.
			if rows[row-1][0] == "Arabic" && col != 0 {
				style = style.Copy().Align(lipgloss.Right)
			}

			return style
		}).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	fmt.Println(t)
}
