package table

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StyleFunc is the style function that determines the style of a Cell.
//
// It takes the row and column of the cell as an input and determines the
// lipgloss Style to use for that cell position.
//
// Example:
//
//	t := table.New().
//	    Headers("Name", "Age").
//	    Row("Kini", 4).
//	    Row("Eli", 1).
//	    Row("Iris", 102).
//	    StyleFunc(func(row, col int) lipgloss.Style {
//	        switch {
//	           case row == 0:
//	               return HeaderStyle
//	           case row%2 == 0:
//	               return EvenRowStyle
//	           default:
//	               return OddRowStyle
//	           }
//	    })
type StyleFunc func(row, col int) lipgloss.Style

// NoTableStyle is a TableStyleFunc that returns a new Style with no attributes.
func NoTableStyle(_, _ int) lipgloss.Style {
	return lipgloss.NewStyle()
}

// Table is a type for rendering tables.
type Table struct {
	styleFunc StyleFunc
	border    lipgloss.Border

	borderTop    bool
	borderBottom bool
	borderLeft   bool
	borderRight  bool
	borderHeader bool
	borderColumn bool
	borderRow    bool

	borderStyle lipgloss.Style
	headers     []any
	rows        [][]any

	width int

	// widths tracks the width of each column.
	widths []int

	// heights tracks the height of each row.
	heights []int
}

// New returns a new Table that can be modified through different
// attributes.
//
// By default, a table has no border, no styling, and no rows.
func New() *Table {
	return &Table{
		styleFunc:    NoTableStyle,
		border:       lipgloss.HiddenBorder(),
		borderBottom: true,
		borderColumn: true,
		borderHeader: true,
		borderLeft:   true,
		borderRight:  true,
		borderTop:    true,
	}
}

// ClearRows clears the table rows.
func (t *Table) ClearRows() *Table {
	t.rows = make([][]any, 0)
	return t
}

// StyleFunc sets the style for a cell based on it's position (row, column).
func (t *Table) StyleFunc(style StyleFunc) *Table {
	t.styleFunc = style
	return t
}

// style returns the style for a cell based on it's position (row, column).
func (t *Table) style(row, col int) lipgloss.Style {
	if t.styleFunc == nil {
		return lipgloss.NewStyle()
	}
	return t.styleFunc(row, col)
}

// Rows sets the table rows.
func (t *Table) Rows(rows ...[]any) *Table {
	t.rows = rows
	return t
}

// Row appends a row of data to the table.
func (t *Table) Row(row ...any) *Table {
	t.rows = append(t.rows, row)
	return t
}

// Headers sets the table headers.
func (t *Table) Headers(headers ...any) *Table {
	t.headers = headers
	return t
}

// Border sets the table border.
func (t *Table) Border(border lipgloss.Border) *Table {
	t.border = border
	return t
}

// BorderTop sets the top border.
func (t *Table) BorderTop(v bool) *Table {
	t.borderTop = v
	return t
}

// BorderBottom sets the bottom border.
func (t *Table) BorderBottom(v bool) *Table {
	t.borderBottom = v
	return t
}

// BorderLeft sets the left border.
func (t *Table) BorderLeft(v bool) *Table {
	t.borderLeft = v
	return t
}

// BorderRight sets the right border.
func (t *Table) BorderRight(v bool) *Table {
	t.borderRight = v
	return t
}

// BorderHeader sets the header separator border.
func (t *Table) BorderHeader(v bool) *Table {
	t.borderHeader = v
	return t
}

// BorderColumn sets the column border separator.
func (t *Table) BorderColumn(v bool) *Table {
	t.borderColumn = v
	return t
}

// BorderRow sets the row border separator.
func (t *Table) BorderRow(v bool) *Table {
	t.borderRow = v
	return t
}

// BorderStyle sets the style for the table border.
func (t *Table) BorderStyle(style lipgloss.Style) *Table {
	t.borderStyle = style
	return t
}

// Width sets the table width, this auto-sizes the columns to fit the width by
// either expanding or contracting the widths of each column as a best effort
// approach.
func (t *Table) Width(w int) *Table {
	t.width = w
	return t
}

