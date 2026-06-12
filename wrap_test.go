package lipgloss

import (
	"bytes"
	"testing"
)

// TestWrapWriterWriteAfterClose verifies that writing to a WrapWriter after it
// has been closed is a safe no-op rather than a nil-pointer panic. Close()
// returns the internal ANSI parser to a sync.Pool and nils it out; out-of-order
// teardown of nested writer chains can route a trailing style/link reset back
// through an already-closed writer, so Write must tolerate a nil parser.
func TestWrapWriterWriteAfterClose(t *testing.T) {
	var buf bytes.Buffer
	w := NewWrapWriter(&buf)

	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	n, err := w.Write([]byte("after close"))
	if err != nil {
		t.Fatalf("write after close: %v", err)
	}
	if n != len("after close") {
		t.Fatalf("write after close: got n=%d, want %d", n, len("after close"))
	}
}
