//go:build wasm
// +build wasm

package table

import (
	"unsafe"

	"github.com/charmbracelet/lipgloss/v2"
)

//go:export TableStyleFunc
func (t *Table) wasmStyleFunc(id int32) *Table {
	t.styleFunc = func(row, col int) lipgloss.Style {
		style := lipgloss.GetStyleFromJS(id, int32(row), int32(col))
		if style == nil {
			return lipgloss.NewStyle()
		}
		return *style
	}
	return t
}

//go:export TableDataNew
func wasmNewStringData() *StringData {
	return NewStringData()
}

//go:export TableDataAppend
func (d *StringData) wasmAppend(ptrArray *uint32, count int) *StringData {
	// Convert to a slice of pointers
	pointerArray := unsafe.Slice(ptrArray, count*2)
	row := []string{}

	for i := range count {
		ptr := pointerArray[i*2]
		length := pointerArray[i*2+1]

		if length > 0 {
			// Convert pointer to byte slice
			bytes := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
			row = append(row, string(bytes))
		} else {
			row = append(row, "")
		}
	}

	d.Append(row)
	return d
}

//go:export TableDataAtPtr
func (d *StringData) wasmAtPtr(row, col int32) uintptr {
	value := d.At(int(row), int(col))
	if value == "" {
		return 0
	}
	
	// Convert string to bytes and return pointer
	bytes := []byte(value)
	if len(bytes) == 0 {
		return 0
	}
	
	return uintptr(unsafe.Pointer(&bytes[0]))
}

//go:export TableDataAtLength
func (d *StringData) wasmAtLength(row, col int32) int {
	value := d.At(int(row), int(col))
	return len(value)
}

//go:export TableDataRows
func (d *StringData) wasmRows() int32 {
	return int32(d.Rows())
}

//go:export TableDataColumns
func (d *StringData) wasmColumns() int32 {
	return int32(d.Columns())
}

//go:export TableSetData
func (t *Table) wasmSetData(data *StringData) *Table {
	t.data = data
	return t
}
