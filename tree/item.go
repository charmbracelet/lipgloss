package tree

// Children is the interface that wraps the basic methods of a list model.
type Children interface {
	// At returns the content item of the given index.
	At(index int) Node

	// Length returns the number of items in the list.
	Length() int

	// Removes the given index returning the Data without it.
	Remove(index int) Children

	// Append adds the given node to the underlying data.
	Append(item Node) Children
}

type NodeData []Node

var _ Children = NodeData(nil)

func (n NodeData) Append(item Node) Children {
	n = append(n, item)
	return n
}

func (n NodeData) Remove(i int) Children {
	if i < 0 || len(n) < i+1 {
		return n
	}
	n = append(n[:i], n[i+1:]...)
	return n
}

func (n NodeData) Length() int {
	return len(n)
}

func (n NodeData) At(i int) Node {
	if i >= 0 && i < len(n) {
		return n[i]
	}
	return nil
}

// NewStringData returns a Data of strings.
func NewStringData(data ...string) Children {
	result := make([]Node, 0, len(data))
	for _, d := range data {
		s := StringNode(d)
		result = append(result, &s)
	}
	return NodeData(result)
}

var _ Children = NewFilter(nil)

// Filter applies a filter on some data. You could use this to create a new
// tree whose values all satisfy the condition provided in the Filter() function.
type Filter struct {
	data   Children
	filter func(index int) bool
}

// Append allows Filter to implement Data. It adds an element to the Tree.
func (m *Filter) Append(item Node) Children {
	m.data = m.data.Append(item)
	return m
}

// NewFilter initializes a new Filter.
func NewFilter(data Children) *Filter {
	return &Filter{data: data}
}

// Remove allows Filter to implement Data. It removes an element from the Tree.
func (m *Filter) Remove(index int) Children {
	m.data = m.data.Remove(index)
	return m
}

// At returns the item at the given index.
// The index is relative to the filtered results.
func (m *Filter) At(index int) Node {
	j := 0
	for i := 0; i < m.data.Length(); i++ {
		if m.filter(i) {
			if j == index {
				return m.data.At(i)
			}
			j++
		}
	}

	return nil
}

// Filter uses a filter function to set a condition that all the data must satisfy to be in the Tree.
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
