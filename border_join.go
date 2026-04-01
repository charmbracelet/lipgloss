package lipgloss

// JoinBordersLeft returns a new Border that represents the left border when
// joining two bordered boxes horizontally. The right side corners and edges
// are replaced with the appropriate middle joining characters.
//
// Example:
//
//	left := lipgloss.JoinBordersLeft(lipgloss.NormalBorder())
//	right := lipgloss.JoinBordersRight(lipgloss.NormalBorder())
//
//	leftBox := lipgloss.NewStyle().Border(left).Render("Left")
//	rightBox := lipgloss.NewStyle().Border(right).Render("Right")
//
//	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox))
//
// This will produce a single shared border between the two boxes instead of
// a doubled border.
func JoinBordersLeft(b Border) Border {
	b.TopRight = b.MiddleTop
	b.BottomRight = b.MiddleBottom
	b.Right = b.Middle
	return b
}

// JoinBordersRight returns a new Border that represents the right border when
// joining two bordered boxes horizontally. The left side corners and edges
// are removed since they are provided by the left box's border.
func JoinBordersRight(b Border) Border {
	b.TopLeft = ""
	b.BottomLeft = ""
	b.Left = ""
	return b
}

// JoinBordersTop returns a new Border that represents the top border when
// joining two bordered boxes vertically. The bottom corners and edges are
// replaced with the appropriate middle joining characters.
func JoinBordersTop(b Border) Border {
	b.BottomLeft = b.MiddleLeft
	b.BottomRight = b.MiddleRight
	b.Bottom = b.Middle
	return b
}

// JoinBordersBottom returns a new Border that represents the bottom border when
// joining two bordered boxes vertically. The top corners and edges are removed
// since they are provided by the top box's border.
func JoinBordersBottom(b Border) Border {
	b.TopLeft = ""
	b.TopRight = ""
	b.Top = ""
	return b
}

// JoinBordersMiddleHorizontal returns a new Border for a box that sits in
// the middle of a horizontal row of bordered boxes. Both the left and right
// borders are replaced with joining characters.
func JoinBordersMiddleHorizontal(b Border) Border {
	b.TopLeft = ""
	b.BottomLeft = ""
	b.Left = ""
	b.TopRight = b.MiddleTop
	b.BottomRight = b.MiddleBottom
	b.Right = b.Middle
	return b
}

// JoinBordersMiddleVertical returns a new Border for a box that sits in
// the middle of a vertical column of bordered boxes. Both the top and bottom
// borders are replaced with joining characters.
func JoinBordersMiddleVertical(b Border) Border {
	b.TopLeft = ""
	b.TopRight = ""
	b.Top = ""
	b.BottomLeft = b.MiddleLeft
	b.BottomRight = b.MiddleRight
	b.Bottom = b.Middle
	return b
}
