package lipgloss

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	"github.com/muesli/termenv"
)

var renderer = NewRenderer()

// Renderer is a lipgloss terminal renderer.
type Renderer struct {
	output            *termenv.Output
	hasDarkBackground bool
}

// RendererOption is a function that can be used to configure a Renderer.
type RendererOption func(r *Renderer)

// DefaultRenderer returns the default renderer.
func DefaultRenderer() *Renderer {
	return renderer
}

// NewRenderer creates a new Renderer.
func NewRenderer(options ...RendererOption) *Renderer {
	r := &Renderer{}
	for _, option := range options {
		option(r)
	}
	if r.output == nil {
		r.output = termenv.DefaultOutput()
	}
	if !r.hasDarkBackground {
		r.hasDarkBackground = r.output.HasDarkBackground()
	}
	return r
}

// WithOutput sets the termenv Output to use for rendering.
func WithOutput(output *termenv.Output) RendererOption {
	return func(r *Renderer) {
		r.output = output
	}
}

// WithDarkBackground forces the renderer to use a dark background.
func WithDarkBackground() RendererOption {
	return func(r *Renderer) {
		r.SetHasDarkBackground(true)
	}
}

// WithColorProfile sets the color profile on the renderer. This function is
// primarily intended for testing. For details, see the note on
// [Renderer.SetColorProfile].
func WithColorProfile(p termenv.Profile) RendererOption {
	return func(r *Renderer) {
		r.SetColorProfile(p)
	}
}

// ColorProfile returns the detected termenv color profile.
func (r *Renderer) ColorProfile() termenv.Profile {
	return r.output.Profile
}

// ColorProfile returns the detected termenv color profile.
func ColorProfile() termenv.Profile {
	return renderer.ColorProfile()
}

// SetColorProfile sets the color profile on the renderer. This function exists
// mostly for testing purposes so that you can assure you're testing against
// a specific profile.
//
// Outside of testing you likely won't want to use this function as the color
// profile will detect and cache the terminal's color capabilities and choose
// the best available profile.
//
// Available color profiles are:
//
// termenv.Ascii (no color, 1-bit)
// termenv.ANSI (16 colors, 4-bit)
// termenv.ANSI256 (256 colors, 8-bit)
// termenv.TrueColor (16,777,216 colors, 24-bit)
//
// This function is thread-safe.
func (r *Renderer) SetColorProfile(p termenv.Profile) {
	r.output.Profile = p
}

// SetColorProfile sets the color profile on the default renderer. This
// function exists mostly for testing purposes so that you can assure you're
// testing against a specific profile.
//
// Outside of testing you likely won't want to use this function as the color
// profile will detect and cache the terminal's color capabilities and choose
// the best available profile.
//
// Available color profiles are:
//
// termenv.Ascii (no color, 1-bit)
// termenv.ANSI (16 colors, 4-bit)
// termenv.ANSI256 (256 colors, 8-bit)
// termenv.TrueColor (16,777,216 colors, 24-bit)
//
// This function is thread-safe.
func SetColorProfile(p termenv.Profile) {
	renderer.SetColorProfile(p)
}

// HasDarkBackground returns whether or not the terminal has a dark background.
func (r *Renderer) HasDarkBackground() bool {
	return r.hasDarkBackground
}

// HasDarkBackground returns whether or not the terminal has a dark background.
func HasDarkBackground() bool {
	return renderer.HasDarkBackground()
}

// SetHasDarkBackground sets the background color detection value on the
// renderer. This function exists mostly for testing purposes so that you can
// assure you're testing against a specific background color setting.
//
// Outside of testing you likely won't want to use this function as the
// backgrounds value will be automatically detected and cached against the
// terminal's current background color setting.
//
// This function is thread-safe.
func (r *Renderer) SetHasDarkBackground(b bool) {
	r.hasDarkBackground = b
}

