package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	ul := lipgloss.NewStyle().MarginLeft(2).Underline(true)
	ulsp := lipgloss.NewStyle().MarginLeft(2).Underline(true).UnderlineSpaces(true)
	nul := lipgloss.NewStyle().MarginLeft(2).Underline(false)
	ulnsp := lipgloss.NewStyle().MarginLeft(2).Underline(true).UnderlineSpaces(false)
	nulsp := lipgloss.NewStyle().MarginLeft(2).Underline(false).UnderlineSpaces(true)
	nulnsp := lipgloss.NewStyle().MarginLeft(2).Underline(false).UnderlineSpaces(false)

	fmt.Println(ul.Render("Underline - abc"))
	fmt.Println(ulsp.Render("Underline Spaces - abc"))
	fmt.Println(nul.Render("No Underline - abc"))
	fmt.Println(ulnsp.Render("Underline No Spaces - abc"))
	fmt.Println(nulsp.Render("No Underline Spaces - abc"))
	fmt.Println(nulnsp.Render("No Underline No Spaces - abc"))
}
