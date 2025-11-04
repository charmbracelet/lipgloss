package lipgloss

import (
	"image/color"
	"strings"
	"sync"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/clipperhouse/displaywidth"
)

// Block represents a rectangular area of styled text.
type Block struct {
	value string

	// Zero means automatic width based on content.
	width, height       int
	maxWidth, maxHeight int

	paddingTop, paddingRight, paddingBottom, paddingLeft int
	paddingChar                                          rune

	marginTop, marginRight, marginBottom, marginLeft int
	marginChar                                       rune
	magrinBg                                         color.Color

	borderStyle                                              Border
	borderTop, borderRight, borderBottom, borderLeft         bool
	borderTopFg, borderRightFg, borderBottomFg, borderLeftFg color.Color
	borderTopBg, borderRightBg, borderBottomBg, borderLeftBg color.Color
	borderFgBlender                                          BorderBlender

	alignHorizontal Position
	alignVertical   Position
}

// NewBlock is a convenience function for creating a new [Block].
func NewBlock() Block {
	return Block{}
}

// SetString sets a default string value for the [Block]. This value will be
// prepended to any [Block.Render] calls.
func (b Block) SetString(s string) Block {
	b.value = s
	return b
}

// Width sets the width of the block before applying margins. This means your
// styled content will exactly equal the size set here. Text will wrap based on
// Padding and Borders set on the style.
func (b Block) Width(width int) Block {
	b.width = width
	return b
}

// Height sets the height of the block before applying margins. If the height of
// the text block is less than this value after applying padding (or not), the
// block will be set to this height.
func (b Block) Height(height int) Block {
	b.height = height
	return b
}

// MaxWidth applies a max width to the block. This is useful in enforcing a
// certain width at render time, particularly with arbitrary strings and
// styles.
//
// Example:
//
//	var userInput string = "..."
//	var userBlock lipgloss.Block
//	fmt.Println(userBlock.MaxWidth(16).Render(userInput))
func (b Block) MaxWidth(maxWidth int) Block {
	b.maxWidth = maxWidth
	return b
}

// MaxHeight applies a max height to the block. This is useful in enforcing a
// certain height at render time, particularly with arbitrary strings and
// styles.
func (b Block) MaxHeight(maxHeight int) Block {
	b.maxHeight = maxHeight
	return b
}

// Align is a shorthand method for setting horizontal and vertical alignment.
//
// With one argument, the position value is applied to the horizontal alignment.
//
// With two arguments, the value is applied to the horizontal and vertical
// alignments, in that order.
func (b Block) Align(positions ...Position) Block {
	switch len(positions) {
	case 1:
		b.alignHorizontal = positions[0]
	case 2:
		b.alignHorizontal = positions[0]
		b.alignVertical = positions[1]
	}
	return b
}

// AlignHorizontal sets the horizontal alignment of the [Block].
func (b Block) AlignHorizontal(pos Position) Block {
	b.alignHorizontal = pos
	return b
}

// AlignVertical sets the vertical alignment of the [Block].
func (b Block) AlignVertical(pos Position) Block {
	b.alignVertical = pos
	return b
}

// Padding is a shorthand method for setting padding on all sides at once.
//
// With one argument, the value is applied to all sides.
//
// With two arguments, the value is applied to the vertical and horizontal
// sides, in that order.
//
// With three arguments, the value is applied to the top side, the horizontal
// sides, and the bottom side, in that order.
//
// With four arguments, the value is applied clockwise starting from the top
// side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments no padding will be added.
func (b Block) Padding(values ...int) Block {
	switch len(values) {
	case 1:
		b.paddingTop = values[0]
		b.paddingRight = values[0]
		b.paddingBottom = values[0]
		b.paddingLeft = values[0]
	case 2:
		b.paddingTop = values[0]
		b.paddingRight = values[1]
		b.paddingBottom = values[0]
		b.paddingLeft = values[1]
	case 3:
		b.paddingTop = values[0]
		b.paddingRight = values[1]
		b.paddingBottom = values[2]
		b.paddingLeft = values[1]
	case 4:
		b.paddingTop = values[0]
		b.paddingRight = values[1]
		b.paddingBottom = values[2]
		b.paddingLeft = values[3]
	}
	return b
}

