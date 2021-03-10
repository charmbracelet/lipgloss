package lipgloss

// UnsetBold removes the bold style rule, if set.
func (s Style) UnsetBold() Style {
	delete(s, boldKey)
	return s
}

// UnsetItalic removes the italic style rule, if set.
func (s Style) UnsetItalic() Style {
	delete(s, italicKey)
	return s
}

// UnsetItalic removes the underline style rule, if set.
func (s Style) UnsetUnderline() Style {
	delete(s, underlineKey)
	return s
}

// UnsetStrikethrough removes the strikethrough style rule, if set.
func (s Style) UnsetStrikethrough() Style {
	delete(s, strikethroughKey)
	return s
}

// UnsetReverse removes the reverse style rule, if set.
func (s Style) UnsetReverse() Style {
	delete(s, reverseKey)
	return s
}

// UnsetBlink removes the bold style rule, if set.
func (s Style) UnsetBlink() Style {
	delete(s, blinkKey)
	return s
}

// UnsetFaint removes the faint style rule, if set.
func (s Style) UnsetFaint() Style {
	delete(s, faintKey)
	return s
}

// UnsetForegroundremoves the foreground style rule, if set.
func (s Style) UnsetForeground() Style {
	delete(s, foregroundKey)
	return s
}

// UnsetBackground removes the background style rule, if set.
func (s Style) UnsetBackground() Style {
	delete(s, backgroundKey)
	return s
}

// UnsetWidth removes the width style rule, if set.
func (s Style) UnsetWidth() Style {
	delete(s, widthKey)
	return s
}

// UnsetAlign removes the text alignment style rule, if set.
func (s Style) UnsetAlign() Style {
	delete(s, alignKey)
	return s
}

// UnsetPadding removes all padding style rules.
func (s Style) UnsetPadding() Style {
	delete(s, leftPaddingKey)
	delete(s, rightPaddingKey)
	delete(s, topPaddingKey)
	delete(s, bottomPaddingKey)
	return s
}

// UnsetLeftPadding removes the left padding style rule, if set.
func (s Style) UnsetLeftPadding() Style {
	delete(s, leftPaddingKey)
	return s
}

// UnsetRightPadding removes the left padding style rule, if set.
func (s Style) UnsetRightPadding() Style {
	delete(s, rightPaddingKey)
	return s
}

// UnsetTopPadding removes the top padding style rule, if set.
func (s Style) UnsetTopPadding() Style {
	delete(s, topPaddingKey)
	return s
}

// UnsetBottomPadding removes the bottom style rule, if set.
func (s Style) UnsetBottomPadding() Style {
	delete(s, bottomPaddingKey)
	return s
}

// UnsetColorWhitespace removes the rule for coloring padding, if set.
func (s Style) UnsetColorWhitespace() Style {
	delete(s, colorWhitespaceKey)
	return s
}

// UnsetMargins removes all margin style rules.
func (s Style) UnsetMargins() Style {
	delete(s, leftMarginKey)
	delete(s, rightMarginKey)
	delete(s, topMarginKey)
	delete(s, bottomMarginKey)
	return s
}

// UnsetLeftMargin removes the left margin style rule, if set.
func (s Style) UnsetLeftMargin() Style {
	delete(s, leftMarginKey)
	return s
}

// UnsetRightMargin removes the right margin style rule, if set.
func (s Style) UnsetRightMargin() Style {
	delete(s, rightMarginKey)
	return s
}

// UnsetTopMargin removes the top margin style rule, if set.
func (s Style) UnsetTopMargin() Style {
	delete(s, topMarginKey)
	return s
}

// UnsetBottomMargin removes the bottom margin style rule, if set.
func (s Style) UnsetBottomMargin() Style {
	delete(s, bottomMarginKey)
	return s
}

// UnsetInline removes the inline style rule, if set.
func (s Style) UnsetInline() Style {
	delete(s, inlineKey)
	return s
}

// UnsetMaxWidth removes the max width style rule, if set.
func (s Style) UnsetMaxWidth() Style {
	delete(s, maxWidthKey)
	return s
}

// UnsetDrawClearTrailingSpaces removes the rule for drawing clear trailing
// spaces, if set.
func (s Style) UnsetDrawClearTrailingSpaces() Style {
	delete(s, drawClearTrailingSpacesKey)
	return s
}

// UnsetUnderlineWhitespace removes the rule for underlining whitespace, if
// set.
func (s Style) UnsetUnderlineWhitespace() Style {
	delete(s, underlineSpacesKey)
	return s
}

// UnsetUnderlineWhitespace removes the rule for strikingn through whitespace,
// if set.
func (s Style) UnsetStrikethroughWhitespace() Style {
	delete(s, strikethroughSpacesKey)
	return s
}

// UnsetUnderlineSpaces removes the value set by UnderlineSpaces.
func (s Style) UnsetUnderlineSpaces(v bool) Style {
	delete(s, underlineSpacesKey)
	return s
}

// UnsetUnderlineSpaces removes the value set by UnsetStrikethroughSpaces.
func (s Style) UnsetStrikethroughSpaces(v bool) Style {
	delete(s, strikethroughSpacesKey)
	return s
}
