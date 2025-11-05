package lipgloss

import (
	"image/color"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/clipperhouse/displaywidth"
)

// Style2 represents a rectangular area of styled text.
type Style2 struct {
	value string

	// Zero means automatic width based on content.
	width, height       int
	maxWidth, maxHeight int

	bg, fg color.Color

	// Text attributes
	style                                                   uv.Style
	noUnderlineSpaces, strikethroughSpaces, colorWhitespace bool

	paddingTop, paddingRight, paddingBottom, paddingLeft int
	paddingChar                                          rune

	marginTop, marginRight, marginBottom, marginLeft int
	marginChar                                       rune
	magrinBg                                         color.Color

	border                                                   Border
	borderTop, borderRight, borderBottom, borderLeft         bool
	borderTopFg, borderRightFg, borderBottomFg, borderLeftFg color.Color
	borderTopBg, borderRightBg, borderBottomBg, borderLeftBg color.Color
	borderFgBlender                                          BorderBlender

	alignHorizontal Position
	alignVertical   Position

	inline bool

	tabWidthOk bool
	tabWidth   int

	transform func(string) string
}

// NewStyle2 is a convenience function for creating a new [Style2].
func NewStyle2() Style2 {
	return Style2{}
}

// SetString sets a default string value for the [Style2]. This value will be
// prepended to any [Style2.Render] calls.
func (b Style2) SetString(s string) Style2 {
	b.value = s
	return b
}

// String renders the block to a string. This is a shorthand for [Style2.Render]
// that implements the [fmt.Stringer] interface.
func (b *Style2) String() string {
	return b.Render()
}

// Background sets the default background color for the block. Any characters
// rendered within the block that do not have a background color set will use
// this color.
func (b Style2) Background(col color.Color) Style2 {
	b.bg = col
	return b
}

// GetBackground returns the background color of the block.
func (b *Style2) GetBackground() color.Color {
	return b.bg
}

// UnsetBackground removes the background color from the block.
func (b Style2) UnsetBackground() Style2 {
	b.bg = nil
	return b
}

// Foreground sets the default foreground color for the block. Any characters
// rendered within the block that do not have a foreground color set will use
// this color.
func (b Style2) Foreground(col color.Color) Style2 {
	b.fg = col
	return b
}

// GetForeground returns the foreground color of the block.
func (b *Style2) GetForeground() color.Color {
	return b.fg
}

// UnsetForeground removes the foreground color from the block.
func (b Style2) UnsetForeground() Style2 {
	b.fg = nil
	return b
}

// Bold sets the bold attribute.
func (b Style2) Bold(v bool) Style2 {
	b.style = b.style.Bold(v)
	return b
}

// GetBold returns the bold attribute.
func (b *Style2) GetBold() bool {
	return b.style.Attrs.Contains(uv.BoldAttr)
}

// UnsetBold removes the bold attribute.
func (b Style2) UnsetBold() Style2 {
	b.style = b.style.Bold(false)
	return b
}

// Italic sets the italic attribute.
func (b Style2) Italic(v bool) Style2 {
	b.style = b.style.Italic(v)
	return b
}

// GetItalic returns the italic attribute.
func (b *Style2) GetItalic() bool {
	return b.style.Attrs.Contains(uv.ItalicAttr)
}

// UnsetItalic removes the italic attribute.
func (b Style2) UnsetItalic() Style2 {
	b.style = b.style.Italic(false)
	return b
}

// Underline sets the underline attribute.
func (b Style2) Underline(v bool) Style2 {
	us := uv.SingleUnderline
	if !v {
		us = uv.NoUnderline
	}
	b.style = b.style.UnderlineStyle(us)
	return b
}

// GetUnderline returns the underline attribute.
func (b *Style2) GetUnderline() bool {
	return b.style.UlStyle > uv.NoUnderline
}

// UnsetUnderline removes the underline attribute.
func (b Style2) UnsetUnderline() Style2 {
	b.style.UlStyle = uv.NoUnderline
	return b
}

// Strikethrough sets the strikethrough attribute.
func (b Style2) Strikethrough(v bool) Style2 {
	b.style = b.style.Strikethrough(v)
	return b
}

// GetStrikethrough returns the strikethrough attribute.
func (b *Style2) GetStrikethrough() bool {
	return b.style.Attrs.Contains(uv.StrikethroughAttr)
}

