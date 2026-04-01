package lipgloss

import "testing"

func TestGetHorizontalFrameSize(t *testing.T) {
	tests := []struct {
		name  string
		style Style
		want  int
	}{
		{"empty style", NewStyle(), 0},
		{"padding only", NewStyle().Padding(0, 2), 4},
		{"margin only", NewStyle().Margin(0, 3), 6},
		{"border only", NewStyle().Border(NormalBorder()), 2},
		{"all combined", NewStyle().Padding(0, 1).Margin(0, 1).Border(NormalBorder()), 6},
		{"border style without sides", NewStyle().BorderStyle(RoundedBorder()), 2},
		{"left border only", NewStyle().BorderStyle(NormalBorder()).BorderLeft(true), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.GetHorizontalFrameSize()
			if got != tt.want {
				t.Errorf("GetHorizontalFrameSize() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGetVerticalFrameSize(t *testing.T) {
	tests := []struct {
		name  string
		style Style
		want  int
	}{
		{"empty style", NewStyle(), 0},
		{"padding only", NewStyle().Padding(2, 0), 4},
		{"margin only", NewStyle().Margin(1, 0), 2},
		{"border only", NewStyle().Border(NormalBorder()), 2},
		{"all combined", NewStyle().Padding(1, 0).Margin(1, 0).Border(NormalBorder()), 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.GetVerticalFrameSize()
			if got != tt.want {
				t.Errorf("GetVerticalFrameSize() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGetFrameSize(t *testing.T) {
	s := NewStyle().Padding(1, 2).Margin(1, 1).Border(RoundedBorder())
	h, v := s.GetFrameSize()
	wantH := 2 + 4 + 2 // margin + padding + border
	wantV := 2 + 2 + 2  // margin + padding + border
	if h != wantH {
		t.Errorf("GetFrameSize() horizontal = %d, want %d", h, wantH)
	}
	if v != wantV {
		t.Errorf("GetFrameSize() vertical = %d, want %d", v, wantV)
	}
}
