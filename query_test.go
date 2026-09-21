package lipgloss

import (
	"os"
	"testing"
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
