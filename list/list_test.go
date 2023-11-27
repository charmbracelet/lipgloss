package list

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestList(t *testing.T) {
	l := New().
		Item("Foo").
		Item("Bar").
		Item("Baz")

	expected := strings.TrimPrefix(`
• Foo
• Bar
• Baz
`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}

func TestListIndent(t *testing.T) {
	l := New().
		Item(New().Indent(4).Item("Foo")). // custom nesting
		Item(New("Bar")).                  // should be nested
		Item("Baz")                        // normal indent

	expected := strings.TrimPrefix(`
    • Foo
  • Bar
• Baz
`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}

func TestHide(t *testing.T) {
	l := New().
		Item("Foo").
		Item(New("Bar").Hide(true)).
		Item("Baz")

	expected := strings.TrimPrefix(`
• Foo
• Baz
`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}

func TestListIntegers(t *testing.T) {
	l := New().
		Item(1).
		Item(2).
		Item(3)

	expected := strings.TrimPrefix(`
• 1
• 2
• 3
`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}

func TestNestedList(t *testing.T) {
	l := New().
		Item("Foo").
		Item("Bar").
		Item(New().
			Item("Qux").
			Item("Quux")).
		Item("Baz")

	expected := strings.TrimPrefix(`
• Foo
• Bar
  • Qux
  • Quux
• Baz
`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}

func TestDeepNestedList(t *testing.T) {
	l := New(
		"Foo",
		"Bar",
		New(
			"Baz",
			New(
				"Qux",
				"Quux",
				"Corge",
			),
			"Grault",
			New(
				"i",
				"ii",
				New(
					"a",
					"b",
				),
				"iii",
			),
		),
		"Garply",
		New(
			"i",
			"ii",
			"iii",
		),
	)

	expected := strings.TrimPrefix(`
• Foo
• Bar
  • Baz
    • Qux
    • Quux
    • Corge
  • Grault
    • i
    • ii
      • a
      • b
    • iii
• Garply
  • i
  • ii
  • iii
`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}

func TestEnumerators(t *testing.T) {
	tests := []struct {
		enumeration Enumerator
		expected    string
	}{
		{
			enumeration: Alphabet,
			expected: `
A. Foo
B. Bar
C. Baz
  A. Qux
  B. Quux
`,
		},
		{
			enumeration: Arabic,
			expected: `
1. Foo
2. Bar
3. Baz
  1. Qux
  2. Quux
`,
		},
		{
			enumeration: Roman,
			expected: `
  I. Foo
 II. Bar
III. Baz
   I. Qux
  II. Quux
`,
		},
		{
			enumeration: Bullet,
			expected: `
• Foo
• Bar
• Baz
  • Qux
  • Quux
`,
		},
		{
			enumeration: Tree,
			expected: `
├─ Foo
├─ Bar
└─ Baz
  ├─ Qux
  └─ Quux
`,
		},
	}

	for _, test := range tests {
		expected := strings.TrimPrefix(test.expected, "\n")

		l := New().
			Enumerator(test.enumeration).
			Item("Foo").
			Item("Bar").
			Item("Baz").
			Item(New("Qux", "Quux").Enumerator(test.enumeration))

		if l.String() != expected {
			t.Errorf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
		}
	}
}

func TestEnumeratorsTransform(t *testing.T) {
	tests := []struct {
		enumeration Enumerator
		style       lipgloss.Style
		expected    string
	}{
		{
			enumeration: Alphabet,
			style:       lipgloss.NewStyle().MarginRight(1).Transform(strings.ToLower),
			expected: `
a. Foo
b. Bar
c. Baz
  a. Qux
  b. Quux
`,
		},
		{
			enumeration: Arabic,
			style: lipgloss.NewStyle().MarginRight(1).Transform(func(s string) string {
				return strings.Replace(s, ".", ")", 1)
			}),
			expected: `
1) Foo
2) Bar
3) Baz
  1) Qux
  2) Quux
`,
		},
		{
			enumeration: Roman,
			style: lipgloss.NewStyle().Transform(func(s string) string {
				return "(" + strings.Replace(strings.ToLower(s), ".", "", 1) + ") "
			}),
			expected: `
  (i) Foo
 (ii) Bar
(iii) Baz
   (i) Qux
  (ii) Quux
`,
		},
		{
			enumeration: Bullet,
			style: lipgloss.NewStyle().Transform(func(s string) string {
				return "- " // this is better done by replacing the enumerator.
			}),
			expected: `
- Foo
- Bar
- Baz
  - Qux
  - Quux
`,
		},
		{
			enumeration: Tree,
			style: lipgloss.NewStyle().MarginRight(1).Transform(func(s string) string {
				return strings.Replace(s, "─", "───", 1)
			}),
			expected: `
├─── Foo
├─── Bar
└─── Baz
  ├─── Qux
  └─── Quux
`,
		},
	}

	for _, test := range tests {
		expected := strings.TrimPrefix(test.expected, "\n")

		l := New().
			Enumerator(test.enumeration).
			EnumeratorStyle(test.style).
			Item("Foo").
			Item("Bar").
			Item("Baz").
			Item(New("Qux", "Quux").Enumerator(test.enumeration).EnumeratorStyle(test.style))

		if l.String() != expected {
			t.Errorf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
		}
	}
}

func TestBullet(t *testing.T) {
	tests := []struct {
		enum Enumerator
		i    int
		exp  string
	}{
		{Alphabet, 1, "A"},
		{Alphabet, 26, "Z"},
		{Alphabet, 27, "AA"},
		{Alphabet, 52, "AZ"},
		{Alphabet, 53, "BA"},
		{Alphabet, 80, "CB"},
		{Alphabet, 702, "ZZ"},
		{Alphabet, 703, "AAA"},
		{Alphabet, 802, "ADV"},
		{Alphabet, 1001, "ALM"},
		{Roman, 1, "I"},
		{Roman, 26, "XXVI"},
		{Roman, 27, "XXVII"},
		{Roman, 51, "LI"},
		{Roman, 101, "CI"},
		{Roman, 702, "DCCII"},
		{Roman, 1001, "MI"},
	}

	for _, test := range tests {
		bullet := strings.TrimSuffix(test.enum(nil, test.i), ".")
		if bullet != test.exp {
			t.Errorf("expected: %s, got: %s\n", test.exp, bullet)
		}
	}
}

func TestEnumeratorsAlign(t *testing.T) {
	fooList := strings.Split(strings.TrimSuffix(strings.Repeat("Foo ", 100), " "), " ")
	l := New().Enumerator(Roman)
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
       C. Foo
`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}
