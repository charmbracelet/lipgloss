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
)

func TestTree(t *testing.T) {
	tr := tree.New().
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items(
					"Qux",
					tree.New().
						Root("Quux").
						Items(
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
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items(
					"Qux",
					tree.New().
						Root("Quux").
						Items("Foo", "Bar").
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
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items(
					"Qux",
					tree.New().
						Root("Quux").
						Items(
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
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items("Qux", "Quuux"),
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
		Items(
			tree.New().
				Root("Bar").
				Items("Qux", "Quuux"),
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
		Items(
			"Bar",
			"Foo",
			tree.New().
				Items(
					"Qux",
					"Qux",
					"Qux",
					"Qux",
					"Qux",
				),
			tree.New().
				Items(
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
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items(
					"Qux",
					tree.New().
						Root("Quux").
						Items("Foo", "Bar"),
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
		Items(
			nil,
			tree.New().
				Root("Bar").
				Items(
					"Qux",
					tree.New().
						Root("Quux").
						Item("Bar"),
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
	quuux := tree.StringNode("Quuux")
	tree := tree.New().
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items(
					tree.StringNode("Qux"),
					tree.New().
						Root("Quux").
						Items("Foo",
							"Bar",
						),
					&quuux,
				),
			"Baz",
		).
		ItemStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))).
		EnumeratorStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			PaddingRight(1)).
		Enumerator(func(tree.Data, int) (indent string, prefix string) {
			return "->", "->"
		})

	want := `
-> Foo
-> Bar
-> -> Qux
-> -> Quux
-> -> -> Foo
-> -> -> Bar
-> -> Quuux
-> Baz
	`
	assertEqual(t, want, tree.String())
}

func TestTreeMultilineNode(t *testing.T) {
	tree := tree.New().
		Root("Big\nRoot\nNode").
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items(
					"Line 1\nLine 2\nLine 3\nLine 4",
					tree.New().
						Root("Quux").
						Items(
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
		Items(
			tree.New().
				Root("Parent").
				Items("child 1", "child 2").
				ItemStyleFunc(func(tree.Data, int) lipgloss.Style {
					return lipgloss.NewStyle().
						SetString("*")
				}).
				EnumeratorStyleFunc(func(_ tree.Data, i int) lipgloss.Style {
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
		Items(
			"Foo",
			"Foo",
			"Foo",
			"Foo",
			"Foo",
		).Enumerator(func(_ tree.Data, i int) (indent string, prefix string) {
		romans := map[int]string{
			1: "I",
			2: "II",
			3: "III",
			4: "IV",
			5: "V",
			6: "VI",
		}
		return "", romans[i+1]
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
		Items("Willy ", "Nilly").
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
		Items(
			"Foo",
			"Baz",
		).Enumerator(func(data tree.Data, i int) (indent string, prefix string) {
		if data.At(i).Name() == "Foo" {
			return "", ">"
		}
		return "", "-"
	})

	want := `
Root
> Foo
- Baz
	`
	assertEqual(t, want, tree.String())
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
			return index != 1
		}).
		Append(tree.StringNode("Qux")).
		Remove(3)

	tree := tree.New().
		Root("Root").
		Item(data)

	want := `
Root
├── Foo
├── Baz
└── Qux
	`

	assertEqual(t, want, tree.String())
	if got := data.At(1); got.Name() != "Baz" {
		t.Errorf("want Baz, got %v", got)
	}
	if got := data.At(10); got != nil {
		t.Errorf("want nil, got %v", got)
	}
}

func TestNodeDataRemoveOutOfBounds(t *testing.T) {
	data := tree.NewStringData("a").Remove(-1).Remove(1)
	if l := data.Length(); l != 1 {
		t.Errorf("want data to contain 1 items, has %d", l)
	}
}

func TestTreeTable(t *testing.T) {
	tree := tree.New().
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items(
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
		Items(
			"Foo",
			"Bar",
			tree.New().
				Item("Baz"),
			"Qux",
		)

	t2 := tree.New().
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Item("Baz"),
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
		Item(list.New("A", "B", "C").
			Enumerator(list.Arabic)).
		Item(list.New("1", "2", "3").
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
		Enumerator(func(_ tree.Data, i int) (string, string) {
			if i == 1 {
				return "", "│\n│"
			}
			return "", " "
		}).
		ItemStyle(paddingsStyle).
		Item("Foo Document\nThe Foo Files").
		Item("Bar Document\nThe Bar Files").
		Item("Baz Document\nThe Baz Files")
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
		Item("Foo").
		Item("Bar").
		Item(
			tree.New().
				Root("Baz").
				Enumerator(func(_ tree.Data, i int) (string, string) {
					if i == 1 {
						return "", "│\n│"
					}
					return "", " "
				}).
				ItemStyle(paddingsStyle).
				Item("Foo Document\nThe Foo Files").
				Item("Bar Document\nThe Bar Files").
				Item("Baz Document\nThe Baz Files"),
		).
		Item("Qux")
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
	glowEnum := func(_ tree.Data, i int) (string, string) {
		if i == 1 {
			return "  ", "│\n│"
		}
		return "  ", " "
	}
	paddingsStyle := lipgloss.NewStyle().PaddingLeft(1).PaddingBottom(1)
	tree := tree.New().
		Enumerator(glowEnum).
		ItemStyle(paddingsStyle).
		Item("Foo Document\nThe Foo Files").
		Item("Bar Document\nThe Bar Files").
		Item(
			tree.New().
				Enumerator(glowEnum).
				ItemStyle(paddingsStyle).
				Item("Qux Document\nThe Qux Files").
				Item("Quux Document\nThe Quux Files").
				Item("Quuux Document\nThe Quuux Files"),
		).
		Item("Baz Document\nThe Baz Files")
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
		Item(0).
		Item(true).
		Item([]any{"Foo", "Bar"}).
		Item([]string{"Qux", "Quux", "Quuux"})

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
