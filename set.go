package lipgloss

// Bold sets a bold formatting rule.
func (s Style) Bold(v bool) Style {
	s.bold = &v
	return s
}

// Italic sets an italic formatting rule. In some terminal emulators this will
// render with "reverse" coloring if not italic font variant is available.
func (s Style) Italic(v bool) Style {
	s.italic = &v
	return s
}

// Underine sets an underline rule. By default, underlines will not be drawn on
// whitespace like margins and padding. To change this behavior set
// renderUnderlinesOnSpaces.
func (s Style) Underline(v bool) Style {
	s.underline = &v
	return s
}

// Strikethrough sets a strikethrough rule. By default, strikes will not be
// drawn on whitespace like margins and padding. To change this behavior set
// renderStrikethroughOnSpaces.
func (s Style) Strikethrough(v bool) Style {
	s.strikethrough = &v
	return s
}

// Reverse sets a rule for inverting foreground and background colors.
func (s Style) Reverse(v bool) Style {
	s.reverse = &v
	return s
}

// Blink sets a rule for blinking forground text.
func (s Style) Blink(v bool) Style {
	s.blink = &v
	return s
}

// Faint sets a rule for rendering the foreground color in a dimmer shade.
func (s Style) Faint(v bool) Style {
	s.faint = &v
	return s
}

// Foreground sets a foreground color.
//
//     // Sets the foreground to blue
//     s := lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
//
//     // Removes the foreground color
//     s.Foreground(lipgloss.NoColor)
//
func (s Style) Foreground(c ColorType) Style {
	s.foreground = &c
	return s
}

// Background sets a background color.
func (s Style) Background(c ColorType) Style {
	s.background = &c
	return s
}

// Width sets the width of the block before applying margins and padding. This
// effects when.
func (s Style) Width(i int) Style {
	s.width = &i
	return s
}

// Align sets a text alignment rule.
func (s Style) Align(a Align) Style {
	s.align = &a
	return s
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
func (s Style) Padding(i ...int) Style {
	top, right, bottom, left, ok := whichSides(i...)
	if !ok {
		return s
	}

	s.topPadding = &top
	s.rightPadding = &right
	s.bottomPadding = &bottom
	s.leftPadding = &left
	return s
}

func (s Style) LeftPadding(i int) Style {
	s.leftPadding = &i
	return s
}

func (s Style) RightPadding(i int) Style {
	s.rightPadding = &i
	return s
}

func (s Style) TopPadding(i int) Style {
	s.topPadding = &i
	return s
}

func (s Style) BottomPadding(i int) Style {
	s.bottomPadding = &i
	return s
}

func (s Style) StylePadding(v bool) Style {
	s.stylePadding = &v
	return s
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
// With more than four arguments no padding will be added.
func (s Style) Margin(i ...int) Style {
	top, right, bottom, left, ok := whichSides(i...)
	if !ok {
		return s
	}

	s.topMargin = &top
	s.rightMargin = &right
	s.bottomMargin = &bottom
	s.leftMargin = &left
	return s
}

func (s Style) LeftMargin(i int) Style {
	s.leftMargin = &i
	return s
}

func (s Style) RightMargin(i int) Style {
	s.rightMargin = &i
	return s
}

func (s Style) TopMargin(i int) Style {
	s.topMargin = &i
	return s
}

func (s Style) BottomMargin(i int) Style {
	s.bottomMargin = &i
	return s
}

// Inline makes rendering output one line. This is useful for enforcing that
// rendering occurs on a single line at render time, particularly with styles
// and strings you may not have control of. Works well with MaxWidth().
func (s Style) Inline(v bool) Style {
	s.inline = &v
	return s
}

// WithMaxWidth applies a max width to a given style. This is useful in
// enforcing a certain width at render time, particularly with aribtrary
// strings and styles.
//
// Example:
//
//     var userInput string = "..."
//     var userStyle = text.Style{ /* ... */ }
//     fmt.Println(userStyle.MaxWidth(16).Render(userInput))
//
func (s Style) MaxWidth(n int) Style {
	s.maxWidth = &n
	return s
}

// Whether or not to draw trailing spaces with no background color. By default
// we leave them in.
func (s Style) DrawClearTrailingSpaces(v bool) Style {
	s.drawClearTrailingSpaces = &v
	return s
}

func (s Style) UnderlineWhitespace(v bool) Style {
	s.underlineWhitespace = &v
	return s
}

func (s Style) StrikethroughWhitespace(v bool) Style {
	s.strikethroughWhitespace = &v
	return s
}

// whichSides is a helper method for setting values on sides of a block based
// on the number of arguments. It follows the CSS shorthand rules for blocks
// like margin, padding. and borders. Here are how the rules work:
//
// 0 args:  do nothing
// 1 arg:   all sides
// 2 args:  top -> bottom
// 3 args:  top -> horizontal -> bottom
// 4 args:  top -> right -> bottom -> left
// 5+ args: do nothing
func whichSides(i ...int) (top, right, bottom, left int, ok bool) {
	switch len(i) {
	case 1:
		top = i[0]
		bottom = i[0]
		left = i[0]
		right = i[0]
		ok = true
	case 2:
		top = i[0]
		bottom = i[0]
		left = i[1]
		right = i[1]
		ok = true
	case 3:
		top = i[0]
		left = i[1]
		right = i[1]
		bottom = i[2]
		ok = true
	case 4:
		top = i[0]
		right = i[1]
		bottom = i[2]
		left = i[3]
		ok = true
	}
	return
}