// PaddingLeft sets the left padding of the [Block].
func (b Block) PaddingLeft(padding int) Block {
	b.paddingLeft = padding
	return b
}

// PaddingRight sets the right padding of the [Block].
func (b Block) PaddingRight(padding int) Block {
	b.paddingRight = padding
	return b
}

// PaddingTop sets the top padding of the [Block].
func (b Block) PaddingTop(padding int) Block {
	b.paddingTop = padding
	return b
}

// PaddingBottom sets the bottom padding of the [Block].
func (b Block) PaddingBottom(padding int) Block {
	b.paddingBottom = padding
	return b
}

// PaddingChar sets the character used for padding. This is useful for
// rendering blocks with a specific character, such as a space or a dot.
// Example of using [NBSP] as padding to prevent line breaks:
//
//	```go
//	s := lipgloss.NewBlock().PaddingChar(lipgloss.NBSP)
//	```
func (b Block) PaddingChar(r rune) Block {
	b.paddingChar = r
	return b
}

// Margin is a shorthand method for setting margins on all sides at once.
//
// With one argument, the value is applied to all sides.
//
// With two arguments, the value is applied to the vertical and horizontal
// sides, in that order.
//
// With three arguments, the value is applied to the top side, the horizontal
// sides, and the bottom side, in that order.
//
// With four arguments, the value is applied clockwise starting from the top
// side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments no margin will be added.
func (b Block) Margin(values ...int) Block {
	switch len(values) {
	case 1:
		b.marginTop = values[0]
		b.marginRight = values[0]
		b.marginBottom = values[0]
		b.marginLeft = values[0]
	case 2:
		b.marginTop = values[0]
		b.marginRight = values[1]
		b.marginBottom = values[0]
		b.marginLeft = values[1]
	case 3:
		b.marginTop = values[0]
		b.marginRight = values[1]
		b.marginBottom = values[2]
		b.marginLeft = values[1]
	case 4:
		b.marginTop = values[0]
		b.marginRight = values[1]
		b.marginBottom = values[2]
		b.marginLeft = values[3]
	}
	return b
}

// MarginLeft sets the left margin of the [Block].
func (b Block) MarginLeft(margin int) Block {
	b.marginLeft = margin
	return b
}

// MarginRight sets the right margin of the [Block].
func (b Block) MarginRight(margin int) Block {
	b.marginRight = margin
	return b
}

// MarginTop sets the top margin of the [Block].
func (b Block) MarginTop(margin int) Block {
	b.marginTop = margin
	return b
}

// MarginBottom sets the bottom margin of the [Block].
func (b Block) MarginBottom(margin int) Block {
	b.marginBottom = margin
	return b
}

// MarginChar sets the character used for margin. This is useful for
// rendering blocks with a specific character, such as a space or a dot.
func (b Block) MarginChar(r rune) Block {
	b.marginChar = r
	return b
}

// MarginBackground sets the background color for the margin area.
func (b Block) MarginBackground(col color.Color) Block {
	b.magrinBg = col
	return b
}

// Border is shorthand for setting the border style and which sides should
// have a border at once. The variadic argument sides works as follows:
//
// With one value, the value is applied to all sides.
//
// With two values, the values are applied to the vertical and horizontal
// sides, in that order.
//
// With three values, the values are applied to the top side, the horizontal
// sides, and the bottom side, in that order.
//
// With four values, the values are applied clockwise starting from the top
// side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments the border will be applied to all sides.
//
// Examples:
//
//	// Applies borders to the top and bottom only
//	lipgloss.NewBlock().Border(lipgloss.NormalBorder(), true, false)
//
//	// Applies rounded borders to the right and bottom only
//	lipgloss.NewBlock().Border(lipgloss.RoundedBorder(), false, true, true, false)
func (b Block) Border(border Border, sides ...bool) Block {
	b.borderStyle = border

	switch len(sides) {
	case 1:
		b.borderTop = sides[0]
		b.borderRight = sides[0]
		b.borderBottom = sides[0]
		b.borderLeft = sides[0]
	case 2:
		b.borderTop = sides[0]
		b.borderRight = sides[1]
		b.borderBottom = sides[0]
		b.borderLeft = sides[1]
	case 3:
		b.borderTop = sides[0]
		b.borderRight = sides[1]
		b.borderBottom = sides[2]
		b.borderLeft = sides[1]
	case 4:
		b.borderTop = sides[0]
		b.borderRight = sides[1]
		b.borderBottom = sides[2]
		b.borderLeft = sides[3]
	default:
		b.borderTop = true
		b.borderRight = true
		b.borderBottom = true
		b.borderLeft = true
	}

	return b
}

