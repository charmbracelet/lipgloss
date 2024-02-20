package require

import (
	"strings"
	"testing"
	"unicode"

	"github.com/aymanbagabas/go-udiff"
)

// Equal verifies the strings are equal, assuming its terminal output.
func Equal(tb testing.TB, expected, got string) {
	tb.Helper()

	cleanExpected := trimSpace(expected)
	cleanGot := trimSpace(got)
	if diff := udiff.Unified("expected", "got", cleanExpected, cleanGot); diff != "" {
		tb.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n\ndiff:\n\n%s\n\n", cleanExpected, cleanGot, diff)
	}
}

func trimSpace(s string) string {
	var result []string
	ss := strings.Split(s, "\n")
	for i, line := range ss {
		// ignore begging and ending empty lines
		if strings.TrimSpace(line) == "" && (i == 0 || i == len(ss)-1) {
			continue
		}
		result = append(result, strings.TrimRightFunc(line, unicode.IsSpace))
	}
	return strings.Join(result, "\n")
}
