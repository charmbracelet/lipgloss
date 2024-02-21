package tree

// Enumerator returns the indent (also known as branch) and tree prefixes of a
// given item.
type Enumerator func(data Data, i int, last bool) (indent string, prefix string)

// DefaultEnumerator enumerates items.
func DefaultEnumerator(_ Data, _ int, last bool) (indent, prefix string) {
	if last {
		return "   ", "└──"
	}
	return "│  ", "├──"
}
