package adaptive

import (
	"fmt"
	"image/color"
	"io"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/input"
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
	// nolint: errcheck
	err = queryTerminal(in, out, defaultQueryTimeout,
		func(events []input.Event) bool {
			for _, e := range events {
				switch e := e.(type) {
				case input.BackgroundColorEvent:
					c = e.Color
					continue // we need to consume the next DA1 event
				case input.PrimaryDeviceAttributesEvent:
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
type queryTerminalFilter func(events []input.Event) bool

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
	rd, err := input.NewDriver(in, "", 0)
	if err != nil {
		return fmt.Errorf("could not create driver: %w", err)
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

	for {
		events, err := rd.ReadEvents()
		if err != nil {
			return err
		}

		if !filter(events) {
			break
		}
	}

	return nil
}