// BorderStyle defines the Border on a style. A Border contains a series of
// definitions for the sides and corners of a border.
//
// Note that if border visibility has not been set for any sides when setting
// the border style, the border will be enabled for all sides during rendering.
//
// You can define border characters as you'd like, though several default
// styles are included: NormalBorder(), RoundedBorder(), BlockBorder(),
// OuterHalfBlockBorder(), InnerHalfBlockBorder(), ThickBorder(),
// and DoubleBorder().
//
// Example:
//
//	lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder())
func (b Block) BorderStyle(border Border) Block {
	b.borderStyle = border
	return b
}

// BorderTop sets whether the top border is visible.
func (b Block) BorderTop(visible bool) Block {
	b.borderTop = visible
	return b
}

// BorderRight sets whether the right border is visible.
func (b Block) BorderRight(visible bool) Block {
	b.borderRight = visible
	return b
}

// BorderBottom sets whether the bottom border is visible.
func (b Block) BorderBottom(visible bool) Block {
	b.borderBottom = visible
	return b
}

// BorderLeft sets whether the left border is visible.
func (b Block) BorderLeft(visible bool) Block {
	b.borderLeft = visible
	return b
}

// BorderForeground is a shorthand function for setting all of the
// foreground colors of the borders at once. The arguments work as follows:
//
// With one argument, the argument is applied to all sides.
//
// With two arguments, the arguments are applied to the vertical and horizontal
// sides, in that order.
//
// With three arguments, the arguments are applied to the top side, the
// horizontal sides, and the bottom side, in that order.
//
// With four arguments, the arguments are applied clockwise starting from the
// top side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments nothing will be set.
func (b Block) BorderForeground(colors ...color.Color) Block {
	switch len(colors) {
	case 1:
		b.borderTopFg = colors[0]
		b.borderRightFg = colors[0]
		b.borderBottomFg = colors[0]
		b.borderLeftFg = colors[0]
	case 2:
		b.borderTopFg = colors[0]
		b.borderRightFg = colors[1]
		b.borderBottomFg = colors[0]
		b.borderLeftFg = colors[1]
	case 3:
		b.borderTopFg = colors[0]
		b.borderRightFg = colors[1]
		b.borderBottomFg = colors[2]
		b.borderLeftFg = colors[1]
	case 4:
		b.borderTopFg = colors[0]
		b.borderRightFg = colors[1]
		b.borderBottomFg = colors[2]
		b.borderLeftFg = colors[3]
	}
	return b
}

// BorderTopForeground sets the foreground color of the top border.
func (b Block) BorderTopForeground(col color.Color) Block {
	b.borderTopFg = col
	return b
}

// BorderRightForeground sets the foreground color of the right border.
func (b Block) BorderRightForeground(col color.Color) Block {
	b.borderRightFg = col
	return b
}

// BorderBottomForeground sets the foreground color of the bottom border.
func (b Block) BorderBottomForeground(col color.Color) Block {
	b.borderBottomFg = col
	return b
}

// BorderLeftForeground sets the foreground color of the left border.
func (b Block) BorderLeftForeground(col color.Color) Block {
	b.borderLeftFg = col
	return b
}

