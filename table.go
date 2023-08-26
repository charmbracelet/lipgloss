package lipgloss

import (
	"fmt"
	"strings"
)

// TableStyleFunc is the style function that determines the style of a Cell.
//
// It takes the row and column of the cell as an input and determines the
// lipgloss Style to use for that cell position.
type TableStyleFunc func(row, col int) Style

// NoTableStyle is a TableStyleFunc that returns a new Style with no attributes.
func NoTableStyle(_, _ int) Style {
	return NewStyle()
}

// Table is a type for rendering tables.
type Table struct {
	styleFunc TableStyleFunc
	border    Border

	borderTop    bool
	borderBottom bool
	borderLeft   bool
	borderRight  bool
	borderHeader bool

	borderStyle Style
	headers     []any
	rows        [][]any
	// widths tracks the width of each column.
	widths []int
}

// NewTable returns a new Table that can be modified through different
// attributes.
//
// By default, a table has no border, no styling, and no rows.
func NewTable() *Table {
	return &Table{
		styleFunc:    NoTableStyle,
		border:       HiddenBorder(),
		borderTop:    true,
		borderHeader: true,
		borderBottom: true,
		borderRight:  true,
		borderLeft:   true,
	}
}

// ClearRows clears the table rows.
func (t *Table) ClearRows() *Table {
	t.rows = make([][]any, 0)
	return t
}

// StyleFunc sets the style for a cell based on it's position (row, column).
func (t *Table) StyleFunc(style TableStyleFunc) *Table {
	t.styleFunc = style
	return t
}

// Rows sets the table rows.
func (t *Table) Rows(rows ...[]any) *Table {
	t.rows = rows

	if len(rows) == 0 {
		return t
	}

	row := rows[0]

	// Update the widths.
	if t.widths == nil {
		t.widths = make([]int, len(row))
	}

	// Update the widths.
	for _, row := range rows {
		for j, cell := range row {
			t.widths[j] = max(t.widths[j], Width(fmt.Sprint(cell)))
		}
	}

	return t
}

// Row appends a row of data to the table.
func (t *Table) Row(row ...any) *Table {
	t.rows = append(t.rows, row)

	// Update the widths.
	if t.widths == nil {
		t.widths = make([]int, len(row))
	}

	for i, cell := range row {
		t.widths[i] = max(t.widths[i], Width(fmt.Sprint(cell)))
	}

	return t
}

// Headers sets the table headers.
func (t *Table) Headers(headers ...any) *Table {
	t.headers = headers

	// Update the widths.
	if t.widths == nil {
		t.widths = make([]int, len(headers))
	}

	for i, cell := range headers {
		t.widths[i] = max(t.widths[i], Width(fmt.Sprint(cell)))
	}

	return t
}

// Border sets the table border.
func (t *Table) Border(border Border) *Table {
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

// BorderHeader sets the right border.
func (t *Table) BorderHeader(v bool) *Table {
	t.borderHeader = v
	return t
}

// BorderStyle sets the style for the table border.
func (t *Table) BorderStyle(style Style) *Table {
	t.borderStyle = style
	return t
}

// String returns the table as a string.
func (t *Table) String() string {
	if (t.headers == nil || len(t.headers) <= 0) && len(t.rows) == 0 {
		return ""
	}

	var s strings.Builder

	// It's possible that the styling affects the width of the table or rows.
	//
	// It's also possible that the styling function was set after the headers
	// and rows.
	//
	// So let's update the widths one last time.
	for i, cell := range t.headers {
		t.widths[i] = max(t.widths[i], Width(t.styleFunc(0, i).Render(fmt.Sprint(cell))))
	}
	for r, row := range t.rows {
		for i, cell := range row {
			t.widths[i] = max(t.widths[i], Width(t.styleFunc(r+1, i).Render(fmt.Sprint(cell))))
		}
	}

	hasHeaders := len(t.headers) > 0
	headers := t.headers

	// Write the top border.
	if t.borderTop {
		s.WriteString(t.borderStyle.Render(t.border.TopLeft))
		if !hasHeaders {
			headers = t.rows[0]
		}
		for i := 0; i < len(headers); i++ {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Top, t.widths[i])))
			if i < len(headers)-1 {
				s.WriteString(t.borderStyle.Render(t.border.TopMiddle))
			}
		}
		s.WriteString(t.borderStyle.Render(t.border.TopRight))
		s.WriteString("\n")
	}

	// Write the headers.
	if hasHeaders && t.borderLeft {
		s.WriteString(t.borderStyle.Render(t.border.Left))
	}
	for i, header := range t.headers {
		s.WriteString(t.styleFunc(0, i).
			UnsetPaddingTop().
			UnsetPaddingBottom().
			UnsetMarginTop().
			UnsetMarginBottom().
			Width(t.widths[i]).
			MaxWidth(t.widths[i]).
			Render(fmt.Sprint(header)))
		if i < len(t.headers)-1 {
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
			if i < len(t.headers)-1 {
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
		if t.borderLeft {
			s.WriteString(t.borderStyle.Render(t.border.Left))
		}
		for c, cell := range row {
			s.WriteString(t.styleFunc(r+1, c).
				UnsetPaddingTop().
				UnsetPaddingBottom().
				UnsetMarginTop().
				UnsetMarginBottom().
				Width(t.widths[c]).
				MaxWidth(t.widths[c]).
				Render(fmt.Sprint(cell)))
			if c < len(row)-1 {
				s.WriteString(t.borderStyle.Render(t.border.Left))
			}
		}
		if t.borderRight {
			s.WriteString(t.borderStyle.Render(t.border.Right))
		}
		s.WriteString("\n")
	}

	// Write the bottom border.
	if t.borderBottom {
		s.WriteString(t.borderStyle.Render(t.border.BottomLeft))
		for i := 0; i < len(headers); i++ {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[i])))
			if i < len(headers)-1 {
				s.WriteString(t.borderStyle.Render(t.border.MiddleBottom))
			}
		}
		s.WriteString(t.borderStyle.Render(t.border.BottomRight))
	}

	return s.String()
}

// Render returns the table as a string.
func (t *Table) Render() string {
	return t.String()
}
