package list

import (
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func TestListSimple(t *testing.T) {
	l := New("A", "B", "C")
	got := l.String()
	want := heredoc.Doc(`
• A
• B
• C`)

	if want != got {
		t.Errorf("\nwant:\n\n%s\n\ngot:\n\n%s", want, got)
	}
}

func TestListEnumerator(t *testing.T) {
	l := New("A", "B", "C").Enumerator(Arabic)
	got := l.String()
	want := heredoc.Doc(`
1. A
2. B
3. C`)

	if want != got {
		t.Errorf("\nwant:\n\n%s\n\ngot:\n\n%s", want, got)
	}
}

func TestListNest(t *testing.T) {
	l := New(
		"A",
		"B",
		New(
			"C",
			"D",
			New(
				"E",
				"F",
			),
		),
	)
	got := l.String()
	want := heredoc.Doc(`
• A
• B
  • C
  • D
    • E
    • F`)

	if want != got {
		t.Errorf("\nwant:\n\n%s\n\ngot:\n\n%s", want, got)
	}
}

func TestListInteger(t *testing.T) {
	l := New(
		1,
		2,
		3,
		4,
	)
	got := l.String()
	want := heredoc.Doc(`
• 1
• 2
• 3
• 4`)

	if want != got {
		t.Errorf("\nwant: %s\ngot: %s", want, got)
	}
}

func TestListMultilineItems(t *testing.T) {
	l := New(
		"Bubble Tea\n"+
			"Milky",
		"Milk Tea\n"+
			"Also milky",
		"Actual milk\n"+
			"Milky as well",
	).
		Enumerator(func(i int) string { return "│\n│" }).
		Separator(lipgloss.NewStyle().SetString("\n\n"))

	got := l.String()
	want := heredoc.Doc(`
│ Bubble Tea
│ Milky     

│ Milk Tea  
│ Also milky

│ Actual milk  
│ Milky as well`)

	if want != got {
		t.Errorf("\nwant:\n\n%s\n\ngot:\n\n%s", want, got)
	}
}

func TestEmbeddedTable(t *testing.T) {
	l := New(
		"Foo",
		"Bar",
		table.New().
			Headers("Foo", "Bar", "Baz").
			StyleFunc(func(r, c int) lipgloss.Style {
				return lipgloss.NewStyle().Padding(0, 1)
			}).
			Row("Foo", "Bar", "Baz"),
		"Baz",
	)

	got := l.String()
	want := heredoc.Doc(`
• Foo
• Bar
• ╭─────┬─────┬─────╮
  │ Foo │ Bar │ Baz │
  ├─────┼─────┼─────┤
  │ Foo │ Bar │ Baz │
  ╰─────┴─────┴─────╯
• Baz`)

	if want != got {
		t.Errorf("\nwant:\n\n%s\n\ngot:\n\n%s", want, got)
	}
}

func TestNestedEmbeddedTable(t *testing.T) {
	l := New(
		"Foo",
		"Bar",
		New(
			table.New().
				StyleFunc(func(r, c int) lipgloss.Style {
					return lipgloss.NewStyle().Padding(0, 1)
				}).
				Headers("Foo", "Bar", "Baz").
				Row("Foo", "Bar", "Baz"),
		).Prefix(""),
		"Baz",
	)

	got := l.String()
	want := heredoc.Doc(`
• Foo
• Bar
   ╭─────┬─────┬─────╮
   │ Foo │ Bar │ Baz │
   ├─────┼─────┼─────┤
   │ Foo │ Bar │ Baz │
   ╰─────┴─────┴─────╯
• Baz`)

	if want != got {
		t.Errorf("\nwant:\n\n%s\n\ngot:\n\n%s", want, got)
	}
}

func TestList(t *testing.T) {
	l := New().
		Item("Foo").
		Item("Bar").
		Item("Baz")

	want := heredoc.Doc(`
• Foo
• Bar
• Baz`)

	if l.String() != want {
		t.Fatalf("want:\n\n%s\n\ngot:\n\n%s\n", want, l.String())
	}
}

func TestHide(t *testing.T) {
	l := New().
		Item("Foo").
		Item("Baz")

	want := heredoc.Doc(`
• Foo
• Baz`)

	if l.String() != want {
		t.Fatalf("want:\n\n%s\n\ngot:\n\n%s\n", want, l.String())
	}
}

func TestListIntegers(t *testing.T) {
	l := New().
		Item("1").
		Item("2").
		Item("3")

	want := heredoc.Doc(`
• 1
• 2
• 3`)

	if l.String() != want {
		t.Fatalf("want:\n\n%s\n\ngot:\n\n%s\n", want, l.String())
	}
}

func TestEnumerators(t *testing.T) {
	tests := []struct {
		enumeration Enumerator
		want        string
	}{
		{
			enumeration: Alphabet,
			want: heredoc.Doc(`
A. Foo
B. Bar
C. Baz`),
		},
		{
			enumeration: Arabic,
			want: heredoc.Doc(`
1. Foo
2. Bar
3. Baz`),
		},
		{
			enumeration: Roman,
			want: heredoc.Doc(`
  I. Foo
 II. Bar
III. Baz`),
		},
		{
			enumeration: Bullet,
			want: heredoc.Doc(`
• Foo
• Bar
• Baz`),
		},
	}

	for _, test := range tests {
		l := New().
			Enumerator(test.enumeration).
			Item("Foo").
			Item("Bar").
			Item("Baz")

		if l.String() != test.want {
			t.Errorf("want:\n\n%s\n\ngot:\n\n%s\n", test.want, l.String())
		}
	}
}

func TestBullet(t *testing.T) {
	tests := []struct {
		enumeration Enumerator
		i           int
		exp         string
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
		bullet := strings.TrimSuffix(test.enumeration(test.i), ".")
		if bullet != test.exp {
			t.Errorf("want: %s, got: %s\n", test.exp, bullet)
		}
	}
}

func TestEnumeratorsAlign(t *testing.T) {
	fooList := strings.Split(strings.TrimSuffix(strings.Repeat("Foo ", 100), " "), " ")
	l := New().Enumerator(Roman)
	for _, f := range fooList {
		l = l.Item(f)
	}

	want := heredoc.Doc(`
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
       C. Foo`)

	if l.String() != want {
		t.Fatalf("want:\n\n%s\n\ngot:\n\n%s\n", want, l.String())
	}
}

func TestIndent(t *testing.T) {
	l := New("foo", "bar", "baz").Enumerator(Arabic)

	want := heredoc.Doc(`
  1. foo
  2. bar
  3. baz`)

	if l.String() != want {
		t.Fatalf("want:\n\n%s\n\ngot:\n\n%s\n", want, l.String())
	}
}
