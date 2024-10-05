package lipgloss

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"github.com/rivo/uniseg"
)

// Border contains a series of values which comprise the various parts of a
// border.
type Border struct {
	Top          string
	Bottom       string
	Left         string
	Right        string
	TopLeft      string
	TopRight     string
	BottomLeft   string
	BottomRight  string
	MiddleLeft   string
	MiddleRight  string
	Middle       string
	MiddleTop    string
	MiddleBottom string
}

// BorderHorizontalFunc is border function that sets horizontal border text
// at the configured position.
//
// It takes the width of the border and the Top/Bottom border string
// and determines the string for that position.
//
// Example:
//
//	reverseStyle := lipgloss.NewStyle().Reverse(true)
//	t := lipgloss.NewStyle().
//	    Border(lipgloss.NormalBorder()).
//	    BorderTopFunc(lipgloss.Center, func(w int, m string) string {
//	        return reverseStyle.Render(" BIG TITLE ")
//	    }).
//	    BorderBottomFunc(lipgloss.Right, func(width int, middle string) string {
//	        return reverseStyle.Render(fmt.Sprintf(" %d/%d ", m.index + 1, m.count)) + middle
//	    }).
//	    BorderBottomFunc(lipgloss.Left, func(width int, middle string) string {
//	        return middle + reverseStyle.Render(fmt.Sprintf("Status: %s", m.status))
//	    })
type BorderHorizontalFunc interface {
	func(width int, middle string) string
}

// BorderDecoration is type used by Border to set text or decorate the border.
type BorderDecoration struct {
	side  Position
	align Position
	st    interface{}
}

// BorderDecorator is constraint type for a string or function that is used
// to decorate a border.
type BorderDecorator interface {
	string | func() string | BorderHorizontalFunc
}

// NewBorderDecoration is function that sets creates a decoration for the border.
//
// It takes the side of the border (Top|Bottom), the alignment (Left|Center|Right) of the
// decoration, and the decoration.
//
// the decoration can be any of
//   - string
//   - func() string
//   - func(width int, middle string) string  where width is the size of the border and middle
//     is the border string
//
// Example:
//
//	reverseStyle := lipgloss.NewStyle().Reverse(true)
//
//	t := lipgloss.NewStyle().
//	    Border(lipgloss.NormalBorder()).
//	    BorderDecoration(lipgloss.NewBorderDecoration(
//	        lipgloss.Top,
//	        lipgloss.Center,
//	        reverseStyle.Padding(0, 1).Render("BIG TITLE"),
//	    )).
//	    BorderDecoration(lipgloss.NewBorderDecoration(
//	        lipgloss.Bottom,
//	        lipgloss.Right,
//	        func(width int, middle string) string {
//	            return reverseStyle.Render(fmt.Sprintf(" %d/%d ", m.index + 1, m.count)) + middle
//	        },
//	    )).
//	    BorderDecoration(lipgloss.NewBorderDecoration(
//	        lipgloss.Bottom,
//	        lipgloss.Left,
//	        reverseStyle.SetString(fmt.Sprintf("Status: %s", m.status)).String,
//	    ))
func NewBorderDecoration[S BorderDecorator](side, align Position, st S) BorderDecoration {
	return BorderDecoration{
		align: align,
		side:  side,
		st:    st,
	}
}

// GetTopSize returns the width of the top border. If borders contain runes of
// varying widths, the widest rune is returned. If no border exists on the top
// edge, 0 is returned.
func (b Border) GetTopSize() int {
	return getBorderEdgeWidth(b.TopLeft, b.Top, b.TopRight)
}

// GetRightSize returns the width of the right border. If borders contain
// runes of varying widths, the widest rune is returned. If no border exists on
// the right edge, 0 is returned.
func (b Border) GetRightSize() int {
	return getBorderEdgeWidth(b.TopRight, b.Right, b.BottomRight)
}

// GetBottomSize returns the width of the bottom border. If borders contain
// runes of varying widths, the widest rune is returned. If no border exists on
// the bottom edge, 0 is returned.
func (b Border) GetBottomSize() int {
	return getBorderEdgeWidth(b.BottomLeft, b.Bottom, b.BottomRight)
}

