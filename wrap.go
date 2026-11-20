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
	err   error
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

// Write writes to the buffer. After a downstream error, further writes return
// that error without forwarding more output.
func (w *WrapWriter) Write(p []byte) (int, error) {
	if w.p == nil {
		// The writer has been closed and its parser returned to the pool.
		// Writing after close can happen during out-of-order teardown of
		// nested writer chains; treat it as a no-op rather than panicking.
		return len(p), nil
	}
	if w.err != nil {
		return 0, w.err
	}

	n := 0
	for n < len(p) {
		i := bytes.IndexByte(p[n:], '\n')
		if i < 0 {
			i = len(p) - n
		}
		if i > 0 {
			m, err := w.write(p[n : n+i])
			for _, b := range p[n : n+m] {
				w.p.Advance(b)
			}
			n += m
			if err != nil {
				return n, err
			}
		}
		if n == len(p) {
			break
		}

		w.p.Advance('\n')
		if !w.style.IsZero() {
			if _, err := w.write([]byte(ansi.ResetStyle)); err != nil {
				return n, err
			}
		}
		if !w.link.IsZero() {
			if _, err := w.write([]byte(ansi.ResetHyperlink())); err != nil {
				return n, err
			}
		}

		m, err := w.write(p[n : n+1])
		n += m
		if err != nil {
			return n, err
		}

		if !w.link.IsZero() {
			if _, err := w.write([]byte(ansi.SetHyperlink(w.link.URL, w.link.Params))); err != nil {
				return n, err
			}
		}
		if !w.style.IsZero() {
			if _, err := w.write([]byte(w.style.String())); err != nil {
				return n, err
			}
		}
	}

	return n, nil
}

// write retains the first failure: a partially written ANSI sequence cannot be
// safely resumed by a later caller write.
func (w *WrapWriter) write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	if n < len(p) && err == nil {
		err = io.ErrShortWrite
	}
	w.err = err
	return n, err
}

// Close closes the writer, resets the style and link if necessary, and releases
// its parser. Calling it is performance critical, but forgetting it does not
// cause safety issues or leaks. It returns the first downstream error.
func (w *WrapWriter) Close() error {
	if w.p == nil {
		return nil
	}
	defer func() {
		ansi.PutParser(w.p)
		w.p = nil
	}()
	if w.err != nil {
		return w.err
	}
	if !w.style.IsZero() {
		if _, err := w.write([]byte(ansi.ResetStyle)); err != nil {
			return err
		}
	}
	if !w.link.IsZero() {
		if _, err := w.write([]byte(ansi.ResetHyperlink())); err != nil {
			return err
		}
	}
	return nil
}
