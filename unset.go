package lipgloss

// UnsetBold removes the bold style rule, if set.
func (s Style) UnsetBold() Style {
	delete(s.rules, boldKey)
	return s
}

// UnsetItalic removes the italic style rule, if set.
func (s Style) UnsetItalic() Style {
	delete(s.rules, italicKey)
	return s
}

// UnsetItalic removes the underline style rule, if set.
func (s Style) UnsetUnderline() Style {
	delete(s.rules, underlineKey)
	return s
}

// UnsetStrikethrough removes the strikethrough style rule, if set.
func (s Style) UnsetStrikethrough() Style {
	delete(s.rules, strikethroughKey)
	return s
}

// UnsetReverse removes the reverse style rule, if set.
func (s Style) UnsetReverse() Style {
	delete(s.rules, reverseKey)
	return s
}

// UnsetBlink removes the bold style rule, if set.
func (s Style) UnsetBlink() Style {
	delete(s.rules, blinkKey)
	return s
}

// UnsetFaint removes the faint style rule, if set.
func (s Style) UnsetFaint() Style {
	delete(s.rules, faintKey)
	return s
}

// UnsetForegroundremoves the foreground style rule, if set.
func (s Style) UnsetForeground() Style {
	delete(s.rules, foregroundKey)
	return s
}

// UnsetBackground removes the background style rule, if set.
func (s Style) UnsetBackground() Style {
	delete(s.rules, backgroundKey)
	return s
}

// UnsetWidth removes the width style rule, if set.
func (s Style) UnsetWidth() Style {
	delete(s.rules, widthKey)
	return s
}

// UnsetAlign removes the text alignment style rule, if set.
func (s Style) UnsetAlign() Style {
	delete(s.rules, alignKey)
	return s
}

// UnsetPadding removes all padding style rules.
func (s Style) UnsetPadding() Style {
	delete(s.rules, paddingLeftKey)
	delete(s.rules, paddingRightKey)
	delete(s.rules, paddingTopKey)
	delete(s.rules, paddingBottomKey)
	return s
}

// UnsetLeftPadding removes the left padding style rule, if set.
func (s Style) UnsetPaddingLeft() Style {
	delete(s.rules, paddingLeftKey)
	return s
}

// UnsetRightPadding removes the left padding style rule, if set.
func (s Style) UnsetPaddingRight() Style {
	delete(s.rules, paddingRightKey)
	return s
}

// UnsetTopPadding removes the top padding style rule, if set.
func (s Style) UnsetPaddingTop() Style {
	delete(s.rules, paddingTopKey)
	return s
}

// UnsetBottomPadding removes the bottom style rule, if set.
func (s Style) UnsetPaddingBottom() Style {
	delete(s.rules, paddingBottomKey)
	return s
}

// UnsetColorWhitespace removes the rule for coloring padding, if set.
func (s Style) UnsetColorWhitespace() Style {
	delete(s.rules, colorWhitespaceKey)
	return s
}

// UnsetMargins removes all margin style rules.
func (s Style) UnsetMargins() Style {
	delete(s.rules, marginLeftKey)
	delete(s.rules, marginRightKey)
	delete(s.rules, marginTopKey)
	delete(s.rules, marginBottomKey)
	return s
}

// UnsetLeftMargin removes the left margin style rule, if set.
func (s Style) UnsetMarginLeft() Style {
	delete(s.rules, marginLeftKey)
	return s
}

// UnsetRightMargin removes the right margin style rule, if set.
func (s Style) UnsetMarginRight() Style {
	delete(s.rules, marginRightKey)
	return s
}

// UnsetTopMargin removes the top margin style rule, if set.
func (s Style) UnsetMarginTop() Style {
	delete(s.rules, marginTopKey)
	return s
}

// UnsetBottomMargin removes the bottom margin style rule, if set.
func (s Style) UnsetMarginBottom() Style {
	delete(s.rules, marginBottomKey)
	return s
}

// UnsetInline removes the inline style rule, if set.
func (s Style) UnsetInline() Style {
	delete(s.rules, inlineKey)
	return s
}

// UnsetMaxWidth removes the max width style rule, if set.
func (s Style) UnsetMaxWidth() Style {
	delete(s.rules, maxWidthKey)
	return s
}

// UnsetDrawClearTrailingSpaces removes the rule for drawing clear trailing
// spaces, if set.
func (s Style) UnsetDrawClearTrailingSpaces() Style {
	delete(s.rules, drawClearTrailingSpacesKey)
	return s
}

// UnsetUnderlineWhitespace removes the rule for underlining whitespace, if
// set.
func (s Style) UnsetUnderlineWhitespace() Style {
	delete(s.rules, underlineSpacesKey)
	return s
}

// UnsetUnderlineWhitespace removes the rule for strikingn through whitespace,
// if set.
func (s Style) UnsetStrikethroughWhitespace() Style {
	delete(s.rules, strikethroughSpacesKey)
	return s
}

// UnsetUnderlineSpaces removes the value set by UnderlineSpaces.
func (s Style) UnsetUnderlineSpaces() Style {
	delete(s.rules, underlineSpacesKey)
	return s
}

// UnsetUnderlineSpaces removes the value set by UnsetStrikethroughSpaces.
func (s Style) UnsetStrikethroughSpaces() Style {
	delete(s.rules, strikethroughSpacesKey)
	return s
}

// UnsetString sets the underlying string value to the empty string.
func (s Style) UnsetString() Style {
	s.value = ""
	return s
}