// BorderBackground is a shorthand function for setting all of the
// background colors of the borders at once. The arguments work as follows:
//
// With one argument, the argument is applied to all sides.
//
// With two arguments, the arguments are applied to the vertical and horizontal
// sides, in that order.
//
// With three arguments, the arguments are applied to the top side, the
// horizontal sides, and the bottom side, in that order.
//
// With four arguments, the arguments are applied clockwise starting from the
// top side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments nothing will be set.
func (b Block) BorderBackground(colors ...color.Color) Block {
	switch len(colors) {
	case 1:
		b.borderTopBg = colors[0]
		b.borderRightBg = colors[0]
		b.borderBottomBg = colors[0]
		b.borderLeftBg = colors[0]
	case 2:
		b.borderTopBg = colors[0]
		b.borderRightBg = colors[1]
		b.borderBottomBg = colors[0]
		b.borderLeftBg = colors[1]
	case 3:
		b.borderTopBg = colors[0]
		b.borderRightBg = colors[1]
		b.borderBottomBg = colors[2]
		b.borderLeftBg = colors[1]
	case 4:
		b.borderTopBg = colors[0]
		b.borderRightBg = colors[1]
		b.borderBottomBg = colors[2]
		b.borderLeftBg = colors[3]
	}
	return b
}

// BorderTopBackground sets the background color of the top border.
func (b Block) BorderTopBackground(col color.Color) Block {
	b.borderTopBg = col
	return b
}

// BorderRightBackground sets the background color of the right border.
func (b Block) BorderRightBackground(col color.Color) Block {
	b.borderRightBg = col
	return b
}

// BorderBottomBackground sets the background color of the bottom border.
func (b Block) BorderBottomBackground(col color.Color) Block {
	b.borderBottomBg = col
	return b
}

// BorderLeftBackground sets the background color of the left border.
func (b Block) BorderLeftBackground(col color.Color) Block {
	b.borderLeftBg = col
	return b
}

// GetHorizontalMargins returns the style's left and right margins. Unset
// values are measured as 0.
func (b *Block) GetHorizontalMargins() int {
	return b.marginLeft + b.marginRight
}

// GetVerticalMargins returns the style's top and bottom margins. Unset values
// are measured as 0.
func (b *Block) GetVerticalMargins() int {
	return b.marginTop + b.marginBottom
}

// GetHorizontalPadding returns the style's left and right padding. Unset
// values are measured as 0.
func (b *Block) GetHorizontalPadding() int {
	return b.paddingLeft + b.paddingRight
}

// GetVerticalPadding returns the style's top and bottom padding. Unset values
// are measured as 0.
func (b *Block) GetVerticalPadding() int {
	return b.paddingTop + b.paddingBottom
}

// GetBorderTopSize returns the width of the top border. If borders contain
// runes of varying widths, the widest rune is returned. If no border exists on
// the top edge, 0 is returned.
func (b *Block) GetBorderTopSize() int {
	return b.getBorderSize(&b.borderTop, &b.borderStyle.Top)
}

// GetBorderLeftSize returns the width of the left border. If borders contain
// runes of varying widths, the widest rune is returned. If no border exists on
// the left edge, 0 is returned.
func (b *Block) GetBorderLeftSize() int {
	return b.getBorderSize(&b.borderLeft, &b.borderStyle.Left)
}

// GetBorderBottomSize returns the width of the bottom border. If borders
// contain runes of varying widths, the widest rune is returned. If no border
// exists on the left edge, 0 is returned.
func (b *Block) GetBorderBottomSize() int {
	return b.getBorderSize(&b.borderBottom, &b.borderStyle.Bottom)
}

// GetBorderRightSize returns the width of the right border. If borders
// contain runes of varying widths, the widest rune is returned. If no border
// exists on the right edge, 0 is returned.
func (b *Block) GetBorderRightSize() int {
	return b.getBorderSize(&b.borderRight, &b.borderStyle.Right)
}

// GetHorizontalBorderSize returns the width of the horizontal borders. If
// borders contain runes of varying widths, the widest rune is returned. If no
// border exists on the horizontal edges, 0 is returned.
func (b *Block) GetHorizontalBorderSize() int {
	return b.GetBorderTopSize() + b.GetBorderBottomSize()
}

