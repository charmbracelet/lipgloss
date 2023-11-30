package list

// Data is the interface that wraps the basic methods of a list model.
type Data interface {
	// At returns the content item of the given index.
	At(index int) string

	// Length returns the number of items in the list.
	Length() int
}

// StringData is a string-based implementation of the Data interface.
type StringData struct {
	items []string
}

// NewStringData creates a new StringData with the given items.
func NewStringData(items ...string) *StringData {
	return &StringData{
		items: items,
	}
}

// Append appends the given item to the list.
func (m *StringData) Append(item string) {
	m.items = append(m.items, item)
}

// At returns the item at the given index.
func (m *StringData) At(index int) string {
	if index < 0 || index >= len(m.items) {
		return ""
	}

	return m.items[index]
}

// Length returns the number of items in the list.
func (m *StringData) Length() int {
	return len(m.items)
}

// Filter applies a filter on some data.
type Filter struct {
	data   Data
	filter func(index int) bool
}

// NewFilter initializes a new Filter.
func NewFilter(data Data) *Filter {
	return &Filter{data: data}
}

// At returns the item at the given index.
func (m *Filter) At(index int) string {
	j := 0
	for i := 0; i < m.data.Length(); i++ {
		if m.filter(i) {
			if j == index {
				return m.data.At(i)
			}
			j++
		}
	}

	return ""
}

// Filter applies the given filter function to the data.
func (m *Filter) Filter(f func(index int) bool) *Filter {
	m.filter = f
	return m
}

// Length returns the number of items in the list.
func (m *Filter) Length() int {
	j := 0
	for i := 0; i < m.data.Length(); i++ {
		if m.filter(i) {
			j++
		}
	}
	return j
}
