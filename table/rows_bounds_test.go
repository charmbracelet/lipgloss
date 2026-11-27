package table

import "testing"

func TestStringDataAtBounds(t *testing.T) {
	data := NewStringData([]string{"a", "b"}, []string{"c"}, nil)
	for _, tc := range []struct {
		name string
		row  int
		cell int
		want string
	}{
		{"first cell", 0, 0, "a"},
		{"second cell", 0, 1, "b"},
		{"last row with data", 1, 0, "c"},
		{"missing cell", 1, 1, ""},
		{"empty row", 2, 0, ""},
		{"row past end", 3, 0, ""},
		{"cell past end", 0, 2, ""},
		{"negative row", -1, 0, ""},
		{"negative cell", 0, -1, ""},
		{"negative row and cell", -1, -1, ""},
		{"negative cell in empty row", 2, -1, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := data.At(tc.row, tc.cell); got != tc.want {
				t.Errorf("At(%d, %d) = %q, want %q", tc.row, tc.cell, got, tc.want)
			}
		})
	}
}
