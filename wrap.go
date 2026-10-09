package lipgloss

import (
	"bytes"
	"io"
	"strconv"
	"strings"

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
	defer w.Close() //nolint:errcheck
	_, _ = io.WriteString(w, s)
	return buf.String()
}

// isolateANSIAtLineBoundaries makes every line self-contained while
// preserving the style and hyperlink state across newlines. Unlike Wrap, it
// closes the final state so content cannot style anything appended to it.
func isolateANSIAtLineBoundaries(s string) string {
	input := s
	var buf strings.Builder
	buf.Grow(len(s) + len(s)/8)

	var state byte
	var style sgrState
	var linkReplay string
	for len(s) > 0 {
		// The pooled parser has a fixed parameter buffer. We only need the raw
		// sequence here, so decode without collecting params to keep long but
		// valid SGR sequences from overflowing that buffer.
		seq, _, n, newState := ansi.DecodeSequence(s, state, nil)
		if n <= 0 {
			return input
		}

		// An invalid or incomplete CSI sequence has no safe boundary at which
		// to inject resets. Preserve the input instead of splitting it.
		_, validCSI := csiCommand(seq)
		if ansi.HasCsiPrefix(seq) && !validCSI {
			return input
		}

		if seq == "\n" {
			styleReplay := style.sequence()
			if styleReplay != "" {
				buf.WriteString(ansi.ResetStyle)
			}
			if linkReplay != "" {
				buf.WriteString(ansi.ResetHyperlink())
			}
			buf.WriteByte('\n')
			buf.WriteString(linkReplay)
			buf.WriteString(styleReplay)
		} else {
			buf.WriteString(seq)
			switch {
			case isSGRSequence(seq):
				if !style.update(seq) {
					return input
				}
			case ansi.HasOscPrefix(seq):
				linkReplay = updateHyperlinkReplay(linkReplay, seq)
			}
		}
		// Replaying a very large active sequence on many short lines can
		// otherwise amplify untrusted input without bound.
		if buf.Len() > maxANSIIsolationSlack &&
			(buf.Len()-maxANSIIsolationSlack)/maxANSIIsolationExpansion > len(input) {
			return input
		}

		state = newState
		s = s[n:]
	}
	if state != ansi.NormalState {
		return input
	}
	if style.sequence() != "" {
		buf.WriteString(ansi.ResetStyle)
	}
	if linkReplay != "" {
		buf.WriteString(ansi.ResetHyperlink())
	}
	return buf.String()
}

func csiCommand(sequence string) (byte, bool) {
	if !ansi.HasCsiPrefix(sequence) || len(sequence) < 2 {
		return 0, false
	}
	command := sequence[len(sequence)-1]
	return command, command >= '@' && command <= '~'
}

func isSGRSequence(sequence string) bool {
	command, ok := csiCommand(sequence)
	if !ok || command != 'm' {
		return false
	}
	start := 1
	if strings.HasPrefix(sequence, "\x1b[") {
		start = 2
	}
	for i := start; i < len(sequence)-1; i++ {
		c := sequence[i]
		if (c < '0' || c > '9') && c != ';' && c != ':' {
			return false
		}
	}
	return true
}

const (
	maxANSIIsolationExpansion = 64
	maxANSIIsolationSlack     = 4096
	maxSGRSettings            = 128
)

type sgrSetting struct {
	key    string
	params string
}

type sgrState []sgrSetting

// update keeps only the currently active SGR settings. This avoids replaying
// an ever-growing history when content repeatedly enables and selectively
// disables attributes between lines.
func (s *sgrState) update(sequence string) bool {
	start := 1
	if strings.HasPrefix(sequence, "\x1b[") {
		start = 2
	}
	if len(sequence) <= start || sequence[len(sequence)-1] != 'm' {
		return false
	}

	params := strings.Split(sequence[start:len(sequence)-1], ";")
	for i := 0; i < len(params); i++ {
		code, ok := sgrCode(params[i])
		if !ok {
			return false
		}
		end := sgrParamGroupEnd(params, i, code)
		if code == 0 {
			*s = (*s)[:0]
			i = end
			continue
		}
		if resets := sgrResetKeys(code); len(resets) > 0 {
			s.remove(resets...)
			i = end
			continue
		}
		s.set(sgrSettingKey(code), strings.Join(params[i:end+1], ";"))
		if len(*s) > maxSGRSettings {
			return false
		}
		i = end
	}
	return true
}

