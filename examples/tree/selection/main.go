package main

import (
	"fmt"
	"path"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/tree"
)

const selected = "/Users/bash/.config/doom-emacs"

type styles struct {
	base,
	container,
	dir,
	selected,
	dimmed,
	toggle lipgloss.Style
}

func defaultStyles() styles {
	var s styles
	s.base = lipgloss.NewStyle()
	s.container = s.base.
		Margin(1, 2).
		Padding(1, 0)
	s.dir = s.base.
		Inline(true)
	s.toggle = s.base.
		Foreground(lipgloss.Color("5")).
		PaddingRight(1)
	s.selected = s.base.
		Background(lipgloss.Color("8")).
		Foreground(lipgloss.Color("207")).
		Bold(true)
	s.dimmed = s.base.
		Foreground(lipgloss.Color("241"))
	return s
}

type dir struct {
	name   string
	open   bool
	styles styles
}

func (d dir) String() string {
	t := d.styles.toggle.PaddingLeft(1).Render
	n := d.styles.dir.Render
	if d.open {
		return t("▼") + n(d.name)
	}
	return t("▶") + n(d.name)
}

// file implements the Node interface.
type file struct {
	name   string
	styles styles
}

func (s file) String() string {
	return path.Base(s.name)
}

func (s file) Hidden() bool {
	return false
}

func (s file) Children() tree.Children {
	return tree.NodeChildren(nil)
}

func (s file) Value() string {
	return s.String()
}

func (s file) SetValue(val any) {
	return
}

func (s file) SetHidden(val bool) {
	return
}

func isItemSelected(children tree.Children, index int) bool {
	child := children.At(index)
	if file, ok := child.(file); ok && file.name == selected {
		return true
	}

	return false
}

func itemStyle(children tree.Children, index int) lipgloss.Style {
	s := defaultStyles()
	if isItemSelected(children, index) {
		return s.selected
	}

	return s.base
}

func indenterStyle(children tree.Children, index int) lipgloss.Style {
	s := defaultStyles()
	if isItemSelected(children, index) {
		return s.dimmed.Background(s.selected.GetBackground())
	}

	return s.dimmed
}

func main() {
	s := defaultStyles()

	t := tree.Root(dir{"~/charm", true, s}).
		Child(
			dir{"ayman", false, s},
			tree.Root(dir{"bash", true, s}).
				Child(
					file{"/Users/bash/.config/doom-emacs", s},
				),
			tree.Root(dir{"carlos", true, s}).
				Child(
					tree.Root(dir{"emotes", true, s}).
						Child(
							file{"/home/caarlos0/Pictures/chefkiss.png", s},
							file{"/home/caarlos0/Pictures/kekw.png", s},
						),
				),
			dir{"maas", false, s},
		).
		Width(30).
		Indenter(Indenter).
		Enumerator(Enumerator).
		EnumeratorStyleFunc(indenterStyle).
		IndenterStyleFunc(indenterStyle).
		ItemStyleFunc(itemStyle)

	fmt.Println(s.container.Render(t.String()))
}

func Enumerator(children tree.Children, index int) string {
	return " │ "
}

func Indenter(children tree.Children, index int) string {
	return " │ "
}
