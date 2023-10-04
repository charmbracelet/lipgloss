package table

// Data is the interface that wraps the basic methods of a table model.
type Data interface {
	Row(row int) Row
	Append(row Row)
	Count() int
	Columns() int
}

// Row represents one line in the table.
type Row interface {
	Column(col int) string
	Length() int
}

// StringData is a string-based implementation of the Data interface.
type StringData struct {
	rows    []Row
	columns int
}

// Rows creates a new StringData with the given number of columns.
func Rows(rows ...[]string) *StringData {
	m := StringData{columns: 0}

	for _, row := range rows {
		m.columns = max(m.columns, len(row))
		m.rows = append(m.rows, StringRow(row))
	}

	return &m
}

// Append appends the given row to the table.
func (m *StringData) Append(row Row) {
	m.columns = max(m.columns, row.Length())
	m.rows = append(m.rows, row)
}

// Row returns the row at the given index.
func (m *StringData) Row(row int) Row {
	return m.rows[row]
}

// Columns returns the number of columns in the table.
func (m *StringData) Columns() int {
	return m.columns
}

// Item appends the given row to the table.
func (m *StringData) Item(rows ...string) *StringData {
	m.columns = max(m.columns, len(rows))
	m.rows = append(m.rows, StringRow(rows))

	return m
}

// Count returns the number of rows in the table.
func (m *StringData) Count() int {
	return len(m.rows)
}

// StringRow is a simple implementation of the Row interface.
type StringRow []string

// Value returns the value of the column at the given index.
func (r StringRow) Column(col int) string {
	if col >= len(r) {
		return ""
	}

	return r[col]
}

// Value returns the value of the column at the given index.
func (r StringRow) Length() int {
	return len(r)
}

// Filter applies a filter on some data.
type Filter struct {
	data   Data
	filter func(row Row) bool
}

// NewFilter initializes a new Filter.
func NewFilter(data Data) *Filter {
	return &Filter{data: data}
}

// Filter applies the given filter function to the data.
func (m *Filter) Filter(f func(row Row) bool) *Filter {
	m.filter = f
	return m
}

// Row returns the row at the given index.
func (m *Filter) Row(row int) Row {
	j := 0
	for i := 0; i < m.data.Count(); i++ {
		if m.filter(m.data.Row(i)) {
			if j == row {
				return m.data.Row(i)
			}

			j++
		}
	}

	return nil
}

// Append appends the given row to the table.
func (m *Filter) Append(row Row) {
	m.data.Append(row)
}

// Columns returns the number of columns in the table.
func (m *Filter) Columns() int {
	return m.data.Columns()
}

// Count returns the number of rows in the table.
func (m *Filter) Count() int {
	j := 0
	for i := 0; i < m.data.Count(); i++ {
		if m.filter(m.data.Row(i)) {
			j++
		}
	}

	return j
}
