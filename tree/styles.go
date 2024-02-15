package tree

import "github.com/charmbracelet/lipgloss"

// StyleFunc allows the list to be styled per item.
type StyleFunc func(i int) lipgloss.Style

// Style is the styling applied to the list.
type Style struct {
	PrefixFunc StyleFunc
	ItemFunc   StyleFunc
}

// DefaultStyles is the default tree styles.
func DefaultStyles() Style {
	return Style{
		PrefixFunc: func(int) lipgloss.Style {
			return lipgloss.NewStyle()
		},
		ItemFunc: func(i int) lipgloss.Style {
			return lipgloss.NewStyle().MarginLeft(1)
		},
	}
}
