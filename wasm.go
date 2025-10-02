//go:build wasm
// +build wasm

package lipgloss

import (
	"unsafe"
)

// We need Style as uintpr instead of *Style because
// it is used as part of callStyleFunc, it isn't allowed
// to send custom struct go pointers (not u)
//
//go:export StyleNewStyle
func wasmNewStyle() *Style {
	return &Style{}
}

// Import the JavaScript function we want to call
//
//go:wasmimport env callStyleFunc
func WasmStyleFunc(funcId, row, col int32) uintptr

func GetStyleFromJS(funcId, row, col int32) *Style {
	// Call the JavaScript function to get the pointer
	ptr := WasmStyleFunc(funcId, row, col)

	// Cast the pointer to our Style struct
	if ptr == 0 {
		return nil // Handle null pointers
	}

	// Convert the raw pointer to a Go struct pointer
	return (*Style)(unsafe.Pointer(ptr))
}
