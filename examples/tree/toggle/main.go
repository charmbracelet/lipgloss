package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type styles struct {
	base,
	block,
	enumerator,
	dir,
	toggle,
	file lipgloss.Style
}

func defaultStyles() styles {
	var s styles
	s.base = lipgloss.NewStyle().
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("225"))
	s.block = s.base.
		Padding(1, 3).
		Margin(1, 3).
		Width(40)
	s.enumerator = s.base.
		Foreground(lipgloss.Color("212")).
		PaddingRight(1)
	s.dir = s.base.
		Inline(true)
	s.toggle = s.base.
		Foreground(lipgloss.Color("207")).
		PaddingRight(1)
	s.file = s.base
	return s
}

type dir struct {
	name   string
	open   bool
	styles styles
}

func (d dir) String() string {
	t := d.styles.toggle.Render
	n := d.styles.dir.Render
	if d.open {
		return t("▼") + n(d.name)
	}
	return t("▶") + n(d.name)
}

type file struct {
	name   string
	styles styles
}

func (s file) String() string {
	return s.styles.file.Render(s.name)
}

func main() {
	s := defaultStyles()

	t := tree.New().
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(s.enumerator).
		Root(dir{"~", true, s}).
		Child(
			dir{"ayman", false, s},
			dir{"bash", false, s},
			dir{"carlos", true, s},
			tree.New().
				Root(dir{"emotes", true, s}).
				Child(
					file{"chefkiss.png", s},
					file{"kekw.png", s},
				),
			dir{"maas", false, s},
		)

	fmt.Println(s.block.Render(t.String()))
}
