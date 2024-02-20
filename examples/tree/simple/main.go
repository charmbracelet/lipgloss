package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	t := tree.New("root", "child 1", "child 2", tree.New("child 3", "child 3.1"))
	fmt.Println(t)
}
