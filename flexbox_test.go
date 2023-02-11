package lipgloss

import (
	"fmt"
	"testing"
)

func TestExample(t *testing.T) {
	container := NewStyle().
		Border(normalBorder).
		FlexWrap(FlexWrapWrap).
		FlexDirection(FlexDirRowReverse).
		// FlexDirection(FlexDirColumn).
		Width(100)
	// Height(40)

	style1 := NewStyle().
		FlexGrow(2).
		Width(75).
		Border(blockBorder, true).
		SetString("content 1\n\nfooo")

	style2 := NewStyle().
		FlexGrow(1).
		Border(doubleBorder, true).
		Padding(2).
		SetString("content 2")

	style3 := NewStyle().
		FlexGrow(1).
		Width(60).
		Border(doubleBorder, true).
		SetString("content 3\n\nfoobar\nbarfoo")

	fmt.Println(Flexbox(container, style1, style2, style3))
}

func TestJustifyContent(t *testing.T) {
	container := NewStyle().
		Border(normalBorder).
		FlexWrap(FlexWrapWrap).
		FlexDirection(FlexDirColumn).
		// FlexJustifyContent(FlexJustifyContentFlexStart).
		// FlexJustifyContent(FlexJustifyContentFlexEnd).
		// FlexJustifyContent(FlexJustifyContentCenter).
		// FlexJustifyContent(FlexJustifyContentSpaceBetween).
		FlexJustifyContent(FlexJustifyContentSpaceAround).
		Width(100).
		Height(5)

	style1 := NewStyle().
		Border(blockBorder, true).
		SetString("content 1\n\nfooo")

	style2 := NewStyle().
		Border(doubleBorder, true).
		Padding(2).
		SetString("content 2")

	style3 := NewStyle().
		Width(40).
		Border(doubleBorder, true).
		SetString("content 3\n\nfoobar\nbarfoo")

	style4 := NewStyle().
		Width(60).
		Border(doubleBorder, true).
		SetString("content 4\n\nfoobar\nbarfoo")

	fmt.Println(Flexbox(container, style1, style2, style3, style4))
}
