package golden

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/aymanbagabas/go-udiff"
)

var update = flag.Bool("update", false, "update .golden files")

// RequireEqual is a helper function to assert the given output is
// the expected from the golden files, printing its diff in case it is not.
//
// You can update the golden files by running your tests with the -update flag.
func RequireEqual(tb testing.TB, out []byte) {
	tb.Helper()

	golden := filepath.Join("testdata", tb.Name()+".golden")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil { //nolint: gomnd
			tb.Fatal(err)
		}
		if err := os.WriteFile(golden, out, 0o600); err != nil { //nolint: gomnd
			tb.Fatal(err)
		}
	}

	path := filepath.Join(tb.TempDir(), tb.Name()+".out")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { //nolint: gomnd
		tb.Fatal(err)
	}
	if err := os.WriteFile(path, out, 0o600); err != nil { //nolint: gomnd
		tb.Fatal(err)
	}

	goldenBts, err := os.ReadFile(golden)
	if err != nil {
		tb.Fatal(err)
	}
	diff := udiff.Unified("golden", "run", string(goldenBts), string(out))
	if diff != "" {
		tb.Fatalf("output does not match, diff:\n\n%s", diff)
	}
	tb.Logf("got:\n%s", string(goldenBts))
}
