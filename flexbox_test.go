package lipgloss

import (
	"fmt"
	"testing"
)

func TestExample(t *testing.T) {
	container := NewStyle().
		Border(normalBorder).
		FlexWrap(FlexWrapWrap).
		Width(100).
		Height(40)

	style1 := NewStyle().
		FlexGrow(2).
		Width(75).
		Border(blockBorder, true).
		SetString("content 1")

	style2 := NewStyle().
		FlexGrow(1).
		Border(doubleBorder, true).
		SetString("content 2")

	fmt.Println(Flexbox(container, style1, style2))
}
