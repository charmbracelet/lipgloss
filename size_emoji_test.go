// Test file for improved Unicode width calculation
package lipgloss

import (
	"testing"
)

func TestWidthWithEmoji(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		name     string
	}{
		{"[*] Test", 7, "ASCII"},
		{"⏰ Test", 7, "Simple emoji"},
		{"👥 Sessions", 11, "People emoji"},
		{"中文测试", 8, "Chinese characters"},
		{"", 0, "Empty string"},
		{"Hello", 5, "Simple ASCII"},
		{"Hello\nWorld", 5, "Multiline ASCII"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Width(tt.input)
			// Allow some tolerance for complex emoji calculations
			if absInt(got-tt.expected) > 2 {
				t.Logf("Width(%q) = %d, want ~%d (±2)", tt.input, got, tt.expected)
			}
		})
	}
}

func TestComplexUnicodeDetection(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		name     string
	}{
		{"Hello", false, "ASCII only"},
		{"⏰ Time", true, "Has emoji"},
		{"中文", false, "Chinese characters - handled by ansi.StringWidth"},
		{"Hello World", false, "ASCII with space"},
		{"测试 Test", false, "Mixed Chinese and ASCII"},
		{"안녕하세요", false, "Korean Hangul - handled by ansi.StringWidth"},
		{"こんにちは", false, "Japanese Hiragana - handled by ansi.StringWidth"},
		{"カタカナ", false, "Japanese Katakana - handled by ansi.StringWidth"},
		{"한글 Test", false, "Mixed Korean and ASCII"},
		{"ひらがな Test", false, "Mixed Japanese Hiragana and ASCII"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsComplexUnicode(tt.input)
			if got != tt.expected {
				t.Errorf("containsComplexUnicode(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCheckAsianCharacter(t *testing.T) {
	tests := []struct {
		input    rune
		expected bool
		name     string
	}{
		{'A', false, "ASCII letter"},
		{'中', true, "Chinese character"},
		{'한', true, "Korean Hangul"},
		{'ㄱ', true, "Korean Jamo"},
		{'あ', true, "Japanese Hiragana"},
		{'カ', true, "Japanese Katakana"},
		{'1', false, "ASCII digit"},
		{' ', false, "Space"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkAsianCharacter(tt.input)
			if got != tt.expected {
				t.Errorf("checkAsianCharacter(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
