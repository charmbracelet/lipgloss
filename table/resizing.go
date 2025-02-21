package table

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (t *Table) resize() {
	strData, ok := t.data.(*StringData)
	if !ok {
		return
	}
	r := newResizer(t.width, t.height, t.headers, strData.rows)
	r.wrap = t.wrap
	r.borderColumn = t.borderColumn
	r.yPaddings = make([][]int, len(r.allRows))

	allRows := append([][]string{t.headers}, strData.rows...)
	r.rowHeights = r.defaultRowHeigths()

	for i, row := range allRows {
		r.yPaddings[i] = make([]int, len(row))

		for j := range row {
			column := &r.columns[j]
			style := t.styleFunc(i, j)

			topMargin, rightMargin, bottomMargin, leftMargin := style.GetMargin()
			topPadding, rightPadding, bottomPadding, leftPadding := style.GetPadding()

			totalHorizontalPadding := leftMargin + rightMargin + leftPadding + rightPadding
			column.xPadding = max(column.xPadding, totalHorizontalPadding)
			column.fixedWidth = max(column.fixedWidth, style.GetWidth())

			r.rowHeights[i] = max(r.rowHeights[i], style.GetHeight())

			totalVerticalPadding := topMargin + bottomMargin + topPadding + bottomPadding
			r.yPaddings[i][j] = totalVerticalPadding
		}
	}

	// A table width wasn't specified. In this case, detect according to
	// content width.
	if r.tableWidth <= 0 {
		r.tableWidth = r.detectTableWidth()
	}

	t.widths, t.heights = r.optimizedWidths()
}

type resizerColumn struct {
	index      int
	min        int
	max        int
	median     int
	rows       [][]string
	xPadding   int // horizontal padding
	fixedWidth int
}

type resizer struct {
	tableWidth  int
	tableHeight int
	headers     []string
	allRows     [][]string
	rowHeights  []int
	columns     []resizerColumn

	wrap         bool
	borderColumn bool
	yPaddings    [][]int // vertical paddings
}

func newResizer(tableWidth, tableHeight int, headers []string, rows [][]string) *resizer {
	r := &resizer{
		tableWidth:  tableWidth,
		tableHeight: tableHeight,
		headers:     headers,
	}

	r.allRows = append([][]string{headers}, rows...)

	for _, row := range r.allRows {
		for i, cell := range row {
			cellLen := lipgloss.Width(cell)

			// Header or first row. Just add as is.
			if len(r.columns) <= i {
				r.columns = append(r.columns, resizerColumn{
					index:  i, // TODO(@andreynering): Save as `-1` if is header?
					min:    cellLen,
					max:    cellLen,
					median: cellLen,
				})
				continue
			}

			r.columns[i].rows = append(r.columns[i].rows, row)
			r.columns[i].min = min(r.columns[i].min, cellLen)
			r.columns[i].max = max(r.columns[i].max, cellLen)
		}
	}
	for j := range r.columns {
		widths := make([]int, len(r.columns[j].rows))
		for i, row := range r.columns[j].rows {
			widths[i] = lipgloss.Width(row[j])
		}
		r.columns[j].median = median(widths)
	}

	return r
}

func (r *resizer) optimizedWidths() (colWidths, rowHeights []int) {
	if r.maxTotal() <= r.tableWidth {
		return r.expandTableWidth()
	}
	return r.shrinkTableWidth()
}

func (r *resizer) detectTableWidth() int {
	return r.maxCharCount() + r.totalPadding() + r.totalBorder()
}

func (r *resizer) expandTableWidth() (colWidths, rowHeights []int) {
	colWidths = r.maxColumnWidths()

	for {
		totalWidth := sum(colWidths) + r.totalBorder()
		if totalWidth >= r.tableWidth {
			break
		}

		shorterColumnIndex := 0
		shorterColumnWidth := math.MaxInt32

		for i, width := range colWidths {
			if width < shorterColumnWidth {
				shorterColumnWidth = width
				shorterColumnIndex = i
			}
		}

		colWidths[shorterColumnIndex]++
	}

	rowHeights = r.expandRowHeigths(colWidths)
	return
}

