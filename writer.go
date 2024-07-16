package lipgloss

import (
	"bytes"
	"fmt"
	"image/color"
	"io"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/ansi/parser"
)

// NewWriter creates a new Lip Gloss writer that writes text to the given writer.
//
// If environ is nil, it will use os.Environ() to get the environment variables.
//
// It queries the given writer to determine if it supports ANSI escape codes.
// If it does, along with the given environment variables, it will determine
// the appropriate color profile to use for color formatting.
//
// This respects the NO_COLOR, CLICOLOR, and CLICOLOR_FORCE environment variables.
func NewWriter(w io.Writer, environ []string) *Writer {
	return &Writer{
		Forward: w,
		Profile: DetectColorProfile(w, environ),
		parser:  ansi.NewParser(parser.MaxParamsSize, 0),
	}
}

// When a Lip Gloss writer is created, it queries the given writer to determine
// if it supports ANSI escape codes. If it does, the Lip Gloss writer will
// write text with ANSI escape codes. If it does not, the Lip Gloss writer will
// write text without ANSI escape codes.
// It also determines the appropriate color profile to use based on the
// capabilities of the underlying writer and environment.

// Writer represents a Lip Gloss writer that writes text to the underlying
// writer.
type Writer struct {
	parser *ansi.Parser

	Forward io.Writer
	Profile Profile
}

// Print writes the given text to the underlying writer.
func (w *Writer) Print(text string) (n int, err error) {
	return fmt.Fprint(w, text)
}

// Println writes the given text to the underlying writer followed by a newline.
func (w *Writer) Println(text string) (n int, err error) {
	return fmt.Fprintln(w, text)
}

// Printf writes the given text to the underlying writer with the given format.
func (w *Writer) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, format, a...)
}

// Write writes the given text to the underlying writer.
func (w *Writer) Write(p []byte) (int, error) {
	if w.Profile == TrueColor {
		return w.Forward.Write(p)
	} else if w.Profile == NoTTY {
		return io.WriteString(w.Forward, ansi.Strip(string(p)))
	}

	if w.parser == nil {
		w.parser = ansi.NewParser(parser.MaxParamsSize, 0)
	}

	convertColorAppend := func(c ansi.Color, sel colorSelector, pen *ansi.CsiSequence) {
		if c := w.Profile.Convert(c); c != nil && c != noColor {
			pen.Params = append(pen.Params, ansiColorToParams(c, sel)...)
		}
	}

	var buf bytes.Buffer
	for i, b := range p {
		w.parser.Advance(func(s ansi.Sequence) {
			pen := ansi.CsiSequence{Cmd: 'm'}
			switch s := s.(type) {
			case ansi.CsiSequence:
				switch s.Cmd {
				case 'm':
					if w.Profile > Ascii {
						return
					}
					for j := 0; j < len(s.Params); j++ {
						param := s.Param(j)
						switch param {
						case 30, 31, 32, 33, 34, 35, 36, 37: // 8-bit foreground color
							if w.Profile > ANSI {
								convertColorAppend(ansi.BasicColor(param-30), foreground, &pen)
								continue
							}
						case 39: // default foreground color
							if w.Profile > ANSI {
								continue
							}
						case 40, 41, 42, 43, 44, 45, 46, 47: // 8-bit background color
							if w.Profile > ANSI {
								convertColorAppend(ansi.BasicColor(param-40), background, &pen)
								continue
							}
						case 49: // default background color
							if w.Profile > ANSI {
								continue
							}
						case 90, 91, 92, 93, 94, 95, 96, 97: // 8-bit bright foreground color
							if w.Profile > ANSI {
								convertColorAppend(ansi.BasicColor(param-90+8), foreground, &pen)
								continue
							}
						case 100, 101, 102, 103, 104, 105, 106, 107: // 8-bit bright background color
							if w.Profile > ANSI {
								convertColorAppend(ansi.BasicColor(param-100+8), background, &pen)
								continue
							}
						case 59: // default underline color
							if w.Profile > ANSI {
								continue
							}
						case 38: // 16 or 24-bit foreground color
							fallthrough
						case 48: // 16 or 24-bit background color
							fallthrough
						case 58: // 16 or 24-bit underline color
							var sel colorSelector
							switch param {
							case 38:
								sel = foreground
							case 48:
								sel = background
							case 58:
								sel = underline
							}
							if c := readColor(&j, &s); c != nil {
								switch c.(type) {
								case ansi.ExtendedColor:
									if w.Profile > ANSI256 {
										convertColorAppend(c, sel, &pen)
										continue
									}
								default:
									if w.Profile > TrueColor {
										convertColorAppend(c, sel, &pen)
										continue
									}
								}
								pen.Params = append(pen.Params, ansiColorToParams(c, sel)...)
								continue
							}
						}
						pen.Params = append(pen.Params, param)
					}
					buf.Write(pen.Bytes())
					return
				}
			}
			buf.Write(s.Bytes())
		}, b, i < len(p)-1)
	}

	return w.Forward.Write(buf.Bytes())
}

// WriteString writes the given text to the underlying writer.
func (w *Writer) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

func readColor(idxp *int, seq *ansi.CsiSequence) (c ansi.Color) {
	i := *idxp
	paramsLen := len(seq.Params)
	// Note: we accept both main and subparams here
	switch seq.Param(i + 1) {
	case 2: // RGB
		if i+2 < paramsLen && i+3 < paramsLen && i+4 < paramsLen {
			c = color.RGBA{
				R: uint8(seq.Param(i + 2)),
				G: uint8(seq.Param(i + 3)),
				B: uint8(seq.Param(i + 4)),
				A: 0xff,
			}
			*idxp += 4
		}
	case 5: // 256 colors
		if i+2 < paramsLen {
			c = ansi.ExtendedColor(seq.Param(i + 2))
			*idxp += 2
		}
	}
	return
}

type colorSelector uint8

const (
	foreground colorSelector = iota
	background
	underline
)

func ansiColorToParams(c ansi.Color, sel colorSelector) []int {
	switch c := c.(type) {
	case ansi.BasicColor:
		offset := 30
		if c >= ansi.BrightBlack {
			offset = 90
			c -= ansi.BrightBlack
		}
		switch sel {
		case foreground:
			return []int{offset + int(c)}
		case background:
			return []int{offset + 10 + int(c)}
		case underline:
			// NOTE: ANSI doesn't have underline colors, use ANSI256.
			return []int{58, 5, int(c)}
		}
	case ansi.ExtendedColor:
		switch sel {
		case foreground:
			return []int{38, 5, int(c)}
		case background:
			return []int{48, 5, int(c)}
		case underline:
			return []int{58, 5, int(c)}
		}
	default:
		r, g, b, _ := c.RGBA()
		r = r >> 8
		g = g >> 8
		b = b >> 8
		switch sel {
		case foreground:
			return []int{38, 2, int(r), int(g), int(b)}
		case background:
			return []int{48, 2, int(r), int(g), int(b)}
		case underline:
			return []int{58, 2, int(r), int(g), int(b)}
		}
	}
	return nil
}
