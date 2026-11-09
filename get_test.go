package lipgloss

import "testing"

func BenchmarkGetters(b *testing.B) {
	s := NewStyle().Bold(true).Foreground(Color("#ff0000"))
	for b.Loop() {
		_, _, _, _ = s.GetPadding()
		_, _, _, _ = s.GetMargin()
		_, _, _, _, _ = s.GetBorder()
		_ = s.GetWidth()
		_ = s.GetHeight()
		_ = s.GetMaxWidth()
		_ = s.GetMaxHeight()
		_ = s.GetBold()
		_ = s.GetItalic()
		_ = s.GetFaint()
		_ = s.GetBlink()
		_ = s.GetReverse()
		_ = s.GetStrikethrough()
		_ = s.GetForeground()
	}
}
