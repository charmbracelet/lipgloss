package lipgloss

import (
	"testing"
	"time"
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
