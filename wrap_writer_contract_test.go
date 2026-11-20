package lipgloss_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
)

var errWrapSink = errors.New("sink failed")

type wrapSink struct {
	bytes.Buffer
	remaining int
	err       error
	calls     int
}

func (s *wrapSink) Write(p []byte) (int, error) {
	s.calls++
	if len(p) <= s.remaining {
		s.remaining -= len(p)
		return s.Buffer.Write(p)
	}
	n, _ := s.Buffer.Write(p[:s.remaining])
	s.remaining = 0
	return n, s.err
}

func TestWrapWriterClosedPipe(t *testing.T) {
	reader, pipe := io.Pipe()
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	defer pipe.Close()
	w := lipgloss.NewWrapWriter(pipe)
	defer w.Close()
	n, err := io.Copy(w, strings.NewReader("customer records"))
	if n != 0 || !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("io.Copy = (%d, %v), want (0, ErrClosedPipe)", n, err)
	}
}

// Each segment records whether its bytes came from the caller. The sink can
// fail at every output byte, including bytes inserted around a newline.
func TestWrapWriterFailureCounts(t *testing.T) {
	const style = "\x1b[31m"
	const link = "\x1b]8;;https://example.com\a"
	const resetStyle = "\x1b[m"
	const resetLink = "\x1b]8;;\a"
	type segment struct {
		text  string
		input bool
	}
	cases := []struct {
		name     string
		segments []segment
	}{
		{"plain", []segment{{"ab", true}, {"\n", true}, {"cd", true}}},
		{"style", []segment{{style + "ab", true}, {resetStyle, false}, {"\n", true}, {style, false}, {"cd", true}}},
		{"link", []segment{{link + "ab", true}, {resetLink, false}, {"\n", true}, {link, false}, {"cd", true}}},
		{"both", []segment{{style + link + "ab", true}, {resetStyle, false}, {resetLink, false}, {"\n", true}, {link, false}, {style, false}, {"cd", true}}},
	}
	for _, tc := range cases {
		var input, output string
		counts := []int{0}
		accepted := 0
		for _, part := range tc.segments {
			output += part.text
			if part.input {
				input += part.text
			}
			for range len(part.text) {
				if part.input {
					accepted++
				}
				counts = append(counts, accepted)
			}
		}
		for _, sinkErr := range []error{errWrapSink, nil} {
			for cutoff := 0; cutoff < len(output); cutoff++ {
				t.Run(fmt.Sprintf("%s/error=%v/cutoff=%d", tc.name, sinkErr, cutoff), func(t *testing.T) {
					sink := &wrapSink{remaining: cutoff, err: sinkErr}
					w := lipgloss.NewWrapWriter(sink)
					n, err := w.Write([]byte(input))
					wantErr := sinkErr
					if wantErr == nil {
						wantErr = io.ErrShortWrite
					}
					if n != counts[cutoff] || !errors.Is(err, wantErr) {
						t.Fatalf("Write = (%d, %v), want (%d, %v)", n, err, counts[cutoff], wantErr)
					}
					if sink.String() != output[:cutoff] {
						t.Fatalf("output = %q, want %q", sink.String(), output[:cutoff])
					}
					calls := sink.calls
					if n, err = w.Write([]byte("later")); n != 0 || !errors.Is(err, wantErr) {
						t.Fatalf("write after failure = (%d, %v)", n, err)
					}
					if err = w.Close(); !errors.Is(err, wantErr) {
						t.Fatalf("Close after failure = %v", err)
					}
					if sink.calls != calls {
						t.Fatalf("forwarded after failure: %d calls became %d", calls, sink.calls)
					}
					if err = w.Close(); err != nil {
						t.Fatalf("second Close = %v", err)
					}
					if n, err = w.Write([]byte("closed")); n != 6 || err != nil {
						t.Fatalf("write after close = (%d, %v)", n, err)
					}
				})
			}
		}
	}
}

func TestWrapWriterCloseFailures(t *testing.T) {
	const input = "\x1b[31m\x1b]8;;https://example.com\x1b\\ok"
	const resets = "\x1b[m\x1b]8;;\a"
	for _, sinkErr := range []error{errWrapSink, nil} {
		for cutoff := 0; cutoff < len(resets); cutoff++ {
			t.Run(fmt.Sprintf("error=%v/cutoff=%d", sinkErr, cutoff), func(t *testing.T) {
				sink := &wrapSink{remaining: len(input) + cutoff, err: sinkErr}
				w := lipgloss.NewWrapWriter(sink)
				if n, err := w.Write([]byte(input)); n != len(input) || err != nil {
					t.Fatalf("setup Write = (%d, %v)", n, err)
				}
				wantErr := sinkErr
				if wantErr == nil {
					wantErr = io.ErrShortWrite
				}
				if err := w.Close(); !errors.Is(err, wantErr) {
					t.Fatalf("Close = %v, want %v", err, wantErr)
				}
				if sink.String() != input+resets[:cutoff] {
					t.Fatalf("output = %q", sink.String())
				}
				calls := sink.calls
				if err := w.Close(); err != nil {
					t.Fatalf("second Close = %v", err)
				}
				if n, err := w.Write([]byte("closed")); n != 6 || err != nil {
					t.Fatalf("write after close = (%d, %v)", n, err)
				}
				if sink.calls != calls {
					t.Fatal("closed writer forwarded output")
				}
			})
		}
	}
}

func TestWrapWriterHealthyStreaming(t *testing.T) {
	const style = "\x1b[31m"
	const link = "\x1b]8;;https://example.com\a"
	const resetStyle = "\x1b[m"
	const resetLink = "\x1b]8;;\a"
	input := style + "\x1b]8;;https://example.com\x1b\\" + "café\n猫\n"
	want := style + "\x1b]8;;https://example.com\x1b\\" + "café" + resetStyle + resetLink + "\n" + link + style + "猫" + resetStyle + resetLink + "\n" + link + style + resetStyle + resetLink
	for _, chunkSize := range []int{1, 2, len(input)} {
		var dst bytes.Buffer
		w := lipgloss.NewWrapWriter(&dst)
		for start := 0; start < len(input); start += chunkSize {
			end := min(start+chunkSize, len(input))
			if n, err := w.Write([]byte(input[start:end])); n != end-start || err != nil {
				t.Fatalf("Write = (%d, %v)", n, err)
			}
		}
		style, link := w.Style(), w.Link()
		if style.IsZero() || link.IsZero() {
			t.Fatal("style or link lost")
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		if dst.String() != want {
			t.Fatalf("chunk %d output = %q, want %q", chunkSize, dst.String(), want)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		if dst.String() != want {
			t.Fatal("second Close emitted duplicate resets")
		}
	}
}
