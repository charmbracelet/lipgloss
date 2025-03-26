package table

import (
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const (
	// HeaderRow denotes the header's row index used when rendering headers. Use
	// this value when looking to customize header styles in StyleFunc.
	HeaderRow int = -1
	// FirstCol stores the index of the first column in a table.
	FirstCol int = 0
	// NoBorder is used to fill a gap with repeated spaces instead of drawing a border.
	NoBorder string = " "
	// EmptyCell is used to clear a cell value.
	EmptyCell string = ""
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

	borderTop          bool
	borderBottom       bool
	borderLeft         bool
	borderRight        bool
	borderHeader       bool
	borderColumn       bool
	borderRow          bool
	borderMergeColumns []int

	borderStyle lipgloss.Style
	headers     []string
	data        Data

	width           int
	height          int
	useManualHeight bool
	offset          int
	wrap            bool

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
		wrap:         true,
		data:         NewStringData(),
	}
}

// ClearRows clears the table rows.
func (t *Table) ClearRows() *Table {
	t.data = NewStringData()
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

// Data sets the table data.
func (t *Table) Data(data Data) *Table {
	t.data = data
	return t
}

// Rows appends rows to the table data.
func (t *Table) Rows(rows ...[]string) *Table {
	for _, row := range rows {
		switch t.data.(type) {
		case *StringData:
			t.data.(*StringData).Append(row)
		}
	}
	return t
}

// Row appends a row to the table data.
func (t *Table) Row(row ...string) *Table {
	switch t.data.(type) {
	case *StringData:
		t.data.(*StringData).Append(row)
	}
	return t
}

// Headers sets the table headers.
func (t *Table) Headers(headers ...string) *Table {
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

// BorderColumnMerge sets the columns where adjacent vertical cells that contain
// the same data get merged into a single cell.
// Prerequesite: BorderRow and BorderColumn must be set to true.
func (t *Table) BorderMergeColumns(cols ...int) *Table {
	if t.borderRow && t.borderColumn {
		t.borderMergeColumns = cols
	}
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

// Height sets the table height.
func (t *Table) Height(h int) *Table {
	t.height = h
	t.useManualHeight = true
	return t
}

// Offset sets the table rendering offset.
//
// Warning: you may declare Offset only after setting Rows. Otherwise it will be
// ignored.
func (t *Table) Offset(o int) *Table {
	t.offset = o
	return t
}

// Wrap dictates whether or not the table content should wrap.
func (t *Table) Wrap(w bool) *Table {
	t.wrap = w
	return t
}

// String returns the table as a string.
func (t *Table) String() string {
	hasHeaders := len(t.headers) > 0
	hasRows := t.data != nil && t.data.Rows() > 0

	if !hasHeaders && !hasRows {
		return ""
	}

	// Add empty cells to the headers, until it's the same length as the longest
	// row (only if there are at headers in the first place).
	if hasHeaders {
		for i := len(t.headers); i < t.data.Columns(); i++ {
			t.headers = append(t.headers, "")
		}
	}

	// Do all the sizing calculations for width and height.
	t.resize()

	var sb strings.Builder

	if t.borderTop {
		sb.WriteString(t.constructTopBorder())
		sb.WriteString("\n")
	}

	if hasHeaders {
		sb.WriteString(t.constructHeaders())
		sb.WriteString("\n")
	}

	var bottom string
	if t.borderBottom {
		bottom = t.constructBottomBorder()
	}

	// If there are no data rows render nothing.
	if t.data.Rows() > 0 {
		switch {
		case t.useManualHeight:
			// The height of the top border. Subtract 1 for the newline.
			topHeight := lipgloss.Height(sb.String()) - 1
			availableLines := t.height - (topHeight + lipgloss.Height(bottom))

			// if the height is larger than the number of rows, use the number
			// of rows.
			if availableLines > t.data.Rows() {
				availableLines = t.data.Rows()
			}
			sb.WriteString(t.constructRows(availableLines))

		default:
			for r := t.offset; r < t.data.Rows(); r++ {
				sb.WriteString(t.constructRow(r, false))
			}
		}
	}

	sb.WriteString(bottom)

	return lipgloss.NewStyle().
		MaxHeight(t.computeHeight()).
		MaxWidth(t.width).
		Render(sb.String())
}

// computeHeight computes the height of the table in it's current configuration.
func (t *Table) computeHeight() int {
	hasHeaders := len(t.headers) > 0
	return sum(t.heights) - 1 + btoi(hasHeaders) +
		btoi(t.borderTop) + btoi(t.borderBottom) +
		btoi(t.borderHeader) + t.data.Rows()*btoi(t.borderRow)
}

// Render returns the table as a string.
func (t *Table) Render() string {
	return t.String()
}

// constructTopBorder constructs the top border for the table given it's current
// border configuration and data.
func (t *Table) constructTopBorder() string {
	var s strings.Builder
	if t.borderLeft {
		s.WriteString(t.borderStyle.Render(t.border.TopLeft))
	}
	for i := 0; i < len(t.widths); i++ {
		s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Top, t.widths[i])))
		if i < len(t.widths)-1 && t.borderColumn {
			s.WriteString(t.borderStyle.Render(t.border.MiddleTop))
		}
	}
	if t.borderRight {
		s.WriteString(t.borderStyle.Render(t.border.TopRight))
	}
	return s.String()
}

