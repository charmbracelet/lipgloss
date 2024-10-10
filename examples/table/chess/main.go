package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	board := [][]string{
		{"♜", "♞", "♝", "♛", "♚", "♝", "♞", "♜"},
		{"♟", "♟", "♟", "♟", "♟", "♟", "♟", "♟"},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{"♙", "♙", "♙", "♙", "♙", "♙", "♙", "♙"},
		{"♖", "♘", "♗", "♕", "♔", "♗", "♘", "♖"},
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(board...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
		})

	ranks := labelStyle.Render(strings.Join([]string{" A", "B", "C", "D", "E", "F", "G", "H  "}, "   "))
	files := labelStyle.Render(strings.Join([]string{" 1", "2", "3", "4", "5", "6", "7", "8 "}, "\n\n "))

	fmt.Println(lipgloss.JoinVertical(lipgloss.Right, lipgloss.JoinHorizontal(lipgloss.Center, files, t.Render()), ranks) + "\n")
}