// GetVerticalBorderSize returns the width of the vertical borders. If
// borders contain runes of varying widths, the widest rune is returned. If no
// border exists on the vertical edges, 0 is returned.
func (b *Block) GetVerticalBorderSize() int {
	return b.GetBorderLeftSize() + b.GetBorderRightSize()
}

// GetHorizontalFrameSize returns the sum of the style's horizontal margins, padding
// and border widths.
//
// Provisional: this method may be renamed.
func (b *Block) GetHorizontalFrameSize() int {
	return b.GetHorizontalMargins() +
		b.GetHorizontalPadding() +
		b.GetHorizontalBorderSize()
}

// GetVerticalFrameSize returns the sum of the style's vertical margins, padding
// and border widths.
//
// Provisional: this method may be renamed.
func (b *Block) GetVerticalFrameSize() int {
	return b.GetVerticalMargins() +
		b.GetVerticalPadding() +
		b.GetVerticalBorderSize()
}

// GetFrameSize returns the sum of the margins, padding and border width for
// both the horizontal and vertical margins.
func (b *Block) GetFrameSize() (horizontal, vertical int) {
	return b.GetHorizontalFrameSize(), b.GetVerticalFrameSize()
}

// Draw renders the block to the screen within the given area. The area defines
// the total space available for the block. Any size constraints defined on the
// block itself are ignored; the block will fill the entire area.
func (b *Block) Draw(scr uv.Screen, area uv.Rectangle) {
	value := b.value

	// Margin area is the full area
	marginArea := area
	hasMargins := b.marginTop > 0 || b.marginRight > 0 || b.marginBottom > 0 || b.marginLeft > 0

	// Inner area after margins is the border area
	borderArea := marginArea
	if b.marginLeft > 0 {
		borderArea.Min.X += b.marginLeft
	}
	if b.marginRight > 0 {
		borderArea.Max.X -= b.marginRight
	}
	if b.marginTop > 0 {
		borderArea.Min.Y += b.marginTop
	}
	if b.marginBottom > 0 {
		borderArea.Max.Y -= b.marginBottom
	}

	border := b.borderStyle
	borderTop := b.borderTop
	borderRight := b.borderRight
	borderBottom := b.borderBottom
	borderLeft := b.borderLeft
	if border != (Border{}) && !borderTop && !borderRight && !borderBottom && !borderLeft {
		// Border set without sides enabled defaults to all sides.
		borderTop = true
		borderRight = true
		borderBottom = true
		borderLeft = true
	}
	hasBorders := border != Border{} && (borderTop || borderRight || borderBottom || borderLeft)
	borderTopEdgeWidth := b.getBorderSize(&borderTop, &border.Top)
	borderBottomEdgeWidth := b.getBorderSize(&borderBottom, &border.Bottom)
	borderLeftEdgeWidth := b.getBorderSize(&borderLeft, &border.Left)
	borderRightEdgeWidth := b.getBorderSize(&borderRight, &border.Right)

	// Padding area is the area after margins and borders
	paddingArea := borderArea
	if borderLeft {
		if border.Left == "" {
			border.Left = " "
		}
		paddingArea.Min.X += borderLeftEdgeWidth
	}
	if borderRight {
		if border.Right == "" {
			border.Right = " "
		}
		paddingArea.Max.X -= borderRightEdgeWidth
	}
	if borderTop {
		if border.Top == "" {
			border.Top = " "
		}
		paddingArea.Min.Y += borderTopEdgeWidth
	}
	if borderBottom {
		if border.Bottom == "" {
			border.Bottom = " "
		}
		paddingArea.Max.Y -= borderBottomEdgeWidth
	}
	if borderTop && borderLeft && border.TopLeft == "" {
		border.TopLeft = " "
	}
	if borderTop && borderRight && border.TopRight == "" {
		border.TopRight = " "
	}
	if borderBottom && borderLeft && border.BottomLeft == "" {
		border.BottomLeft = " "
	}
	if borderBottom && borderRight && border.BottomRight == "" {
		border.BottomRight = " "
	}

	hasPadding := b.paddingTop > 0 || b.paddingRight > 0 || b.paddingBottom > 0 || b.paddingLeft > 0
	// Content area is the area after margins, borders, and padding
	contentArea := paddingArea
	if b.paddingLeft > 0 {
		contentArea.Min.X += b.paddingLeft
	}
	if b.paddingRight > 0 {
		contentArea.Max.X -= b.paddingRight
	}
	if b.paddingTop > 0 {
		contentArea.Min.Y += b.paddingTop
	}
	if b.paddingBottom > 0 {
		contentArea.Max.Y -= b.paddingBottom
	}

	// Render loop.
	var wg sync.WaitGroup

	// Margins
	if hasMargins {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for y := marginArea.Min.Y; y < marginArea.Max.Y; y++ {
				for x := marginArea.Min.X; x < marginArea.Max.X; x++ {
					if x >= marginArea.Min.X && x < marginArea.Max.X &&
						y >= marginArea.Min.Y && y < marginArea.Max.Y {
						continue
					}
					scr.SetCell(x, y, &uv.Cell{
						Content: " ",
						Style: uv.Style{
							Bg: b.magrinBg,
						},
						Width: 1,
					})
				}
			}
		}()
	}

	// Borders
	if hasBorders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.drawBorder(scr, borderArea, &border,
				borderTop, borderRight,
				borderBottom, borderLeft,
				borderTopEdgeWidth, borderBottomEdgeWidth,
				borderLeftEdgeWidth, borderRightEdgeWidth,
			)
		}()
	}

	// Padding
	if hasPadding {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fill := " "
			if b.paddingChar != 0 {
				fill = string(b.paddingChar)
			}
			for y := paddingArea.Min.Y; y < paddingArea.Max.Y; y++ {
				for x := paddingArea.Min.X; x < paddingArea.Max.X; x++ {
					if x >= contentArea.Min.X && x < contentArea.Max.X && y >= contentArea.Min.Y && y < contentArea.Max.Y {
						continue
					}
					scr.SetCell(x, y, &uv.Cell{Content: fill, Width: 1})
				}
			}
		}()
	}

	// Align text
	value = alignTextHorizontalWithMethod(value, b.alignHorizontal,
		contentArea.Dx(), nil, scr.WidthMethod())
	value = alignTextVertical(value, b.alignVertical,
		contentArea.Dy(), nil)

	uv.NewStyledString(value).Draw(scr, contentArea)

	wg.Wait()
}