// String returns the table as a string.
func (t *Table) String() string {
	if (t.headers == nil || len(t.headers) <= 0) && len(t.rows) == 0 {
		return ""
	}

	var s strings.Builder

	hasHeaders := len(t.headers) > 0
	longestRow := t.headers
	if !hasHeaders {
		longestRow = t.rows[0]
	}

	// Find longest row.
	for _, row := range t.rows {
		if len(row) > len(longestRow) {
			longestRow = row
		}
	}

	// Add empty cells to the headers, until it's the same length as the longest
	// row (only if there are at headers in the first place).
	if hasHeaders {
		for i := len(t.headers); i < len(longestRow); i++ {
			t.headers = append(t.headers, "")
		}
	}

	// Initialize the widths.
	t.widths = make([]int, len(longestRow))
	t.heights = make([]int, boolToInt(hasHeaders)+len(t.rows))

	// It's possible that the styling affects the width of the table or rows.
	//
	// It's also possible that the styling function was set after the headers
	// and rows.
	//
	// So let's update the widths one last time.
	for i, cell := range t.headers {
		t.widths[i] = max(t.widths[i], lipgloss.Width(t.style(0, i).Render(fmt.Sprint(cell))))
		t.heights[0] = max(t.heights[0], lipgloss.Height(t.style(0, i).Render(fmt.Sprint(cell))))
	}
	for r, row := range t.rows {
		for i, cell := range row {
			rendered := t.style(r+1, i).Render(fmt.Sprint(cell))
			t.heights[r+boolToInt(hasHeaders)] = max(t.heights[r+boolToInt(hasHeaders)], lipgloss.Height(rendered))
			t.widths[i] = max(t.widths[i], lipgloss.Width(rendered))
		}
	}

	tableWidth := sum(t.widths)
	if t.borderColumn {
		tableWidth += (len(t.widths) - 1)
	}
	tableWidth += boolToInt(t.borderLeft)
	tableWidth += boolToInt(t.borderRight)

	if tableWidth < t.width && t.width > 0 {
		// The table is too narrow, so we need to expand it.
		for tableWidth < t.width {
			// Add an equal amount to each column.
			for i := range t.widths {
				t.widths[i]++
				tableWidth++
			}
		}
	} else if tableWidth > t.width && t.width > 0 {
		// The table is too wide, so we need to shrink it.
		for tableWidth > t.width {
			// Subtract an equal amount from each column.
			for i := range t.widths {
				t.widths[i]--
				tableWidth--
			}
		}
	}

	// Write the top border.
	if t.borderTop {
		if t.borderLeft {
			s.WriteString(t.borderStyle.Render(t.border.TopLeft))
		}
		for i := 0; i < len(longestRow); i++ {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Top, t.widths[i])))
			if i < len(longestRow)-1 && t.borderColumn {
				s.WriteString(t.borderStyle.Render(t.border.MiddleTop))
			}
		}
		if t.borderRight {
			s.WriteString(t.borderStyle.Render(t.border.TopRight))
		}
		s.WriteString("\n")
	}

	// Write the headers.
	if hasHeaders && t.borderLeft {
		s.WriteString(t.borderStyle.Render(t.border.Left))
	}
	for i, header := range t.headers {
		s.WriteString(t.style(0, i).
			MaxHeight(1).
			Width(t.widths[i]).
			MaxWidth(t.widths[i]).
			Render(fmt.Sprint(header)))
		if i < len(t.headers)-1 && t.borderColumn {
			s.WriteString(t.borderStyle.Render(t.border.Left))
		}
	}
	if hasHeaders && t.borderHeader {
		if t.borderRight {
			s.WriteString(t.borderStyle.Render(t.border.Right))
		}
		s.WriteString("\n")
		if t.borderLeft {
			s.WriteString(t.borderStyle.Render(t.border.MiddleLeft))
		}
	}
	if t.borderHeader {
		for i := 0; i < len(t.headers); i++ {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Top, t.widths[i])))
			if i < len(t.headers)-1 && t.borderColumn {
				s.WriteString(t.borderStyle.Render(t.border.Middle))
			}
		}
	}
	if hasHeaders && t.borderRight {
		if t.borderHeader {
			s.WriteString(t.borderStyle.Render(t.border.MiddleRight))
		} else {
			s.WriteString(t.borderStyle.Render(t.border.Right))
		}
	}
	if hasHeaders {
		s.WriteString("\n")
	}

	// Write the data.
	for r, row := range t.rows {
		height := t.heights[r+boolToInt(hasHeaders)]

		left := strings.Repeat(t.borderStyle.Render(t.border.Left)+"\n", height)
		right := strings.Repeat(t.borderStyle.Render(t.border.Right)+"\n", height)

		// Append empty cells to the row, until it's the same length as the
		// longest row.
		for i := len(row); i < len(longestRow); i++ {
			row = append(row, "")
		}

		var cells []string
		if t.borderLeft {
			cells = append(cells, left)
		}

		for c, cell := range row {
			cells = append(cells, t.style(r+1, c).
				Height(height).
				MaxHeight(height).
				Width(t.widths[c]).
				MaxWidth(t.widths[c]).
				Render(fmt.Sprint(cell)))

			if c < len(row)-1 && t.borderColumn {
				cells = append(cells, left)
			}
		}

		if t.borderRight {
			cells = append(cells, right)
		}

		for i, cell := range cells {
			cells[i] = strings.TrimRight(cell, "\n")
		}

		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cells...) + "\n")

		if t.borderRow && r < len(t.rows)-1 {
			s.WriteString(t.borderStyle.Render(t.border.MiddleLeft))
			for i := 0; i < len(longestRow); i++ {
				s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[i])))
				if i < len(longestRow)-1 && t.borderColumn {
					s.WriteString(t.borderStyle.Render(t.border.Middle))
				}
			}
			s.WriteString(t.borderStyle.Render(t.border.MiddleRight) + "\n")
		}
	}

	// Write the bottom border.
	if t.borderBottom {
		if t.borderLeft {
			s.WriteString(t.borderStyle.Render(t.border.BottomLeft))
		}
		for i := 0; i < len(longestRow); i++ {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[i])))
			if i < len(longestRow)-1 && t.borderColumn {
				s.WriteString(t.borderStyle.Render(t.border.MiddleBottom))
			}
		}
		if t.borderRight {
			s.WriteString(t.borderStyle.Render(t.border.BottomRight))
		}
	}

	return s.String()
}

// Render returns the table as a string.
func (t *Table) Render() string {
	return t.String()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sum(n []int) int {
	var sum int
	for _, i := range n {
		sum += i
	}
	return sum
}
