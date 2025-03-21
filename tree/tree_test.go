package tree_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/aymanbagabas/go-udiff"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/muesli/termenv"
)

func TestTree(t *testing.T) {
	tr := tree.New().
		Child(
			"Foo",
			tree.Root("Bar").
				Child(
					"Qux",
					tree.Root("Quux").
						Child(
							"Foo",
							"Bar",
						),
					"Quuux",
				),
			"Baz",
		)

	want := `
├── Foo
├── Bar
│   ├── Qux
│   ├── Quux
│   │   ├── Foo
│   │   └── Bar
│   └── Quuux
└── Baz
	`
	assertEqual(t, want, tr.String())

	tr.Enumerator(tree.RoundedEnumerator)

	want = `
├── Foo
├── Bar
│   ├── Qux
│   ├── Quux
│   │   ├── Foo
│   │   ╰── Bar
│   ╰── Quuux
╰── Baz
`
	assertEqual(t, want, tr.String())
}

func TestTreeHidden(t *testing.T) {
	tree := tree.New().
		Child(
			"Foo",
			tree.Root("Bar").
				Child(
					"Qux",
					tree.Root("Quux").
						Child("Foo", "Bar").
						Hide(true),
					"Quuux",
				),
			"Baz",
		)

	want := `
├── Foo
├── Bar
│   ├── Qux
│   └── Quuux
└── Baz
	`
	assertEqual(t, want, tree.String())
}

func TestTreeAllHidden(t *testing.T) {
	tree := tree.New().
		Child(
			"Foo",
			tree.Root("Bar").
				Child(
					"Qux",
					tree.Root("Quux").
						Child(
							"Foo",
							"Bar",
						),
					"Quuux",
				),
			"Baz",
		).Hide(true)

	want := ``

	assertEqual(t, want, tree.String())
}

func TestTreeRoot(t *testing.T) {
	tree := tree.New().
		Root("Root").
		Child(
			"Foo",
			tree.Root("Bar").
				Child("Qux", "Quuux"),
			"Baz",
		)
	want := `
Root
├── Foo
├── Bar
│   ├── Qux
│   └── Quuux
└── Baz
	`
	assertEqual(t, want, tree.String())
}

func TestTreeStartsWithSubtree(t *testing.T) {
	tree := tree.New().
		Child(
			tree.New().
				Root("Bar").
				Child("Qux", "Quuux"),
			"Baz",
		)

	want := `
├── Bar
│   ├── Qux
│   └── Quuux
└── Baz
	`
	assertEqual(t, want, tree.String())
}

func TestTreeAddTwoSubTreesWithoutName(t *testing.T) {
	tree := tree.New().
		Child(
			"Bar",
			"Foo",
			tree.New().
				Child(
					"Qux",
					"Qux",
					"Qux",
					"Qux",
					"Qux",
				),
			tree.New().
				Child(
					"Quux",
					"Quux",
					"Quux",
					"Quux",
					"Quux",
				),
			"Baz",
		)

	want := `
├── Bar
├── Foo
│   ├── Qux
│   ├── Qux
│   ├── Qux
│   ├── Qux
│   ├── Qux
│   ├── Quux
│   ├── Quux
│   ├── Quux
│   ├── Quux
│   └── Quux
└── Baz
	`
	assertEqual(t, want, tree.String())
}

func TestTreeLastNodeIsSubTree(t *testing.T) {
	tree := tree.New().
		Child(
			"Foo",
			tree.Root("Bar").
				Child("Qux",
					tree.Root("Quux").Child("Foo", "Bar"),
					"Quuux",
				),
		)

	want := `
├── Foo
└── Bar
    ├── Qux
    ├── Quux
    │   ├── Foo
    │   └── Bar
    └── Quuux
	`
	assertEqual(t, want, tree.String())
}

func TestTreeNil(t *testing.T) {
	tree := tree.New().
		Child(
			nil,
			tree.Root("Bar").
				Child(
					"Qux",
					tree.Root("Quux").
						Child("Bar"),
					"Quuux",
				),
			"Baz",
		)

	want := `
├── Bar
│   ├── Qux
│   ├── Quux
│   │   └── Bar
│   └── Quuux
└── Baz
	`
	assertEqual(t, want, tree.String())
}

