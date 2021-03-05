package lipgloss

// UnsetBold removes the bold style rule, if set.
func (s Style) UnsetBold() Style {
	s.bold = nil
	return s
}

// UnsetItalic removes the italic style rule, if set.
func (s Style) UnsetItalic() Style {
	s.italic = nil
	return s
}

// UnsetItalic removes the underline style rule, if set.
func (s Style) UnsetUnderline() Style {
	s.underline = nil
	return s
}

// UnsetStrikethrough removes the strikethrough style rule, if set.
func (s Style) UnsetStrikethrough() Style {
	s.strikethrough = nil
	return s
}

// UnsetReverse removes the reverse style rule, if set.
func (s Style) UnsetReverse() Style {
	s.reverse = nil
	return s
}

// UnsetBlink removes the bold style rule, if set.
func (s Style) UnsetBlink() Style {
	s.blink = nil
	return s
}

// UnsetFaint removes the faint style rule, if set.
func (s Style) UnsetFaint() Style {
	s.faint = nil
	return s
}

// UnsetForegroundremoves the foreground style rule, if set.
func (s Style) UnsetForeground() Style {
	s.foreground = nil
	return s
}

// UnsetBackground removes the background style rule, if set.
func (s Style) UnsetBackground() Style {
	s.background = nil
	return s
}

// UnsetWidth removes the width style rule, if set.
func (s Style) UnsetWidth() Style {
	s.width = nil
	return s
}

// UnsetAlign removes the text alignment style rule, if set.
func (s Style) UnsetAlign() Style {
	s.align = nil
	return s
}

// UnsetPadding removes all padding style rules.
func (s Style) UnsetPadding() Style {
	s.leftPadding = nil
	s.rightPadding = nil
	s.topPadding = nil
	s.bottomPadding = nil
	return s
}

// UnsetLeftPadding removes the left padding style rule, if set.
func (s Style) UnsetLeftPadding() Style {
	s.leftPadding = nil
	return s
}

// UnsetRightPadding removes the left padding style rule, if set.
func (s Style) UnsetRightPadding() Style {
	s.rightPadding = nil
	return s
}

// UnsetTopPadding removes the top padding style rule, if set.
func (s Style) UnsetTopPadding() Style {
	s.topPadding = nil
	return s
}

// UnsetBottomPadding removes the bottom style rule, if set.
func (s Style) UnsetBottomPadding() Style {
	s.bottomPadding = nil
	return s
}

// UnsetColorWhitespace removes the rule for coloring padding, if set.
func (s Style) UnsetColorWhitespace() Style {
	s.colorWhitespace = nil
	return s
}

// UnsetMargins removes all margin style rules.
func (s Style) UnsetMargins() Style {
	s.leftMargin = nil
	s.rightMargin = nil
	s.topMargin = nil
	s.bottomMargin = nil
	return s
}

// UnsetLeftMargin removes the left margin style rule, if set.
func (s Style) UnsetLeftMargin() Style {
	s.leftMargin = nil
	return s
}

// UnsetRightMargin removes the right margin style rule, if set.
func (s Style) UnsetRightMargin() Style {
	s.rightMargin = nil
	return s
}

// UnsetTopMargin removes the top margin style rule, if set.
func (s Style) UnsetTopMargin() Style {
	s.topMargin = nil
	return s
}

// UnsetBottomMargin removes the bottom margin style rule, if set.
func (s Style) UnsetBottomMargin() Style {
	s.bottomMargin = nil
	return s
}

// UnsetInline removes the inline style rule, if set.
func (s Style) UnsetInline() Style {
	s.inline = nil
	return s
}

// UnsetMaxWidth removes the max width style rule, if set.
func (s Style) UnsetMaxWidth() Style {
	s.maxWidth = nil
	return s
}

// UnsetDrawClearTrailingSpaces removes the rule for drawing clear trailing
// spaces, if set.
func (s Style) UnsetDrawClearTrailingSpaces() Style {
	s.drawClearTrailingSpaces = nil
	return s
}

// UnsetUnderlineWhitespace removes the rule for underlining whitespace, if
// set.
func (s Style) UnsetUnderlineWhitespace() Style {
	s.underlineWhitespace = nil
	return s
}

// UnsetUnderlineWhitespace removes the rule for strikingn through whitespace,
// if set.
func (s Style) UnsetStrikethroughWhitespace() Style {
	s.strikethroughWhitespace = nil
	return s
}

func (s Style) UnsetUnderlineSpaces(v bool) Style {
	s.underlineSpaces = nil
	return s
}

func (s Style) UnsetStrikethroughSpaces(v bool) Style {
	s.strikethroughSpaces = nil
	return s
}
