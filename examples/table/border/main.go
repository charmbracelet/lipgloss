package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	re := lipgloss.NewRenderer(os.Stdout)

	var (
		HeaderStyle         = re.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Align(lipgloss.Center)
		CellStyle           = re.NewStyle().Padding(0, 1)
		DefaultCellStyle    = CellStyle.Copy().Foreground(lipgloss.Color("6"))
		SelectedCellStyle   = CellStyle.Copy().Foreground(lipgloss.Color("5")).Bold(true)
		DefaultBorderStyle  = re.NewStyle().Foreground(lipgloss.Color("8")).Faint(true)
		SelectedBorderStyle = re.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	)

	rows := [][]string{
		{"English", "Hello", "Hi"},
		{"Chinese", "您好", "你好"},
		{"Japanese", "こんにちは", "やあ"},
		{"Arabic", "أهلين", "أهلا"},
		{"Russian", "Здравствуйте", "Привет"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	for idx := range rows {
		t := table.New().
			BorderRow(true).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == 0 {
					return HeaderStyle
				}
				if row == idx+1 {
					return SelectedCellStyle
				}
				return DefaultCellStyle
			}).
			BorderStyleFunc(func(row, col int, borderType table.BorderType) lipgloss.Style {
				if row == idx {
					switch borderType {
					case table.BorderBottom:
						return SelectedBorderStyle
					}
				} else if row == idx+1 {
					switch borderType {
					case table.BorderLeft:
						if col == 0 {
							return SelectedBorderStyle
						}
					case table.BorderRight, table.BorderBottom:
						return SelectedBorderStyle
					}
				}
				return DefaultBorderStyle
			}).
			Headers("LANGUAGE", "FORMAL", "INFORMAL").
			Rows(rows...)

		fmt.Println(t)
		fmt.Println()
	}
}
