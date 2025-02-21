package table

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

type shrinkStrategy string

const (
	biggestDiffToMedian shrinkStrategy = "biggestDiffToMedian"
	biggestColumn       shrinkStrategy = "biggestColumn"
)

func (t *Table) resize() {
	strData, ok := t.data.(*StringData)
	if !ok {
		return
	}
	r := newResizer(t.width, t.height, t.headers, strData.rows)
	r.wrap = t.wrap
	r.borderColumn = t.borderColumn

	allRows := append([][]string{t.headers}, strData.rows...)
	r.rowHeights = r.defaultRowHeigths()

	for i, row := range allRows {
		for j := range row {
			column := &r.columns[j]
			style := t.styleFunc(i, j)

			_, rightMargin, _, leftMargin := style.GetMargin()
			_, rightPadding, _, leftPadding := style.GetPadding()

			totalPadding := leftMargin + rightMargin + leftPadding + rightPadding
			column.padding = max(column.padding, totalPadding)
			column.fixedWidth = max(column.fixedWidth, style.GetWidth())

			r.rowHeights[i] = max(r.rowHeights[i], style.GetHeight())
		}
	}

	// A table width wasn't specified. In this case, detect according to
	// content width.
	if r.tableWidth <= 0 {
		r.tableWidth = r.detectTableWidth()
	}

	t.widths, t.heights = r.optimizedWidths()

	// fmt.Printf("t.heights: %+v\n", t.heights)
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
	headers     []string
	allRows     [][]string
	rowHeights  []int
	columns     []resizerColumn

	wrap         bool
	borderColumn bool
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
	return r.shrinkTableWidth(biggestDiffToMedian)
	// shrinkedWidths, ok := r.shrinkTableWidth()
	// if ok {
	// 	return shrinkedWidths
	// }
	// return r.shrinkTableAndContentWidth()
}

func (r *resizer) detectTableWidth() int {
	return r.maxCharCount() + r.totalPadding() + r.totalBorder()
}

func (r *resizer) expandTableWidth() (colWidths, rowHeights []int) {
	colWidths = r.maxColumnWidths()
	rowHeights = r.defaultRowHeigths()

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

	return
}

// func (r *resizer) shrinkTableWidth() (colWidths []int, ok bool) {
func (r *resizer) shrinkTableWidth(strategy shrinkStrategy) (colWidths, rowHeights []int) {
	colWidths = r.maxColumnWidths()
	rowHeights = r.defaultRowHeigths()
	// startCuttingContent := false

	for {
		totalWidth := sum(colWidths) + r.totalBorder()
		if totalWidth <= r.tableWidth {
			break
		}

		biggestDiffToMedian := -math.MaxInt32
		biggestColumn := -math.MaxInt32
		biggestDiffToMedianIndex := -math.MaxInt32

		for j, width := range colWidths {
			// Do not truncate headers.
			// if len(r.headers) > j && lipgloss.Width(r.headers[j]) <= width {
			// 	continue
			// }

			switch strategy {
			case "biggestDiffToMedian":
				diffToMedian := width - r.columns[j].median
				if diffToMedian > 0 && diffToMedian > biggestDiffToMedian {
					biggestDiffToMedian = diffToMedian
					biggestDiffToMedianIndex = j
				}
			case "biggestColumn":
				if width > biggestColumn {
					biggestColumn = width
					biggestDiffToMedianIndex = j
				}
			}
		}

		// println("biggestDiffToMedianIndex", biggestDiffToMedianIndex, "biggestDiffToMedian", biggestDiffToMedian)

		// Shirinking would cut content.
		// if biggestDiffToMedianIndex < 0 {
		if strategy == "biggestDiffToMedian" && biggestDiffToMedian <= 0 {
			println("shirinking would cut content")
			strategy = "biggestColumn"
			// startCuttingContent = true
			continue
		}

		// No more space to shrink.
		if colWidths[biggestDiffToMedianIndex] == 0 {
			println("no more space to shrink")
			// startCuttingContent = true
			break
		}

		colWidths[biggestDiffToMedianIndex]--
	}

	// If wrap is enabled (the default), increase the height of the rows
	// that have content that would not fit into a single line.
	if r.wrap {
		for i, row := range r.allRows {
			for j, cell := range row {
				if lipgloss.Width(cell) > colWidths[j] {
					// panic("here!")
					rowHeights[i] += 5
					// biggestDiffToMedianIndex = j
					break
				}
			}
		}
	}

	// if !shinkWouldCutContent {
	// 	return colWidths
	// }

	// if !r.wrap {
	// 	panic("TODO: implement")
	// }

	// for j := range r.columns {

	// }

	// for i, row := range r.allRows {
	// 	for j, col := range row {
	// 		recordWidth, _ := lipgloss.Size(col)

	// 		println("recordWidth", recordWidth, "colWidths[j]", colWidths[j])

	// 		if recordWidth > colWidths[j] {
	// 			r.rowHeights[i]++
	// 			biggestDiffToMedianIndex = j

	// 			// if r.rowHeights[i] > r.tableHeight {
	// 			// 	break
	// 			// }

	// 			// colWidths[j] = max(colWidths[j], recordWidth)

	// 			// totalWidth := sum(colWidths) + r.totalBorder()
	// 			// if totalWidth >= r.tableWidth {
	// 			// 	break
	// 			// }

	// 			// break
	// 		}
	// 	}
	// }

	return
}

// func (r *resizer) shrinkTableAndContentWidth() []int {

// }

func (r *resizer) defaultRowHeigths() (rowHeights []int) {
	rowHeights = make([]int, len(r.allRows))
	for i := range rowHeights {
		rowHeights[i] = 1
	}
	return
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

// func (r *resizer) medianColumnWidths() []int {
// 	medianColumnWidths := make([]int, len(r.columns))
// 	for i, col := range r.columns {
// 		if col.fixedWidth > 0 {
// 			medianColumnWidths[i] = col.fixedWidth
// 		} else {
// 			medianColumnWidths[i] = col.median + r.paddingForCol(col.index)
// 		}
// 	}
// 	return medianColumnWidths
// }

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
