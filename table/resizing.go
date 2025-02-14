package table

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

func (t *Table) resize() {
	strData, ok := t.data.(*StringData)
	if !ok {
		return
	}
	r := newResizer(t.width, t.height, t.headers, strData.rows)
	r.borderColumn = t.borderColumn

	allRows := append([][]string{t.headers}, strData.rows...)

	for i, row := range allRows {
		for j := range row {
			column := &r.columns[j]
			style := t.styleFunc(i, j)

			_, rightMargin, _, leftMargin := style.GetMargin()
			_, rightPadding, _, leftPadding := style.GetPadding()

			totalPadding := leftMargin + rightMargin + leftPadding + rightPadding
			column.padding = max(column.padding, totalPadding)
			column.fixedWidth = max(column.fixedWidth, style.GetWidth())
		}
	}

	// A table width wasn't specified. In this case, detect according to
	// content width.
	if r.tableWidth <= 0 {
		r.tableWidth = r.detectTableWidth()
	}

	t.widths = r.optimizedWidths()
}

type resizerColumn struct {
	index      int
	min        int
	max        int
	median     int
	rows       [][]string
	padding    int
	fixedWidth int
}

type resizer struct {
	tableWidth  int
	tableHeight int
	columns     []resizerColumn

	borderColumn bool
}

func newResizer(tableWidth, tableHeight int, headers []string, rows [][]string) *resizer {
	r := &resizer{
		tableWidth:  tableWidth,
		tableHeight: tableHeight,
	}

	allRows := append([][]string{headers}, rows...)

	for _, row := range allRows {
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
	for i := range r.columns {
		widths := make([]int, len(r.columns[i].rows))
		for j, row := range r.columns[i].rows {
			widths[j] = lipgloss.Width(row[i])
		}
		r.columns[i].median = median(widths)
	}

	return r
}

func (r *resizer) optimizedWidths() []int {
	if r.maxTotal() <= r.tableWidth {
		return r.expandTableWidth()
	}
	return r.shrinkTableWidth()
}

func (r *resizer) detectTableWidth() int {
	return r.maxCharCount() + r.totalPadding() + r.totalBorder()
}

func (r *resizer) expandTableWidth() []int {
	colWidths := r.maxColumnWidths()

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

	return colWidths
}

func (r *resizer) shrinkTableWidth() []int {
	colWidths := r.maxColumnWidths()

	for {
		totalWidth := sum(colWidths) + r.totalBorder()
		if totalWidth <= r.tableWidth {
			break
		}

		biggestDiffToMedian := -math.MaxInt32
		biggestDiffToMedianIndex := -math.MaxInt32

		for i, width := range colWidths {
			diffToMedian := width - r.columns[i].median

			if diffToMedian > biggestDiffToMedian {
				biggestDiffToMedian = diffToMedian
				biggestDiffToMedianIndex = i
			}
		}

		// No more space to shrink.
		if colWidths[biggestDiffToMedianIndex] == 0 {
			break
		}

		colWidths[biggestDiffToMedianIndex]--
	}

	return colWidths
}

func (r *resizer) maxColumnWidths() []int {
	maxColumnWidths := make([]int, len(r.columns))
	for i, col := range r.columns {
		if col.fixedWidth > 0 {
			maxColumnWidths[i] = col.fixedWidth
		} else {
			maxColumnWidths[i] = col.max + r.paddingForCol(col.index)
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
			count += col.fixedWidth - r.paddingForCol(col.index)
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
			maxTotal += column.max + r.paddingForCol(j)
		}
	}
	return
}

func (r *resizer) totalPadding() (totalPadding int) {
	for _, col := range r.columns {
		totalPadding += col.padding
	}
	return
}

func (r *resizer) paddingForCol(j int) int {
	if j >= len(r.columns) {
		return 0
	}
	return r.columns[j].padding
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
