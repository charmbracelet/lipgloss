package table

// Model is the interface that wraps the basic methods of a table model.
type Model interface {
	Row(row int) Row
	Count() int
	Columns() int
}

// Row represents one line in the table.
type Row interface {
	Column(col int) string
}

// StringModel is a string-based implementation of the Model interface.
type StringModel struct {
	rows    []Row
	columns int
}

// NewStringModel creates a new StringModel with the given number of columns.
func NewStringModel(columns int) *StringModel {
	return &StringModel{columns: columns}
}

// Row returns the row at the given index.
func (m *StringModel) Row(row int) Row {
	return m.rows[row]
}

// Columns returns the number of columns in the table.
func (m *StringModel) Columns() int {
	return m.columns
}

// AppendRows appends the given rows to the table.
func (m *StringModel) AppendRows(rows ...[]string) *StringModel {
	for _, row := range rows {
		m.rows = append(m.rows, StringRow(row))
	}

	return m
}

// AppendRow appends the given row to the table.
func (m *StringModel) AppendRow(rows ...string) *StringModel {
	m.rows = append(m.rows, StringRow(rows))

	return m
}

// Count returns the number of rows in the table.
func (m *StringModel) Count() int {
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

type FilterModel struct {
	model  Model
	filter func(row Row) bool
}

func NewFilterModel(model Model) *FilterModel {
	return &FilterModel{model: model}
}

func (m *FilterModel) Filter(f func(row Row) bool) *FilterModel {
	m.filter = f
	return m
}

func (m *FilterModel) Row(row int) Row {
	j := 0
	for i := 0; i < m.model.Count(); i++ {
		if m.filter(m.model.Row(i)) {
			if j == row {
				return m.model.Row(i)
			}

			j++
		}
	}

	return nil
}

func (m *FilterModel) Columns() int {
	return m.model.Columns()
}

func (m *FilterModel) Count() int {
	j := 0
	for i := 0; i < m.model.Count(); i++ {
		if m.filter(m.model.Row(i)) {
			j++
		}
	}

	return j
}
