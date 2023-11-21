package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/list"
)

func main() {
	l := list.New("A", "B", "C")
	fmt.Println(l)
}
