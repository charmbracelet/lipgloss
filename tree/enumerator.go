package tree

// Enumerator returns the indent (also known as branch) and tree prefixes of a
// given item.
type Enumerator func(children Children, index int) (indent string, prefix string)

// DefaultEnumerator enumerates items.
//
//	root
//	├── foo
//	└── bar
func DefaultEnumerator(children Children, index int) (indent, prefix string) {
	if children.Length()-1 == index {
		return "   ", "└──"
	}
	return "│  ", "├──"
}

// RoundedEnumerator enumerates items.
//
//	root
//	├── foo
//	╰── bar
func RoundedEnumerator(children Children, index int) (indent, prefix string) {
	if children.Length()-1 == index {
		return "   ", "╰──"
	}
	return "│  ", "├──"
}
