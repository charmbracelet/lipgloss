package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func main() {
	str := "1 "
	red := lipgloss.Color("9")

	style1 := lipgloss.NewStyle().
		Foreground(red).
		Render(str)

	style2 := lipgloss.NewStyle().
		Width(1).
		Render(style1)

	fmt.Println("Lipgloss")
	fmt.Printf("Before Width: '%v'\n", style1)
	fmt.Printf("After Width : '%v'\n", style2)
	fmt.Println()

	fmt.Println("ANSI")
	fmt.Printf("ASCII Width : '%v'\n", ansi.Wrap(str, 1, ""))
	fmt.Printf("ANSI Width  : '%v'\n", ansi.Wrap(style1, 1, ""))
}
