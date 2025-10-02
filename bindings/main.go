package main

import (
	"os"
	"unsafe"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/charmbracelet/lipgloss/v2/list"
	_ "github.com/charmbracelet/lipgloss/v2/table"
	"github.com/charmbracelet/lipgloss/v2/tree"
)

// Note: this function is used on wasm module initalization.
func main() {}

//go:export DetectFromEnvVars
func wasmDetectFromEnvVars(ptrArray *uint32, count int) {
	// Convert to a slice of pointers
	pointerArray := unsafe.Slice(ptrArray, count*2)
	vars := []string{}

	for i := range count {
		strPtr := uintptr(pointerArray[i*2])
		strLen := int(pointerArray[i*2+1])

		// Create a byte slice and convert to string
		bytes := unsafe.Slice((*byte)(unsafe.Pointer(strPtr)), strLen)
		str := string(bytes)
		vars = append(vars, str)
	}

	compat.Profile = colorprofile.Env(vars) // truecolor
	lipgloss.Writer = &colorprofile.Writer{
		Forward: os.Stdout,
		Profile: colorprofile.Env(vars),
	}
}

var buf [1024]byte

//go:export getBuffer
func GetBuffer() *byte {
	return &buf[0]
}

// TinyGo-specific memory management
// Use a simple memory pool approach instead of complex tracking
var memoryPool [][]byte
var poolIndex int

//go:export wasmMalloc
func wasmMalloc(size int) uintptr {
	if size <= 0 {
		return 0
	}

	// For TinyGo, use a simpler approach with pre-allocated pools
	// This avoids the string length overflow issues
	data := make([]byte, size)

	// Store in pool to prevent GC
	if poolIndex >= len(memoryPool) {
		memoryPool = append(memoryPool, make([][]byte, 100)...)
	}
	memoryPool[poolIndex] = data
	poolIndex++

	return uintptr(unsafe.Pointer(&data[0]))
}

//go:export wasmFree
func wasmFree(ptr uintptr) {
	// For TinyGo, we'll rely on periodic cleanup rather than immediate freeing
	// This is more compatible with TinyGo's GC
}

//go:export wasmGC
func wasmGC() {
	// Periodic cleanup for TinyGo
	if poolIndex > 1000 {
		// Reset the pool periodically to prevent memory buildup
		memoryPool = memoryPool[:0]
		poolIndex = 0
	}
}

//go:export getMemorySize
func wasmGetMemorySize() int {
	// Return a reasonable estimate for TinyGo
	return 16 * 1024 * 1024 // 16MB
}

//go:export getAllocatedSize
func wasmGetAllocatedSize() int {
	return poolIndex * 1024 // Rough estimate
}

// List exports

//go:export ListEnumeratorAlphabet
func wasmListEnumeratorAlphabet() int32 {
	return 0 // Alphabet
}

//go:export ListEnumeratorArabic
func wasmListEnumeratorArabic() int32 {
	return 1 // Arabic
}

//go:export ListEnumeratorBullet
func wasmListEnumeratorBullet() int32 {
	return 2 // Bullet
}

//go:export ListEnumeratorDash
func wasmListEnumeratorDash() int32 {
	return 3 // Dash
}

//go:export ListEnumeratorRoman
func wasmListEnumeratorRoman() int32 {
	return 4 // Roman
}

//go:export ListEnumeratorAsterisk
func wasmListEnumeratorAsterisk() int32 {
	return 5 // Asterisk
}

//go:export ListItem
func wasmListItem(listPtr uintptr, strPtr uintptr, strLen int) {
	l := (*list.List)(unsafe.Pointer(listPtr))
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(strPtr)), strLen)
	str := string(bytes)
	l.Item(str)
}

//go:export ListItemList
func wasmListItemList(listPtr uintptr, itemListPtr uintptr) {
	l := (*list.List)(unsafe.Pointer(listPtr))
	itemList := (*list.List)(unsafe.Pointer(itemListPtr))
	l.Item(itemList)
}

//go:export ListOffset
func wasmListOffset(listPtr uintptr, start int, end int) {
	l := (*list.List)(unsafe.Pointer(listPtr))
	l.Offset(start, end)
}

//go:export ListEnumeratorStyle
func wasmListEnumeratorStyle(listPtr uintptr, stylePtr uintptr) {
	l := (*list.List)(unsafe.Pointer(listPtr))
	style := (*lipgloss.Style)(unsafe.Pointer(stylePtr))
	l.EnumeratorStyle(*style)
}

//go:export ListItemStyle
func wasmListItemStyle(listPtr uintptr, stylePtr uintptr) {
	l := (*list.List)(unsafe.Pointer(listPtr))
	style := (*lipgloss.Style)(unsafe.Pointer(stylePtr))
	l.ItemStyle(*style)
}

