package main

import (
	"os"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/list"
	"charm.land/lipgloss/v2/table"
)

func main() {
	hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(hasDarkBG)

	purple := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		MarginRight(1)

	pink := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		MarginRight(1)

	base := lipgloss.NewStyle().
		MarginBottom(1).
		MarginLeft(1)

	faint := lipgloss.NewStyle().Faint(true)

	dim := lipgloss.Color("250")
	highlight := lipgloss.Color("#EE6FF8")

	special := lightDark(lipgloss.Color("#43BF6D"), lipgloss.Color("#73F59F"))

	checklistEnumStyle := func(items list.Items, index int) lipgloss.Style {
		switch index {
		case 1, 2, 4:
			return lipgloss.NewStyle().
				Foreground(special).
				PaddingRight(1)
		default:
			return lipgloss.NewStyle().PaddingRight(1)
		}
	}

	checklistEnum := func(items list.Items, index int) string {
		switch index {
		case 1, 2, 4:
			return "✓"
		default:
			return "•"
		}
	}

	checklistStyle := func(items list.Items, index int) lipgloss.Style {
		switch index {
		case 1, 2, 4:
			return lipgloss.NewStyle().
				Strikethrough(true).
				Foreground(lightDark(lipgloss.Color("#969B86"), lipgloss.Color("#696969")))
		default:
			return lipgloss.NewStyle()
		}
	}

	gradient := lipgloss.Blend1D(5, lipgloss.Color("#F25D94"), lipgloss.Color("#643AFF"))

	titleStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#FFF7DB"))

	lipglossStyleFunc := func(items list.Items, index int) lipgloss.Style {
		if index == items.Length()-1 {
			return titleStyle.Padding(1, 2).Margin(0, 0, 1, 0).MaxWidth(20).Background(gradient[index])
		}
		return titleStyle.Padding(0, 5-index, 0, index+2).MaxWidth(20).Background(gradient[index])
	}

	history := "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."

	l := list.New().
		EnumeratorStyle(purple).
		Item("Lip Gloss").
		Item("Blush").
		Item("Eye Shadow").
		Item("Mascara").
		Item("Foundation").
		Item(
			list.New().
				EnumeratorStyle(pink).
				Item("Citrus Fruits to Try").
				Item(
					list.New().
						ItemStyleFunc(checklistStyle).
						EnumeratorStyleFunc(checklistEnumStyle).
						Enumerator(checklistEnum).
						Item("Grapefruit").
						Item("Yuzu").
						Item("Citron").
						Item("Kumquat").
						Item("Pomelo"),
				).
				Item("Actual Lip Gloss Vendors").
				Item(
					list.New().
						ItemStyleFunc(checklistStyle).
						EnumeratorStyleFunc(checklistEnumStyle).
						Enumerator(checklistEnum).
						Item("Glossier").
						Item("Claire‘s Boutique").
						Item("Nyx").
						Item("Mac").
						Item("Milk").
						Item(
							list.New().
								EnumeratorStyle(purple).
								Enumerator(list.Dash).
								ItemStyleFunc(lipglossStyleFunc).
								Item("Lip Gloss").
								Item("Lip Gloss").
								Item("Lip Gloss").
								Item("Lip Gloss").
								Item(
									list.New().
										EnumeratorStyle(lipgloss.NewStyle().Foreground(gradient[4]).MarginRight(1)).
										Item("\nStyle Definitions for Nice Terminal Layouts\n─────").
										Item("From Charm").
										Item("https://github.com/charmbracelet/lipgloss").
										Item(
											list.New().
												EnumeratorStyle(lipgloss.NewStyle().Foreground(gradient[3]).MarginRight(1)).
												Item("Emperors: Julio-Claudian dynasty").
												Item(
													lipgloss.NewStyle().Padding(1).Render(
														list.New(
															"Augustus",
															"Tiberius",
															"Caligula",
															"Claudius",
															"Nero",
														).Enumerator(list.Roman).String(),
													),
												).
												Item(
													lipgloss.NewStyle().
														Bold(true).
														Foreground(lipgloss.Color("#FAFAFA")).
														Background(lipgloss.Color("#7D56F4")).
														AlignHorizontal(lipgloss.Center).
														AlignVertical(lipgloss.Center).
														Padding(1, 3).
														Margin(0, 1, 1, 1).
														Width(40).
														Render(history),
												).
												Item(
													table.New().
														Width(30).
														BorderStyle(purple.MarginRight(0)).
														StyleFunc(func(row, col int) lipgloss.Style {
															style := lipgloss.NewStyle()

															if col == 0 {
																style = style.Align(lipgloss.Center)
															} else {
																style = style.Align(lipgloss.Right).PaddingRight(2)
															}
															if row == 0 {
																return style.Bold(true).Align(lipgloss.Center).PaddingRight(0)
															}
															return style.Faint(true)
														}).
														Headers("ITEM", "QUANTITY").
														Row("Apple", "6").
														Row("Banana", "10").
														Row("Orange", "2").
														Row("Strawberry", "12"),
												).
												Item("Documents").
												Item(
													list.New().
														Enumerator(func(_ list.Items, i int) string {
															if i == 1 {
																return "│\n│"
															}
															return " "
														}).
														ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
															if i == 1 {
																return base.Foreground(highlight)
															}
															return base.Foreground(dim)
														}).
														EnumeratorStyleFunc(func(_ list.Items, i int) lipgloss.Style {
															if i == 1 {
																return lipgloss.NewStyle().Foreground(highlight)
															}
															return lipgloss.NewStyle().Foreground(dim)
														}).
														Item("Foo Document\n" + faint.Render("1 day ago")).
														Item("Bar Document\n" + faint.Render("2 days ago")).
														Item("Baz Document\n" + faint.Render("10 minutes ago")).
														Item("Qux Document\n" + faint.Render("1 month ago")),
												).
												Item("EOF"),
										).
										Item("go get github.com/charmbracelet/lipgloss/list\n"),
								).
								Item("See ya later"),
						),
				).
				Item("List"),
		).
		Item("xoxo, Charm_™")

	lipgloss.Println(l)
}