// SetHasDarkBackground sets the background color detection value for the
// default renderer. This function exists mostly for testing purposes so that
// you can assure you're testing against a specific background color setting.
//
// Outside of testing you likely won't want to use this function as the
// backgrounds value will be automatically detected and cached against the
// terminal's current background color setting.
//
// This function is thread-safe.
func SetHasDarkBackground(b bool) {
	renderer.SetHasDarkBackground(b)
}

// Render formats a string according to the given style.
func (r *Renderer) Render(s Style, str string) string {
	var (
		te           = r.ColorProfile().String()
		teSpace      = r.ColorProfile().String()
		teWhitespace = r.ColorProfile().String()

		bold          = s.getAsBool(boldKey, false)
		italic        = s.getAsBool(italicKey, false)
		underline     = s.getAsBool(underlineKey, false)
		strikethrough = s.getAsBool(strikethroughKey, false)
		reverse       = s.getAsBool(reverseKey, false)
		blink         = s.getAsBool(blinkKey, false)
		faint         = s.getAsBool(faintKey, false)

		fg = s.getAsColor(foregroundKey)
		bg = s.getAsColor(backgroundKey)

		width           = s.getAsInt(widthKey)
		height          = s.getAsInt(heightKey)
		horizontalAlign = s.getAsPosition(alignHorizontalKey)
		verticalAlign   = s.getAsPosition(alignVerticalKey)

		topPadding    = s.getAsInt(paddingTopKey)
		rightPadding  = s.getAsInt(paddingRightKey)
		bottomPadding = s.getAsInt(paddingBottomKey)
		leftPadding   = s.getAsInt(paddingLeftKey)

		colorWhitespace = s.getAsBool(colorWhitespaceKey, true)
		inline          = s.getAsBool(inlineKey, false)
		maxWidth        = s.getAsInt(maxWidthKey)
		maxHeight       = s.getAsInt(maxHeightKey)

		underlineSpaces     = underline && s.getAsBool(underlineSpacesKey, true)
		strikethroughSpaces = strikethrough && s.getAsBool(strikethroughSpacesKey, true)

		// Do we need to style whitespace (padding and space outside
		// paragraphs) separately?
		styleWhitespace = reverse

		// Do we need to style spaces separately?
		useSpaceStyler = underlineSpaces || strikethroughSpaces
	)

	if len(s.rules) == 0 {
		return str
	}

	// Enable support for ANSI on the legacy Windows cmd.exe console. This is a
	// no-op on non-Windows systems and on Windows runs only once.
	enableLegacyWindowsANSI()

	if bold {
		te = te.Bold()
	}
	if italic {
		te = te.Italic()
	}
	if underline {
		te = te.Underline()
	}
	if reverse {
		if reverse {
			teWhitespace = teWhitespace.Reverse()
		}
		te = te.Reverse()
	}
	if blink {
		te = te.Blink()
	}
	if faint {
		te = te.Faint()
	}

	if fg != noColor {
		fgc := r.color(fg)
		te = te.Foreground(fgc)
		if styleWhitespace {
			teWhitespace = teWhitespace.Foreground(fgc)
		}
		if useSpaceStyler {
			teSpace = teSpace.Foreground(fgc)
		}
	}

	if bg != noColor {
		bgc := r.color(bg)
		te = te.Background(bgc)
		if colorWhitespace {
			teWhitespace = teWhitespace.Background(bgc)
		}
		if useSpaceStyler {
			teSpace = teSpace.Background(bgc)
		}
	}

	if underline {
		te = te.Underline()
	}
	if strikethrough {
		te = te.CrossOut()
	}

	if underlineSpaces {
		teSpace = teSpace.Underline()
	}
	if strikethroughSpaces {
		teSpace = teSpace.CrossOut()
	}

	// Strip newlines in single line mode
	if inline {
		str = strings.ReplaceAll(str, "\n", "")
	}

	// Word wrap
	if !inline && width > 0 {
		wrapAt := width - leftPadding - rightPadding
		str = wordwrap.String(str, wrapAt)
		str = wrap.String(str, wrapAt) // force-wrap long strings
	}

	// Render core text
	{
		var b strings.Builder

		l := strings.Split(str, "\n")
		for i := range l {
			if useSpaceStyler {
				// Look for spaces and apply a different styler
				for _, r := range l[i] {
					if unicode.IsSpace(r) {
						b.WriteString(teSpace.Styled(string(r)))
						continue
					}
					b.WriteString(te.Styled(string(r)))
				}
			} else {
				b.WriteString(te.Styled(l[i]))
			}
			if i != len(l)-1 {
				b.WriteRune('\n')
			}
		}

		str = b.String()
	}

	// Padding
	if !inline {
		if leftPadding > 0 {
			var st *termenv.Style
			if colorWhitespace || styleWhitespace {
				st = &teWhitespace
			}
			str = padLeft(str, leftPadding, st)
		}

		if rightPadding > 0 {
			var st *termenv.Style
			if colorWhitespace || styleWhitespace {
				st = &teWhitespace
			}
			str = padRight(str, rightPadding, st)
		}

		if topPadding > 0 {
			str = strings.Repeat("\n", topPadding) + str
		}

		if bottomPadding > 0 {
			str += strings.Repeat("\n", bottomPadding)
		}
	}

	// Height
	if height > 0 {
		str = alignTextVertical(str, verticalAlign, height, nil)
	}

	// Set alignment. This will also pad short lines with spaces so that all
	// lines are the same length, so we run it under a few different conditions
	// beyond alignment.
	{
		numLines := strings.Count(str, "\n")

		if !(numLines == 0 && width == 0) {
			var st *termenv.Style
			if colorWhitespace || styleWhitespace {
				st = &teWhitespace
			}
			str = alignTextHorizontal(str, horizontalAlign, width, st)
		}
	}

	if !inline {
		str = s.applyBorder(r, str)
		str = s.applyMargins(r, str, inline)
	}

	// Truncate according to MaxWidth
	if maxWidth > 0 {
		lines := strings.Split(str, "\n")

		for i := range lines {
			lines[i] = truncate.String(lines[i], uint(maxWidth))
		}

		str = strings.Join(lines, "\n")
	}

	// Truncate according to MaxHeight
	if maxHeight > 0 {
		lines := strings.Split(str, "\n")
		str = strings.Join(lines[:min(maxHeight, len(lines))], "\n")
	}

	return str
}

