package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func main() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true).MarginRight(1)

	t := tree.New().Root(".").EnumeratorStyle(enumeratorStyle).ItemStyle(itemStyle)
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			t.Item(tree.New().Root(path))
		}
		return nil
	})

	fmt.Println(t)
}
