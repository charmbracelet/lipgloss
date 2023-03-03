package lipgloss

func (s Style) unset(key propKey) Style {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.rules, key)
	return s
}

// UnsetBold removes the bold style rule, if set.
func (s Style) UnsetBold() Style {
	return s.unset(boldKey)
}

// UnsetItalic removes the italic style rule, if set.
func (s Style) UnsetItalic() Style {
	return s.unset(italicKey)
}

// UnsetUnderline removes the underline style rule, if set.
func (s Style) UnsetUnderline() Style {
	return s.unset(underlineKey)
}

// UnsetStrikethrough removes the strikethrough style rule, if set.
func (s Style) UnsetStrikethrough() Style {
	return s.unset(strikethroughKey)
}

// UnsetReverse removes the reverse style rule, if set.
func (s Style) UnsetReverse() Style {
	return s.unset(reverseKey)
}

// UnsetBlink removes the bold style rule, if set.
func (s Style) UnsetBlink() Style {
	return s.unset(blinkKey)
}

// UnsetFaint removes the faint style rule, if set.
func (s Style) UnsetFaint() Style {
	return s.unset(faintKey)
}

// UnsetForeground removes the foreground style rule, if set.
func (s Style) UnsetForeground() Style {
	return s.unset(foregroundKey)
}

// UnsetBackground removes the background style rule, if set.
func (s Style) UnsetBackground() Style {
	return s.unset(backgroundKey)
}

// UnsetWidth removes the width style rule, if set.
func (s Style) UnsetWidth() Style {
	return s.unset(widthKey)
}

// UnsetHeight removes the height style rule, if set.
func (s Style) UnsetHeight() Style {
	return s.unset(heightKey)
}

// UnsetAlign removes the horizontal and vertical text alignment style rule, if set.
func (s Style) UnsetAlign() Style {
	s = s.unset(alignHorizontalKey)
	return s.unset(alignVerticalKey)
}

// UnsetAlignHorizontal removes the horizontal text alignment style rule, if set.
func (s Style) UnsetAlignHorizontal() Style {
	return s.unset(alignHorizontalKey)
}

// UnsetAlignVertical removes the vertical text alignment style rule, if set.
func (s Style) UnsetAlignVertical() Style {
	return s.unset(alignVerticalKey)
}

// UnsetPadding removes all padding style rules.
func (s Style) UnsetPadding() Style {
	s = s.unset(paddingLeftKey)
	s = s.unset(paddingRightKey)
	s = s.unset(paddingTopKey)
	return s.unset(paddingBottomKey)
}

// UnsetPaddingLeft removes the left padding style rule, if set.
func (s Style) UnsetPaddingLeft() Style {
	return s.unset(paddingLeftKey)
}

// UnsetPaddingRight removes the right padding style rule, if set.
func (s Style) UnsetPaddingRight() Style {
	return s.unset(paddingRightKey)
}

// UnsetPaddingTop removes the top padding style rule, if set.
func (s Style) UnsetPaddingTop() Style {
	return s.unset(paddingTopKey)
}

// UnsetPaddingBottom removes the bottom style rule, if set.
func (s Style) UnsetPaddingBottom() Style {
	return s.unset(paddingBottomKey)
}

// UnsetColorWhitespace removes the rule for coloring padding, if set.
func (s Style) UnsetColorWhitespace() Style {
	return s.unset(colorWhitespaceKey)
}

// UnsetMargins removes all margin style rules.
func (s Style) UnsetMargins() Style {
	s = s.unset(marginLeftKey)
	s = s.unset(marginRightKey)
	s = s.unset(marginTopKey)
	return s.unset(marginBottomKey)
}

// UnsetMarginLeft removes the left margin style rule, if set.
func (s Style) UnsetMarginLeft() Style {
	return s.unset(marginLeftKey)
}

// UnsetMarginRight removes the right margin style rule, if set.
func (s Style) UnsetMarginRight() Style {
	return s.unset(marginRightKey)
}

// UnsetMarginTop removes the top margin style rule, if set.
func (s Style) UnsetMarginTop() Style {
	return s.unset(marginTopKey)
}

// UnsetMarginBottom removes the bottom margin style rule, if set.
func (s Style) UnsetMarginBottom() Style {
	return s.unset(marginBottomKey)
}

