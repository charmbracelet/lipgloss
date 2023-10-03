package table

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"golang.org/x/exp/slices"
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

// DefaultStyles is a TableStyleFunc that returns a new Style with no attributes.
func DefaultStyles(_, _ int) lipgloss.Style {
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
		styleFunc:    DefaultStyles,
		border:       lipgloss.RoundedBorder(),
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
	hasHeaders := t.headers != nil && len(t.headers) > 0
	hasRows := t.rows != nil && len(t.rows) > 0

	if !hasHeaders && !hasRows {
		return ""
	}

	var s strings.Builder

	// Find the longest row length.
	longestRowLen := len(t.headers)
	for _, row := range t.rows {
		longestRowLen = max(longestRowLen, len(row))
	}

	// Add empty cells to the headers, until it's the same length as the longest
	// row (only if there are at headers in the first place).
	if hasHeaders {
		for i := len(t.headers); i < longestRowLen; i++ {
			t.headers = append(t.headers, "")
		}
	}

	// Initialize the widths.
	t.widths = make([]int, longestRowLen)
	t.heights = make([]int, btoi(hasHeaders)+len(t.rows))

	// The style function may affect width of the table. It's possible to set
	// the StyleFunc after the headers and rows. Update the widths for a final
	// time.
	for i, cell := range t.headers {
		t.widths[i] = max(t.widths[i], lipgloss.Width(t.style(0, i).Render(fmt.Sprint(cell))))
		t.heights[0] = max(t.heights[0], lipgloss.Height(t.style(0, i).Render(fmt.Sprint(cell))))
	}

	for r, row := range t.rows {
		for i, cell := range row {
			rendered := t.style(r+1, i).Render(fmt.Sprint(cell))
			t.heights[r+btoi(hasHeaders)] = max(t.heights[r+btoi(hasHeaders)], lipgloss.Height(rendered))
			t.widths[i] = max(t.widths[i], lipgloss.Width(rendered))
		}
	}

	// Table Resizing Logic.
	//
	// Given a user defined table width, we must ensure the table is exactly that
	// width. This must account for all borders, column, separators, and column
	// data.
	//
	// In the case where the table is narrower than the specified table width,
	// we simply expand the columns evenly to fit the width.
	// For example, a table with 3 columns takes up 50 characters total, and the
	// width specified is 80, we expand each column by 10 characters, adding 30
	// to the total width.
	//
	// In the case where the table is wider than the specified table width, we
	// _could_ simply shrink the columns evenly but this would result in data
	// being truncated (perhaps unnecessarily). The naive approach could result
	// in very poor cropping of the table data. So, instead of shrinking columns
	// evenly, we calculate the median non-whitespace length of each column, and
	// shrink the columns based on the largest median.
	//
	// For example,
	//  ┌──────┬───────────────┬──────────┐
	//  │ Name │ Age of Person │ Location │
	//  ├──────┼───────────────┼──────────┤
	//  │ Kini │ 40            │ New York │
	//  │ Eli  │ 30            │ London   │
	//  │ Iris │ 20            │ Paris    │
	//  └──────┴───────────────┴──────────┘
	//
	// Median non-whitespace length  vs column width of each column:
	//
	// Name: 4 / 5
	// Age of Person: 2 / 15
	// Location: 6 / 10
	//
	// The biggest difference is 15 - 2, so we can shrink the 2nd column by 13.

	width := t.computeWidth()

	if width < t.width && t.width > 0 {
		// Table is too narrow, expand the columns evenly until it reaches the
		// desired width.
		var i int
		for width < t.width {
			t.widths[i]++
			width++
			i = (i + 1) % len(t.widths)
		}
	} else if width > t.width && t.width > 0 {
		// Table is too wide, narrow the columns until it reaches the desired
		// width. We choose which columns to narrow based on the median
		// non-whitespace length of each column.
		columnMedians := make([]int, len(t.widths))
		for c := range t.widths {
			trimmedWidth := make([]int, len(t.rows))
			for r := range t.rows {
				renderedCell := t.style(r+btoi(hasHeaders), c).Render(fmt.Sprint(t.rows[r][c]))
				nonWhitespaceChars := lipgloss.Width(strings.TrimRight(renderedCell, " ")) + 1
				trimmedWidth[r] = nonWhitespaceChars + 1 // +1 for some padding or truncation character
			}
			columnMedians[c] = median(trimmedWidth)
		}

		for width > t.width {
			// Find the column with the largest median.
			largestDifference, largestDifferenceIndex := 0, 0
			for i, median := range columnMedians {
				difference := (t.widths[i] - median)
				if median > largestDifference {
					largestDifference = difference
					largestDifferenceIndex = i
				}
			}

			if largestDifference <= 0 {
				break
			}

			if width-largestDifference < t.width {
				largestDifference = width - t.width
			}

			width -= largestDifference
			// Set column width to the median.
			newWidth := t.widths[largestDifferenceIndex] - largestDifference
			t.widths[largestDifferenceIndex] = max(newWidth, 1)
			columnMedians[largestDifferenceIndex] = 0
		}

		// If the table is still too wide, we need to shrink it further. This
		// time, we shrink the columns evenly.
		for width > t.width {
			// Is the width unreasonably small?
			if t.width <= (len(t.widths) + (len(t.widths) * btoi(t.borderColumn)) + btoi(t.borderLeft) + btoi(t.borderRight)) {
				// Give up.
				break
			}

			for i := range t.widths {
				if t.widths[i] > 1 {
					t.widths[i]--
					width--
				}
				if width <= t.width {
					break
				}
			}
		}
	}

	// Write the top border.
	if t.borderTop {
		if t.borderLeft {
			s.WriteString(t.borderStyle.Render(t.border.TopLeft))
		}
		for i := 0; i < longestRowLen; i++ {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Top, t.widths[i])))
			if i < longestRowLen-1 && t.borderColumn {
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
			Render(runewidth.Truncate(fmt.Sprint(header), t.widths[i], "…")))
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
		height := t.heights[r+btoi(hasHeaders)]

		left := strings.Repeat(t.borderStyle.Render(t.border.Left)+"\n", height)
		right := strings.Repeat(t.borderStyle.Render(t.border.Right)+"\n", height)

		// Append empty cells to the row, until it's the same length as the
		// longest row.
		for i := len(row); i < longestRowLen; i++ {
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
				Render(runewidth.Truncate(fmt.Sprint(cell), t.widths[c]*height, "…")))

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
			for i := 0; i < longestRowLen; i++ {
				s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[i])))
				if i < longestRowLen-1 && t.borderColumn {
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
		for i := 0; i < longestRowLen; i++ {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[i])))
			if i < longestRowLen-1 && t.borderColumn {
				s.WriteString(t.borderStyle.Render(t.border.MiddleBottom))
			}
		}
		if t.borderRight {
			s.WriteString(t.borderStyle.Render(t.border.BottomRight))
		}
	}

	height := sum(t.heights) - 1 + btoi(hasHeaders) +
		btoi(t.borderHeader) + btoi(t.borderTop) + btoi(t.borderBottom) +
		len(t.rows)*btoi(t.borderRow)

	return lipgloss.NewStyle().MaxHeight(height).MaxWidth(t.width).Render(s.String())
}

// compute the width of the table in it's current configuration.
func (t *Table) computeWidth() int {
	width := sum(t.widths) + btoi(t.borderLeft) + btoi(t.borderRight)
	if t.borderColumn {
		width += len(t.widths) - 1
	}
	return width
}

// Render returns the table as a string.
func (t *Table) Render() string {
	return t.String()
}

// btoi converts a boolean to an integer, 1 if true, 0 if false.
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// max returns the greater of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// sum returns the sum of all integers in a slice.
func sum(n []int) int {
	var sum int
	for _, i := range n {
		sum += i
	}
	return sum
}

// median returns the median of a slice of integers.
func median(n []int) int {
	slices.Sort(n)

	if len(n) <= 0 {
		return 0
	}
	if len(n)%2 == 0 {
		return (n[len(n)/2-1] + n[len(n)/2]) / 2
	}
	return n[len(n)/2]
}
