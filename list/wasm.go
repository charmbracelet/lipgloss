//go:build wasm
// +build wasm

package list

import "github.com/charmbracelet/lipgloss/v2"

//go:export ListStyleFunc
func (l *List) wasmStyleFunc(id int32) *List {
	l.ItemStyleFunc(func(items Items, index int) lipgloss.Style {
		style := lipgloss.GetStyleFromJS(id, int32(index), 0)
		if style == nil {
			return lipgloss.NewStyle()
		}
		return *style
	})
	return l
}

//go:export ListEnumeratorStyleFunc
func (l *List) wasmEnumeratorStyleFunc(id int32) *List {
	l.EnumeratorStyleFunc(func(items Items, index int) lipgloss.Style {
		style := lipgloss.GetStyleFromJS(id, int32(index), 0)
		if style == nil {
			return lipgloss.NewStyle()
		}
		return *style
	})
	return l
}

//go:export ListRenderPtr
func (l *List) wasmRenderPtr() *byte {
	// Return a pointer to the first byte of the string
	str := l.String()
	if len(str) == 0 {
		return nil
	}
	return &([]byte(str)[0])
}

//go:export ListRenderLength
func (l *List) wasmRenderLength() int {
	// Return the length of the string
	return len(l.String())
}