package list_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/internal/require"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
)

func TestList(t *testing.T) {
	l := list.New().
		Item("Foo").
		Item("Bar").
		Item("Baz")

	expected := `
• Foo
• Bar
• Baz
	`
	require.Equal(t, expected, l.String())
}

func TestListItems(t *testing.T) {
	l := list.New().
		Items([]string{"Foo", "Bar", "Baz"})

	expected := `
• Foo
• Bar
• Baz
	`
	require.Equal(t, expected, l.String())
}

func TestSublist(t *testing.T) {
	l := list.New().
		Item("Foo").
		Item("Bar").
		Item(list.New("Hi", "Hello")).Enumerator(list.Roman).
		Item("Qux")

	expected := `
  I. Foo
 II. Bar
  • Hi
  • Hello
III. Qux
	`
	require.Equal(t, expected, l.String())
}

func TestComplexSublist(t *testing.T) {
	style1 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		MarginRight(1)
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		MarginRight(1)

	l := list.New().
		Item("Foo").
		Item("Bar").
		Item(list.New("foo2", "bar2")).
		Item("Qux").
		Item(
			list.New("aaa", "bbb").
				EnumeratorStyle(style1).
				Enumerator(list.Roman),
		).
		Item("Deep").
		Item(
			list.New().
				EnumeratorStyle(style2).
				Enumerator(list.Alphabet).
				Item("foo").
				Item("Deeper").
				Item(
					list.New().
						EnumeratorStyle(style1).
						Enumerator(list.Arabic).
						Item("a").
						Item("b").
						Item("Even Deeper, inherit parent renderer").
						Item(
							list.New().
								Enumerator(list.Asterisk).
								EnumeratorStyle(style2).
								Item("sus").
								Item("d minor").
								Item("f#").
								Item("One ore level, with another renderer").
								Item(
									list.New().
										EnumeratorStyle(style1).
										Enumerator(list.Dash).
										Item("a\nmultine\nstring").
										Item("hoccus poccus").
										Item("abra kadabra").
										Item("And finally, a tree within all this").
										Item(

											tree.New().
												EnumeratorStyle(style2).
												Item("another\nmultine\nstring").
												Item("something").
												Item("a subtree").
												Item(
													tree.New().
														EnumeratorStyle(style2).
														Item("yup").
														Item("many itens").
														Item("another"),
												).
												Item("hallo").
												Item("wunderbar!"),
										).
										Item("this is a tree\nand other obvious statements"),
								),
						),
				).
				Item("bar"),
		).
		Item("Baz")

	expected := `
• Foo
• Bar
  • foo2
  • bar2
• Qux
   I. aaa
  II. bbb
• Deep
  A. foo
  B. Deeper
    1. a
    2. b
    3. Even Deeper, inherit parent renderer
      * sus
      * d minor
      * f#
      * One ore level, with another renderer
        - a
          multine
          string
        - hoccus poccus
        - abra kadabra
        - And finally, a tree within all this
          ├── another
          │   multine
          │   string
          ├── something
          ├── a subtree
          │   ├── yup
          │   ├── many itens
          │   └── another
          ├── hallo
          └── wunderbar!
        - this is a tree
          and other obvious statements
  C. bar
• Baz
	`
	require.Equal(t, expected, l.String())
}

func TestMultiline(t *testing.T) {
	l := list.New().
		Item("Item1\nline 2\nline 3").
		Item("Item2\nline 2\nline 3").
		Item("3")

	expected := `
• Item1
  line 2
  line 3
• Item2
  line 2
  line 3
• 3
	`
	require.Equal(t, expected, l.String())
}

func TestListIntegers(t *testing.T) {
	l := list.New().
		Item("1").
		Item("2").
		Item("3")

	expected := `
• 1
• 2
• 3
	`
	require.Equal(t, expected, l.String())
}

func TestEnumerators(t *testing.T) {
	tests := map[string]struct {
		enumerator list.Enumerator
		expected   string
	}{
		"alphabet": {
			enumerator: list.Alphabet,
			expected: `
A. Foo
B. Bar
C. Baz
			`,
		},
		"arabic": {
			enumerator: list.Arabic,
			expected: `
1. Foo
2. Bar
3. Baz
			`,
		},
		"roman": {
			enumerator: list.Roman,
			expected: `
  I. Foo
 II. Bar
III. Baz
			`,
		},
		"bullet": {
			enumerator: list.Bullet,
			expected: `
• Foo
• Bar
• Baz
			`,
		},
		"asterisk": {
			enumerator: list.Asterisk,
			expected: `
* Foo
* Bar
* Baz
			`,
		},
		"dash": {
			enumerator: list.Dash,
			expected: `
- Foo
- Bar
- Baz
			`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			l := list.New().
				Enumerator(test.enumerator).
				Item("Foo").
				Item("Bar").
				Item("Baz")

			require.Equal(t, test.expected, l.String())
		})
	}
}

