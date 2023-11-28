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
func Alphabet(_ *List, i int) string {
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
func Arabic(_ *List, i int) string {
	return fmt.Sprintf("%d.", i)
}

// Roman is the enumeration for roman numerals listing.
//
// /   I. Foo
// /  II. Bar
// / III. Baz
// /  IV. Qux.
func Roman(_ *List, i int) string {
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
func Bullet(_ *List, _ int) string {
	return "•"
}

// Tree is the enumeration for the tree listing.
//
// ├─ Foo
// ├─ Bar
// ├─ Baz
// └─ Qux.
func Tree(l *List, index int) string {
	// out of bounds?
	if index < 0 || index > len(l.Items) {
		return ""
	}

	switch index {
	// is last item of list.
	case len(l.Items):
		return "└─"
	default:
		switch l.Items[index].(type) {
		case *List:
			return "└─"
		default:
			return "├─"
		}
	}
}