// UnsetStrikethrough removes the strikethrough attribute.
func (b Style2) UnsetStrikethrough() Style2 {
	b.style = b.style.Strikethrough(false)
	return b
}

// Reverse sets the reverse attribute.
func (b Style2) Reverse(v bool) Style2 {
	b.style = b.style.Reverse(v)
	return b
}

// GetReverse returns the reverse attribute.
func (b *Style2) GetReverse() bool {
	return b.style.Attrs.Contains(uv.ReverseAttr)
}

// UnsetReverse removes the reverse attribute.
func (b Style2) UnsetReverse() Style2 {
	b.style = b.style.Reverse(false)
	return b
}

// Blink sets the blink attribute.
func (b Style2) Blink(v bool) Style2 {
	b.style = b.style.SlowBlink(v)
	return b
}

// GetBlink returns the blink attribute.
func (b *Style2) GetBlink() bool {
	return b.style.Attrs.Contains(uv.SlowBlinkAttr)
}

// UnsetBlink removes the blink attribute.
func (b Style2) UnsetBlink() Style2 {
	b.style = b.style.SlowBlink(false)
	return b
}

// Faint sets the faint attribute.
func (b Style2) Faint(v bool) Style2 {
	b.style = b.style.Faint(v)
	return b
}

// GetFaint returns the faint attribute.
func (b *Style2) GetFaint() bool {
	return b.style.Attrs.Contains(uv.FaintAttr)
}

// UnsetFaint removes the faint attribute.
func (b Style2) UnsetFaint() Style2 {
	b.style = b.style.Faint(false)
	return b
}

// UnderlineSpaces sets whether to underline spaces.
func (b Style2) UnderlineSpaces(v bool) Style2 {
	b.noUnderlineSpaces = !v
	return b
}

// GetUnderlineSpaces returns whether spaces are underlined.
func (b *Style2) GetUnderlineSpaces() bool {
	return !b.noUnderlineSpaces
}

// UnsetUnderlineSpaces removes the underline spaces setting.
func (b Style2) UnsetUnderlineSpaces() Style2 {
	b.noUnderlineSpaces = true
	return b
}

// StrikethroughSpaces sets whether to strikethrough spaces.
func (b Style2) StrikethroughSpaces(v bool) Style2 {
	b.strikethroughSpaces = v
	return b
}

// GetStrikethroughSpaces returns whether spaces are struck through.
func (b *Style2) GetStrikethroughSpaces() bool {
	return b.strikethroughSpaces
}

// UnsetStrikethroughSpaces removes the strikethrough spaces setting.
func (b Style2) UnsetStrikethroughSpaces() Style2 {
	b.strikethroughSpaces = false
	return b
}

// ColorWhitespace sets whether to color whitespace.
func (b Style2) ColorWhitespace(v bool) Style2 {
	b.colorWhitespace = v
	return b
}

// GetColorWhitespace returns whether whitespace is colored.
func (b *Style2) GetColorWhitespace() bool {
	return b.colorWhitespace
}

// UnsetColorWhitespace removes the color whitespace setting.
func (b Style2) UnsetColorWhitespace() Style2 {
	b.colorWhitespace = false
	return b
}

// Inline sets whether the style is inline.
func (b Style2) Inline(v bool) Style2 {
	b.inline = v
	return b
}

// GetInline returns whether the style is inline.
func (b *Style2) GetInline() bool {
	return b.inline
}

// UnsetInline removes the inline setting.
func (b Style2) UnsetInline() Style2 {
	b.inline = false
	return b
}

// TabWidth sets the tab width.
func (b Style2) TabWidth(width int) Style2 {
	b.tabWidthOk = true
	b.tabWidth = width
	return b
}

// GetTabWidth returns the tab width.
func (b *Style2) GetTabWidth() int {
	if !b.tabWidthOk {
		return tabWidthDefault
	}
	return b.tabWidth
}

// UnsetTabWidth removes the tab width setting.
func (b Style2) UnsetTabWidth() Style2 {
	b.tabWidthOk = false
	b.tabWidth = 0
	return b
}

// Transform sets a transform function to apply to the rendered text.
func (b Style2) Transform(fn func(string) string) Style2 {
	b.transform = fn
	return b
}

