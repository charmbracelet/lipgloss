package lipgloss

import "testing"

// TestInlineReplacesNewlinesWithSpace verifies that Inline mode replaces
// newlines with spaces so that words from separate lines remain separated,
// rather than being concatenated. This is the behavior described in #116.
func TestInlineReplacesNewlinesWithSpace(t *testing.T) {
	style := NewStyle().Inline(true)

	tests := []struct {
		input    string
		expected string
	}{
		{"hello\nworld", "hello world"},
		{"hello\r\nworld", "hello world"},
		{"a\nb\nc", "a b c"},
		{"single", "single"},
	}

	for _, tt := range tests {
		actual := style.Render(tt.input)
		if actual != tt.expected {
			t.Errorf("inline Render(%q): got %q, want %q", tt.input, actual, tt.expected)
		}
	}
}