// constructBottomBorder constructs the bottom border for the table given it's current
// border configuration and data.
func (t *Table) constructBottomBorder() string {
	var s strings.Builder
	if t.borderLeft {
		s.WriteString(t.borderStyle.Render(t.border.BottomLeft))
	}
	for i := 0; i < len(t.widths); i++ {
		s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[i])))
		if i < len(t.widths)-1 && t.borderColumn {
			s.WriteString(t.borderStyle.Render(t.border.MiddleBottom))
		}
	}
	if t.borderRight {
		s.WriteString(t.borderStyle.Render(t.border.BottomRight))
	}
	return s.String()
}

// constructHeaders constructs the headers for the table given it's current
// header configuration and data.
func (t *Table) constructHeaders() string {
	var s strings.Builder
	if t.borderLeft {
		s.WriteString(t.borderStyle.Render(t.border.Left))
	}
	for i, header := range t.headers {
		s.WriteString(t.style(HeaderRow, i).
			MaxHeight(1).
			Width(t.widths[i]).
			MaxWidth(t.widths[i]).
			Render(t.truncateCell(header, -1, i)))
		if i < len(t.headers)-1 && t.borderColumn {
			s.WriteString(t.borderStyle.Render(t.border.Left))
		}
	}
	if t.borderHeader {
		if t.borderRight {
			s.WriteString(t.borderStyle.Render(t.border.Right))
		}
		s.WriteString("\n")
		if t.borderLeft {
			s.WriteString(t.borderStyle.Render(t.border.MiddleLeft))
		}
		for i := 0; i < len(t.headers); i++ {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Top, t.widths[i])))
			if i < len(t.headers)-1 && t.borderColumn {
				s.WriteString(t.borderStyle.Render(t.border.Middle))
			}
		}
		if t.borderRight {
			s.WriteString(t.borderStyle.Render(t.border.MiddleRight))
		}
	}
	if t.borderRight && !t.borderHeader {
		s.WriteString(t.borderStyle.Render(t.border.Right))
	}
	return s.String()
}

func (t *Table) constructRows(availableLines int) string {
	var sb strings.Builder

	// The number of rows to render after removing the offset.
	offsetRowCount := t.data.Rows() - t.offset

	// The number of rows to render. We always render at least one row.
	rowsToRender := availableLines
	rowsToRender = max(rowsToRender, 1)

	// Check if we need to render an overflow row.
	needsOverflow := rowsToRender < offsetRowCount

	// only use the offset as the starting value if there is overflow.
	rowIdx := t.offset
	if !needsOverflow {
		// if there is no overflow, just render to the height of the table
		// check there's enough content to fill the table
		rowIdx = t.data.Rows() - rowsToRender
	}
	for rowsToRender > 0 && rowIdx < t.data.Rows() {
		// Whenever the height is too small to render all rows, the bottom row will be an overflow row (ellipsis).
		isOverflow := needsOverflow && rowsToRender == 1

		sb.WriteString(t.constructRow(rowIdx, isOverflow))

		rowIdx++
		rowsToRender--
	}
	return sb.String()
}