// GetLeftSize returns the width of the left border. If borders contain runes
// of varying widths, the widest rune is returned. If no border exists on the
// left edge, 0 is returned.
func (b Border) GetLeftSize() int {
	return getBorderEdgeWidth(b.TopLeft, b.Left, b.BottomLeft)
}

func getBorderEdgeWidth(borderParts ...string) (maxWidth int) {
	for _, piece := range borderParts {
		w := maxRuneWidth(piece)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

var (
	noBorder = Border{}

	normalBorder = Border{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "┌",
		TopRight:     "┐",
		BottomLeft:   "└",
		BottomRight:  "┘",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}

	roundedBorder = Border{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "╰",
		BottomRight:  "╯",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}

	blockBorder = Border{
		Top:         "█",
		Bottom:      "█",
		Left:        "█",
		Right:       "█",
		TopLeft:     "█",
		TopRight:    "█",
		BottomLeft:  "█",
		BottomRight: "█",
	}

	outerHalfBlockBorder = Border{
		Top:         "▀",
		Bottom:      "▄",
		Left:        "▌",
		Right:       "▐",
		TopLeft:     "▛",
		TopRight:    "▜",
		BottomLeft:  "▙",
		BottomRight: "▟",
	}

	innerHalfBlockBorder = Border{
		Top:         "▄",
		Bottom:      "▀",
		Left:        "▐",
		Right:       "▌",
		TopLeft:     "▗",
		TopRight:    "▖",
		BottomLeft:  "▝",
		BottomRight: "▘",
	}

	thickBorder = Border{
		Top:          "━",
		Bottom:       "━",
		Left:         "┃",
		Right:        "┃",
		TopLeft:      "┏",
		TopRight:     "┓",
		BottomLeft:   "┗",
		BottomRight:  "┛",
		MiddleLeft:   "┣",
		MiddleRight:  "┫",
		Middle:       "╋",
		MiddleTop:    "┳",
		MiddleBottom: "┻",
	}

	doubleBorder = Border{
		Top:          "═",
		Bottom:       "═",
		Left:         "║",
		Right:        "║",
		TopLeft:      "╔",
		TopRight:     "╗",
		BottomLeft:   "╚",
		BottomRight:  "╝",
		MiddleLeft:   "╠",
		MiddleRight:  "╣",
		Middle:       "╬",
		MiddleTop:    "╦",
		MiddleBottom: "╩",
	}

	hiddenBorder = Border{
		Top:          " ",
		Bottom:       " ",
		Left:         " ",
		Right:        " ",
		TopLeft:      " ",
		TopRight:     " ",
		BottomLeft:   " ",
		BottomRight:  " ",
		MiddleLeft:   " ",
		MiddleRight:  " ",
		Middle:       " ",
		MiddleTop:    " ",
		MiddleBottom: " ",
	}
)

// NormalBorder returns a standard-type border with a normal weight and 90
// degree corners.
func NormalBorder() Border {
	return normalBorder
}

// RoundedBorder returns a border with rounded corners.
func RoundedBorder() Border {
	return roundedBorder
}

// BlockBorder returns a border that takes the whole block.
func BlockBorder() Border {
	return blockBorder
}

// OuterHalfBlockBorder returns a half-block border that sits outside the frame.
func OuterHalfBlockBorder() Border {
	return outerHalfBlockBorder
}

// InnerHalfBlockBorder returns a half-block border that sits inside the frame.
func InnerHalfBlockBorder() Border {
	return innerHalfBlockBorder
}

// ThickBorder returns a border that's thicker than the one returned by
// NormalBorder.
func ThickBorder() Border {
	return thickBorder
}

// DoubleBorder returns a border comprised of two thin strokes.
func DoubleBorder() Border {
	return doubleBorder
}

// HiddenBorder returns a border that renders as a series of single-cell
// spaces. It's useful for cases when you want to remove a standard border but
// maintain layout positioning. This said, you can still apply a background
// color to a hidden border.
func HiddenBorder() Border {
	return hiddenBorder
}

func (s Style) applyBorder(str string) string {
	var (
		topSet    = s.isSet(borderTopKey)
		rightSet  = s.isSet(borderRightKey)
		bottomSet = s.isSet(borderBottomKey)
		leftSet   = s.isSet(borderLeftKey)

		border    = s.getBorderStyle()
		hasTop    = s.getAsBool(borderTopKey, false)
		hasRight  = s.getAsBool(borderRightKey, false)
		hasBottom = s.getAsBool(borderBottomKey, false)
		hasLeft   = s.getAsBool(borderLeftKey, false)

		topFG    = s.getAsColor(borderTopForegroundKey)
		rightFG  = s.getAsColor(borderRightForegroundKey)
		bottomFG = s.getAsColor(borderBottomForegroundKey)
		leftFG   = s.getAsColor(borderLeftForegroundKey)

		topBG    = s.getAsColor(borderTopBackgroundKey)
		rightBG  = s.getAsColor(borderRightBackgroundKey)
		bottomBG = s.getAsColor(borderBottomBackgroundKey)
		leftBG   = s.getAsColor(borderLeftBackgroundKey)

		topFuncs    = s.borderTopFunc
		bottomFuncs = s.borderBottomFunc
	)

	// If a border is set and no sides have been specifically turned on or off
	// render borders on all sides.
	if border != noBorder && !(topSet || rightSet || bottomSet || leftSet) {
		hasTop = true
		hasRight = true
		hasBottom = true
		hasLeft = true
	}

	// If no border is set or all borders are been disabled, abort.
	if border == noBorder || (!hasTop && !hasRight && !hasBottom && !hasLeft) {
		return str
	}

	lines, width := getLines(str)

	if hasLeft {
		if border.Left == "" {
			border.Left = " "
		}
	}

	if hasRight && border.Right == "" {
		border.Right = " "
	}

	// If corners should be rendered but are set with the empty string, fill them
	// with a single space.
	if hasTop && hasLeft && border.TopLeft == "" {
		border.TopLeft = " "
	}
	if hasTop && hasRight && border.TopRight == "" {
		border.TopRight = " "
	}
	if hasBottom && hasLeft && border.BottomLeft == "" {
		border.BottomLeft = " "
	}
	if hasBottom && hasRight && border.BottomRight == "" {
		border.BottomRight = " "
	}

	// Figure out which corners we should actually be using based on which
	// sides are set to show.
	if hasTop {
		switch {
		case !hasLeft && !hasRight:
			border.TopLeft = ""
			border.TopRight = ""
		case !hasLeft:
			border.TopLeft = ""
		case !hasRight:
			border.TopRight = ""
		}
	}
	if hasBottom {
		switch {
		case !hasLeft && !hasRight:
			border.BottomLeft = ""
			border.BottomRight = ""
		case !hasLeft:
			border.BottomLeft = ""
		case !hasRight:
			border.BottomRight = ""
		}
	}

	// For now, limit corners to one rune.
	border.TopLeft = getFirstRuneAsString(border.TopLeft)
	border.TopRight = getFirstRuneAsString(border.TopRight)
	border.BottomRight = getFirstRuneAsString(border.BottomRight)
	border.BottomLeft = getFirstRuneAsString(border.BottomLeft)

	var out strings.Builder

	// Render top
	if hasTop {
		var top string
		if len(topFuncs) > 0 {
			top = renderAnnotatedHorizontalEdge(border.TopLeft, border.Top, border.TopRight, topFuncs, width)
		} else {
			top = renderHorizontalEdge(border.TopLeft, border.Top, border.TopRight, width)
		}
		top = s.styleBorder(top, topFG, topBG)
		out.WriteString(top)
		out.WriteRune('\n')
	}

	leftRunes := []rune(border.Left)
	leftIndex := 0

	rightRunes := []rune(border.Right)
	rightIndex := 0

	// Render sides
	for i, l := range lines {
		if hasLeft {
			r := string(leftRunes[leftIndex])
			leftIndex++
			if leftIndex >= len(leftRunes) {
				leftIndex = 0
			}
			out.WriteString(s.styleBorder(r, leftFG, leftBG))
		}
		out.WriteString(l)
		if hasRight {
			r := string(rightRunes[rightIndex])
			rightIndex++
			if rightIndex >= len(rightRunes) {
				rightIndex = 0
			}
			out.WriteString(s.styleBorder(r, rightFG, rightBG))
		}
		if i < len(lines)-1 {
			out.WriteRune('\n')
		}
	}

	// Render bottom
	if hasBottom {
		var bottom string
		if len(bottomFuncs) > 0 {
			bottom = renderAnnotatedHorizontalEdge(border.BottomLeft, border.Bottom, border.BottomRight, bottomFuncs, width)
		} else {
			bottom = renderHorizontalEdge(border.BottomLeft, border.Bottom, border.BottomRight, width)
		}
		bottom = s.styleBorder(bottom, bottomFG, bottomBG)
		out.WriteRune('\n')
		out.WriteString(bottom)
	}

	return out.String()
}

// Render the horizontal (top or bottom) portion of a border.
func renderAnnotatedHorizontalEdge(left, middle, right string, bFuncs []interface{}, width int) string {
	if middle == "" {
		middle = " "
	}

	ts := make([]string, 3)
	ws := make([]int, 3)
	for i, f := range bFuncs {
		if f == nil {
			continue
		}
		remainingWidth := width
		if remainingWidth < 1 {
			break
		}
		switch f := f.(type) {
		case string:
			ts[i] = ansi.Truncate(f, remainingWidth, "")
			ws[i] = ansi.StringWidth(ts[i])
			remainingWidth -= ws[i]
		case func(int, string) string:
			ts[i] = f(remainingWidth, middle)
			ts[i] = ansi.Truncate(ts[i], remainingWidth, "")
			ws[i] = ansi.StringWidth(ts[i])
			remainingWidth -= ws[i]
		case func() string:
			ts[i] = ansi.Truncate(f(), remainingWidth, "")
			ws[i] = ansi.StringWidth(ts[i])
			remainingWidth -= ws[i]
		}
	}

	runes := []rune(middle)
	j := 0

	out := strings.Builder{}
	out.WriteString(left)
	out.WriteString(ts[0])

	for i := ws[0]; i < width-ws[2]; {
		if ws[1] > 0 && i == (width-ws[1])/2 {
			out.WriteString(ts[1])
			i += ws[1]
		}
		out.WriteRune(runes[j])
		j++
		if j >= len(runes) {
			j = 0
		}
		i += ansi.StringWidth(string(runes[j]))
	}
	out.WriteString(ts[2])
	out.WriteString(right)

	return out.String()
}

// Render the horizontal (top or bottom) portion of a border.
func renderHorizontalEdge(left, middle, right string, width int) string {
	if middle == "" {
		middle = " "
	}

	runes := []rune(middle)
	j := 0

	out := strings.Builder{}
	out.WriteString(left)
	for i := 0; i < width; {
		out.WriteRune(runes[j])
		j++
		if j >= len(runes) {
			j = 0
		}
		i += ansi.StringWidth(string(runes[j]))
	}
	out.WriteString(right)

	return out.String()
}

// Apply foreground and background styling to a border.
func (s Style) styleBorder(border string, fg, bg TerminalColor) string {
	if fg == noColor && bg == noColor {
		return border
	}

	style := termenv.Style{}

	if fg != noColor {
		style = style.Foreground(fg.color(s.r))
	}
	if bg != noColor {
		style = style.Background(bg.color(s.r))
	}

	return style.Styled(border)
}

func maxRuneWidth(str string) int {
	var width int

	state := -1
	for len(str) > 0 {
		var w int
		_, str, w, state = uniseg.FirstGraphemeClusterInString(str, state)
		if w > width {
			width = w
		}
	}

	return width
}

func getFirstRuneAsString(str string) string {
	if str == "" {
		return str
	}
	r := []rune(str)
	return string(r[0])
}
