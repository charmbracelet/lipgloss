package lipgloss

import (
	"math"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// JoinBorderVertical is a utility function for vertically joining two
// potentially multi-lined strings along a horizontal axis with borders.
// The first argument is the style used for the borders. The second argument
// is the alignment position, with 0 being all the way to the left and 1 being
// all the way to the right.
//
// If you just want to align to the left, right or center you may as well just
// use the helper constants Left, Center, and Right.
//
// Example:
//
//	blockB := "...\n...\n..."
//	blockA := "...\n...\n...\n...\n..."
//
//	// Join 20% from the top
//	str := JoinBorderVertical(s, 0.2, blockA, blockB)
//
//	// Join on the right edge
//	str := JoinBorderVertical(s, Right, blockA, blockB)
func JoinBorderVertical(borderStyle Style, pos Position, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	var (
		maxWidth int
		blocks   = make([][]string, len(strs))

		stt, str, stb, stl = getBorderStyles(borderStyle)

		border, top, right, bottom, left = borderStyle.GetBorder()
	)

	for i := range strs {
		var w int
		blocks[i], w = getLines(strs[i])
		if w > maxWidth {
			maxWidth = w
		}
	}

	var sb strings.Builder

	if top {
		var buff strings.Builder
		if left {
			buff.WriteString(border.TopLeft)
		}
		buff.WriteString(strings.Repeat(border.Top, maxWidth))
		if right {
			buff.WriteString(border.TopRight)
		}
		sb.WriteString(stt.Render(buff.String()))
		sb.WriteRune('\n')
	}

	for i, block := range blocks {
		for j, line := range block {
			w := maxWidth - ansi.StringWidth(line)

			if left {
				sb.WriteString(stl.Render(border.Left))
			}

			switch pos { //nolint:exhaustive
			case Left:
				sb.WriteString(line)
				sb.WriteString(strings.Repeat(" ", w))

			case Right:
				sb.WriteString(strings.Repeat(" ", w))
				sb.WriteString(line)

			default: // Somewhere in the middle
				if w < 1 {
					sb.WriteString(line)
					break
				}

				split := int(math.Round(float64(w) * pos.value()))
				right := w - split
				left := w - right

				sb.WriteString(strings.Repeat(" ", left))
				sb.WriteString(line)
				sb.WriteString(strings.Repeat(" ", right))
			}

			if right {
				sb.WriteString(str.Render(border.Right))
			}

			// Write a newline as long as we're not on the last line of the
			// last block.
			if !(i == len(blocks)-1 && j == len(block)-1) {
				sb.WriteRune('\n')
			}
		}
		if i < len(blocks)-1 {
			var buff strings.Builder
			if left {
				buff.WriteString(border.MiddleLeft)
			}
			buff.WriteString(strings.Repeat(border.Bottom, maxWidth))
			if right {
				buff.WriteString(border.MiddleRight)
			}
			sb.WriteString(stt.Render(buff.String()))
			sb.WriteRune('\n')
		}
	}

	if bottom {
		sb.WriteRune('\n')
		var buff strings.Builder
		if left {
			buff.WriteString(border.BottomLeft)
		}
		buff.WriteString(strings.Repeat(border.Bottom, maxWidth))
		if right {
			buff.WriteString(border.BottomRight)
		}
		sb.WriteString(stb.Render(buff.String()))
	}

	return sb.String()
}

// JoinBorderHorizontal is a utility function for horizontally joining two
// potentially multi-lined strings along a vertical axis. The first argument is
// the position, with 0 being all the way at the top and 1 being all the way
// at the bottom.
//
// If you just want to align to the top, center or bottom you may as well just
// use the helper constants Top, Center, and Bottom.
//
// Example:
//
//	blockB := "...\n...\n..."
//	blockA := "...\n...\n...\n...\n..."
//
//	// Join 20% from the top
//	str := JoinBorderHorizontal(s, 0.2, blockA, blockB)
//
//	// Join on the top edge
//	str := JoinBorderHorizontal(s, Top, blockA, blockB)
func JoinBorderHorizontal(borderStyle Style, pos Position, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	var (
		// Groups of strings broken into multiple lines
		blocks = make([][]string, len(strs))

		stt, str, stb, stl = getBorderStyles(borderStyle)

		border, top, right, bottom, left = borderStyle.GetBorder()

		// Max line widths for the above text blocks
		maxWidths = make([]int, len(strs))

		// Height of the tallest block
		maxHeight int
	)

	// Break text blocks into lines and get max widths for each text block
	for i, str := range strs {
		blocks[i], maxWidths[i] = getLines(str)
		if len(blocks[i]) > maxHeight {
			maxHeight = len(blocks[i])
		}
	}

	// Add extra lines to make each side the same height
	for i := range blocks {
		if len(blocks[i]) >= maxHeight {
			continue
		}

		extraLines := make([]string, maxHeight-len(blocks[i]))

		switch pos { //nolint:exhaustive
		case Top:
			blocks[i] = append(blocks[i], extraLines...)

		case Bottom:
			blocks[i] = append(extraLines, blocks[i]...)

		default: // Somewhere in the middle
			n := len(extraLines)
			split := int(math.Round(float64(n) * pos.value()))
			top := n - split
			bottom := n - top

			blocks[i] = append(extraLines[top:], blocks[i]...)
			blocks[i] = append(blocks[i], extraLines[bottom:]...)
		}
	}

	// Merge lines
	var sb strings.Builder

	// write top border
	if top {
		var buff strings.Builder
		buff.WriteString(border.TopLeft)
		for j := range blocks {
			if j > 0 {
				buff.WriteString(border.MiddleTop)
			}
			buff.WriteString(strings.Repeat(border.Top, maxWidths[j]))
		}
		buff.WriteString(border.TopRight)
		sb.WriteString(stt.Render(buff.String()))
		sb.WriteRune('\n')
	}

	for i := range blocks[0] { // remember, all blocks have the same number of members now
		for j, block := range blocks {
			if left || j > 0 {
				sb.WriteString(stl.Render(border.Left))
			}
			sb.WriteString(block[i])

			// Also make lines the same length
			sb.WriteString(strings.Repeat(" ", maxWidths[j]-ansi.StringWidth(block[i])))
		}
		if right {
			sb.WriteString(str.Render(border.Right))
		}
		if i < len(blocks[0])-1 {
			sb.WriteRune('\n')
		}
	}

	// write bottom border
	if bottom {
		sb.WriteRune('\n')

		var buff strings.Builder
		buff.WriteString(border.BottomLeft)
		for j := range blocks {
			if j > 0 {
				buff.WriteString(border.MiddleBottom)
			}
			buff.WriteString(strings.Repeat(border.Bottom, maxWidths[j]))
		}
		buff.WriteString(border.BottomRight)
		sb.Write([]byte(stb.Render(buff.String())))
	}

	return sb.String()
}

// getBorderStyles gets the styles for each side of the border
// returns Top, Right, Bottom, and Left Styles
func getBorderStyles(borderStyle Style) (Style, Style, Style, Style) {

	var (
		styleT = NewStyle().
			Foreground(borderStyle.GetBorderTopForeground()).
			Background(borderStyle.GetBorderTopBackground())

		styleR = NewStyle().
			Foreground(borderStyle.GetBorderRightForeground()).
			Background(borderStyle.GetBorderRightBackground())

		styleB = NewStyle().
			Foreground(borderStyle.GetBorderBottomForeground()).
			Background(borderStyle.GetBorderBottomBackground())

		styleL = NewStyle().
			Foreground(borderStyle.GetBorderLeftForeground()).
			Background(borderStyle.GetBorderLeftBackground())
	)

	return styleT, styleR, styleB, styleL
}
