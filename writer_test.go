package lipgloss

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/charmbracelet/colorprofile"
)

func TestMain(m *testing.M) {
	switch os.Getenv("LIPGLOSS_WRITER_PROBE") {
	case "render":
		NewStyle().Bold(true).Render("hi")
		os.Exit(0)
	case "print":
		_, _ = Print("hi")
		os.Exit(0)
	case "assign":
		Writer = &colorprofile.Writer{Forward: io.Discard, Profile: colorprofile.TrueColor}
		_, _ = Print("hi")
		os.Exit(0)
	default:
		os.Exit(m.Run())
	}
}

func TestRenderDoesNotQueryTmux(t *testing.T) {
	t.Parallel()
	if writerProbeCalledTmux(t, "render") {
		t.Fatal("rendering a string queried tmux")
	}
}

func TestPrintQueriesTmux(t *testing.T) {
	t.Parallel()
	if !writerProbeCalledTmux(t, "print") {
		t.Fatal("printing did not query tmux")
	}
}

func TestAssignedWriterSkipsTmuxQuery(t *testing.T) {
	t.Parallel()
	if writerProbeCalledTmux(t, "assign") {
		t.Fatal("an assigned Writer still queried tmux")
	}
}

func writerProbeCalledTmux(t *testing.T, probe string) bool {
	t.Helper()

	dir := t.TempDir()
	marker := filepath.Join(dir, "called")
	script := filepath.Join(dir, "tmux")
	body := "#!/bin/sh\nprintf x > " + strconv.Quote(marker) + "\n"
	if err := os.WriteFile(script, []byte(body), 0o700); err != nil {
		t.Fatalf("write tmux: %v", err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = []string{
		"LIPGLOSS_WRITER_PROBE=" + probe,
		"PATH=" + dir,
		"TMUX=" + dir,
		"TERM=tmux-256color",
		"TTY_FORCE=1",
		"COLORTERM=truecolor",
		"TMPDIR=" + os.TempDir(),
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("probe %s: %v\n%s", probe, err, out)
	}

	_, statErr := os.Stat(marker)
	return statErr == nil
}
