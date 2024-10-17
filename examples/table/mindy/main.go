package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	labelStyle := lipgloss.NewStyle().Width(3).Align(lipgloss.Right)
	swatchStyle := lipgloss.NewStyle().Width(6)

	data := [][]string{}
	for i := 0; i < 13; i += 8 {
		data = append(data, makeRow(i, i+5))
	}
	data = append(data, makeEmptyRow())
	for i := 6; i < 15; i += 8 {
		data = append(data, makeRow(i, i+1))
	}
	data = append(data, makeEmptyRow())
	for i := 16; i < 231; i += 6 {
		data = append(data, makeRow(i, i+5))
	}
	data = append(data, makeEmptyRow())
	for i := 232; i < 256; i += 6 {
		data = append(data, makeRow(i, i+5))
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			color := lipgloss.Color(fmt.Sprint(data[row-1][col-col%2]))
			switch {
			case col%2 == 0:
				return labelStyle.Foreground(color)
			default:
				return swatchStyle.Background(color)
			}
		})

	lipgloss.Println(t)
}

const rowLength = 12

func makeRow(start, end int) []string {
	var row []string
	for i := start; i <= end; i++ {
		row = append(row, fmt.Sprint(i))
		row = append(row, "")
	}
	for i := len(row); i < rowLength; i++ {
		row = append(row, "")
	}
	return row
}

func makeEmptyRow() []string {
	return makeRow(0, -1)
}
