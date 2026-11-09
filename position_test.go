package lipgloss

import "testing"

func TestPlaceHorizontal(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		pos      Position
		str      string
		expected string
	}{
		{
			name:     "left constant",
			width:    10,
			pos:      Left,
			str:      "Hello",
			expected: "Hello     ",
		},
		{
			name:     "left zero float",
			width:    10,
			pos:      0.0,
			str:      "Hello",
			expected: "Hello     ",
		},
		{
			name:     "near left float",
			width:    10,
			pos:      0.000000001,
			str:      "Hello",
			expected: "Hello     ",
		},
		{
			name:     "right constant",
			width:    10,
			pos:      Right,
			str:      "Hello",
			expected: "     Hello",
		},
		{
			name:     "right one float",
			width:    10,
			pos:      1.0,
			str:      "Hello",
			expected: "     Hello",
		},
		{
			name:     "near right float",
			width:    10,
			pos:      0.999999999,
			str:      "Hello",
			expected: "     Hello",
		},
		{
			name:     "pos 0.2",
			width:    10,
			pos:      0.2,
			str:      "Hello",
			expected: " Hello    ",
		},
		{
			name:     "pos 0.5 center constant",
			width:    10,
			pos:      Center,
			str:      "Hello",
			expected: "   Hello  ",
		},
		{
			name:     "pos 0.8",
			width:    10,
			pos:      0.8,
			str:      "Hello",
			expected: "    Hello ",
		},
		{
			name:     "multiline",
			width:    6,
			pos:      0.25,
			str:      "Hi\nWorld",
			expected: " Hi   \nWorld ",
		},
		{
			name:     "no-op when width <= content width",
			width:    3,
			pos:      Center,
			str:      "Hello",
			expected: "Hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := PlaceHorizontal(tc.width, tc.pos, tc.str)
			if got != tc.expected {
				t.Errorf("PlaceHorizontal(%d, %v, %q) = %q, want %q", tc.width, tc.pos, tc.str, got, tc.expected)
			}
		})
	}
}

func TestPlaceVertical(t *testing.T) {
	tests := []struct {
		name     string
		height   int
		pos      Position
		str      string
		expected string
	}{
		{
			name:     "top constant",
			height:   3,
			pos:      Top,
			str:      "Hello",
			expected: "Hello\n     \n     ",
		},
		{
			name:     "top zero float",
			height:   3,
			pos:      0.0,
			str:      "Hello",
			expected: "Hello\n     \n     ",
		},
		{
			name:     "near top float",
			height:   3,
			pos:      0.000000001,
			str:      "Hello",
			expected: "Hello\n     \n     ",
		},
		{
			name:     "bottom constant",
			height:   3,
			pos:      Bottom,
			str:      "Hello",
			expected: "     \n     \nHello",
		},
		{
			name:     "bottom one float",
			height:   3,
			pos:      1.0,
			str:      "Hello",
			expected: "     \n     \nHello",
		},
		{
			name:     "near bottom float",
			height:   3,
			pos:      0.999999999,
			str:      "Hello",
			expected: "     \n     \nHello",
		},
		{
			name:     "center constant",
			height:   3,
			pos:      Center,
			str:      "Hello",
			expected: "     \nHello\n     ",
		},
		{
			name:     "pos 0.25 with 4 lines",
			height:   5,
			pos:      0.25,
			str:      "Hello",
			expected: "     \nHello\n     \n     \n     ",
		},
		{
			name:     "no-op when height <= content height",
			height:   1,
			pos:      Center,
			str:      "Hello",
			expected: "Hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := PlaceVertical(tc.height, tc.pos, tc.str)
			if got != tc.expected {
				t.Errorf("PlaceVertical(%d, %v, %q) = %q, want %q", tc.height, tc.pos, tc.str, got, tc.expected)
			}
		})
	}
}

func TestPlace(t *testing.T) {
	got := Place(5, 3, Left, Top, "X")
	expected := "X    \n     \n     "
	if got != expected {
		t.Errorf("Place(5, 3, Left, Top, %q) = %q, want %q", "X", got, expected)
	}

	got = Place(5, 3, Right, Bottom, "X")
	expected = "     \n     \n    X"
	if got != expected {
		t.Errorf("Place(5, 3, Right, Bottom, %q) = %q, want %q", "X", got, expected)
	}
}