func TestEnumeratorsTransform(t *testing.T) {
	tests := map[string]struct {
		enumeration list.Enumerator
		style       lipgloss.Style
		expected    string
	}{
		"alphabet lower": {
			enumeration: list.Alphabet,
			style:       lipgloss.NewStyle().MarginRight(1).Transform(strings.ToLower),
			expected: `
a. Foo
b. Bar
c. Baz
			`,
		},
		"arabic)": {
			enumeration: list.Arabic,
			style: lipgloss.NewStyle().MarginRight(1).Transform(func(s string) string {
				return strings.Replace(s, ".", ")", 1)
			}),
			expected: `
1) Foo
2) Bar
3) Baz
			`,
		},
		"roman within ()": {
			enumeration: list.Roman,
			style: lipgloss.NewStyle().Transform(func(s string) string {
				return "(" + strings.Replace(strings.ToLower(s), ".", "", 1) + ") "
			}),
			expected: `
  (i) Foo
 (ii) Bar
(iii) Baz
			`,
		},
		"bullet is dash": {
			enumeration: list.Bullet,
			style: lipgloss.NewStyle().Transform(func(s string) string {
				return "- " // this is better done by replacing the enumerator.
			}),
			expected: `
- Foo
- Bar
- Baz
			`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			l := list.New().
				EnumeratorStyle(test.style).
				Enumerator(test.enumeration).
				Item("Foo").
				Item("Bar").
				Item("Baz")

			require.Equal(t, test.expected, l.String())
		})
	}
}

func TestBullet(t *testing.T) {
	tests := []struct {
		enum list.Enumerator
		i    int
		exp  string
	}{
		{list.Alphabet, 0, "A"},
		{list.Alphabet, 25, "Z"},
		{list.Alphabet, 26, "AA"},
		{list.Alphabet, 51, "AZ"},
		{list.Alphabet, 52, "BA"},
		{list.Alphabet, 79, "CB"},
		{list.Alphabet, 701, "ZZ"},
		{list.Alphabet, 702, "AAA"},
		{list.Alphabet, 801, "ADV"},
		{list.Alphabet, 1000, "ALM"},
		{list.Roman, 0, "I"},
		{list.Roman, 25, "XXVI"},
		{list.Roman, 26, "XXVII"},
		{list.Roman, 50, "LI"},
		{list.Roman, 100, "CI"},
		{list.Roman, 701, "DCCII"},
		{list.Roman, 1000, "MI"},
	}

	for _, test := range tests {
		prefix := test.enum(nil, test.i)
		bullet := strings.TrimSuffix(prefix, ".")
		if bullet != test.exp {
			t.Errorf("expected: %s, got: %s\n", test.exp, bullet)
		}
	}
}

func TestEnumeratorsAlign(t *testing.T) {
	fooList := strings.Split(strings.TrimSuffix(strings.Repeat("Foo ", 100), " "), " ")
	l := list.New().Enumerator(list.Roman)
	for _, f := range fooList {
		l.Item(f)
	}

	expected := strings.TrimPrefix(`
       I. Foo
      II. Foo
     III. Foo
      IV. Foo
       V. Foo
      VI. Foo
     VII. Foo
    VIII. Foo
      IX. Foo
       X. Foo
      XI. Foo
     XII. Foo
    XIII. Foo
     XIV. Foo
      XV. Foo
     XVI. Foo
    XVII. Foo
   XVIII. Foo
     XIX. Foo
      XX. Foo
     XXI. Foo
    XXII. Foo
   XXIII. Foo
    XXIV. Foo
     XXV. Foo
    XXVI. Foo
   XXVII. Foo
  XXVIII. Foo
    XXIX. Foo
     XXX. Foo
    XXXI. Foo
   XXXII. Foo
  XXXIII. Foo
   XXXIV. Foo
    XXXV. Foo
   XXXVI. Foo
  XXXVII. Foo
 XXXVIII. Foo
   XXXIX. Foo
      XL. Foo
     XLI. Foo
    XLII. Foo
   XLIII. Foo
    XLIV. Foo
     XLV. Foo
    XLVI. Foo
   XLVII. Foo
  XLVIII. Foo
    XLIX. Foo
       L. Foo
      LI. Foo
     LII. Foo
    LIII. Foo
     LIV. Foo
      LV. Foo
     LVI. Foo
    LVII. Foo
   LVIII. Foo
     LIX. Foo
      LX. Foo
     LXI. Foo
    LXII. Foo
   LXIII. Foo
    LXIV. Foo
     LXV. Foo
    LXVI. Foo
   LXVII. Foo
  LXVIII. Foo
    LXIX. Foo
     LXX. Foo
    LXXI. Foo
   LXXII. Foo
  LXXIII. Foo
   LXXIV. Foo
    LXXV. Foo
   LXXVI. Foo
  LXXVII. Foo
 LXXVIII. Foo
   LXXIX. Foo
    LXXX. Foo
   LXXXI. Foo
  LXXXII. Foo
 LXXXIII. Foo
  LXXXIV. Foo
   LXXXV. Foo
  LXXXVI. Foo
 LXXXVII. Foo
LXXXVIII. Foo
  LXXXIX. Foo
      XC. Foo
     XCI. Foo
    XCII. Foo
   XCIII. Foo
    XCIV. Foo
     XCV. Foo
    XCVI. Foo
   XCVII. Foo
  XCVIII. Foo
    XCIX. Foo
       C. Foo`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}