func TestTreeCustom(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	tree := tree.New().
		Child(
			"Foo",
			tree.New().
				Root("Bar").
				Child(
					"Qux",
					tree.New().
						Root("Quux").
						Child("Foo",
							"Bar",
						),
					"Quuux",
				),
			"Baz",
		).
		ItemStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))).
		EnumeratorStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			PaddingRight(1)).
		Enumerator(func(tree.Children, int) string {
			return "->"
		}).
		Indenter(func(tree.Children, int) string {
			return "->"
		})

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestTreeMultilineNode(t *testing.T) {
	tree := tree.New().
		Root("Big\nRoot\nNode").
		Child(
			"Foo",
			tree.New().
				Root("Bar").
				Child(
					"Line 1\nLine 2\nLine 3\nLine 4",
					tree.New().
						Root("Quux").
						Child(
							"Foo",
							"Bar",
						),
					"Quuux",
				),
			"Baz\nLine 2",
		)

	want := `
Big
Root
Node
├── Foo
├── Bar
│   ├── Line 1
│   │   Line 2
│   │   Line 3
│   │   Line 4
│   ├── Quux
│   │   ├── Foo
│   │   └── Bar
│   └── Quuux
└── Baz
    Line 2
	`
	assertEqual(t, want, tree.String())
}

func TestTreeSubTreeWithCustomEnumerator(t *testing.T) {
	tree := tree.New().
		Root("The Root Node™").
		Child(
			tree.New().
				Root("Parent").
				Child("child 1", "child 2").
				ItemStyleFunc(func(tree.Children, int) lipgloss.Style {
					return lipgloss.NewStyle().
						SetString("*")
				}).
				EnumeratorStyleFunc(func(_ tree.Children, i int) lipgloss.Style {
					return lipgloss.NewStyle().
						SetString("+").
						PaddingRight(1)
				}),
			"Baz",
		)

	want := `
The Root Node™
├── Parent
│   + ├── * child 1
│   + └── * child 2
└── Baz
	`
	assertEqual(t, want, tree.String())
}

func TestTreeMixedEnumeratorSize(t *testing.T) {
	tree := tree.New().
		Root("The Root Node™").
		Child(
			"Foo",
			"Foo",
			"Foo",
			"Foo",
			"Foo",
		).Enumerator(func(_ tree.Children, i int) string {
		romans := map[int]string{
			1: "I",
			2: "II",
			3: "III",
			4: "IV",
			5: "V",
			6: "VI",
		}
		return romans[i+1]
	})

	want := `
The Root Node™
  I Foo
 II Foo
III Foo
 IV Foo
  V Foo
	`
	assertEqual(t, want, tree.String())
}

func TestTreeStyleNilFuncs(t *testing.T) {
	tree := tree.New().
		Root("Silly").
		Child("Willy ", "Nilly").
		ItemStyleFunc(nil).
		EnumeratorStyleFunc(nil)

	want := `
Silly
├──Willy
└──Nilly
	`
	assertEqual(t, want, tree.String())
}

func TestTreeStyleAt(t *testing.T) {
	tree := tree.New().
		Root("Root").
		Child(
			"Foo",
			"Baz",
		).Enumerator(func(data tree.Children, i int) string {
		if data.At(i).Value() == "Foo" {
			return ">"
		}
		return "-"
	})

	want := `
Root
> Foo
- Baz
	`
	assertEqual(t, want, tree.String())
}

