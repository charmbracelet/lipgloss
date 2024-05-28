package list

import (
	"fmt"
	"strings"
)

// Enumerator is the type of enumeration to use for the list styling.
// It is the prefix for the list.
type Enumerator func(items Items, index int) string

// Alphabet is the enumeration for alphabetical listing.
//
//	a. Foo
//	b. Bar
//	c. Baz
//	d. Qux.
func Alphabet(_ Items, i int) string {
	if i >= alphabetLen*alphabetLen+alphabetLen {
		return fmt.Sprintf("%c%c%c.", 'A'+i/alphabetLen/alphabetLen-1, 'A'+(i/alphabetLen)%alphabetLen-1, 'A'+i%alphabetLen)
	}
	if i >= alphabetLen {
		return fmt.Sprintf("%c%c.", 'A'+i/alphabetLen-1, 'A'+(i)%alphabetLen)
	}
	return fmt.Sprintf("%c.", 'A'+i%alphabetLen)
}

const alphabetLen = 26

// Arabic is the enumeration for arabic numerals listing.
//
//  1. Foo
//  2. Bar
//  3. Baz
//  4. Qux.
func Arabic(_ Items, i int) string {
	return fmt.Sprintf("%d.", i+1)
}

// Roman is the enumeration for roman numerals listing.
//
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
//	\• Foo
//	\• Bar
//	\• Baz
//	\• Qux.
func Bullet(Items, int) string {
	return "•"
}

// Asterisk is an enumeration using asterisks.
//
//	\* Foo
//	\* Bar
//	\* Baz
//	\* Qux.
func Asterisk(Items, int) string {
	return "*"
}

// Dash is an enumeration using dashes.
//
//	\- Foo
//	\- Bar
//	\- Baz
//	\- Qux.
func Dash(Items, int) string {
	return "-"
}
