package lipgloss

import (
	"io"
	"testing"
	"time"
)

// TestPlaceWithTabCharacterHang tests that Place doesn't hang when using
// tab characters with WithWhitespaceChars option.
func TestPlaceWithTabCharacterHang(t *testing.T) {
	// Set up a timeout to prevent test from hanging indefinitely
	done := make(chan bool)

	go func() {
		// This should not hang
		_ = Place(10, 3, Center, Center, "hello",
			WithWhitespaceChars("\t"),
		)
		done <- true
	}()

	select {
	case <-done:
		// Test passed - function completed
	case <-time.After(2 * time.Second):
		t.Fatal("Place() hung when using tab character in WithWhitespaceChars - issue #108")
	}
}

// TestPlaceWithTabVariations tests various tab character combinations
func TestPlaceWithTabVariations(t *testing.T) {
	testCases := []struct {
		name  string
		chars string
	}{
		{"tab character", "\t"},
		{"multiple tabs", "\t\t"},
		{"tab and space", "\t "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			done := make(chan bool)

			go func() {
				_ = Place(10, 3, Center, Center, "hello",
					WithWhitespaceChars(tc.chars),
				)
				done <- true
			}()

			select {
			case <-done:
				// Test passed
			case <-time.After(2 * time.Second):
				t.Fatalf("Place() hung with whitespace chars %q", tc.chars)
			}
		})
	}
}

// TestWhitespaceRenderWithTabChar tests the whitespace.render method directly with tabs
func TestWhitespaceRenderWithTabChar(t *testing.T) {
	r := NewRenderer(io.Discard)

	testCases := []struct {
		name  string
		chars string
		width int
	}{
		{"tab character small width", "\t", 5},
		{"tab character large width", "\t", 20},
		{"multiple tabs", "\t\t", 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			done := make(chan bool)

			go func() {
				ws := newWhitespace(r, WithWhitespaceChars(tc.chars))
				_ = ws.render(tc.width)
				done <- true
			}()

			select {
			case <-done:
				// Test passed
			case <-time.After(2 * time.Second):
				t.Fatalf("whitespace.render() hung with chars %q and width %d", tc.chars, tc.width)
			}
		})
	}
}
