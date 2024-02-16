package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/tree"
)

const abcLen = 26

const indent = " "

// Alphabet is the enumeration for alphabetical listing.
//
// a. Foo
// b. Bar
// c. Baz
// d. Qux.
func Alphabet(_ tree.Atter, i int, _ bool) (string, string) {
	if i >= abcLen*abcLen+abcLen {
		return indent, fmt.Sprintf("%c%c%c.", 'A'+i/abcLen/abcLen-1, 'A'+(i/abcLen)%abcLen-1, 'A'+i%abcLen)
	}
	if i >= abcLen {
		return indent, fmt.Sprintf("%c%c.", 'A'+i/abcLen-1, 'A'+(i)%abcLen)
	}
	return indent, fmt.Sprintf("%c.", 'A'+i%abcLen)
}

// Arabic is the enumeration for arabic numerals listing.
//
// 1. Foo
// 2. Bar
// 3. Baz
// 4. Qux.
func Arabic(_ tree.Atter, i int, _ bool) (string, string) {
	return indent, fmt.Sprintf("%d.", i+1)
}

// Roman is the enumeration for roman numerals listing.
//
// /   I. Foo
// /  II. Bar
// / III. Baz
// /  IV. Qux.
func Roman(_ tree.Atter, i int, _ bool) (string, string) {
	var (
		roman  = []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
		arabic = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
		result strings.Builder
	)
	for v, value := range arabic {
		for i >= value-1 {
			i -= value
			result.WriteString(roman[v])
		}
	}
	result.WriteRune('.')
	return indent, result.String()
}

// Bullet is the enumeration for bullet listing.
//
// • Foo
// • Bar
// • Baz
// • Qux.
func Bullet(tree.Atter, int, bool) (string, string) {
	return indent, "•"
}

// Asterisk is an enumeration using asterisks.
//
// * Foo
// * Bar
// * Baz
// * Qux.
func Asterisk(tree.Atter, int, bool) (string, string) {
	return indent, "*"
}

// Dash is an enumeration using dashes.
//
// * Foo
// * Bar
// * Baz
// * Qux.
func Dash(tree.Atter, int, bool) (string, string) {
	return indent, "-"
}
