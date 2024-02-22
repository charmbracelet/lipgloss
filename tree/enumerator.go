package tree

// Enumerator returns the indent (also known as branch) and tree prefixes of a
// given item.
type Enumerator func(data Data, i int) (indent string, prefix string)

// DefaultEnumerator enumerates items.
func DefaultEnumerator(data Data, i int) (indent, prefix string) {
	if data.Length()-1 == i {
		return "   ", "└──"
	}
	return "│  ", "├──"
}
