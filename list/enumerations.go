package list

import (
	"fmt"
	"strings"
)

// Enumerator is the type of enumeration to use for the list styling.
// It is the prefix for the list.
type Enumerator func(l *List, i int) string

const abcLen = 26

// Alphabet is the enumeration for alphabetical listing.
//
// a. Foo
// b. Bar
// c. Baz
// d. Qux.
var Alphabet Enumerator = func(_ *List, i int) string {
	if i > abcLen*abcLen+abcLen {
		return fmt.Sprintf("%c%c%c.", 'A'+(i-1)/abcLen/abcLen-1, 'A'+((i-1)/abcLen)%abcLen-1, 'A'+(i-1)%abcLen)
	}
	if i > abcLen {
		return fmt.Sprintf("%c%c.", 'A'+(i-1)/abcLen-1, 'A'+(i-1)%abcLen)
	}
	return fmt.Sprintf("%c.", 'A'+(i-1)%abcLen)
}

// Arabic is the enumeration for arabic numerals listing.
//
// 1. Foo
// 2. Bar
// 3. Baz
// 4. Qux.
var Arabic Enumerator = func(_ *List, i int) string {
	return fmt.Sprintf("%d.", i)
}

// Roman is the enumeration for roman numerals listing.
//
// i. Foo
// ii. Bar
// iii. Baz
// iv. Qux.
var Roman Enumerator = func(_ *List, i int) string {
	var (
		roman  = []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
		arabic = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
		result strings.Builder
	)
	for v, value := range arabic {
		for i >= value {
			i -= value
			result.WriteString(roman[v])
		}
	}
	result.WriteRune('.')
	return result.String()
}

// Bullet is the enumeration for bullet listing.
//
// • Foo
// • Bar
// • Baz
// • Qux.
var Bullet Enumerator = func(_ *List, _ int) string {
	return "•"
}

// Tree is the enumeration for the tree listing.
//
// ├─ Foo
// ├─ Bar
// ├─ Baz
// └─ Qux.
var Tree Enumerator = func(l *List, index int) string {
	// out of bounds?
	if index < 0 || index > len(l.items) {
		return ""
	}

	switch index {
	// is last item of list.
	case len(l.items):
		return "└─"
	default:
		switch l.items[index].(type) {
		case *List:
			return "└─"
		default:
			return "├─"
		}
	}
}
