package lipgloss

import (
	"os"
	"testing"
	"time"
)

// Pipes are not terminals, so there is nothing to answer a background color
// query. BackgroundColor must say so and return, not wait.
//
// On Windows it used to open CONIN$/CONOUT$ and query the console instead, which
// left any program with redirected stdio -- a service, a CI job -- waiting on a
// console that never replies.
func TestBackgroundColorDoesNotQueryNonTerminals(t *testing.T) {
	r, w := pipe(t)

	done := make(chan error, 1)
	go func() {
		_, err := BackgroundColor(r, w)
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Error("want an error for handles that are not terminals, got nil")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("BackgroundColor blocked on handles that are not terminals")
	}
}

// The documented fallback: assume a dark background when detection is not possible.
func TestHasDarkBackgroundDefaultsToDarkOnNonTerminals(t *testing.T) {
	r, w := pipe(t)

	done := make(chan bool, 1)
	go func() {
		done <- HasDarkBackground(r, w)
	}()

	select {
	case isDark := <-done:
		if !isDark {
			t.Error("want the dark default when detection is not possible, got false")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("HasDarkBackground blocked on handles that are not terminals")
	}
}

func pipe(t *testing.T) (r, w *os.File) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	t.Cleanup(func() {
		_ = r.Close()
		_ = w.Close()
	})

	return r, w
}
