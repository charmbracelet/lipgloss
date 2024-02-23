package tree_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/internal/require"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss/tree"
)

func TestTree(t *testing.T) {
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
		)

	expected := `
├── Foo
├── Bar
│   ├── Qux
│   ├── Quux
│   │   ├── Foo
│   │   └── Bar
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
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

	expected := `
├── Foo
├── Bar
│   ├── Qux
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
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

	expected := ``
	require.Equal(t, expected, tree.String())
}

func TestTreeRoot(t *testing.T) {
	tree := tree.New().
		Root("The Root").
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items("Qux", "Quuux"),
			"Baz",
		)
	expected := `
The Root
├── Foo
├── Bar
│   ├── Qux
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeStartsWithSubtree(t *testing.T) {
	tree := tree.New().Items(
		tree.New().
			Root("Bar").
			Items("Qux", "Quuux"),
		"Baz",
	)

	expected := `
├── Bar
│   ├── Qux
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeAddTwoSubTreesWithoutName(t *testing.T) {
	tree := tree.New().Items(
		"bar",
		"foo",
		tree.New().Items(
			"Bar 11",
			"Bar 12",
			"Bar 13",
			"Bar 14",
			"Bar 15",
		),
		tree.New().Items(
			"Bar 21",
			"Bar 22",
			"Bar 23",
			"Bar 24",
			"Bar 25",
		),
		"Baz",
	)

	expected := `
├── bar
├── foo
│   ├── Bar 11
│   ├── Bar 12
│   ├── Bar 13
│   ├── Bar 14
│   ├── Bar 15
│   ├── Bar 21
│   ├── Bar 22
│   ├── Bar 23
│   ├── Bar 24
│   └── Bar 25
└── Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeLastNodeIsSubTree(t *testing.T) {
	tree := tree.New().Items(
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

	expected := `
├── Foo
└── Bar
    ├── Qux
    ├── Quux
    │   ├── Foo
    │   └── Bar
    └── Quuux
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeNil(t *testing.T) {
	tree := tree.New().Items(
		nil,
		tree.New().
			Root("Bar").Items(
			"Qux",
			tree.New().Root("Quux").Item("Bar"),
			"Quuux",
		),
		"Baz",
	)

	expected := `
├── Bar
│   ├── Qux
│   ├── Quux
│   │   └── Bar
│   └── Quuux
└── Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeCustom(t *testing.T) {
	quuux := tree.StringNode("Quuux")
	tree := tree.New().Items(
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
		ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("9"))).
		EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).MarginRight(1)).
		Enumerator(func(tree.Data, int) (indent string, prefix string) {
			return "->", "->"
		})

	expected := `
-> Foo
-> Bar
-> -> Qux
-> -> Quux
-> -> -> Foo
-> -> -> Bar
-> -> Quuux
-> Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeMultilineNode(t *testing.T) {
	tree := tree.New().
		Root("Multiline\nRoot\nNode").
		Items(
			"Foo",
			tree.New().
				Root("Bar").
				Items(
					"Qux\nLine 2\nLine 3\nLine 4",
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

	expected := `
Multiline
Root
Node
├── Foo
├── Bar
│   ├── Qux
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
	require.Equal(t, expected, tree.String())
}

func TestTreeSubTreeWithCustomRenderer(t *testing.T) {
	tree := tree.New().
		Root("The Root Node(tm)").
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
						MarginRight(1)
				}),
			"Baz",
		)

	expected := `
The Root Node(tm)
├── Parent
│   + ├── * child 1
│   + └── * child 2
└── Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeMixedEnumeratorSize(t *testing.T) {
	tree := tree.New().
		Root("The Root Node(tm)").
		Items(
			"child 1",
			"child 2",
			"child 3",
			"child 4",
			"child 5",
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

	expected := `
The Root Node(tm)
  I child 1
 II child 2
III child 3
 IV child 4
  V child 5
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeStyleNilFuncs(t *testing.T) {
	tree := tree.New().
		Root("Multiline").
		Items(
			"Foo",
			"Baz",
		).ItemStyleFunc(nil).
		EnumeratorStyleFunc(nil)

	expected := `
Multiline
├──Foo
└──Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeStyleAt(t *testing.T) {
	tree := tree.New().
		Root("Multiline").
		Items(
			"Foo",
			"Baz",
		).Enumerator(func(data tree.Data, i int) (indent string, prefix string) {
		if data.At(i).Name() == "Foo" {
			return "", ">"
		}
		return "", "-"
	})

	expected := `
Multiline
> Foo
- Baz
	`
	require.Equal(t, expected, tree.String())
}

func TestAt(t *testing.T) {
	nodes := tree.NewStringData("foo", "bar")
	t.Run("0", func(t *testing.T) {
		if s := nodes.At(0).String(); s != "foo" {
			t.Errorf("expected 'foo', got '%s'", s)
		}
	})
	t.Run("10", func(t *testing.T) {
		if n := nodes.At(10); n != nil {
			t.Errorf("expected nil, got '%s'", n)
		}
	})
	t.Run("-1", func(t *testing.T) {
		if n := nodes.At(10); n != nil {
			t.Errorf("expected nil, got '%s'", n)
		}
	})
}

func TestFilter(t *testing.T) {
	data := tree.NewFilter(tree.NewStringData("Foo", "Bar", "Baz", "Nope")).
		Filter(func(index int) bool {
			return index != 1
		}).
		Append(tree.StringNode("Qux")).
		Remove(3)
	tree := tree.New().Root("Root").Item(data)

	expected := `
Root
├── Foo
├── Baz
└── Qux
	`

	require.Equal(t, expected, tree.String())
	if got := data.At(1); got.Name() != "Baz" {
		t.Errorf("expected to get Baz, got %v", got)
	}
	if got := data.At(10); got != nil {
		t.Errorf("expected to get nil, got %v", got)
	}
}

func TestNodeDataRemoveOutOfBounds(t *testing.T) {
	data := tree.NewStringData("a").Remove(-1).Remove(1)
	if l := data.Length(); l != 1 {
		t.Errorf("expected data to contain 1 items, has %d", l)
	}
}

func TestTreeTable(t *testing.T) {
	tree := tree.New().
		Items(
			"a",
			tree.New().
				Root("b").
				Items(
					"c",
					"d",
					table.New().
						Width(40).
						Headers("a", "b").
						Row("1", "2").
						Row("3", "4"),
					"e",
				),
			"c",
		)
	expected := `
├── a
├── b
│   ├── c
│   ├── d
│   ├── ╭───────────────────┬──────────────────╮
│   │   │a                  │b                 │
│   │   ├───────────────────┼──────────────────┤
│   │   │1                  │2                 │
│   │   │3                  │4                 │
│   │   ╰───────────────────┴──────────────────╯
│   └── e
└── c
	`
	require.Equal(t, expected, tree.String())
}

func TestTreeOffset(t *testing.T) {
	enum := func(tree.Data, int) (string, string) {
		return "", "*"
	}

	t.Run("start", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			tree := tree.New().
				Root("root").
				Items("a", "b", "c", "d").
				OffsetStart(0).
				Enumerator(enum)

			expected := `
root
* a
* b
* c
* d
		`
			require.Equal(t, expected, tree.String())
		})

		t.Run("in bounds", func(t *testing.T) {
			tree := tree.New().
				Root("root").
				Items("a", "b", "c", "d").
				OffsetStart(2).
				Enumerator(enum)

			expected := `
root
* c
* d
		`
			require.Equal(t, expected, tree.String())
		})

		t.Run("max", func(t *testing.T) {
			tree := tree.New().
				Root("root").
				Items("a", "b", "c", "d").
				OffsetStart(4).
				Enumerator(enum)

			expected := `
root
		`
			require.Equal(t, expected, tree.String())
		})
		t.Run("out bounds", func(t *testing.T) {
			tree := tree.New().
				Root("root").
				Items("a", "b", "c", "d").
				OffsetStart(6).
				Enumerator(enum)

			expected := `
root
		`
			require.Equal(t, expected, tree.String())
		})
	})
	t.Run("end", func(t *testing.T) {
		t.Run("min", func(t *testing.T) {
			tree := tree.New().
				Root("root").
				Items("a", "b", "c", "d").
				OffsetEnd(0).
				Enumerator(enum)

			expected := `
root
* a
* b
* c
* d
		`
			require.Equal(t, expected, tree.String())
		})

		t.Run("in bounds", func(t *testing.T) {
			tree := tree.New().
				Root("root").
				Items("a", "b", "c", "d").
				OffsetEnd(2).
				Enumerator(enum)

			expected := `
root
* a
* b
		`
			require.Equal(t, expected, tree.String())
		})

		t.Run("max", func(t *testing.T) {
			tree := tree.New().
				Root("root").
				Items("a", "b", "c", "d").
				OffsetEnd(4).
				Enumerator(enum)

			expected := `
root
		`
			require.Equal(t, expected, tree.String())
		})
		t.Run("out bounds", func(t *testing.T) {
			tree := tree.New().
				Root("root").
				Items("a", "b", "c", "d").
				OffsetEnd(6).
				Enumerator(enum)

			expected := `
root
		`
			require.Equal(t, expected, tree.String())
		})
	})

	t.Run("start and end", func(t *testing.T) {
		tree := tree.New().
			Root("root").
			Items("a", "b", "c", "d").
			OffsetStart(1).
			OffsetEnd(2).
			Enumerator(enum)

		expected := `
root
* b
		`
		require.Equal(t, expected, tree.String())
	})
}

func TestAddItemWithoutRootAndWithRoot(t *testing.T) {
	t1 := tree.New().
		Items(
			"foo",
			"bar",
			tree.New().
				Item("zaz"),
			"qux",
		)
	t2 := tree.New().
		Items(
			"foo",
			tree.New().
				Root("bar").
				Item("zaz"),
			"qux",
		)
	expected := `
├── foo
├── bar
│   └── zaz
└── qux
	`
	require.Equal(t, expected, t1.String())
	require.Equal(t, expected, t2.String())
}

func TestEmbedListWithinTree(t *testing.T) {
	t1 := tree.New().
		Item(list.New("A", "B", "C").Enumerator(list.Arabic)).
		Item(list.New("1", "2", "3").Enumerator(list.Alphabet))
	expected := `
├── 1. A
│   2. B
│   3. C
└── A. 1
    B. 2
    C. 3
	`
	require.Equal(t, expected, t1.String())
}

func TestMultilinePrefix(t *testing.T) {
	marginsStyle := lipgloss.NewStyle().MarginLeft(1).MarginBottom(1)
	tree := tree.New().
		Enumerator(func(_ tree.Data, i int) (string, string) {
			if i == 1 {
				return "", "│\n│"
			}
			return "", " "
		}).
		ItemStyle(marginsStyle).
		Item("Document 0\nSome tagline").
		Item("Document 1\nHello world").
		Item("Document 2\nSome other tagline")
	expected := `
   Document 0
   Some tagline

│  Document 1
│  Hello world

   Document 2
   Some other tagline
	`
	require.Equal(t, expected, tree.String())
}

func TestMultilinePrefixSingleLineItem(t *testing.T) {
	marginsStyle := lipgloss.NewStyle().MarginLeft(1).MarginBottom(1)
	tree := tree.New().
		Enumerator(func(_ tree.Data, i int) (string, string) {
			if i == 1 {
				return "", "│\n│"
			}
			return "", " "
		}).
		ItemStyle(marginsStyle).
		Item("Document 0\nhello").
		Item("Document 1\n").
		Item("Document 2\nhello again")
	expected := `
   Document 0
   hello

│  Document 1
│

   Document 2
   hello again
	`
	require.Equal(t, expected, tree.String())
}

func TestMultilinePrefixSubtree(t *testing.T) {
	marginsStyle := lipgloss.NewStyle().MarginLeft(1).MarginBottom(1)
	tree := tree.New().
		Item("Hello").
		Item("Foo").
		Item(
			tree.New().
				Root("Bar").
				Enumerator(func(_ tree.Data, i int) (string, string) {
					if i == 1 {
						return "", "│\n│"
					}
					return "", " "
				}).
				ItemStyle(marginsStyle).
				Item("Document 0\nSome tagline").
				Item("Document 1\nHello world").
				Item("Document 2\nSome other tagline"),
		).
		Item("Fuss")
	expected := `
├── Hello
├── Foo
├── Bar
│      Document 0
│      Some tagline
│
│   │  Document 1
│   │  Hello world
│
│      Document 2
│      Some other tagline
│
└── Fuss
	`
	require.Equal(t, expected, tree.String())
}

func TestMultilinePrefixInception(t *testing.T) {
	glowEnum := func(_ tree.Data, i int) (string, string) {
		if i == 1 {
			return "  ", "│\n│"
		}
		return "  ", " "
	}
	marginsStyle := lipgloss.NewStyle().MarginLeft(1).MarginBottom(1)
	tree := tree.New().
		Enumerator(glowEnum).
		ItemStyle(marginsStyle).
		Item("Document 0\nSome tagline").
		Item("Document 1\nHello world").
		Item(
			tree.New().
				Enumerator(glowEnum).
				ItemStyle(marginsStyle).
				Item("Document 1a\nnothing important").
				Item("Document 1b\nsomething").
				Item("Document 1c\nsomething else"),
		).
		Item("Document 2\nSome other tagline")
	expected := `
    Document 0
    Some tagline

│   Document 1
│   Hello world

       Document 1a
       nothing important

   │   Document 1b
   │   something

       Document 1c
       something else

    Document 2
    Some other tagline
	`
	require.Equal(t, expected, tree.String())
}