// UnsetMarginBackground removes the margin's background color. Note that the
// margin's background color can be set from the background color of another
// style during inheritance.
func (s Style) UnsetMarginBackground() Style {
	return s.unset(marginBackgroundKey)
}

// UnsetBorderStyle removes the border style rule, if set.
func (s Style) UnsetBorderStyle() Style {
	return s.unset(borderStyleKey)
}

// UnsetBorderTop removes the border top style rule, if set.
func (s Style) UnsetBorderTop() Style {
	return s.unset(borderTopKey)
}

// UnsetBorderRight removes the border right style rule, if set.
func (s Style) UnsetBorderRight() Style {
	return s.unset(borderRightKey)
}

// UnsetBorderBottom removes the border bottom style rule, if set.
func (s Style) UnsetBorderBottom() Style {
	return s.unset(borderBottomKey)
}

// UnsetBorderLeft removes the border left style rule, if set.
func (s Style) UnsetBorderLeft() Style {
	return s.unset(borderLeftKey)
}

// UnsetBorderForeground removes all border foreground color styles, if set.
func (s Style) UnsetBorderForeground() Style {
	s = s.unset(borderTopForegroundKey)
	s = s.unset(borderRightForegroundKey)
	s = s.unset(borderBottomForegroundKey)
	return s.unset(borderLeftForegroundKey)
}

// UnsetBorderTopForeground removes the top border foreground color rule,
// if set.
func (s Style) UnsetBorderTopForeground() Style {
	return s.unset(borderTopForegroundKey)
}

// UnsetBorderRightForeground removes the right border foreground color rule,
// if set.
func (s Style) UnsetBorderRightForeground() Style {
	return s.unset(borderRightForegroundKey)
}

// UnsetBorderBottomForeground removes the bottom border foreground color
// rule, if set.
func (s Style) UnsetBorderBottomForeground() Style {
	return s.unset(borderBottomForegroundKey)
}

// UnsetBorderLeftForeground removes the left border foreground color rule,
// if set.
func (s Style) UnsetBorderLeftForeground() Style {
	return s.unset(borderLeftForegroundKey)
}

// UnsetBorderBackground removes all border background color styles, if
// set.
func (s Style) UnsetBorderBackground() Style {
	s = s.unset(borderTopBackgroundKey)
	s = s.unset(borderRightBackgroundKey)
	s = s.unset(borderBottomBackgroundKey)
	return s.unset(borderLeftBackgroundKey)
}

// UnsetBorderTopBackgroundColor removes the top border background color rule,
// if set.
func (s Style) UnsetBorderTopBackgroundColor() Style {
	return s.unset(borderTopBackgroundKey)
}

// UnsetBorderRightBackground removes the right border background color
// rule, if set.
func (s Style) UnsetBorderRightBackground() Style {
	return s.unset(borderRightBackgroundKey)
}

// UnsetBorderBottomBackground removes the bottom border background color
// rule, if set.
func (s Style) UnsetBorderBottomBackground() Style {
	return s.unset(borderBottomBackgroundKey)
}

// UnsetBorderLeftBackground removes the left border color rule, if set.
func (s Style) UnsetBorderLeftBackground() Style {
	return s.unset(borderLeftBackgroundKey)
}

// UnsetInline removes the inline style rule, if set.
func (s Style) UnsetInline() Style {
	return s.unset(inlineKey)
}

// UnsetMaxWidth removes the max width style rule, if set.
func (s Style) UnsetMaxWidth() Style {
	return s.unset(maxWidthKey)
}

// UnsetMaxHeight removes the max height style rule, if set.
func (s Style) UnsetMaxHeight() Style {
	return s.unset(maxHeightKey)
}

// UnsetUnderlineSpaces removes the value set by UnderlineSpaces.
func (s Style) UnsetUnderlineSpaces() Style {
	return s.unset(underlineSpacesKey)
}

// UnsetStrikethroughSpaces removes the value set by StrikethroughSpaces.
func (s Style) UnsetStrikethroughSpaces() Style {
	return s.unset(strikethroughSpacesKey)
}

// UnsetString sets the underlying string value to the empty string.
func (s Style) UnsetString() Style {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.value = ""
	return s
}
