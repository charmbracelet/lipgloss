package lipgloss

import (
	"strings"
	"testing"
)

func TestNoTabConversionMultiline(t *testing.T) {
	s := NewStyle().TabWidth(NoTabConversion)

	got := s.Render("a\tb")
	if got != "a\tb" {
		t.Fatalf("single-line: got %q want %q", got, "a\tb")
	}

	got = s.Render("a\tb\nc\td")
	if got != "a\tb\nc\td" {
		t.Fatalf("multi-line: got %q want %q", got, "a\tb\nc\td")
	}
}

func TestNoTabConversionMultilineWithColor(t *testing.T) {
	s := NewStyle().TabWidth(NoTabConversion).Foreground(Color("2"))
	got := s.Render("a\tb\nc\td")
	if !strings.Contains(got, "\t") {
		t.Fatalf("colored multi-line lost tabs: %q", got)
	}
	if strings.Contains(got, "    ") {
		t.Fatalf("colored multi-line expanded tabs to spaces: %q", got)
	}
}

func TestDefaultTabConversionStillAppliesOnMultiline(t *testing.T) {
	s := NewStyle()
	got := s.Render("a\tb\nc\td")
	want := "a    b\nc    d"
	if got != want {
		t.Fatalf("default multi-line tabs: got %q want %q", got, want)
	}
}