func (r *resizer) shrinkTableWidth() (colWidths, rowHeights []int) {
	colWidths = r.maxColumnWidths()

	// Cut width of columns that are way too big.
	shrinkBiggestColumns := func(veryBigOnly bool) {
		for {
			totalWidth := sum(colWidths) + r.totalBorder()
			if totalWidth <= r.tableWidth {
				break
			}

			bigColumnIndex := -math.MaxInt32
			bigColumnWidth := -math.MaxInt32

			for j, width := range colWidths {
				if veryBigOnly {
					if width >= (r.tableWidth/2) && width > bigColumnWidth {
						bigColumnWidth = width
						bigColumnIndex = j
					}
				} else {
					if width > bigColumnWidth {
						bigColumnWidth = width
						bigColumnIndex = j
					}
				}
			}

			if bigColumnIndex < 0 || colWidths[bigColumnIndex] == 0 {
				break
			}
			colWidths[bigColumnIndex]--
		}
	}

	// Cut width of columns that differ the most from the median.
	shrinkToMedian := func() {
		for {
			totalWidth := sum(colWidths) + r.totalBorder()
			if totalWidth <= r.tableWidth {
				break
			}

			biggestDiffToMedian := -math.MaxInt32
			biggestDiffToMedianIndex := -math.MaxInt32

			for j, width := range colWidths {
				diffToMedian := width - r.columns[j].median
				if diffToMedian > 0 && diffToMedian > biggestDiffToMedian {
					biggestDiffToMedian = diffToMedian
					biggestDiffToMedianIndex = j
				}
			}

			if biggestDiffToMedianIndex <= 0 || colWidths[biggestDiffToMedianIndex] == 0 {
				break
			}
			colWidths[biggestDiffToMedianIndex]--
		}
	}

	shrinkBiggestColumns(true)
	shrinkToMedian()
	shrinkBiggestColumns(false)

	rowHeights = r.expandRowHeigths(colWidths)
	return
}

func (r *resizer) expandRowHeigths(colWidths []int) (rowHeights []int) {
	rowHeights = r.defaultRowHeigths()
	if !r.wrap {
		return rowHeights
	}
	for i, row := range r.allRows {
		for j, cell := range row {
			height := r.detectContentHeight(cell, colWidths[j]-r.xPaddingForCol(j)) + r.xPaddingForCell(i, j)
			if height > rowHeights[i] {
				rowHeights[i] = height
			}
		}
	}
	return
}

func (r *resizer) defaultRowHeigths() (rowHeights []int) {
	rowHeights = make([]int, len(r.allRows))
	for i := range rowHeights {
		if i < len(r.rowHeights) {
			rowHeights[i] = r.rowHeights[i]
		}
		if rowHeights[i] < 1 {
			rowHeights[i] = 1
		}
	}
	return
}

func (r *resizer) maxColumnWidths() []int {
	maxColumnWidths := make([]int, len(r.columns))
	for i, col := range r.columns {
		if col.fixedWidth > 0 {
			maxColumnWidths[i] = col.fixedWidth
		} else {
			maxColumnWidths[i] = col.max + r.xPaddingForCol(col.index)
		}
	}
	return maxColumnWidths
}

func (r *resizer) columnCount() int {
	return len(r.columns)
}

func (r *resizer) maxCharCount() int {
	var count int
	for _, col := range r.columns {
		if col.fixedWidth > 0 {
			count += col.fixedWidth - r.xPaddingForCol(col.index)
		} else {
			count += col.max
		}
	}
	return count
}

func (r *resizer) maxTotal() (maxTotal int) {
	for j, column := range r.columns {
		if column.fixedWidth > 0 {
			maxTotal += column.fixedWidth
		} else {
			maxTotal += column.max + r.xPaddingForCol(j)
		}
	}
	return
}

func (r *resizer) totalPadding() (totalPadding int) {
	for _, col := range r.columns {
		totalPadding += col.xPadding
	}
	return
}

func (r *resizer) xPaddingForCol(j int) int {
	if j >= len(r.columns) {
		return 0
	}
	return r.columns[j].xPadding
}

func (r *resizer) xPaddingForCell(i, j int) int {
	if i >= len(r.yPaddings) || j >= len(r.yPaddings[i]) {
		return 0
	}
	return r.yPaddings[i][j]
}

func (r *resizer) totalBorder() int {
	return (r.columnCount() * r.borderPerCell()) + r.extraBorder()
}

func (r *resizer) borderPerCell() int {
	if r.borderColumn {
		return 1
	}
	return 0
}

func (r *resizer) extraBorder() int {
	if r.borderColumn {
		return 1
	}
	return 0
}

func (r *resizer) detectContentHeight(content string, width int) (height int) {
	if width == 0 {
		return 1
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")
	for _, line := range strings.Split(content, "\n") {
		height += strings.Count(ansi.Wrap(line, width, ""), "\n") + 1
	}
	return
}