func TestRootStyle(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	tree := tree.New().
		Root("Root").
		Child(
			"Foo",
			"Baz",
		).
		RootStyle(lipgloss.NewStyle().Background(lipgloss.Color("#5A56E0"))).
		ItemStyle(lipgloss.NewStyle().Background(lipgloss.Color("#04B575")))

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestInheritedStyles(t *testing.T) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	tr := tree.
		Root("⁜ Makeup").
		Child(
			"Glossier",
			"Fenty Beauty",
			tree.New().Child(
				"Gloss Bomb Universal Lip Luminizer",
				"Hot Cheeks Velour Blushlighter",
			),
			"Nyx",
			"Mac",
			"Milk",
		).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyle(itemStyle)
	// Add a Tree as a Child of "Glossier"
	tr.Replace(0, tr.Children().At(0).Child(
		tree.Root("Apparel").Child("Pink Hoodie", "Baseball Cap"),
	))

	// Add a Leaf as a Child of "Glossier"
	tr.Children().At(0).Child("Makeup")
	golden.RequireEqual(t, []byte(tr.String()))
}

func TestAt(t *testing.T) {
	data := tree.NewStringData("Foo", "Bar")

	if s := data.At(0).String(); s != "Foo" {
		t.Errorf("want 'Foo', got '%s'", s)
	}

	if n := data.At(10); n != nil {
		t.Errorf("want nil, got '%s'", n)
	}

	if n := data.At(-1); n != nil {
		t.Errorf("want nil, got '%s'", n)
	}
}

func TestFilter(t *testing.T) {
	data := tree.NewFilter(tree.NewStringData(
		"Foo",
		"Bar",
		"Baz",
		"Nope",
	)).
		Filter(func(index int) bool {
			return index != 3
		})

	tree := tree.New().
		Root("Root").
		Child(data)

	want := `
Root
├── Foo
├── Bar
└── Baz
	`

	assertEqual(t, want, tree.String())
	if got := data.At(1); got.Value() != "Bar" {
		t.Errorf("want Bar, got %v", got)
	}
	if got := data.At(10); got != nil {
		t.Errorf("want nil, got %v", got)
	}
}

func TestNodeDataRemoveOutOfBounds(t *testing.T) {
	data := tree.NewStringData("a")
	if l := data.Length(); l != 1 {
		t.Errorf("want data to contain 1 items, has %d", l)
	}
}

func TestTreeTable(t *testing.T) {
	tree := tree.New().
		Child(
			"Foo",
			tree.New().
				Root("Bar").
				Child(
					"Baz",
					"Baz",
					table.New().
						Width(20).
						StyleFunc(func(row, col int) lipgloss.Style {
							return lipgloss.NewStyle().Padding(0, 1)
						}).
						Headers("Foo", "Bar").
						Row("Qux", "Baz").
						Row("Qux", "Baz"),
					"Baz",
				),
			"Qux",
		)
	want := `
├── Foo
├── Bar
│   ├── Baz
│   ├── Baz
│   ├── ╭─────────┬────────╮
│   │   │ Foo     │ Bar    │
│   │   ├─────────┼────────┤
│   │   │ Qux     │ Baz    │
│   │   │ Qux     │ Baz    │
│   │   ╰─────────┴────────╯
│   └── Baz
└── Qux
	`
	assertEqual(t, want, tree.String())
}

func TestAddItemWithAndWithoutRoot(t *testing.T) {
	t1 := tree.New().
		Child(
			"Foo",
			"Bar",
			tree.New().
				Child("Baz"),
			"Qux",
		)

	t2 := tree.New().
		Child(
			"Foo",
			tree.New().
				Root("Bar").
				Child("Baz"),
			"Qux",
		)

	want := `
├── Foo
├── Bar
│   └── Baz
└── Qux
	`
	assertEqual(t, want, t1.String())
	assertEqual(t, want, t2.String())
}

func TestEmbedListWithinTree(t *testing.T) {
	t1 := tree.New().
		Child(list.New("A", "B", "C").
			Enumerator(list.Arabic)).
		Child(list.New("1", "2", "3").
			Enumerator(list.Alphabet))

	want := `
├── 1. A
│   2. B
│   3. C
└── A. 1
    B. 2
    C. 3
	`
	assertEqual(t, want, t1.String())
}

