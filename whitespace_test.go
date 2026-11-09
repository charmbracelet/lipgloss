package lipgloss

import (
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
)

func TestWhitespaceRenderWithTab(t *testing.T) {
	// This test verifies that rendering whitespace with tab characters
	// doesn't cause an infinite loop (issue #108)
	done := make(chan bool, 1)

	go func() {
		ws := newWhitespace(WithWhitespaceChars("\t"))
		_ = ws.render(10)
		done <- true
	}()

	select {
	case <-done:
		// Success - render completed
	case <-time.After(2 * time.Second):
		t.Fatal("whitespace.render() with tab character caused infinite loop")
	}
}

func TestWhitespaceRenderWithZeroWidthChar(t *testing.T) {
	// Test with zero-width joiner (another zero-width character)
	done := make(chan bool, 1)

	go func() {
		ws := newWhitespace(WithWhitespaceChars("\u200d")) // zero-width joiner
		_ = ws.render(5)
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("whitespace.render() with zero-width character caused infinite loop")
	}
}

func TestWhitespaceRenderNormal(t *testing.T) {
	// Verify normal behavior still works
	ws := newWhitespace(WithWhitespaceChars("*"))
	result := ws.render(5)
	if len(result) != 5 {
		t.Errorf("expected 5 characters, got %d", len(result))
	}
}

func TestWhitespaceRenderWideChars(t *testing.T) {
	// Rendering with wide (multi-cell) characters must never exceed the
	// requested width. Any leftover cells are padded with spaces instead.
	for _, tc := range []struct {
		chars string
		width int
	}{
		{"橋", 1},
		{"橋", 3},
		{"橋", 5},
		{"a橋", 4},
		{"橋", 0},
	} {
		ws := newWhitespace(WithWhitespaceChars(tc.chars))
		got := ansi.StringWidth(ws.render(tc.width))
		if got != tc.width {
			t.Errorf("render(%d) with chars %q produced width %d, want %d",
				tc.width, tc.chars, got, tc.width)
		}
	}
}