func (s sgrState) sequence() string {
	if len(s) == 0 {
		return ""
	}
	var buf strings.Builder
	buf.WriteString("\x1b[")
	for i, setting := range s {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.WriteString(setting.params)
	}
	buf.WriteByte('m')
	return buf.String()
}

func (s *sgrState) set(key, params string) {
	s.remove(key)
	*s = append(*s, sgrSetting{key: key, params: params})
}

func (s *sgrState) remove(keys ...string) {
	settings := *s
	for _, key := range keys {
		for i := 0; i < len(settings); i++ {
			if settings[i].key == key {
				settings = append(settings[:i], settings[i+1:]...)
				i--
			}
		}
	}
	*s = settings
}

func sgrCode(param string) (int, bool) {
	if param == "" {
		return 0, true
	}
	if i := strings.IndexByte(param, ':'); i >= 0 {
		param = param[:i]
	}
	code, err := strconv.Atoi(param)
	return code, err == nil
}

// sgrParamGroupEnd skips color components so zero-valued RGB channels are not
// mistaken for a full SGR reset.
func sgrParamGroupEnd(params []string, i int, code int) int {
	if (code != 38 && code != 48 && code != 58) || strings.ContainsRune(params[i], ':') || i+1 >= len(params) {
		return i
	}
	mode, ok := sgrCode(params[i+1])
	if !ok {
		return i
	}
	span := 0
	switch mode {
	case 0, 1:
		span = 1
	case 2, 3:
		span = 4
	case 4, 6:
		span = 5
	case 5:
		span = 2
	}
	return min(i+span, len(params)-1)
}

func sgrSettingKey(code int) string {
	switch {
	case code == 1:
		return "bold"
	case code == 2:
		return "faint"
	case code == 4:
		return "underline"
	case code == 7:
		return "inverse"
	case code == 8:
		return "conceal"
	case code == 9:
		return "strike"
	case code >= 10 && code <= 19:
		return "font"
	case code == 26:
		return "proportional"
	case (code >= 30 && code <= 37) || code == 38 || (code >= 90 && code <= 97):
		return "foreground"
	case (code >= 40 && code <= 47) || code == 48 || (code >= 100 && code <= 107):
		return "background"
	case code == 53:
		return "overline"
	case code == 58:
		return "underline-color"
	default:
		return strconv.Itoa(code)
	}
}

func sgrResetKeys(code int) []string {
	switch code {
	case 10:
		return []string{"font"}
	case 22:
		return []string{"bold", "faint"}
	case 23:
		return []string{"3", "20"}
	case 24:
		return []string{"underline", "21"}
	case 25:
		return []string{"5", "6"}
	case 27:
		return []string{"inverse"}
	case 28:
		return []string{"conceal"}
	case 29:
		return []string{"strike"}
	case 39:
		return []string{"foreground"}
	case 49:
		return []string{"background"}
	case 50:
		return []string{"proportional"}
	case 54:
		return []string{"51", "52"}
	case 55:
		return []string{"overline"}
	case 59:
		return []string{"underline-color"}
	case 65:
		return []string{"60", "61", "62", "63", "64"}
	case 75:
		return []string{"73", "74"}
	default:
		return nil
	}
}

func updateHyperlinkReplay(current, sequence string) string {
	start := 1
	if strings.HasPrefix(sequence, "\x1b]") {
		start = 2
	}
	end := len(sequence)
	switch {
	case end > start && sequence[end-1] == ansi.BEL:
		end--
	case end > start && sequence[end-1] == ansi.ST:
		end--
	case end >= start+2 && sequence[end-2:] == "\x1b\\":
		end -= 2
	default:
		return current
	}
	data := []byte(sequence[start:end])
	params := bytes.SplitN(data, []byte{';'}, 3)
	if len(params) != 3 {
		return current
	}
	command, err := strconv.Atoi(string(params[0]))
	if err != nil || command != 8 {
		return current
	}
	if len(params[2]) == 0 {
		return ""
	}
	return sequence
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