// GetTransform returns the transform function.
func (b *Style2) GetTransform() func(string) string {
	return b.transform
}

// UnsetTransform removes the transform function.
func (b Style2) UnsetTransform() Style2 {
	b.transform = nil
	return b
}

// Width sets the width of the block before applying margins. This means your
// styled content will exactly equal the size set here. Text will wrap based on
// Padding and Borders set on the style.
func (b Style2) Width(width int) Style2 {
	b.width = width
	return b
}

// GetWidth returns the width of the block.
func (b *Style2) GetWidth() int {
	return b.width
}

// UnsetWidth removes the width constraint from the block.
func (b Style2) UnsetWidth() Style2 {
	b.width = 0
	return b
}

// Height sets the height of the block before applying margins. If the height of
// the text block is less than this value after applying padding (or not), the
// block will be set to this height.
func (b Style2) Height(height int) Style2 {
	b.height = height
	return b
}

// GetHeight returns the height of the block.
func (b *Style2) GetHeight() int {
	return b.height
}

// UnsetHeight removes the height constraint from the block.
func (b Style2) UnsetHeight() Style2 {
	b.height = 0
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
func (b Style2) MaxWidth(maxWidth int) Style2 {
	b.maxWidth = maxWidth
	return b
}

// GetMaxWidth returns the max width of the block.
func (b *Style2) GetMaxWidth() int {
	return b.maxWidth
}

// UnsetMaxWidth removes the max width constraint from the block.
func (b Style2) UnsetMaxWidth() Style2 {
	b.maxWidth = 0
	return b
}

// MaxHeight applies a max height to the block. This is useful in enforcing a
// certain height at render time, particularly with arbitrary strings and
// styles.
func (b Style2) MaxHeight(maxHeight int) Style2 {
	b.maxHeight = maxHeight
	return b
}

// GetMaxHeight returns the max height of the block.
func (b *Style2) GetMaxHeight() int {
	return b.maxHeight
}

// UnsetMaxHeight removes the max height constraint from the block.
func (b Style2) UnsetMaxHeight() Style2 {
	b.maxHeight = 0
	return b
}

// Align is a shorthand method for setting horizontal and vertical alignment.
//
// With one argument, the position value is applied to the horizontal alignment.
//
// With two arguments, the value is applied to the horizontal and vertical
// alignments, in that order.
func (b Style2) Align(positions ...Position) Style2 {
	switch len(positions) {
	case 1:
		b.alignHorizontal = positions[0]
	case 2:
		b.alignHorizontal = positions[0]
		b.alignVertical = positions[1]
	}
	return b
}

// GetAlign returns the horizontal and vertical alignment of the block.
func (b *Style2) GetAlign() (horizontal, vertical Position) {
	return b.alignHorizontal, b.alignVertical
}

// UnsetAlign removes both horizontal and vertical alignment from the block.
func (b Style2) UnsetAlign() Style2 {
	b.alignHorizontal = Position(0)
	b.alignVertical = Position(0)
	return b
}

// AlignHorizontal sets the horizontal alignment of the [Style2].
func (b Style2) AlignHorizontal(pos Position) Style2 {
	b.alignHorizontal = pos
	return b
}

// GetAlignHorizontal returns the horizontal alignment of the block.
func (b *Style2) GetAlignHorizontal() Position {
	return b.alignHorizontal
}

// UnsetAlignHorizontal removes the horizontal alignment from the block.
func (b Style2) UnsetAlignHorizontal() Style2 {
	b.alignHorizontal = Position(0)
	return b
}

// AlignVertical sets the vertical alignment of the [Style2].
func (b Style2) AlignVertical(pos Position) Style2 {
	b.alignVertical = pos
	return b
}

// GetAlignVertical returns the vertical alignment of the block.
func (b *Style2) GetAlignVertical() Position {
	return b.alignVertical
}

// UnsetAlignVertical removes the vertical alignment from the block.
func (b Style2) UnsetAlignVertical() Style2 {
	b.alignVertical = Position(0)
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
func (b Style2) Padding(values ...int) Style2 {
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

// GetPadding returns the padding for all sides of the block.
func (b *Style2) GetPadding() (top, right, bottom, left int) {
	return b.paddingTop, b.paddingRight, b.paddingBottom, b.paddingLeft
}

// UnsetPadding removes padding from all sides of the block.
func (b Style2) UnsetPadding() Style2 {
	b.paddingTop = 0
	b.paddingRight = 0
	b.paddingBottom = 0
	b.paddingLeft = 0
	return b
}

// PaddingLeft sets the left padding of the [Style2].
func (b Style2) PaddingLeft(padding int) Style2 {
	b.paddingLeft = padding
	return b
}

// GetPaddingLeft returns the left padding of the block.
func (b *Style2) GetPaddingLeft() int {
	return b.paddingLeft
}

// UnsetPaddingLeft removes the left padding from the block.
func (b Style2) UnsetPaddingLeft() Style2 {
	b.paddingLeft = 0
	return b
}

// PaddingRight sets the right padding of the [Style2].
func (b Style2) PaddingRight(padding int) Style2 {
	b.paddingRight = padding
	return b
}

// GetPaddingRight returns the right padding of the block.
func (b *Style2) GetPaddingRight() int {
	return b.paddingRight
}

// UnsetPaddingRight removes the right padding from the block.
func (b Style2) UnsetPaddingRight() Style2 {
	b.paddingRight = 0
	return b
}

// PaddingTop sets the top padding of the [Style2].
func (b Style2) PaddingTop(padding int) Style2 {
	b.paddingTop = padding
	return b
}

// GetPaddingTop returns the top padding of the block.
func (b *Style2) GetPaddingTop() int {
	return b.paddingTop
}

// UnsetPaddingTop removes the top padding from the block.
func (b Style2) UnsetPaddingTop() Style2 {
	b.paddingTop = 0
	return b
}

// PaddingBottom sets the bottom padding of the [Style2].
func (b Style2) PaddingBottom(padding int) Style2 {
	b.paddingBottom = padding
	return b
}

// GetPaddingBottom returns the bottom padding of the block.
func (b *Style2) GetPaddingBottom() int {
	return b.paddingBottom
}

// UnsetPaddingBottom removes the bottom padding from the block.
func (b Style2) UnsetPaddingBottom() Style2 {
	b.paddingBottom = 0
	return b
}

// PaddingChar sets the character used for padding. This is useful for
// rendering blocks with a specific character, such as a space or a dot.
// Example of using [NBSP] as padding to prevent line breaks:
//
//	```go
//	s := lipgloss.NewBlock().PaddingChar(lipgloss.NBSP)
//	```
func (b Style2) PaddingChar(r rune) Style2 {
	b.paddingChar = r
	return b
}

// GetPaddingChar returns the character used for padding.
func (b *Style2) GetPaddingChar() rune {
	return b.paddingChar
}

// UnsetPaddingChar removes the custom padding character from the block.
func (b Style2) UnsetPaddingChar() Style2 {
	b.paddingChar = 0
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
func (b Style2) Margin(values ...int) Style2 {
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

// GetMargin returns the margin for all sides of the block.
func (b *Style2) GetMargin() (top, right, bottom, left int) {
	return b.marginTop, b.marginRight, b.marginBottom, b.marginLeft
}

// UnsetMargins removes margins from all sides of the block.
func (b Style2) UnsetMargins() Style2 {
	b.marginTop = 0
	b.marginRight = 0
	b.marginBottom = 0
	b.marginLeft = 0
	return b
}

// MarginLeft sets the left margin of the [Style2].
func (b Style2) MarginLeft(margin int) Style2 {
	b.marginLeft = margin
	return b
}

// GetMarginLeft returns the left margin of the block.
func (b *Style2) GetMarginLeft() int {
	return b.marginLeft
}

// UnsetMarginLeft removes the left margin from the block.
func (b Style2) UnsetMarginLeft() Style2 {
	b.marginLeft = 0
	return b
}

// MarginRight sets the right margin of the [Style2].
func (b Style2) MarginRight(margin int) Style2 {
	b.marginRight = margin
	return b
}

// GetMarginRight returns the right margin of the block.
func (b *Style2) GetMarginRight() int {
	return b.marginRight
}

// UnsetMarginRight removes the right margin from the block.
func (b Style2) UnsetMarginRight() Style2 {
	b.marginRight = 0
	return b
}

// MarginTop sets the top margin of the [Style2].
func (b Style2) MarginTop(margin int) Style2 {
	b.marginTop = margin
	return b
}

// GetMarginTop returns the top margin of the block.
func (b *Style2) GetMarginTop() int {
	return b.marginTop
}

// UnsetMarginTop removes the top margin from the block.
func (b Style2) UnsetMarginTop() Style2 {
	b.marginTop = 0
	return b
}

// MarginBottom sets the bottom margin of the [Style2].
func (b Style2) MarginBottom(margin int) Style2 {
	b.marginBottom = margin
	return b
}

// GetMarginBottom returns the bottom margin of the block.
func (b *Style2) GetMarginBottom() int {
	return b.marginBottom
}

// UnsetMarginBottom removes the bottom margin from the block.
func (b Style2) UnsetMarginBottom() Style2 {
	b.marginBottom = 0
	return b
}

// MarginChar sets the character used for margin. This is useful for
// rendering blocks with a specific character, such as a space or a dot.
func (b Style2) MarginChar(r rune) Style2 {
	b.marginChar = r
	return b
}

// GetMarginChar returns the character used for margin.
func (b *Style2) GetMarginChar() rune {
	return b.marginChar
}

// UnsetMarginChar removes the custom margin character from the block.
func (b Style2) UnsetMarginChar() Style2 {
	b.marginChar = 0
	return b
}

// MarginBackground sets the background color for the margin area.
func (b Style2) MarginBackground(col color.Color) Style2 {
	b.magrinBg = col
	return b
}

// GetMarginBackground returns the background color for the margin area.
func (b *Style2) GetMarginBackground() color.Color {
	return b.magrinBg
}

// UnsetMarginBackground removes the margin background color from the block.
func (b Style2) UnsetMarginBackground() Style2 {
	b.magrinBg = nil
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
func (b Style2) Border(border Border, sides ...bool) Style2 {
	b.border = border

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
//	lipgloss.NewBlock().BorderStyle(lipgloss.ThickBorder())
func (b Style2) BorderStyle(border Border) Style2 {
	b.border = border
	return b
}

// GetBorderStyle returns the border style of the block.
func (b *Style2) GetBorderStyle() Border {
	return b.border
}

// GetBorder returns the border style and which sides have borders enabled.
func (b *Style2) GetBorder() (border Border, top, right, bottom, left bool) {
	return b.border, b.borderTop, b.borderRight, b.borderBottom, b.borderLeft
}

// UnsetBorderStyle removes the border style from the block.
func (b Style2) UnsetBorderStyle() Style2 {
	b.border = Border{}
	return b
}

// BorderTop sets whether the top border is visible.
func (b Style2) BorderTop(visible bool) Style2 {
	b.borderTop = visible
	return b
}

// GetBorderTop returns whether the top border is visible.
func (b *Style2) GetBorderTop() bool {
	return b.borderTop
}

// UnsetBorderTop disables the top border.
func (b Style2) UnsetBorderTop() Style2 {
	b.borderTop = false
	return b
}

// BorderRight sets whether the right border is visible.
func (b Style2) BorderRight(visible bool) Style2 {
	b.borderRight = visible
	return b
}

// GetBorderRight returns whether the right border is visible.
func (b *Style2) GetBorderRight() bool {
	return b.borderRight
}

// UnsetBorderRight disables the right border.
func (b Style2) UnsetBorderRight() Style2 {
	b.borderRight = false
	return b
}

// BorderBottom sets whether the bottom border is visible.
func (b Style2) BorderBottom(visible bool) Style2 {
	b.borderBottom = visible
	return b
}

// GetBorderBottom returns whether the bottom border is visible.
func (b *Style2) GetBorderBottom() bool {
	return b.borderBottom
}

// UnsetBorderBottom disables the bottom border.
func (b Style2) UnsetBorderBottom() Style2 {
	b.borderBottom = false
	return b
}

// BorderLeft sets whether the left border is visible.
func (b Style2) BorderLeft(visible bool) Style2 {
	b.borderLeft = visible
	return b
}

// GetBorderLeft returns whether the left border is visible.
func (b *Style2) GetBorderLeft() bool {
	return b.borderLeft
}

// UnsetBorderLeft disables the left border.
func (b Style2) UnsetBorderLeft() Style2 {
	b.borderLeft = false
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
func (b Style2) BorderForeground(colors ...color.Color) Style2 {
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
func (b Style2) BorderTopForeground(col color.Color) Style2 {
	b.borderTopFg = col
	return b
}

// GetBorderTopForeground returns the foreground color of the top border.
func (b *Style2) GetBorderTopForeground() color.Color {
	return b.borderTopFg
}

// UnsetBorderTopForeground removes the foreground color from the top border.
func (b Style2) UnsetBorderTopForeground() Style2 {
	b.borderTopFg = nil
	return b
}

// BorderRightForeground sets the foreground color of the right border.
func (b Style2) BorderRightForeground(col color.Color) Style2 {
	b.borderRightFg = col
	return b
}

// GetBorderRightForeground returns the foreground color of the right border.
func (b *Style2) GetBorderRightForeground() color.Color {
	return b.borderRightFg
}

// UnsetBorderRightForeground removes the foreground color from the right border.
func (b Style2) UnsetBorderRightForeground() Style2 {
	b.borderRightFg = nil
	return b
}

// BorderBottomForeground sets the foreground color of the bottom border.
func (b Style2) BorderBottomForeground(col color.Color) Style2 {
	b.borderBottomFg = col
	return b
}

// GetBorderBottomForeground returns the foreground color of the bottom border.
func (b *Style2) GetBorderBottomForeground() color.Color {
	return b.borderBottomFg
}

// UnsetBorderBottomForeground removes the foreground color from the bottom border.
func (b Style2) UnsetBorderBottomForeground() Style2 {
	b.borderBottomFg = nil
	return b
}

// BorderLeftForeground sets the foreground color of the left border.
func (b Style2) BorderLeftForeground(col color.Color) Style2 {
	b.borderLeftFg = col
	return b
}

// GetBorderLeftForeground returns the foreground color of the left border.
func (b *Style2) GetBorderLeftForeground() color.Color {
	return b.borderLeftFg
}

// UnsetBorderLeftForeground removes the foreground color from the left border.
func (b Style2) UnsetBorderLeftForeground() Style2 {
	b.borderLeftFg = nil
	return b
}

// UnsetBorderForeground removes the foreground color from all borders.
func (b Style2) UnsetBorderForeground() Style2 {
	b.borderTopFg = nil
	b.borderRightFg = nil
	b.borderBottomFg = nil
	b.borderLeftFg = nil
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
func (b Style2) BorderBackground(colors ...color.Color) Style2 {
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
func (b Style2) BorderTopBackground(col color.Color) Style2 {
	b.borderTopBg = col
	return b
}

// GetBorderTopBackground returns the background color of the top border.
func (b *Style2) GetBorderTopBackground() color.Color {
	return b.borderTopBg
}

// UnsetBorderTopBackground removes the background color from the top border.
func (b Style2) UnsetBorderTopBackground() Style2 {
	b.borderTopBg = nil
	return b
}

// BorderRightBackground sets the background color of the right border.
func (b Style2) BorderRightBackground(col color.Color) Style2 {
	b.borderRightBg = col
	return b
}

// GetBorderRightBackground returns the background color of the right border.
func (b *Style2) GetBorderRightBackground() color.Color {
	return b.borderRightBg
}

// UnsetBorderRightBackground removes the background color from the right border.
func (b Style2) UnsetBorderRightBackground() Style2 {
	b.borderRightBg = nil
	return b
}

// BorderBottomBackground sets the background color of the bottom border.
func (b Style2) BorderBottomBackground(col color.Color) Style2 {
	b.borderBottomBg = col
	return b
}

// GetBorderBottomBackground returns the background color of the bottom border.
func (b *Style2) GetBorderBottomBackground() color.Color {
	return b.borderBottomBg
}

// UnsetBorderBottomBackground removes the background color from the bottom border.
func (b Style2) UnsetBorderBottomBackground() Style2 {
	b.borderBottomBg = nil
	return b
}

// BorderLeftBackground sets the background color of the left border.
func (b Style2) BorderLeftBackground(col color.Color) Style2 {
	b.borderLeftBg = col
	return b
}

// GetBorderLeftBackground returns the background color of the left border.
func (b *Style2) GetBorderLeftBackground() color.Color {
	return b.borderLeftBg
}

// UnsetBorderLeftBackground removes the background color from the left border.
func (b Style2) UnsetBorderLeftBackground() Style2 {
	b.borderLeftBg = nil
	return b
}

// UnsetBorderBackground removes the background color from all borders.
func (b Style2) UnsetBorderBackground() Style2 {
	b.borderTopBg = nil
	b.borderRightBg = nil
	b.borderBottomBg = nil
	b.borderLeftBg = nil
	return b
}

// GetHorizontalMargins returns the style's left and right margins. Unset
// values are measured as 0.
func (b *Style2) GetHorizontalMargins() int {
	return b.marginLeft + b.marginRight
}

// GetVerticalMargins returns the style's top and bottom margins. Unset values
// are measured as 0.
func (b *Style2) GetVerticalMargins() int {
	return b.marginTop + b.marginBottom
}

// GetHorizontalPadding returns the style's left and right padding. Unset
// values are measured as 0.
func (b *Style2) GetHorizontalPadding() int {
	return b.paddingLeft + b.paddingRight
}

// GetVerticalPadding returns the style's top and bottom padding. Unset values
// are measured as 0.
func (b *Style2) GetVerticalPadding() int {
	return b.paddingTop + b.paddingBottom
}

// GetBorderTopSize returns the width of the top border. If borders contain
// runes of varying widths, the widest rune is returned. If no border exists on
// the top edge, 0 is returned.
func (b *Style2) GetBorderTopSize() int {
	return b.getBorderSize(&b.borderTop, &b.border.Top)
}

// GetBorderLeftSize returns the width of the left border. If borders contain
// runes of varying widths, the widest rune is returned. If no border exists on
// the left edge, 0 is returned.
func (b *Style2) GetBorderLeftSize() int {
	return b.getBorderSize(&b.borderLeft, &b.border.Left)
}

// GetBorderBottomSize returns the width of the bottom border. If borders
// contain runes of varying widths, the widest rune is returned. If no border
// exists on the left edge, 0 is returned.
func (b *Style2) GetBorderBottomSize() int {
	return b.getBorderSize(&b.borderBottom, &b.border.Bottom)
}

// GetBorderRightSize returns the width of the right border. If borders
// contain runes of varying widths, the widest rune is returned. If no border
// exists on the right edge, 0 is returned.
func (b *Style2) GetBorderRightSize() int {
	return b.getBorderSize(&b.borderRight, &b.border.Right)
}

// GetHorizontalBorderSize returns the width of the horizontal borders. If
// borders contain runes of varying widths, the widest rune is returned. If no
// border exists on the horizontal edges, 0 is returned.
func (b *Style2) GetHorizontalBorderSize() int {
	return b.GetBorderTopSize() + b.GetBorderBottomSize()
}

// GetVerticalBorderSize returns the width of the vertical borders. If
// borders contain runes of varying widths, the widest rune is returned. If no
// border exists on the vertical edges, 0 is returned.
func (b *Style2) GetVerticalBorderSize() int {
	return b.GetBorderLeftSize() + b.GetBorderRightSize()
}

// GetHorizontalFrameSize returns the sum of the style's horizontal margins, padding
// and border widths.
//
// Provisional: this method may be renamed.
func (b *Style2) GetHorizontalFrameSize() int {
	return b.GetHorizontalMargins() +
		b.GetHorizontalPadding() +
		b.GetHorizontalBorderSize()
}

// GetVerticalFrameSize returns the sum of the style's vertical margins, padding
// and border widths.
//
// Provisional: this method may be renamed.
func (b *Style2) GetVerticalFrameSize() int {
	return b.GetVerticalMargins() +
		b.GetVerticalPadding() +
		b.GetVerticalBorderSize()
}

// GetFrameSize returns the sum of the margins, padding and border width for
// both the horizontal and vertical margins.
func (b *Style2) GetFrameSize() (horizontal, vertical int) {
	return b.GetHorizontalFrameSize(), b.GetVerticalFrameSize()
}

// Draw renders the block to the screen within the given area. The area defines
// the total space available for the block. Any size constraints defined on the
// block itself are ignored; the block will fill the entire area.
func (b *Style2) Draw(scr uv.Screen, area uv.Rectangle) {
	value := b.value

	// Potentially convert tabs to spaces
	value = b.maybeConvertTabs(value)
	// carriage returns can cause strange behaviour when rendering.
	value = strings.ReplaceAll(value, "\r\n", "\n")

	// Strip newlines in single line mode
	if b.inline {
		value = strings.ReplaceAll(value, "\n", "")
	}

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

	border := b.border
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

	hasPadding := b.inline && b.paddingTop > 0 || b.paddingRight > 0 || b.paddingBottom > 0 || b.paddingLeft > 0
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

	// Render content
	var sb strings.Builder
	var state byte
	var seq string
	var w, n int

	spStyle := b.style
	if !b.strikethroughSpaces && spStyle.Attrs.Contains((uv.StrikethroughAttr)) {
		spStyle = spStyle.Strikethrough(false)
	}
	if b.noUnderlineSpaces && spStyle.UlStyle > 0 {
		spStyle = spStyle.UnderlineStyle(uv.NoUnderline)
	}

	// Skip per-sequence styling if the space-style and the base-style are
	// the same.
	if spStyle != b.style { //nolint:nestif
		str := value
		for len(str) > 0 {
			seq, w, n, state = ansi.DecodeSequence(str, state, nil)
			switch w {
			case 0:
				sb.WriteString(seq)
			case 1:
				if len(seq) == 1 && seq == " " {
					sb.WriteString(spStyle.Sequence() + seq + ansi.ResetStyle)
				} else {
					r, _ := utf8.DecodeRuneInString(seq)
					if unicode.IsSpace(r) {
						sb.WriteString(spStyle.Sequence() + seq + ansi.ResetStyle)
					} else {
						sb.WriteString(b.style.Sequence() + seq + ansi.ResetStyle)
					}
				}
			default:
				sb.WriteString(b.style.Sequence() + seq + ansi.ResetStyle)
			}
			str = str[n:]
		}

		value = sb.String()
	} else {
		value = b.style.Sequence() + value + ansi.ResetStyle
	}

	// Wrap text. This is going to take care of wrapping styles with
	// newlines as well.
	value = Wrap(value, contentArea.Dx(), "")

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Align text
		value = alignTextHorizontalWithMethod(value, b.alignHorizontal,
			contentArea.Dx(), nil, scr.WidthMethod())
		value = alignTextVertical(value, b.alignVertical,
			contentArea.Dy(), nil)

		// Content
		uv.NewStyledString(value).Draw(scr, contentArea)

		// Apply default foreground and background colors for the block
		// content.
		for y := contentArea.Min.Y; y < contentArea.Max.Y; y++ {
			for x := contentArea.Min.X; x < contentArea.Max.X; {
				cell := scr.CellAt(x, y)
				if cell == nil || cell.Width <= 0 {
					x++
					continue
				}
				if b.fg != nil && cell.Style.Fg == nil {
					cell.Style.Fg = b.fg
				}
				if b.bg != nil && cell.Style.Bg == nil {
					cell.Style.Bg = b.bg
				}
				x += cell.Width
			}
		}
	}()

	wg.Wait()
}

// Render renders the block to a string.
func (b Style2) Render(strs ...string) string {
	if len(strs) > 0 {
		b.value += strings.Join(strs, " ")

		// Truncate based on b.maxWidth and b.maxHeight if necessary.
		if b.maxWidth > 0 {
			lines := strings.Split(b.value, "\n")
			for i, line := range lines {
				lines[i] = ansi.Truncate(line, b.maxWidth, "")
			}
			b.value = strings.Join(lines, "\n")
		}
		if b.maxHeight > 0 {
			lines := strings.Split(b.value, "\n")
			height := min(b.maxHeight, len(lines))
			if len(lines) > 0 {
				b.value = strings.Join(lines[:height], "\n")
			}
		}
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
func (b *Style2) drawBorder(
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
func (b *Style2) getBorderSize(isSet *bool, content *string) int {
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
func (b *Style2) isBorderSetWithoutSides() bool {
	return b.border != (Border{}) && !b.borderTop && !b.borderRight && !b.borderBottom && !b.borderLeft
}

func (s Style2) maybeConvertTabs(str string) string {
	tw := tabWidthDefault
	if s.tabWidthOk {
		tw = s.tabWidth
	}
	switch tw {
	case -1:
		return str
	case 0:
		return strings.ReplaceAll(str, "\t", "")
	default:
		return strings.ReplaceAll(str, "\t", strings.Repeat(" ", tw))
	}
}
