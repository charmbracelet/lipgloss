package lipgloss

import (
	"bytes"
	"math"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
)

// Position represents a position along a horizontal or vertical axis. It's in
// situations where an axis is involved, like alignment, joining, placement and
// so on.
//
// A value of 0 represents the start (the left or top) and 1 represents the end
// (the right or bottom). 0.5 represents the center.
//
// There are constants Top, Bottom, Center, Left and Right in this package that
// can be used to aid readability.
type Position float64

func (p Position) value() float64 {
	return math.Min(1, math.Max(0, float64(p)))
}

// Position aliases.
const (
	Top    Position = 0.0
	Bottom Position = 1.0
	Center Position = 0.5
	Left   Position = 0.0
	Right  Position = 1.0
)

// Place places a string or text block vertically in an unstyled box of a given
// width or height.
func Place(width, height int, hPos, vPos Position, str string, opts ...WhitespaceOption) string {
	return renderer.Place(width, height, hPos, vPos, str, opts...)
}

// Place places a string or text block vertically in an unstyled box of a given
// width or height.
func (r *Renderer) Place(width, height int, hPos, vPos Position, str string, opts ...WhitespaceOption) string {
	return r.PlaceVertical(height, vPos, r.PlaceHorizontal(width, hPos, str, opts...), opts...)
}

// PlaceHorizontal places a string or text block horizontally in an unstyled
// block of a given width. If the given width is shorter than the max width of
// the string (measured by its longest line) this will be a noop.
func PlaceHorizontal(width int, pos Position, str string, opts ...WhitespaceOption) string {
	return renderer.PlaceHorizontal(width, pos, str, opts...)
}

// PlaceHorizontal places a string or text block horizontally in an unstyled
// block of a given width. If the given width is shorter than the max width of
// the string (measured by its longest line) this will be a noöp.
func (r *Renderer) PlaceHorizontal(width int, pos Position, str string, opts ...WhitespaceOption) string {
	lines, contentWidth := getLines(str)
	gap := width - contentWidth

	if gap <= 0 {
		return str
	}

	ws := newWhitespace(r, opts...)

	var b strings.Builder
	for i, l := range lines {
		// Is this line shorter than the longest line?
		short := max(0, contentWidth-ansi.PrintableRuneWidth(l))

		switch pos { //nolint:exhaustive
		case Left:
			b.WriteString(l)
			b.WriteString(ws.render(gap + short))

		case Right:
			b.WriteString(ws.render(gap + short))
			b.WriteString(l)

		default: // somewhere in the middle
			totalGap := gap + short

			split := int(math.Round(float64(totalGap) * pos.value()))
			left := totalGap - split
			right := totalGap - left

			b.WriteString(ws.render(left))
			b.WriteString(l)
			b.WriteString(ws.render(right))
		}

		if i < len(lines)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

// PlaceVertical places a string or text block vertically in an unstyled block
// of a given height. If the given height is shorter than the height of the
// string (measured by its newlines) then this will be a noop.
func PlaceVertical(height int, pos Position, str string, opts ...WhitespaceOption) string {
	return renderer.PlaceVertical(height, pos, str, opts...)
}

// PlaceVertical places a string or text block vertically in an unstyled block
// of a given height. If the given height is shorter than the height of the
// string (measured by its newlines) then this will be a noöp.
func (r *Renderer) PlaceVertical(height int, pos Position, str string, opts ...WhitespaceOption) string {
	contentHeight := strings.Count(str, "\n") + 1
	gap := height - contentHeight

	if gap <= 0 {
		return str
	}

	ws := newWhitespace(r, opts...)

	_, width := getLines(str)
	emptyLine := ws.render(width)
	b := strings.Builder{}

	switch pos { //nolint:exhaustive
	case Top:
		b.WriteString(str)
		b.WriteRune('\n')
		for i := 0; i < gap; i++ {
			b.WriteString(emptyLine)
			if i < gap-1 {
				b.WriteRune('\n')
			}
		}

	case Bottom:
		b.WriteString(strings.Repeat(emptyLine+"\n", gap))
		b.WriteString(str)

	default: // Somewhere in the middle
		split := int(math.Round(float64(gap) * pos.value()))
		top := gap - split
		bottom := gap - top

		b.WriteString(strings.Repeat(emptyLine+"\n", top))
		b.WriteString(str)

		for i := 0; i < bottom; i++ {
			b.WriteRune('\n')
			b.WriteString(emptyLine)
		}
	}

	return b.String()
}

// PlaceOverlay places fg on top of bg.
func PlaceOverlay(x, y int, fg, bg string, opts ...WhitespaceOption) string {
	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		// FIXME: return fg or bg?
		return fg
	}
	// TODO: allow placement outside of the bg box?
	x = clamp(x, 0, bgWidth-fgWidth)
	y = clamp(y, 0, bgHeight-fgHeight)

	ws := &whitespace{}
	for _, opt := range opts {
		opt(ws)
	}

	var b strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := truncate.String(bgLine, uint(x))
			pos = ansi.PrintableRuneWidth(left)
			b.WriteString(left)
			if pos < x {
				b.WriteString(ws.render(x - pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		b.WriteString(fgLine)
		pos += ansi.PrintableRuneWidth(fgLine)

		right := cutLeft(bgLine, pos)
		bgWidth := ansi.PrintableRuneWidth(bgLine)
		rightWidth := ansi.PrintableRuneWidth(right)
		if rightWidth <= bgWidth-pos {
			b.WriteString(ws.render(bgWidth - rightWidth - pos))
		}

		b.WriteString(right)
	}

	return b.String()
}

// cutLeft cuts printable characters from the left.
// This function is heavily based on muesli's ansi and truncate packages.
func cutLeft(s string, cutWidth int) string {
	var (
		pos    int
		isAnsi bool
		ab     bytes.Buffer
		b      bytes.Buffer
	)
	for _, c := range s {
		var w int
		if c == ansi.Marker || isAnsi {
			isAnsi = true
			ab.WriteRune(c)
			if ansi.IsTerminator(c) {
				isAnsi = false
				if bytes.HasSuffix(ab.Bytes(), []byte("[0m")) {
					ab.Reset()
				}
			}
		} else {
			w = runewidth.RuneWidth(c)
		}

		if pos >= cutWidth {
			if b.Len() == 0 {
				if ab.Len() > 0 {
					b.Write(ab.Bytes())
				}
				if pos-cutWidth > 1 {
					b.WriteByte(' ')
					continue
				}
			}
			b.WriteRune(c)
		}
		pos += w
	}
	return b.String()
}

func clamp(v, lower, upper int) int {
	return min(max(v, lower), upper)
}