// Render renders the block to a string.
func (b *Block) Render(strs ...string) string {
	if len(strs) > 0 {
		b.value += strings.Join(strs, " ")
	}

	blockWidth, blockHeight := b.width, b.height
	contentWidth, contentHeight := Width(b.value), Height(b.value)

	// Adjust block size based on content size if necessary.
	blockWidth = max(blockWidth, contentWidth+b.GetHorizontalPadding()+b.GetVerticalBorderSize())
	blockHeight = max(blockHeight, contentHeight+b.GetVerticalPadding()+b.GetHorizontalBorderSize())

	// Compute total dimensions including margins. The sizes of borders and
	// padding are included in contentWidth and contentHeight.
	totalWidth := blockWidth + b.GetHorizontalMargins()
	totalHeight := blockHeight + b.GetVerticalMargins()
	if b.maxWidth > 0 {
		totalWidth = min(totalWidth, b.maxWidth)
	}
	if b.maxHeight > 0 {
		totalHeight = min(totalHeight, b.maxHeight)
	}

	scr := uv.NewScreenBuffer(totalWidth, totalHeight)
	b.Draw(scr, scr.Bounds())

	var buf strings.Builder
	for i, l := range scr.Lines {
		buf.WriteString(l.Render())
		if i < len(scr.Lines)-1 {
			_, _ = buf.WriteString("\n")
		}
	}

	return buf.String()
}

