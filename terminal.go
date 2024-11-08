package lipgloss

import (
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/ansi/parser"
	"github.com/muesli/cancelreader"
)

// queryBackgroundColor queries the terminal for the background color.
// If the terminal does not support querying the background color, nil is
// returned.
//
// Note: you will need to set the input to raw mode before calling this
// function.
//
//	state, _ := term.MakeRaw(in.Fd())
//	defer term.Restore(in.Fd(), state)
//
// copied from x/term@v0.1.3.
func queryBackgroundColor(in io.Reader, out io.Writer) (c color.Color, err error) {
	//nolint: errcheck
	err = queryTerminal(in, out, defaultQueryTimeout,
		func(seq string, pa *ansi.Parser) bool {
			switch {
			case ansi.HasOscPrefix(seq):
				switch pa.Cmd {
				case 11: // OSC 11
					parts := strings.Split(string(pa.Data[:pa.DataLen]), ";")
					if len(parts) != 2 {
						break // invalid, but we still need to parse the next sequence
					}
					c = xParseColor(parts[1])
				}
			case ansi.HasCsiPrefix(seq):
				switch pa.Cmd {
				case 'c' | '?'<<parser.MarkerShift: // DA1
					return false
				}
			}
			return true
		}, ansi.RequestBackgroundColor+ansi.RequestPrimaryDeviceAttributes)
	return
}

const defaultQueryTimeout = time.Second * 2

// queryTerminalFilter is a function that filters input events using a type
// switch. If false is returned, the QueryTerminal function will stop reading
// input.
type queryTerminalFilter func(seq string, pa *ansi.Parser) bool

// queryTerminal queries the terminal for support of various features and
// returns a list of response events.
// Most of the time, you will need to set stdin to raw mode before calling this
// function.
// Note: This function will block until the terminal responds or the timeout
// is reached.
// copied from x/term@v0.1.3.
func queryTerminal(
	in io.Reader,
	out io.Writer,
	timeout time.Duration,
	filter queryTerminalFilter,
	query string,
) error {
	rd, err := cancelreader.NewReader(in)
	if err != nil {
		return fmt.Errorf("could not create cancel reader: %w", err)
	}

	defer rd.Close() //nolint: errcheck

	done := make(chan struct{}, 1)
	defer close(done)
	go func() {
		select {
		case <-done:
		case <-time.After(timeout):
			rd.Cancel()
		}
	}()

	if _, err := io.WriteString(out, query); err != nil {
		return fmt.Errorf("could not write query: %w", err)
	}

	pa := ansi.GetParser()
	defer ansi.PutParser(pa)

	var buf [256]byte // 256 bytes should be enough for most responses
	for {
		n, err := rd.Read(buf[:])
		if err != nil {
			return fmt.Errorf("could not read from input: %w", err)
		}

		var state byte
		p := buf[:]
		for n > 0 {
			seq, _, read, newState := ansi.DecodeSequence(p[:n], state, pa)
			if !filter(string(seq), pa) {
				return nil
			}

			state = newState
			n -= read
			p = p[read:]
		}
	}
}

func shift(x uint64) uint64 {
	if x > 0xff {
		x >>= 8
	}
	return x
}

func xParseColor(s string) color.Color {
	switch {
	case strings.HasPrefix(s, "rgb:"):
		parts := strings.Split(s[4:], "/")
		if len(parts) != 3 {
			return color.Black
		}

		r, _ := strconv.ParseUint(parts[0], 16, 32)
		g, _ := strconv.ParseUint(parts[1], 16, 32)
		b, _ := strconv.ParseUint(parts[2], 16, 32)

		return color.RGBA{uint8(shift(r)), uint8(shift(g)), uint8(shift(b)), 255} //nolint:gosec
	case strings.HasPrefix(s, "rgba:"):
		parts := strings.Split(s[5:], "/")
		if len(parts) != 4 {
			return color.Black
		}

		r, _ := strconv.ParseUint(parts[0], 16, 32)
		g, _ := strconv.ParseUint(parts[1], 16, 32)
		b, _ := strconv.ParseUint(parts[2], 16, 32)
		a, _ := strconv.ParseUint(parts[3], 16, 32)

		return color.RGBA{uint8(shift(r)), uint8(shift(g)), uint8(shift(b)), uint8(shift(a))} //nolint:gosec
	}
	return nil
}
