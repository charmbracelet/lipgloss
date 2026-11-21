package lipgloss

import "testing"

// Bool props are stored as bits of Style.attrs indexed by propKey. Keys
// beyond bit 31 (borderLeftKey, inlineKey) only survive if attrs is 64-bit:
// on 32-bit platforms an int-typed attrs truncates int(1<<32) to zero and
// the value is silently dropped. Run with GOARCH=386 to exercise that case.
func TestBoolPropsBeyond32Bits(t *testing.T) {
	s := NewStyle().BorderStyle(NormalBorder()).BorderLeft(true).Inline(true)
	if !s.GetBorderLeft() {
		t.Error("BorderLeft(true) reads back false")
	}
	if !s.GetInline() {
		t.Error("Inline(true) reads back false")
	}
	s = s.BorderLeft(false).Inline(false)
	if s.GetBorderLeft() {
		t.Error("BorderLeft(false) reads back true")
	}
	if s.GetInline() {
		t.Error("Inline(false) reads back true")
	}
}
