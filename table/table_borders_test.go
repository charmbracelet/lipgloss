package table

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

var TableStyleRegular = func(row, col int) lipgloss.Style {
	switch {
	case row == HeaderRow:
		return lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center)
	default:
		return lipgloss.NewStyle().Padding(0, 1)
	}
}

var (
	basicRows = [][]string{
		{"Apple", "Green", "8cm", "7/10"},
		{"Strawberry", "Red", "3cm", "8/10"},
		{"Plum", "Purple", "6cm", "7/10"},
		{"Blueberry", "Blue", "1cm", "9/10"},
	}

	duplicateRows = [][]string{
		{"Apple", "Green", "8cm", "7/10"},
		{"Apple", "Red", "8cm", "8/10"},
		{"Blueberry", "Blue", "1cm", "8/10"},
		{"Blackberry", "Blue", "1.5cm", "9/10"},
	}
)

func TestTableMergeColumnsSingleRow(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "Green", "7/10"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderMergeColumns(0, 1, 2, 3).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Color2", "Rating").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌───────┬───────┬────────┬────────┐
│ Fruit │ Color │ Color2 │ Rating │
├───────┼───────┼────────┼────────┤
│ Apple │ Green │ Green  │ 7/10   │
└───────┴───────┴────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableMergeColumnsSingleCol(t *testing.T) {
	rows := [][]string{
		{"Apple"},
		{"Apple"},
		{"Strawberry"},
		{"Plum"},
		{"Plum"},
		{"Blueberry"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderMergeColumns(0).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┐
│   Fruit    │
├────────────┤
│ Apple      │
│            │
│            │
├────────────┤
│ Strawberry │
├────────────┤
│ Plum       │
│            │
│            │
├────────────┤
│ Blueberry  │
└────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableMergeColumnsDifferentColumns(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm", "7/10"},
		{"Apple", "Red", "8cm", "8/10"},
		{"Peach", "Orange", "7cm", "8/10"},
		{"Blueberry", "Blue", "1cm", "8/10"},
		{"Blackberry", "Blue", "1.5cm", "9/10"},
		{"Plum", "Purple", "6cm", "9/10"},
		{"Plum", "Yellow", "6cm", "7/10"},
		{"Raspberry", "Pink", "1cm", "7/10"},
	}

	tests := []struct {
		name               string
		borderMergeColumns []int
		expected           string
	}{
		{
			name:               "Merge Column 0",
			borderMergeColumns: []int{0},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
│            ├────────┼──────────────┼────────┤
│            │ Red    │ 8cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Peach      │ Orange │ 7cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blackberry │ Blue   │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 9/10   │
│            ├────────┼──────────────┼────────┤
│            │ Yellow │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Raspberry  │ Pink   │ 1cm          │ 7/10   │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge Column 1",
			borderMergeColumns: []int{1},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Red    │ 8cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Peach      │ Orange │ 7cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 8/10   │
├────────────┤        ├──────────────┼────────┤
│ Blackberry │        │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Yellow │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Raspberry  │ Pink   │ 1cm          │ 7/10   │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge Column 2",
			borderMergeColumns: []int{2},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┤              ├────────┤
│ Apple      │ Red    │              │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Peach      │ Orange │ 7cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blackberry │ Blue   │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 9/10   │
├────────────┼────────┤              ├────────┤
│ Plum       │ Yellow │              │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Raspberry  │ Pink   │ 1cm          │ 7/10   │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge Column 3",
			borderMergeColumns: []int{3},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Red    │ 8cm          │ 8/10   │
├────────────┼────────┼──────────────┤        │
│ Peach      │ Orange │ 7cm          │        │
├────────────┼────────┼──────────────┤        │
│ Blueberry  │ Blue   │ 1cm          │        │
├────────────┼────────┼──────────────┼────────┤
│ Blackberry │ Blue   │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┤        │
│ Plum       │ Purple │ 6cm          │        │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Yellow │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┤        │
│ Raspberry  │ Pink   │ 1cm          │        │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge Column 0 2",
			borderMergeColumns: []int{0, 2},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
│            ├────────┤              ├────────┤
│            │ Red    │              │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Peach      │ Orange │ 7cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blackberry │ Blue   │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 9/10   │
│            ├────────┤              ├────────┤
│            │ Yellow │              │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Raspberry  │ Pink   │ 1cm          │ 7/10   │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge Column 0 1",
			borderMergeColumns: []int{0, 1},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
│            ├────────┼──────────────┼────────┤
│            │ Red    │ 8cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Peach      │ Orange │ 7cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 8/10   │
├────────────┤        ├──────────────┼────────┤
│ Blackberry │        │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 9/10   │
│            ├────────┼──────────────┼────────┤
│            │ Yellow │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Raspberry  │ Pink   │ 1cm          │ 7/10   │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge Column 2 3",
			borderMergeColumns: []int{2, 3},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┤              ├────────┤
│ Apple      │ Red    │              │ 8/10   │
├────────────┼────────┼──────────────┤        │
│ Peach      │ Orange │ 7cm          │        │
├────────────┼────────┼──────────────┤        │
│ Blueberry  │ Blue   │ 1cm          │        │
├────────────┼────────┼──────────────┼────────┤
│ Blackberry │ Blue   │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┤        │
│ Plum       │ Purple │ 6cm          │        │
├────────────┼────────┤              ├────────┤
│ Plum       │ Yellow │              │ 7/10   │
├────────────┼────────┼──────────────┤        │
│ Raspberry  │ Pink   │ 1cm          │        │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge Columns Order Switch 3 1",
			borderMergeColumns: []int{3, 1},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Red    │ 8cm          │ 8/10   │
├────────────┼────────┼──────────────┤        │
│ Peach      │ Orange │ 7cm          │        │
├────────────┼────────┼──────────────┤        │
│ Blueberry  │ Blue   │ 1cm          │        │
├────────────┤        ├──────────────┼────────┤
│ Blackberry │        │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┤        │
│ Plum       │ Purple │ 6cm          │        │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Yellow │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┤        │
│ Raspberry  │ Pink   │ 1cm          │        │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Out of Scope : Merge Column 4 5",
			borderMergeColumns: []int{4, 5},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Red    │ 8cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Peach      │ Orange │ 7cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blackberry │ Blue   │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Yellow │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Raspberry  │ Pink   │ 1cm          │ 7/10   │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Out of Scope All Neg : Merge Column -2 -1",
			borderMergeColumns: []int{-2, -1},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Red    │ 8cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Peach      │ Orange │ 7cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blackberry │ Blue   │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Yellow │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Raspberry  │ Pink   │ 1cm          │ 7/10   │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Out of Scope : Merge Column -1 0",
			borderMergeColumns: []int{-1, 0},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
│            ├────────┼──────────────┼────────┤
│            │ Red    │ 8cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Peach      │ Orange │ 7cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blackberry │ Blue   │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 9/10   │
│            ├────────┼──────────────┼────────┤
│            │ Yellow │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Raspberry  │ Pink   │ 1cm          │ 7/10   │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge Column 1 2 3",
			borderMergeColumns: []int{1, 2, 3},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┤              ├────────┤
│ Apple      │ Red    │              │ 8/10   │
├────────────┼────────┼──────────────┤        │
│ Peach      │ Orange │ 7cm          │        │
├────────────┼────────┼──────────────┤        │
│ Blueberry  │ Blue   │ 1cm          │        │
├────────────┤        ├──────────────┼────────┤
│ Blackberry │        │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┤        │
│ Plum       │ Purple │ 6cm          │        │
├────────────┼────────┤              ├────────┤
│ Plum       │ Yellow │              │ 7/10   │
├────────────┼────────┼──────────────┤        │
│ Raspberry  │ Pink   │ 1cm          │        │
└────────────┴────────┴──────────────┴────────┘
`),
		},
		{
			name:               "Merge All Columns 0 1 2 3",
			borderMergeColumns: []int{0, 1, 2, 3},
			expected: strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
│            ├────────┤              ├────────┤
│            │ Red    │              │ 8/10   │
├────────────┼────────┼──────────────┤        │
│ Peach      │ Orange │ 7cm          │        │
├────────────┼────────┼──────────────┤        │
│ Blueberry  │ Blue   │ 1cm          │        │
├────────────┤        ├──────────────┼────────┤
│ Blackberry │        │ 1.5cm        │ 9/10   │
├────────────┼────────┼──────────────┤        │
│ Plum       │ Purple │ 6cm          │        │
│            ├────────┤              ├────────┤
│            │ Yellow │              │ 7/10   │
├────────────┼────────┼──────────────┤        │
│ Raspberry  │ Pink   │ 1cm          │        │
└────────────┴────────┴──────────────┴────────┘
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := New().
				Border(lipgloss.NormalBorder()).
				BorderColumn(true).
				BorderRow(true).
				BorderMergeColumns(tt.borderMergeColumns...).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
				StyleFunc(TableStyleRegular).
				Headers("Fruit", "Color", "Avg Diameter", "Rating").
				Rows(rows...)

			if table.String() != tt.expected {
				t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", tt.expected, table.String())
			}
		})
	}
}

func TestTableEmptyMergeColumnsAllBorders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderMergeColumns().
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(duplicateRows...)

	expected := strings.TrimSpace(`
┌────────────┬───────┬──────────────┬────────┐
│   Fruit    │ Color │ Avg Diameter │ Rating │
├────────────┼───────┼──────────────┼────────┤
│ Apple      │ Green │ 8cm          │ 7/10   │
├────────────┼───────┼──────────────┼────────┤
│ Apple      │ Red   │ 8cm          │ 8/10   │
├────────────┼───────┼──────────────┼────────┤
│ Blueberry  │ Blue  │ 1cm          │ 8/10   │
├────────────┼───────┼──────────────┼────────┤
│ Blackberry │ Blue  │ 1.5cm        │ 9/10   │
└────────────┴───────┴──────────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableEmptyMergeColumnsNoBorders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(false).
		BorderMergeColumns().
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(duplicateRows...)

	expected := strings.TrimSpace(`
┌─────────────────────────────────────────┐
│   Fruit     Color  Avg Diameter  Rating │
├─────────────────────────────────────────┤
│ Apple       Green  8cm           7/10   │
│ Apple       Red    8cm           8/10   │
│ Blueberry   Blue   1cm           8/10   │
│ Blackberry  Blue   1.5cm         9/10   │
└─────────────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableMergeColumnsNoRowNoCol(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(false).
		BorderMergeColumns(0, 1, 2, 3).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(duplicateRows...)

	expected := strings.TrimSpace(`
┌─────────────────────────────────────────┐
│   Fruit     Color  Avg Diameter  Rating │
├─────────────────────────────────────────┤
│ Apple       Green  8cm           7/10   │
│ Apple       Red    8cm           8/10   │
│ Blueberry   Blue   1cm           8/10   │
│ Blackberry  Blue   1.5cm         9/10   │
└─────────────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableMergeColumnsNoCol(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(true).
		BorderMergeColumns(0, 1, 2, 3).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(duplicateRows...)

	expected := strings.TrimSpace(`
┌─────────────────────────────────────────┐
│   Fruit     Color  Avg Diameter  Rating │
├─────────────────────────────────────────┤
│ Apple       Green  8cm           7/10   │
├─────────────────────────────────────────┤
│ Apple       Red    8cm           8/10   │
├─────────────────────────────────────────┤
│ Blueberry   Blue   1cm           8/10   │
├─────────────────────────────────────────┤
│ Blackberry  Blue   1.5cm         9/10   │
└─────────────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableMergeColumnsNoRow(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(false).
		BorderMergeColumns(0, 1, 2, 3).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(duplicateRows...)

	expected := strings.TrimSpace(`
┌────────────┬───────┬──────────────┬────────┐
│   Fruit    │ Color │ Avg Diameter │ Rating │
├────────────┼───────┼──────────────┼────────┤
│ Apple      │ Green │ 8cm          │ 7/10   │
│ Apple      │ Red   │ 8cm          │ 8/10   │
│ Blueberry  │ Blue  │ 1cm          │ 8/10   │
│ Blackberry │ Blue  │ 1.5cm        │ 9/10   │
└────────────┴───────┴──────────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableTooManyHeadersWithBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm"},
		{"Strawberry", "Red", "3cm"},
		{"Plum", "Purple", "6cm"},
		{"Blueberry", "Blue", "1cm"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │         
├────────────┼────────┼──────────────┼────────┤
│ Strawberry │ Red    │ 3cm          │         
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │         
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │         
└────────────┴────────┴──────────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

// ///
func TestTableTooManyHeadersColBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm"},
		{"Apple", "Red", "8cm"},
		{"Blueberry", "Blue", "1cm"},
		{"Blackberry", "Blue", "1.5cm"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┬───────┬──────────────┬────────┐
│   Fruit    │ Color │ Avg Diameter │ Rating │
├────────────┼───────┼──────────────┼────────┤
│ Apple      │ Green │ 8cm          │         
│ Apple      │ Red   │ 8cm          │         
│ Blueberry  │ Blue  │ 1cm          │         
│ Blackberry │ Blue  │ 1.5cm        │         
└────────────┴───────┴──────────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableTooFewHeadersWithBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm"},
		{"Apple", "Red", "8cm"},
		{"Blueberry", "Blue", "1cm"},
		{"Blackberry", "Blue", "1.5cm"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┬───────┬───────┐
│   Fruit    │       │       │
├────────────┼───────┼───────┤
│ Apple      │ Green │ 8cm   │
├────────────┼───────┼───────┤
│ Apple      │ Red   │ 8cm   │
├────────────┼───────┼───────┤
│ Blueberry  │ Blue  │ 1cm   │
├────────────┼───────┼───────┤
│ Blackberry │ Blue  │ 1.5cm │
└────────────┴───────┴───────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableTooFewHeadersColBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm"},
		{"Apple", "Red", "8cm"},
		{"Blueberry", "Blue", "1cm"},
		{"Blackberry", "Blue", "1.5cm"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┬───────┬───────┐
│   Fruit    │       │       │
├────────────┼───────┼───────┤
│ Apple      │ Green │ 8cm   │
│ Apple      │ Red   │ 8cm   │
│ Blueberry  │ Blue  │ 1cm   │
│ Blackberry │ Blue  │ 1.5cm │
└────────────┴───────┴───────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableTooFewHeadersNoBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm"},
		{"Apple", "Red", "8cm"},
		{"Blueberry", "Blue", "1cm"},
		{"Blackberry", "Blue", "1.5cm"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌──────────────────────────┐
│   Fruit                  │
├──────────────────────────┤
│ Apple       Green  8cm   │
│ Apple       Red    8cm   │
│ Blueberry   Blue   1cm   │
│ Blackberry  Blue   1.5cm │
└──────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicNoBorders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(basicRows...)

	expected := strings.TrimSpace(`
┌──────────────────────────────────────────┐
│   Fruit     Color   Avg Diameter  Rating │
├──────────────────────────────────────────┤
│ Apple       Green   8cm           7/10   │
│ Strawberry  Red     3cm           8/10   │
│ Plum        Purple  6cm           7/10   │
│ Blueberry   Blue    1cm           9/10   │
└──────────────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicRowBorders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(basicRows...)

	expected := strings.TrimSpace(`
┌──────────────────────────────────────────┐
│   Fruit     Color   Avg Diameter  Rating │
├──────────────────────────────────────────┤
│ Apple       Green   8cm           7/10   │
├──────────────────────────────────────────┤
│ Strawberry  Red     3cm           8/10   │
├──────────────────────────────────────────┤
│ Plum        Purple  6cm           7/10   │
├──────────────────────────────────────────┤
│ Blueberry   Blue    1cm           9/10   │
└──────────────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicColBorders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(basicRows...)

	expected := strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
│ Strawberry │ Red    │ 3cm          │ 8/10   │
│ Plum       │ Purple │ 6cm          │ 7/10   │
│ Blueberry  │ Blue   │ 1cm          │ 9/10   │
└────────────┴────────┴──────────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicAllBorders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(basicRows...)

	expected := strings.TrimSpace(`
┌────────────┬────────┬──────────────┬────────┐
│   Fruit    │ Color  │ Avg Diameter │ Rating │
├────────────┼────────┼──────────────┼────────┤
│ Apple      │ Green  │ 8cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Strawberry │ Red    │ 3cm          │ 8/10   │
├────────────┼────────┼──────────────┼────────┤
│ Plum       │ Purple │ 6cm          │ 7/10   │
├────────────┼────────┼──────────────┼────────┤
│ Blueberry  │ Blue   │ 1cm          │ 9/10   │
└────────────┴────────┴──────────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicNoBordersNoHeaders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(basicRows...)

	expected := strings.TrimSpace(`
┌───────────────────────────────┐
│ Apple       Green   8cm  7/10 │
│ Strawberry  Red     3cm  8/10 │
│ Plum        Purple  6cm  7/10 │
│ Blueberry   Blue    1cm  9/10 │
└───────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}
func TestTableBasicRowBordersNoHeaders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(basicRows...)

	expected := strings.TrimSpace(`
┌───────────────────────────────┐
│ Apple       Green   8cm  7/10 │
├───────────────────────────────┤
│ Strawberry  Red     3cm  8/10 │
├───────────────────────────────┤
│ Plum        Purple  6cm  7/10 │
├───────────────────────────────┤
│ Blueberry   Blue    1cm  9/10 │
└───────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}
func TestTableBasicColBordersNoHeaders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(basicRows...)

	expected := strings.TrimSpace(`
┌────────────┬────────┬─────┬──────┐
│ Apple      │ Green  │ 8cm │ 7/10 │
│ Strawberry │ Red    │ 3cm │ 8/10 │
│ Plum       │ Purple │ 6cm │ 7/10 │
│ Blueberry  │ Blue   │ 1cm │ 9/10 │
└────────────┴────────┴─────┴──────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicAllBordersNoHeaders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(basicRows...)

	expected := strings.TrimSpace(`
┌────────────┬────────┬─────┬──────┐
│ Apple      │ Green  │ 8cm │ 7/10 │
├────────────┼────────┼─────┼──────┤
│ Strawberry │ Red    │ 3cm │ 8/10 │
├────────────┼────────┼─────┼──────┤
│ Plum       │ Purple │ 6cm │ 7/10 │
├────────────┼────────┼─────┼──────┤
│ Blueberry  │ Blue   │ 1cm │ 9/10 │
└────────────┴────────┴─────┴──────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleCell(t *testing.T) {
	rows := [][]string{
		{"Apple"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(rows...)

	expected := strings.TrimSpace(`
┌───────┐
│ Apple │
└───────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleCellNoInnerBorders(t *testing.T) {
	rows := [][]string{
		{"Apple"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(rows...)

	expected := strings.TrimSpace(`
┌───────┐
│ Apple │
└───────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicEmptyTable(t *testing.T) {
	rows := [][]string{}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌───────┬───────┬──────────────┬────────┐
│ Fruit │ Color │ Avg Diameter │ Rating │
├───────┼───────┼──────────────┼────────┤
└───────┴───────┴──────────────┴────────┘`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicEmptyTableNoHeaders(t *testing.T) {
	rows := [][]string{}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(rows...)

	expected := strings.TrimSpace(``)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleRowWithRowBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm", "7/10"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────────────────────────────┐
│ Fruit  Color  Avg Diameter  Rating │
├────────────────────────────────────┤
│ Apple  Green  8cm           7/10   │
└────────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleRowWithColBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm", "7/10"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌───────┬───────┬──────────────┬────────┐
│ Fruit │ Color │ Avg Diameter │ Rating │
├───────┼───────┼──────────────┼────────┤
│ Apple │ Green │ 8cm          │ 7/10   │
└───────┴───────┴──────────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleRowNoBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm", "7/10"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────────────────────────────┐
│ Fruit  Color  Avg Diameter  Rating │
├────────────────────────────────────┤
│ Apple  Green  8cm           7/10   │
└────────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleRowAllBorders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm", "7/10"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit", "Color", "Avg Diameter", "Rating").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌───────┬───────┬──────────────┬────────┐
│ Fruit │ Color │ Avg Diameter │ Rating │
├───────┼───────┼──────────────┼────────┤
│ Apple │ Green │ 8cm          │ 7/10   │
└───────┴───────┴──────────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleColumnWithRowBorders(t *testing.T) {
	rows := [][]string{
		{"Apple"},
		{"Strawberry"},
		{"Plum"},
		{"Blueberry"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┐
│   Fruit    │
├────────────┤
│ Apple      │
├────────────┤
│ Strawberry │
├────────────┤
│ Plum       │
├────────────┤
│ Blueberry  │
└────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleColumnWithColBorders(t *testing.T) {
	rows := [][]string{
		{"Apple"},
		{"Strawberry"},
		{"Plum"},
		{"Blueberry"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┐
│   Fruit    │
├────────────┤
│ Apple      │
│ Strawberry │
│ Plum       │
│ Blueberry  │
└────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleColumnNoBorders(t *testing.T) {
	rows := [][]string{
		{"Apple"},
		{"Strawberry"},
		{"Plum"},
		{"Blueberry"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		BorderRow(false).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┐
│   Fruit    │
├────────────┤
│ Apple      │
│ Strawberry │
│ Plum       │
│ Blueberry  │
└────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleColumnAllBorders(t *testing.T) {
	rows := [][]string{
		{"Apple"},
		{"Strawberry"},
		{"Plum"},
		{"Blueberry"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Headers("Fruit").
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┐
│   Fruit    │
├────────────┤
│ Apple      │
├────────────┤
│ Strawberry │
├────────────┤
│ Plum       │
├────────────┤
│ Blueberry  │
└────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleRowNoHeaders(t *testing.T) {
	rows := [][]string{
		{"Apple", "Green", "8cm", "7/10"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(rows...)

	expected := strings.TrimSpace(`
┌───────┬───────┬─────┬──────┐
│ Apple │ Green │ 8cm │ 7/10 │
└───────┴───────┴─────┴──────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableBasicSingleColumnNoHeaders(t *testing.T) {
	rows := [][]string{
		{"Apple"},
		{"Strawberry"},
		{"Plum"},
		{"Blueberry"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(TableStyleRegular).
		Rows(rows...)

	expected := strings.TrimSpace(`
┌────────────┐
│ Apple      │
├────────────┤
│ Strawberry │
├────────────┤
│ Plum       │
├────────────┤
│ Blueberry  │
└────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}