func (t *Table) constructRow(row int, isOverflow bool) string {
	var s strings.Builder

	hasHeaders := len(t.headers) > 0
	height := t.heights[row+btoi(hasHeaders)]
	if isOverflow {
		height = 1
	}

	var cells []string
	left := strings.Repeat(t.borderStyle.Render(t.border.Left)+"\n", height)
	if t.borderLeft {
		cells = append(cells, left)
	}

	verticalCellsEqual := func(row, col int) bool { return t.data.At(row-1, col) == t.data.At(row, col) }
	modifyCell := func(col int, cell string) string {
		if row > 0 && slices.Contains(t.borderMergeColumns, col) && verticalCellsEqual(row, col) {
			return EmptyCell // if merging column, clear cells with identical data excluding the first occurrence
		}
		return cell
	}

	for col := 0; col < t.data.Columns(); col++ {
		cell := "…"
		if !isOverflow {
			cell = t.data.At(row, col)
		}

		cellStyle := t.style(row, col)
		if !t.wrap {
			cell = t.truncateCell(cell, row, col)
		}
		cells = append(cells, cellStyle.
			// Account for the margins in the cell sizing.
			Height(height-cellStyle.GetVerticalMargins()).
			MaxHeight(height).
			Width(t.widths[col]-cellStyle.GetHorizontalMargins()).
			MaxWidth(t.widths[col]).
			Render(modifyCell(col, cell)))

		if col < t.data.Columns()-1 && t.borderColumn {
			cells = append(cells, left)
		}
	}

	if t.borderRight {
		right := strings.Repeat(t.borderStyle.Render(t.border.Right)+"\n", height)
		cells = append(cells, right)
	}

	for i, cell := range cells {
		cells[i] = strings.TrimRight(cell, "\n")
	}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cells...) + "\n")

	t.drawRowBorders(&s, row)

	return s.String()
}

// Draws the borders separating rows for singular row.
func (t *Table) drawRowBorders(s *strings.Builder, row int) {
	if t.borderRow && row < t.data.Rows()-1 {
		t.drawLeftmostBorder(s, row)
		t.drawMiddleBorders(s, row)
		t.drawRightmostBorders(s, row)
	}
}

// Draws the leftmost border intersection for a singular row.
func (t *Table) drawLeftmostBorder(s *strings.Builder, row int) {
	if slices.Contains(t.borderMergeColumns, FirstCol) && t.verticalCellsEqual(row, FirstCol) {
		s.WriteString(t.borderStyle.Render(t.border.Left))
	} else {
		s.WriteString(t.borderStyle.Render(t.border.MiddleLeft))
	}
}

// Draws all the horizontal and intersection borders for a singular row, exluding the last column.
func (t *Table) drawMiddleBorders(s *strings.Builder, row int) {
	for col := 0; col < len(t.widths)-1; col++ {
		if !t.borderColumn {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[col])))
		} else if slices.Contains(t.borderMergeColumns, col) && slices.Contains(t.borderMergeColumns, col+1) &&
			t.verticalCellsEqual(row, col) && t.verticalCellsEqual(row, col+1) {
			s.WriteString(t.borderStyle.Render(strings.Repeat(NoBorder, t.widths[col])))
			s.WriteString(t.borderStyle.Render(t.border.Left)) // |
		} else if slices.Contains(t.borderMergeColumns, col) && t.verticalCellsEqual(row, col) {
			s.WriteString(t.borderStyle.Render(strings.Repeat(NoBorder, t.widths[col])))
			s.WriteString(t.borderStyle.Render(t.border.MiddleLeft)) // |-
		} else if slices.Contains(t.borderMergeColumns, col+1) && t.verticalCellsEqual(row, col+1) {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[col])))
			s.WriteString(t.borderStyle.Render(t.border.MiddleRight)) // -|
		} else {
			s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[col])))
			s.WriteString(t.borderStyle.Render(t.border.Middle)) // -|-
		}
	}
}

// Draws the rightmost column's horizontal and intersection border for a singular row.
func (t *Table) drawRightmostBorders(s *strings.Builder, row int) {
	lastCol := len(t.widths) - 1
	if slices.Contains(t.borderMergeColumns, lastCol) && t.verticalCellsEqual(row, lastCol) {
		s.WriteString(t.borderStyle.Render(strings.Repeat(NoBorder, t.widths[lastCol])))
		s.WriteString(t.borderStyle.Render(t.border.Right) + "\n") // |
	} else {
		s.WriteString(t.borderStyle.Render(strings.Repeat(t.border.Bottom, t.widths[lastCol])))
		s.WriteString(t.borderStyle.Render(t.border.MiddleRight) + "\n") // -|
	}
}

// Returns true if the cell at (row, col) and the cell below it contain the same data.
func (t *Table) verticalCellsEqual(row, col int) bool {
	return t.data.At(row, col) == t.data.At(row+1, col)
}

func (t *Table) truncateCell(cell string, rowIndex, colIndex int) string {
	hasHeaders := len(t.headers) > 0
	height := t.heights[rowIndex+btoi(hasHeaders)]
	cellWidth := t.widths[colIndex]
	cellStyle := t.style(rowIndex, colIndex)

	length := (cellWidth * height) - cellStyle.GetHorizontalPadding() - cellStyle.GetHorizontalMargins()
	return ansi.Truncate(cell, length, "…")
}
