package lipgloss_test

import (
	"errors"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
)

type wrapFullErrorWriter struct{}

func (wrapFullErrorWriter) Write(p []byte) (int, error) { return len(p), errWrapSink }

func TestWrapWriterFullWriteError(t *testing.T) {
	w := lipgloss.NewWrapWriter(wrapFullErrorWriter{})
	defer w.Close()
	n, err := w.Write([]byte("complete"))
	if n != 8 || !errors.Is(err, errWrapSink) {
		t.Fatalf("Write = (%d, %v), want (8, sink failed)", n, err)
	}
}

func TestWrapWriterTracksAcceptedStyle(t *testing.T) {
	const input = "\x1b[31m\x1b]8;;https://example.com\a"
	for _, cutoff := range []int{0, 4, 5, len(input) - 1, len(input)} {
		sink := &wrapSink{remaining: cutoff, err: errWrapSink}
		w := lipgloss.NewWrapWriter(sink)
		n, err := w.Write([]byte(input))
		if n != cutoff {
			t.Fatalf("Write count = %d, want %d", n, cutoff)
		}
		if cutoff < len(input) && !errors.Is(err, errWrapSink) {
			t.Fatalf("Write error = %v", err)
		}
		style, link := w.Style(), w.Link()
		if style.IsZero() != (cutoff < 5) {
			t.Fatalf("style at accepted byte %d = %v", cutoff, style)
		}
		if link.IsZero() != (cutoff < len(input)) {
			t.Fatalf("link at accepted byte %d = %v", cutoff, link)
		}
		_ = w.Close()
	}
}