func TestMultilinePrefix(t *testing.T) {
	paddingsStyle := lipgloss.NewStyle().PaddingLeft(1).PaddingBottom(1)
	tree := tree.New().
		Enumerator(func(_ tree.Children, i int) string {
			if i == 1 {
				return "│\n│"
			}
			return " "
		}).
		Indenter(func(_ tree.Children, i int) string {
			return " "
		}).
		ItemStyle(paddingsStyle).
		Child("Foo Document\nThe Foo Files").
		Child("Bar Document\nThe Bar Files").
		Child("Baz Document\nThe Baz Files")
	want := `
   Foo Document
   The Foo Files

│  Bar Document
│  The Bar Files

   Baz Document
   The Baz Files
	`
	assertEqual(t, want, tree.String())
}

func TestMultilinePrefixSubtree(t *testing.T) {
	paddingsStyle := lipgloss.NewStyle().
		Padding(0, 0, 1, 1)

	tree := tree.New().
		Child("Foo").
		Child("Bar").
		Child(
			tree.New().
				Root("Baz").
				Enumerator(func(_ tree.Children, i int) string {
					if i == 1 {
						return "│\n│"
					}
					return " "
				}).
				Indenter(func(tree.Children, int) string {
					return " "
				}).
				ItemStyle(paddingsStyle).
				Child("Foo Document\nThe Foo Files").
				Child("Bar Document\nThe Bar Files").
				Child("Baz Document\nThe Baz Files"),
		).
		Child("Qux")
	want := `
├── Foo
├── Bar
├── Baz
│      Foo Document
│      The Foo Files
│
│   │  Bar Document
│   │  The Bar Files
│
│      Baz Document
│      The Baz Files
│
└── Qux
	`
	assertEqual(t, want, tree.String())
}

func TestMultilinePrefixInception(t *testing.T) {
	glowEnum := func(_ tree.Children, i int) string {
		if i == 1 {
			return "│\n│"
		}
		return " "
	}
	glowIndenter := func(_ tree.Children, i int) string {
		return "  "
	}
	paddingsStyle := lipgloss.NewStyle().PaddingLeft(1).PaddingBottom(1)
	tree := tree.New().
		Enumerator(glowEnum).
		Indenter(glowIndenter).
		ItemStyle(paddingsStyle).
		Child("Foo Document\nThe Foo Files").
		Child("Bar Document\nThe Bar Files").
		Child(
			tree.New().
				Enumerator(glowEnum).
				Indenter(glowIndenter).
				ItemStyle(paddingsStyle).
				Child("Qux Document\nThe Qux Files").
				Child("Quux Document\nThe Quux Files").
				Child("Quuux Document\nThe Quuux Files"),
		).
		Child("Baz Document\nThe Baz Files")
	want := `
    Foo Document
    The Foo Files

│   Bar Document
│   The Bar Files

       Qux Document
       The Qux Files

   │   Quux Document
   │   The Quux Files

       Quuux Document
       The Quuux Files

    Baz Document
    The Baz Files
	`
	assertEqual(t, want, tree.String())
}

func TestTypes(t *testing.T) {
	tree := tree.New().
		Child(0).
		Child(true).
		Child([]any{"Foo", "Bar"}).
		Child([]string{"Qux", "Quux", "Quuux"})

	want := `
├── 0
├── true
├── Foo
├── Bar
├── Qux
├── Quux
└── Quuux
	`
	assertEqual(t, want, tree.String())
}

// assertEqual verifies the strings are equal, assuming its terminal output.
func assertEqual(tb testing.TB, want, got string) {
	tb.Helper()

	want = trimSpace(want)
	got = trimSpace(got)

	diff := udiff.Unified("want", "got", want, got)
	if diff != "" {
		tb.Fatalf("\nwant:\n\n%s\n\ngot:\n\n%s\n\ndiff:\n\n%s\n\n", want, got, diff)
	}
}

func trimSpace(s string) string {
	var result []string //nolint: prealloc
	ss := strings.Split(s, "\n")
	for i, line := range ss {
		if strings.TrimSpace(line) == "" && (i == 0 || i == len(ss)-1) {
			continue
		}
		result = append(result, strings.TrimRightFunc(line, unicode.IsSpace))
	}
	return strings.Join(result, "\n")
}
