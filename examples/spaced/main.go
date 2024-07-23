package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	fancy := "Fancy"
	fmt.Println(fancy)

	fancyBold := lipgloss.NewStyle().Bold(true).Render(fancy)
	fmt.Printf("%q\n", fancyBold)

	fancyBoldStrikethrough := lipgloss.NewStyle().Strikethrough(true).Render(fancyBold)
	fmt.Printf("%q\n", fancyBoldStrikethrough)
}
