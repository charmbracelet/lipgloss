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
