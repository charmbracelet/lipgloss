package list

import (
	"strings"
	"testing"
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

func TestEnumeration(t *testing.T) {
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
