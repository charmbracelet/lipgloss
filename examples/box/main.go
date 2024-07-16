package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var box = lipgloss.NewStyle().Width(8).Height(8).Border(lipgloss.RoundedBorder())

func main() {
	// fmt.Println(box.Render("สวัสดีสวัสดี" + "สวัสดีสวัสดี"))
	fmt.Println(box.Render("สวัสดีสวัสดี" + ansi.SetHyperlink("http://example.com") + "สวัสดีสวัสดี" + ansi.ResetHyperlink()))
}