// Render formats a string according to the given style using the default
// renderer. This is syntactic sugar for rendering with a DefaultRenderer.
func Render(s Style, str string) string {
	return renderer.Render(s, str)
}

func (r *Renderer) colorValue(c TerminalColor) string {
	switch c := c.(type) {
	case ANSIColor:
		return fmt.Sprint(c)
	case Color:
		return string(c)
	case AdaptiveColor:
		if r.HasDarkBackground() {
			return c.Dark
		}
		return c.Light
	case CompleteColor:
		switch r.ColorProfile() {
		case termenv.TrueColor:
			return c.TrueColor
		case termenv.ANSI256:
			return c.ANSI256
		case termenv.ANSI:
			return c.ANSI
		default:
			return ""
		}
	case CompleteAdaptiveColor:
		col := c.Light
		if r.HasDarkBackground() {
			col = c.Dark
		}
		switch r.ColorProfile() {
		case termenv.TrueColor:
			return col.TrueColor
		case termenv.ANSI256:
			return col.ANSI256
		case termenv.ANSI:
			return col.ANSI
		default:
			return ""
		}

	default:
		return ""
	}
}

// color returns a termenv.color for the given TerminalColor.
func (r *Renderer) color(c TerminalColor) termenv.Color {
	return r.ColorProfile().Color(r.colorValue(c))
}
