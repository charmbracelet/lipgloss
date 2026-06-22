package tree_test

import (
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/list"
	"charm.land/lipgloss/v2/table"
	"charm.land/lipgloss/v2/tree"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
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

	t.Run("before", func(t *testing.T) {
		golden.RequireEqual(t, []byte(tr.String()))
	})

	tr.Enumerator(tree.RoundedEnumerator)

	t.Run("after", func(t *testing.T) {
		golden.RequireEqual(t, []byte(tr.String()))
	})
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestTreeStartsWithSubtree(t *testing.T) {
	tree := tree.New().
		Child(
			tree.New().
				Root("Bar").
				Child("Qux", "Quuux"),
			"Baz",
		)

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestTreeCustom(t *testing.T) {
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
		IndenterStyle(lipgloss.NewStyle().
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestTreeStyleNilFuncs(t *testing.T) {
	tree := tree.New().
		Root("Silly").
		Child("Willy ", "Nilly").
		ItemStyleFunc(nil).
		EnumeratorStyleFunc(nil)

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestRootStyle(t *testing.T) {
	tree := tree.New().
		Root("Root").
		Child(
			"Foo",
			"Baz",
		).
		RootStyle(lipgloss.NewStyle().Background(lipgloss.Color("#5A56E0"))).
		ItemStyle(lipgloss.NewStyle().Background(lipgloss.Color("#04B575")))

	golden.RequireEqual(t, []byte(ansi.Strip(tree.String())))
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

func TestTreeOffset(t *testing.T) {
	build := func() *tree.Tree {
		return tree.New().
			Root("Root").
			Child("A", "B", "C", "D", "E")
	}

	t.Run("zero_offset", func(t *testing.T) {
		// Offset(0,0) should show all children.
		tr := build()
		tr.Offset(0, 0)
		got := tr.Children()
		want := []string{"A", "B", "C", "D", "E"}
		assertChildren(t, got, want)
	})

	t.Run("skip_from_start", func(t *testing.T) {
		// Offset(2,0) should skip A,B → show C,D,E.
		tr := build()
		tr.Offset(2, 0)
		got := tr.Children()
		want := []string{"C", "D", "E"}
		assertChildren(t, got, want)
	})

	t.Run("skip_from_end", func(t *testing.T) {
		// Offset(0,2) should skip D,E → show A,B,C.
		tr := build()
		tr.Offset(0, 2)
		got := tr.Children()
		want := []string{"A", "B", "C"}
		assertChildren(t, got, want)
	})

	t.Run("skip_from_both", func(t *testing.T) {
		// Offset(2,1) should skip A,B from start and E from end → show C,D.
		tr := build()
		tr.Offset(2, 1)
		got := tr.Children()
		want := []string{"C", "D"}
		assertChildren(t, got, want)
	})

	t.Run("start_greater_than_end", func(t *testing.T) {
		// Offset(3,1) — start > end is valid, no swap should happen.
		// Skip A,B,C from start, skip E from end → show D.
		tr := build()
		tr.Offset(3, 1)
		got := tr.Children()
		want := []string{"D"}
		assertChildren(t, got, want)
	})

	t.Run("offset_larger_than_children", func(t *testing.T) {
		// Offset(10,0) should produce empty children.
		tr := build()
		tr.Offset(10, 0)
		got := tr.Children()
		if got.Length() != 0 {
			t.Errorf("expected 0 children, got %d", got.Length())
		}
	})

	t.Run("negative_values_clamped_to_zero", func(t *testing.T) {
		// Negative values should be treated as zero.
		tr := build()
		tr.Offset(-1, -2)
		got := tr.Children()
		want := []string{"A", "B", "C", "D", "E"}
		assertChildren(t, got, want)
	})

	t.Run("offset_swallows_all", func(t *testing.T) {
		// Offset(2,3) on 5 items: skip A,B + C,D,E = nothing left.
		tr := build()
		tr.Offset(2, 3)
		got := tr.Children()
		if got.Length() != 0 {
			t.Errorf("expected 0 children, got %d", got.Length())
		}
	})

	t.Run("equal_offset", func(t *testing.T) {
		// Offset(1,1) on 5 items: skip A + skip E → show B,C,D.
		tr := build()
		tr.Offset(1, 1)
		got := tr.Children()
		want := []string{"B", "C", "D"}
		assertChildren(t, got, want)
	})
}

func assertChildren(t *testing.T, got tree.Children, want []string) {
	t.Helper()
	if got.Length() != len(want) {
		t.Errorf("expected %d children, got %d", len(want), got.Length())
		return
	}
	for i, w := range want {
		if got.At(i).Value() != w {
			t.Errorf("expected child %d to be %q, got %q", i, w, got.At(i).Value())
		}
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestAddItemWithAndWithoutRoot(t *testing.T) {
	t.Run("with root", func(t *testing.T) {
		t1 := tree.New().
			Child(
				"Foo",
				"Bar",
				tree.New().
					Child("Baz"),
				"Qux",
			)
		golden.RequireEqual(t, []byte(t1.String()))
	})

	t.Run("without root", func(t *testing.T) {
		t2 := tree.New().
			Child(
				"Foo",
				tree.New().
					Root("Bar").
					Child("Baz"),
				"Qux",
			)
		golden.RequireEqual(t, []byte(t2.String()))
	})
}

func TestEmbedListWithinTree(t *testing.T) {
	t1 := tree.New().
		Child(list.New("A", "B", "C").
			Enumerator(list.Arabic)).
		Child(list.New("1", "2", "3").
			Enumerator(list.Alphabet))

	golden.RequireEqual(t, []byte(t1.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
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

	golden.RequireEqual(t, []byte(tree.String()))
}

func TestTypes(t *testing.T) {
	tree := tree.New().
		Child(0).
		Child(true).
		Child([]any{"Foo", "Bar"}).
		Child([]string{"Qux", "Quux", "Quuux"})

	golden.RequireEqual(t, []byte(tree.String()))
}
