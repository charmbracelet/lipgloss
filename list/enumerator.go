package list

import (
	"fmt"
	"strings"
)

// Enumerator enumerates a list. Given a list of items and the index of the
// current enumeration, it returns the prefix that should be displayed for the
// current item.
//
// For example, a simple Arabic numeral enumeration would be:
//
//	func Arabic(_ Items, i int) string {
//	  return fmt.Sprintf("%d.", i+1)
//	}
//
// There are several predefined enumerators:
//   - Alphabet
//   - Arabic
//   - Bullet
//   - Dash
//   - Roman
//
// Or, define your own.
type Enumerator func(items Items, index int) string

// Indenter indents the children of a tree.
//
// Indenters allow for displaying nested tree items with connecting borders
// to sibling nodes.
//
// For example, the default indenter would be:
//
//	func TreeIndenter(children Children, index int) string {
//		if children.Length()-1 == index {
//			return "│  "
//		}
//
//		return "   "
//	}
type Indenter func(items Items, index int) string

// Alphabet is the enumeration for alphabetical listing.
//
//	Example:
//	  a. Foo
//	  b. Bar
//	  c. Baz
//	  d. Qux.
func Alphabet(_ Items, i int) string {
	if i >= abcLen*abcLen+abcLen {
		return fmt.Sprintf("%c%c%c.", 'A'+i/abcLen/abcLen-1, 'A'+(i/abcLen)%abcLen-1, 'A'+i%abcLen)
	}
	if i >= abcLen {
		return fmt.Sprintf("%c%c.", 'A'+i/abcLen-1, 'A'+(i)%abcLen)
	}
	return fmt.Sprintf("%c.", 'A'+i%abcLen)
}

const abcLen = 26

// Arabic is the enumeration for arabic numerals listing.
//
//	Example:
//	  1. Foo
//	  2. Bar
//	  3. Baz
//	  4. Qux.
func Arabic(_ Items, i int) string {
	return fmt.Sprintf("%d.", i+1)
}

// Roman is the enumeration for roman numerals listing.
//
//	Example:
//	  I. Foo
//	 II. Bar
//	III. Baz
//	 IV. Qux.
func Roman(_ Items, i int) string {
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
	return result.String()
}

// Bullet is the enumeration for bullet listing.
//
//	Example:
//	  • Foo
//	  • Bar
//	  • Baz
//	  • Qux.
func Bullet(Items, int) string {
	return "•"
}

// Asterisk is an enumeration using asterisks.
//
//	Example:
//	  * Foo
//	  * Bar
//	  * Baz
//	  * Qux.
func Asterisk(Items, int) string {
	return "*"
}

// Dash is an enumeration using dashes.
//
//	Example:
//	  - Foo
//	  - Bar
//	  - Baz
//	  - Qux.
func Dash(Items, int) string {
	return "-"
}
