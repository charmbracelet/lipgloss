package lipgloss

// Set a value on the underlying rules map.
func (s *Style) set(key propKey, value interface{}) {
	if s.rules == nil {
		s.rules = make(rules)
	}
	s.rules[key] = value
}

// Bold sets a bold formatting rule.
func (s Style) Bold(v bool) Style {
	s.set(boldKey, v)
	return s
}

// Italic sets an italic formatting rule. In some terminal emulators this will
// render with "reverse" coloring if not italic font variant is available.
func (s Style) Italic(v bool) Style {
	s.set(italicKey, v)
	return s
}

// Underine sets an underline rule. By default, underlines will not be drawn on
// whitespace like margins and padding. To change this behavior set
// renderUnderlinesOnSpaces.
func (s Style) Underline(v bool) Style {
	s.set(underlineKey, v)
	return s
}

// Strikethrough sets a strikethrough rule. By default, strikes will not be
// drawn on whitespace like margins and padding. To change this behavior set
// renderStrikethroughOnSpaces.
func (s Style) Strikethrough(v bool) Style {
	s.set(strikethroughKey, v)
	return s
}

// Reverse sets a rule for inverting foreground and background colors.
func (s Style) Reverse(v bool) Style {
	s.set(reverseKey, v)
	return s
}

// Blink sets a rule for blinking forground text.
func (s Style) Blink(v bool) Style {
	s.set(blinkKey, v)
	return s
}

// Faint sets a rule for rendering the foreground color in a dimmer shade.
func (s Style) Faint(v bool) Style {
	s.set(faintKey, v)
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
	s.set(foregroundKey, c)
	return s
}

// Background sets a background color.
func (s Style) Background(c ColorType) Style {
	s.set(backgroundKey, c)
	return s
}

// Width sets the width of the block before applying margins and padding.
func (s Style) Width(i int) Style {
	s.set(widthKey, i)
	return s
}

// Align sets a text alignment rule.
func (s Style) Align(a Align) Style {
	s.set(alignKey, a)
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

	s.set(paddingTopKey, top)
	s.set(paddingRightKey, right)
	s.set(paddingBottomKey, bottom)
	s.set(paddingLeftKey, left)
	return s
}

// LeftPadding adds padding on the left.
func (s Style) PaddingLeft(i int) Style {
	s.set(paddingLeftKey, i)
	return s
}

// Right Padding adds padding on the right.
func (s Style) RightPadding(i int) Style {
	s.set(paddingRightKey, i)
	return s
}

// TopPadding addds padding to the top of the block.
func (s Style) TopPadding(i int) Style {
	s.set(paddingTopKey, i)
	return s
}

// BottomPadding adds padding to the bottom of the block.
func (s Style) BottomPadding(i int) Style {
	s.set(paddingBottomKey, i)
	return s
}

// ColorWhitespace determines whether or not the background color should be
// applied to the padding. This is true by default as it's more than likely the
// desired and expected behavior, but it can be disabled for certain graphic
// effects.
func (s Style) ColorWhitespace(v bool) Style {
	s.set(colorWhitespaceKey, v)
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

	s.set(marginTopKey, top)
	s.set(marginRightKey, right)
	s.set(marginBottomKey, bottom)
	s.set(marginLeftKey, left)
	return s
}

// LeftMargin sets the value of the left margin.
func (s Style) LeftMargin(i int) Style {
	s.set(marginLeftKey, i)
	return s
}

// RightMargin sets the value of the right margin.
func (s Style) RightMargin(i int) Style {
	s.set(marginRightKey, i)
	return s
}

// TopMargin sets the value of the top margin.
func (s Style) TopMargin(i int) Style {
	s.set(marginTopKey, i)
	return s
}

// BottomMargin sets the value of the bottom margin.
func (s Style) BottomMargin(i int) Style {
	s.set(marginBottomKey, i)
	return s
}

// Inline makes rendering output one line. This is useful for enforcing that
// rendering occurs on a single line at render time, particularly with styles
// and strings you may not have control of. Works well with MaxWidth().
//
// Because this in intended to be used at the time of render, this method will
// not mutate the style and instead return a copy.
//
// Example:
//
//     var userInput string = "..."
//     var userStyle = text.Style{ /* ... */ }
//     fmt.Println(userStyle.Inline().Render(userInput))
//
func (s Style) Inline(v bool) Style {
	o := s.Copy()
	o.set(inlineKey, v)
	return o
}

// MaxWidth applies a max width to a given style. This is useful in enforcing
// a certain width at render time, particularly with arbitrary strings and
// styles.
//
// Because this in intended to be used at the time of render, this method will
// not mutate the style and instead return a copy.
//
// Example:
//
//     var userInput string = "..."
//     var userStyle = text.Style{ /* ... */ }
//     fmt.Println(userStyle.MaxWidth(16).Render(userInput))
//
func (s Style) MaxWidth(n int) Style {
	o := s.Copy()
	o.set(maxWidthKey, n)
	return o
}

// Whether or not to draw trailing spaces with no background color. By default
// we leave them in.
func (s Style) DrawClearTrailingSpaces(v bool) Style {
	s.set(drawClearTrailingSpacesKey, v)
	return s
}

// UnderlineWhitespace determines whether to draw underlines in areas outside
// text blocks. By default this is false, but it can be enabled for certain
// graphic effects. Also note that whitespace can be underlined without
// underlining text.
func (s Style) UnderlineWhitespace(v bool) Style {
	s.set(underlineWhitespaceKey, v)
	return s
}

// StrikethroughWhitespace determines whether to draw strikethroughs in areas
// outside text blocks. By default this is false, but it can be enabled for
// certain graphic effects. Also note that whitespace can be struck without
// striking through text.
func (s Style) StrikethroughWhitespace(v bool) Style {
	s.set(strikethroughWhitespaceKey, v)
	return s
}

// UnderlineSpaces determines whether to underline spaces between words. By
// default this is true. Spaces can also be underlined without underlining the
// text itself.
func (s Style) UnderlineSpaces(v bool) Style {
	s.set(underlineSpacesKey, v)
	return s
}

// StrikethroughSpaces determines whether to apply strikethroughs to spaces
// between words. By default this is true. Spaces can also be struck without
// underlining the text itself.
func (s Style) StrikethroughSpaces(v bool) Style {
	s.set(strikethroughSpacesKey, v)
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
	return top, right, bottom, left, ok
}
