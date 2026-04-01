package tree

import "fmt"

// DynamicLeaf is a leaf node whose value is computed dynamically each time
// it is rendered. This is useful for items whose content changes over time,
// such as a counter or a timestamp.
//
// Example:
//
//	counter := 0
//	item := tree.NewDynamicLeaf(tree.StringerFunc(func() string {
//		counter++
//		return fmt.Sprintf("count: %d", counter)
//	}))
//
//	l := list.New().Item(item)
type DynamicLeaf struct {
	stringer fmt.Stringer
	hidden   bool
}

// NewDynamicLeaf returns a new DynamicLeaf that evaluates its value lazily
// from the provided Stringer.
func NewDynamicLeaf(s fmt.Stringer) *DynamicLeaf {
	return &DynamicLeaf{stringer: s}
}

// Children of a DynamicLeaf are always empty.
func (d DynamicLeaf) Children() Children {
	return NodeChildren(nil)
}

// Value returns the dynamically computed value of this node.
func (d DynamicLeaf) Value() string {
	if d.stringer == nil {
		return ""
	}
	return d.stringer.String()
}

// SetValue sets the stringer for this dynamic leaf. If the value is a
// fmt.Stringer, it is stored directly. Otherwise, the value is converted
// to a string and wrapped in a static stringer.
func (d *DynamicLeaf) SetValue(value any) {
	switch item := value.(type) {
	case fmt.Stringer:
		d.stringer = item
	case string:
		d.stringer = StringerFunc(func() string { return item })
	default:
		s := fmt.Sprintf("%v", value)
		d.stringer = StringerFunc(func() string { return s })
	}
}

// Hidden returns whether this node is hidden.
func (d DynamicLeaf) Hidden() bool {
	return d.hidden
}

// SetHidden hides or shows this node.
func (d *DynamicLeaf) SetHidden(hidden bool) {
	d.hidden = hidden
}

// String returns the dynamically computed string representation.
func (d DynamicLeaf) String() string {
	return d.Value()
}

// StringerFunc is an adapter to allow the use of ordinary functions as
// fmt.Stringer. It works like http.HandlerFunc.
//
// Example:
//
//	s := tree.StringerFunc(func() string {
//		return time.Now().Format(time.RFC3339)
//	})
type StringerFunc func() string

// String calls f().
func (f StringerFunc) String() string {
	return f()
}
