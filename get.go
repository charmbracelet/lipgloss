package lipgloss

import (
	"strings"

	"github.com/muesli/reflow/ansi"
)

// Returns whether or not the given property is set.
func (s Style) isSet(k propKey) bool {
	_, exists := s.rules[k]
	return exists
}

func (s Style) getAsBool(k propKey, defaultVal bool) bool {
	v, ok := s.rules[k]
	if !ok {
		return defaultVal
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return defaultVal
}

func (s Style) getAsColor(k propKey) ColorType {
	v, ok := s.rules[k]
	if !ok {
		return NoColor{}
	}
	if c, ok := v.(ColorType); ok {
		return c
	}
	return NoColor{}
}

func (s Style) getAsInt(k propKey) int {
	v, ok := s.rules[k]
	if !ok {
		return 0
	}
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

func (s Style) getAsPosition(k propKey) Position {
	v, ok := s.rules[k]
	if !ok {
		return Position(0)
	}
	if p, ok := v.(Position); ok {
		return p
	}
	return Position(0)
}

func (s Style) getAsBorderStyle(k propKey) Border {
	v, ok := s.rules[k]
	if !ok {
		return noBorder
	}
	if b, ok := v.(Border); ok {
		return b
	}
	return noBorder
}

// Split a string into lines, additionally returning the size of the widest
// line.
func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := ansi.PrintableRuneWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}
