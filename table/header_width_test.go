package table

import (
	"strings"
	"testing"
)

// A column whose header is wider than its cells must not be squeezed below
// its header width: the header has to stay readable, otherwise the column
// disappears from the table.
//
// This is the table from charmbracelet/glow#941, where the "Status" column
// ended up one cell wide and its header was truncated away.
func TestColumnKeepsHeaderWidthWhenShrinking(t *testing.T) {
	tbl := New().
		Headers("Feature", "Status", "identifier", "Notes").
		Rows(
			[]string{
				"feature one with a moderately long name",
				"Y",
				"module_alpha::test_long_identifier_name_first",
				"comment one. Some explanatory text describing behavior here.",
			},
			[]string{
				"feature two with another moderately long name",
				"Y",
				"module_alpha::test_long_identifier_name_second",
				"comment two. More explanatory text describing behavior here.",
			},
		).
		Width(80)

	out := tbl.Render()
	t.Log("\n" + out)

	if !strings.Contains(out, "Status") {
		t.Errorf("expected the Status header to be visible, got:\n%s", out)
	}
	if !strings.Contains(out, "Y") {
		t.Errorf("expected the Status cells to be visible, got:\n%s", out)
	}
}