//go:export ListEnumerator
func wasmListEnumerator(listPtr uintptr, enumType int32) {
	l := (*list.List)(unsafe.Pointer(listPtr))
	switch enumType {
	case 0:
		l.Enumerator(list.Alphabet)
	case 1:
		l.Enumerator(list.Arabic)
	case 2:
		l.Enumerator(list.Bullet)
	case 3:
		l.Enumerator(list.Dash)
	case 4:
		l.Enumerator(list.Roman)
	case 5:
		l.Enumerator(list.Asterisk)
	}
}

// Tree exports

//go:export TreeNew
func wasmTreeNew() uintptr {
	t := tree.New()
	return uintptr(unsafe.Pointer(t))
}

//go:export TreeRoot
func wasmTreeRoot(treePtr uintptr, strPtr uintptr, strLen int) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(strPtr)), strLen)
	str := string(bytes)
	t.Root(str)
}

//go:export TreeChild
func wasmTreeChild(treePtr uintptr, strPtr uintptr, strLen int) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(strPtr)), strLen)
	str := string(bytes)
	t.Child(str)
}

//go:export TreeChildTree
func wasmTreeChildTree(treePtr uintptr, childTreePtr uintptr) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	childTree := (*tree.Tree)(unsafe.Pointer(childTreePtr))
	t.Child(childTree)
}

//go:export TreeHidden
func wasmTreeHidden(treePtr uintptr) bool {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	return t.Hidden()
}

//go:export TreeHide
func wasmTreeHide(treePtr uintptr, hide bool) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	t.Hide(hide)
}

//go:export TreeOffset
func wasmTreeOffset(treePtr uintptr, start int, end int) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	t.Offset(start, end)
}

//go:export TreeEnumeratorStyle
func wasmTreeEnumeratorStyle(treePtr uintptr, stylePtr uintptr) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	style := (*lipgloss.Style)(unsafe.Pointer(stylePtr))
	t.EnumeratorStyle(*style)
}

//go:export TreeItemStyle
func wasmTreeItemStyle(treePtr uintptr, stylePtr uintptr) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	style := (*lipgloss.Style)(unsafe.Pointer(stylePtr))
	t.ItemStyle(*style)
}

//go:export TreeRootStyle
func wasmTreeRootStyle(treePtr uintptr, stylePtr uintptr) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	style := (*lipgloss.Style)(unsafe.Pointer(stylePtr))
	t.RootStyle(*style)
}

//go:export TreeEnumerator
func wasmTreeEnumerator(treePtr uintptr, enumType int32) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	switch enumType {
	case 0:
		t.Enumerator(tree.DefaultEnumerator)
	case 1:
		t.Enumerator(tree.RoundedEnumerator)
	}
}

//go:export TreeIndenter
func wasmTreeIndenter(treePtr uintptr, indenterType int32) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	switch indenterType {
	case 0:
		t.Indenter(tree.DefaultIndenter)
	}
}

//go:export TreeNewLeaf
func wasmTreeNewLeaf(strPtr uintptr, strLen int, hidden bool) uintptr {
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(strPtr)), strLen)
	str := string(bytes)
	leaf := tree.NewLeaf(str, hidden)
	return uintptr(unsafe.Pointer(leaf))
}

//go:export TreeLeafValue
func wasmTreeLeafValue(leafPtr uintptr) *byte {
	leaf := (*tree.Leaf)(unsafe.Pointer(leafPtr))
	value := leaf.Value()
	bytes := []byte(value)
	if len(bytes) == 0 {
		return nil
	}
	return &bytes[0]
}

//go:export TreeLeafValueLength
func wasmTreeLeafValueLength(leafPtr uintptr) int {
	leaf := (*tree.Leaf)(unsafe.Pointer(leafPtr))
	return len(leaf.Value())
}

//go:export TreeLeafHidden
func wasmTreeLeafHidden(leafPtr uintptr) bool {
	leaf := (*tree.Leaf)(unsafe.Pointer(leafPtr))
	return leaf.Hidden()
}

//go:export TreeLeafSetHidden
func wasmTreeLeafSetHidden(leafPtr uintptr, hidden bool) {
	leaf := (*tree.Leaf)(unsafe.Pointer(leafPtr))
	leaf.SetHidden(hidden)
}

//go:export TreeLeafSetValue
func wasmTreeLeafSetValue(leafPtr uintptr, strPtr uintptr, strLen int) {
	leaf := (*tree.Leaf)(unsafe.Pointer(leafPtr))
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(strPtr)), strLen)
	str := string(bytes)
	leaf.SetValue(str)
}

//go:export TreeChildLeaf
func wasmTreeChildLeaf(treePtr uintptr, leafPtr uintptr) {
	t := (*tree.Tree)(unsafe.Pointer(treePtr))
	leaf := (*tree.Leaf)(unsafe.Pointer(leafPtr))
	t.Child(leaf)
}
