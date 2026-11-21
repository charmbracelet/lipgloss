package lipgloss

import (
	"bytes"
	"io"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
)

// Wrap wraps the given string to the given width, preserving ANSI styles and links.
func Wrap(s string, width int, breakpoints string) string {
	var buf bytes.Buffer
	s = ansi.Wrap(s, width, breakpoints)
	// The writer only inserts resets around newlines, so the output is
	// the same length as the input plus a little. Size it once.
	buf.Grow(len(s) + len(s)/8)
	w := NewWrapWriter(&buf)
	_, _ = io.WriteString(w, s)
	_ = w.Close()
	return buf.String()
}

// WrapWriter is a writer that writes to a buffer and keeps track of the
// current pen style and link state for the purpose of wrapping with newlines.
//
// When it encounters a newline, it resets the style and link, writes the
// newline, and then reapplies the style and link to the next line.
type WrapWriter struct {
	w     io.Writer
	p     *ansi.Parser
	style uv.Style
	link  uv.Link
}

// NewWrapWriter returns a new [WrapWriter].
func NewWrapWriter(w io.Writer) *WrapWriter {
	pw := &WrapWriter{w: w}
	pw.p = ansi.GetParser()
	handleCsi := func(cmd ansi.Cmd, params ansi.Params) {
		if cmd == 'm' {
			uv.ReadStyle(params, &pw.style)
		}
	}
	handleOsc := func(cmd int, data []byte) {
		if cmd == 8 {
			uv.ReadLink(data, &pw.link)
		}
	}
	pw.p.SetHandler(ansi.Handler{
		HandleCsi: handleCsi,
		HandleOsc: handleOsc,
	})
	return pw
}

// Style returns the current pen style.
func (w *WrapWriter) Style() uv.Style {
	return w.style
}

// Link returns the current pen link.
func (w *WrapWriter) Link() uv.Link {
	return w.link
}

// Write writes to the buffer.
func (w *WrapWriter) Write(p []byte) (int, error) {
	if w.p == nil {
		// The writer has been closed and its parser returned to the pool.
		// Writing after close can happen during out-of-order teardown of
		// nested writer chains; treat it as a no-op rather than panicking.
		return len(p), nil
	}
	// Bytes are forwarded in runs rather than one at a time: only a
	// newline needs anything written around it, so everything between
	// newlines is a single downstream Write instead of one per byte.
	start := 0
	for i := range p {
		b := p[i]
		w.p.Advance(b)
		if b != '\n' {
			continue
		}

		if i > start {
			_, _ = w.w.Write(p[start:i])
		}
		if !w.style.IsZero() {
			_, _ = w.w.Write([]byte(ansi.ResetStyle))
		}
		if !w.link.IsZero() {
			_, _ = w.w.Write([]byte(ansi.ResetHyperlink()))
		}

		_, _ = w.w.Write(p[i : i+1])
		start = i + 1

		if !w.link.IsZero() {
			_, _ = w.w.Write([]byte(ansi.SetHyperlink(w.link.URL, w.link.Params)))
		}
		if !w.style.IsZero() {
			_, _ = w.w.Write([]byte(w.style.String()))
		}
	}
	if start < len(p) {
		_, _ = w.w.Write(p[start:])
	}

	return len(p), nil
}

// Close closes the writer, resets the style and link if necessary, and releases
// its parser. Calling it is performance critical, but forgetting it does not
// cause safety issues or leaks.
func (w *WrapWriter) Close() error {
	if !w.style.IsZero() {
		_, _ = w.w.Write([]byte(ansi.ResetStyle))
	}
	if !w.link.IsZero() {
		_, _ = w.w.Write([]byte(ansi.ResetHyperlink()))
	}
	if w.p != nil {
		ansi.PutParser(w.p)
		w.p = nil
	}
	return nil
}
