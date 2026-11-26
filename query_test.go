package lipgloss

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
)

func TestBackgroundColorRejectsNonTerminal(t *testing.T) {
	in, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	out, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	if _, err := BackgroundColor(in, out); err == nil {
		t.Fatal("BackgroundColor succeeded with non-terminal files")
	}
}

func TestQueryTerminalTimeout(t *testing.T) {
	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()

	start := time.Now()
	err := queryTerminal(pr, io.Discard, 50*time.Millisecond, func(string, *ansi.Parser) bool {
		return true
	}, ansi.RequestBackgroundColor)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("queryTerminal took too long to time out: %v", elapsed)
	}
}

func TestQueryBackgroundColorResponse(t *testing.T) {
	in := strings.NewReader("\x1b]11;rgb:1111/2222/3333\x1b\\\x1b[?62;c")
	var out bytes.Buffer

	bg, err := queryBackgroundColor(in, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bg == nil {
		t.Fatal("expected non-nil background color")
	}
	if out.Len() == 0 {
		t.Fatal("expected query to be written to output")
	}
}
