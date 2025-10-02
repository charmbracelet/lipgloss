//go:build wasm
// +build wasm

package tree

import "github.com/charmbracelet/lipgloss/v2"

//go:export TreeStyleFunc
func (t *Tree) wasmStyleFunc(id int32) *Tree {
	t.ItemStyleFunc(func(children Children, index int) lipgloss.Style {
		style := lipgloss.GetStyleFromJS(id, int32(index), 0)
		if style == nil {
			return lipgloss.NewStyle()
		}
		return *style
	})
	return t
}

//go:export TreeEnumeratorStyleFunc
func (t *Tree) wasmEnumeratorStyleFunc(id int32) *Tree {
	t.EnumeratorStyleFunc(func(children Children, index int) lipgloss.Style {
		style := lipgloss.GetStyleFromJS(id, int32(index), 0)
		if style == nil {
			return lipgloss.NewStyle()
		}
		return *style
	})
	return t
}

//go:export TreeRenderPtr
func (t *Tree) wasmRenderPtr() *byte {
	// Return a pointer to the first byte of the string
	str := t.String()
	if len(str) == 0 {
		return nil
	}
	return &([]byte(str)[0])
}

//go:export TreeRenderLength
func (t *Tree) wasmRenderLength() int {
	// Return the length of the string
	return len(t.String())
}

//go:export TreeEnumeratorDefault
func wasmTreeEnumeratorDefault() int32 {
	return 0 // DefaultEnumerator
}

//go:export TreeEnumeratorRounded
func wasmTreeEnumeratorRounded() int32 {
	return 1 // RoundedEnumerator
}

//go:export TreeIndenterDefault
func wasmTreeIndenterDefault() int32 {
	return 0 // DefaultIndenter
}