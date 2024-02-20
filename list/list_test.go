package list

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/internal/golden"
	"github.com/charmbracelet/lipgloss/tree"
)

func TestList(t *testing.T) {
	l := New().
		Item("Foo").
		Item("Bar").
		Item("Baz")

	golden.RequireEqual(t, []byte(l.String()))
}

func TestSublist(t *testing.T) {
	l := New().
		Item("Foo").
		Item("Bar").
		Item(New("Hi", "Hello")).Enumerator(Roman).
		Item("Qux")
	golden.RequireEqual(t, []byte(l.String()))
}

func TestComplexSublist(t *testing.T) {
	style1 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		MarginRight(1)
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		MarginRight(1)

	l := New().
		Item("Foo").
		Item("Bar").
		Item(New("foo2", "bar2")).
		Item("Qux").
		Item(
			New("Qux", "aaa", "bbb").
				EnumeratorStyle(style1).
				Enumerator(Roman),
		).
		Item("Deep").
		Item(
			New().
				EnumeratorStyle(style2).
				Enumerator(Alphabet).
				Item("foo").
				Item("Deeper").
				Item(
					New().
						EnumeratorStyle(style1).
						Enumerator(Arabic).
						Item("a").
						Item("b").
						Item("Even Deeper, inherit parent renderer").
						Item(
							New().
								Enumerator(Asterisk).
								EnumeratorStyle(style2).
								Item("sus").
								Item("d minor").
								Item("f#").
								Item("One ore level, with another renderer").
								Item(
									New().
										EnumeratorStyle(style1).
										Enumerator(Dash).
										Item("a\nmultine\nstring").
										Item("hoccus poccus").
										Item("abra kadabra").
										Item("And finally, a tree within all this").
										Item(

											tree.New("").
												EnumeratorStyle(style2).
												Item("another\nmultine\nstring").
												Item("something").
												Item("a subtree").
												Item(
													tree.New("").
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

	golden.RequireEqual(t, []byte(l.String()))
}

func TestMultiline(t *testing.T) {
	l := New().
		Item("Item1\nline 2\nline 3").
		Item("Item2\nline 2\nline 3").
		Item("3")

	golden.RequireEqual(t, []byte(l.String()))
}

func TestListIntegers(t *testing.T) {
	l := New().
		Item("1").
		Item("2").
		Item("3")

	golden.RequireEqual(t, []byte(l.String()))
}

func TestEnumerators(t *testing.T) {
	tests := map[string]tree.Enumerator{
		"alphabet": Alphabet,
		"arabic":   Arabic,
		"roman":    Roman,
		"bullet":   Bullet,
		"asterisk": Asterisk,
		"dash":     Dash,
	}

	for name, enum := range tests {
		t.Run(name, func(t *testing.T) {
			l := New().
				Enumerator(enum).
				Item("Foo").
				Item("Bar").
				Item("Baz")

			golden.RequireEqual(t, []byte(l.String()))
		})
	}
}

func TestEnumeratorsTransform(t *testing.T) {
	tests := map[string]struct {
		enumeration tree.Enumerator
		style       lipgloss.Style
	}{
		"alphabet lower": {
			enumeration: Alphabet,
			style:       lipgloss.NewStyle().MarginRight(1).Transform(strings.ToLower),
		},
		"arabic)": {
			enumeration: Arabic,
			style: lipgloss.NewStyle().MarginRight(1).Transform(func(s string) string {
				return strings.Replace(s, ".", ")", 1)
			}),
		},
		"roman within ()": {
			enumeration: Roman,
			style: lipgloss.NewStyle().Transform(func(s string) string {
				return "(" + strings.Replace(strings.ToLower(s), ".", "", 1) + ") "
			}),
		},
		"bullet is dash": {
			enumeration: Bullet,
			style: lipgloss.NewStyle().Transform(func(s string) string {
				return "- " // this is better done by replacing the enumerator.
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			l := New().
				EnumeratorStyle(test.style).
				Enumerator(test.enumeration).
				Item("Foo").
				Item("Bar").
				Item("Baz")

			golden.RequireEqual(t, []byte(l.String()))
		})
	}
}

func TestBullet(t *testing.T) {
	tests := []struct {
		enum tree.Enumerator
		i    int
		exp  string
	}{
		{Alphabet, 0, "A"},
		{Alphabet, 25, "Z"},
		{Alphabet, 26, "AA"},
		{Alphabet, 51, "AZ"},
		{Alphabet, 52, "BA"},
		{Alphabet, 79, "CB"},
		{Alphabet, 701, "ZZ"},
		{Alphabet, 702, "AAA"},
		{Alphabet, 801, "ADV"},
		{Alphabet, 1000, "ALM"},
		{Roman, 0, "I"},
		{Roman, 25, "XXVI"},
		{Roman, 26, "XXVII"},
		{Roman, 50, "LI"},
		{Roman, 100, "CI"},
		{Roman, 701, "DCCII"},
		{Roman, 1000, "MI"},
	}

	for _, test := range tests {
		_, prefix := test.enum(nil, test.i, false)
		bullet := strings.TrimSuffix(prefix, ".")
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
       C. Foo`, "\n")

	if l.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s\n", expected, l.String())
	}
}
