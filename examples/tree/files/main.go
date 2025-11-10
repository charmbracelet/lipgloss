package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
)

func addBranches(root *tree.Tree, path string) error {
	items, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			// It's a directory.

			// Skip directories that start with a dot.
			if strings.HasPrefix(item.Name(), ".") {
				continue
			}

			treeBranch := tree.Root(item.Name())
			root.Child(treeBranch)

			// Recurse.
			branchPath := filepath.Join(path, item.Name())
			if err := addBranches(treeBranch, branchPath); err != nil {
				return err
			}
		} else {
			// It's a file.

			// Skip files that start with a dot.
			if strings.HasPrefix(item.Name(), ".") {
				continue
			}

			root.Child(item.Name())
		}
	}

	return nil
}

func main() {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true).PaddingRight(1)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	t := tree.Root(pwd).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(itemStyle).
		ItemStyle(itemStyle)

	if err := addBranches(t, "."); err != nil {
		fmt.Fprintf(os.Stderr, "Error building tree: %v\n", err)
		os.Exit(1)
	}

	lipgloss.Println(t)
}
