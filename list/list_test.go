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

func TestSublists(t *testing.T) {
	l := New().
		Item("Foo").
		Item(NewSublist("Bar", "foo2", "bar2")).
		Item(
			NewSublist("Qux", "aaa", "bbb").
				Renderer(tree.NewDefaultRenderer().Enumerator(Roman)),
		).
		Item(
			NewSublist("Deep").
				Renderer(tree.NewDefaultRenderer().Enumerator(Alphabet)).
				Item("foo").
				Item(
					NewSublist("Deeper").
						Renderer(tree.NewDefaultRenderer().Enumerator(Arabic)).
						Item("a").
						Item("b").
						Item(
							NewSublist("Even Deeper, inherit parent renderer").
								Item("sus").
								Item("d minor").
								Item("f#").
								Item(

									NewSublist("One ore level, with another renderer").
										Renderer(tree.NewDefaultRenderer().Enumerator(Bullet)).
										Item("a\nmultine\nstring").
										Item("hoccus poccus").
										Item("abra kadabra").
										Item(

											NewSublist("And finally, a tree within all this").
												Renderer(tree.NewDefaultRenderer()).
												Item("another\nmultine\nstring").
												Item("something").
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
	}

	for name, enum := range tests {
		t.Run(name, func(t *testing.T) {
			l := New().
				Renderer(NewDefaultRenderer().Enumerator(enum)).
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
				Renderer(
					NewDefaultRenderer().
						EnumeratorStyle(test.style).
						Enumerator(test.enumeration),
				).
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
	l := New().Renderer(NewDefaultRenderer().Enumerator(Roman))
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