// drawBorder draws the border onto the given screen within the given area.
func (b *Block) drawBorder(
	scr uv.Screen, borderArea uv.Rectangle, border *Border,
	borderTop, borderRight, borderBottom, borderLeft bool,
	borderTopEdgeWidth, borderBottomEdgeWidth, borderLeftEdgeWidth, borderRightEdgeWidth int,
) {
	borderTopLeftWidth := b.getBorderSize(&borderTop, &border.TopLeft)
	borderTopRightWidth := b.getBorderSize(&borderRight, &border.TopRight)
	borderBottomLeftWidth := b.getBorderSize(&borderLeft, &border.BottomLeft)
	borderBottomRightWidth := b.getBorderSize(&borderBottom, &border.BottomRight)

	// Precompute styles for each side
	topStyle := uv.Style{Fg: b.borderTopFg, Bg: b.borderTopBg}
	rightStyle := uv.Style{Fg: b.borderRightFg, Bg: b.borderRightBg}
	bottomStyle := uv.Style{Fg: b.borderBottomFg, Bg: b.borderBottomBg}
	leftStyle := uv.Style{Fg: b.borderLeftFg, Bg: b.borderLeftBg}

	for y := borderArea.Min.Y; y < borderArea.Max.Y; y++ {
		for x := borderArea.Min.X; x < borderArea.Max.X; {
			switch {
			case x == borderArea.Min.X && y == borderArea.Min.Y && borderTop && borderLeft:
				// Top-left corner (use top colors)
				scr.SetCell(x, y, &uv.Cell{Content: border.TopLeft, Style: topStyle, Width: borderTopLeftWidth})
				x += borderTopLeftWidth
			case x == borderArea.Max.X-1 && y == borderArea.Min.Y && borderTop && borderRight:
				// Top-right corner (use top colors)
				scr.SetCell(x, y, &uv.Cell{Content: border.TopRight, Style: topStyle, Width: borderTopRightWidth})
				x += borderTopRightWidth
			case x == borderArea.Min.X && y == borderArea.Max.Y-1 && borderBottom && borderLeft:
				// Bottom-left corner (use bottom colors)
				scr.SetCell(x, y, &uv.Cell{Content: border.BottomLeft, Style: bottomStyle, Width: borderBottomLeftWidth})
				x += borderBottomLeftWidth
			case x == borderArea.Max.X-1 && y == borderArea.Max.Y-1 && borderBottom && borderRight:
				// Bottom-right corner (use bottom colors)
				scr.SetCell(x, y, &uv.Cell{Content: border.BottomRight, Style: bottomStyle, Width: borderBottomRightWidth})
				x += borderBottomRightWidth
			case y == borderArea.Min.Y && borderTop:
				// Top edge
				scr.SetCell(x, y, &uv.Cell{Content: border.Top, Style: topStyle, Width: borderTopEdgeWidth})
				x += borderTopEdgeWidth
			case y == borderArea.Max.Y-1 && borderBottom:
				// Bottom edge
				scr.SetCell(x, y, &uv.Cell{Content: border.Bottom, Style: bottomStyle, Width: borderBottomEdgeWidth})
				x += borderBottomEdgeWidth
			case x == borderArea.Min.X && borderLeft:
				// Left edge
				scr.SetCell(x, y, &uv.Cell{Content: border.Left, Style: leftStyle, Width: borderLeftEdgeWidth})
				x += borderLeftEdgeWidth
			case x == borderArea.Max.X-1 && borderRight:
				// Right edge
				scr.SetCell(x, y, &uv.Cell{Content: border.Right, Style: rightStyle, Width: borderRightEdgeWidth})
				x += borderRightEdgeWidth
			default:
				x++
			}
		}
	}
}

// getBorderSize returns the width of the border side.
func (b *Block) getBorderSize(isSet *bool, content *string) int {
	if b.isBorderSetWithoutSides() {
		return 1
	}
	if !*isSet {
		return 0
	}
	if len(*content) == 0 {
		return 1
	}
	if len(*content) > 1 {
		return displaywidth.String(*content)
	}
	return 1
}

// isBorderStyleSetWithoutSides returns true if the border style is set but no
// sides are set. This is used to determine if the border should be rendered by
// default.
func (b *Block) isBorderSetWithoutSides() bool {
	return b.borderStyle != (Border{}) && !b.borderTop && !b.borderRight && !b.borderBottom && !b.borderLeft
}
