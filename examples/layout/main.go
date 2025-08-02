package main

// This example demonstrates various Lip Gloss style and layout features.

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/uniseg"
)

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96

	// How wide to render various columns in the layout.
	columnWidth = 30
)

var (
	// Whether the detected background color is dark. We detect this in init().
	hasDarkBG bool

	// A helper function for choosing either a light or dark color based on the
	// detected background color. We create this in init().
	lightDark lipgloss.LightDarkFunc
)

func init() {
	// Detect the background color.
	hasDarkBG = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)

	// Create a new helper function for choosing either a light or dark color
	// based on the detected background color.
	lightDark = lipgloss.LightDark(hasDarkBG)
}

func main() {
	// Style definitions.
	var (

		// General.

		subtle    = lightDark(lipgloss.Color("#D9DCCF"), lipgloss.Color("#383838"))
		highlight = lightDark(lipgloss.Color("#874BFD"), lipgloss.Color("#7D56F4"))
		special   = lightDark(lipgloss.Color("#43BF6D"), lipgloss.Color("#73F59F"))

		divider = lipgloss.NewStyle().
			SetString("•").
			Padding(0, 1).
			Foreground(subtle).
			String()

		url = lipgloss.NewStyle().Foreground(special).Render

		// Tabs.

		activeTabBorder = lipgloss.Border{
			Top:         "─",
			Bottom:      " ",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┘",
			BottomRight: "└",
		}

		tabBorder = lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┴",
			BottomRight: "┴",
		}

		tab = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(highlight).
			Padding(0, 1)

		activeTab = tab.Border(activeTabBorder, true)

		tabGap = tab.
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)

		// Title.

		titleStyle = lipgloss.NewStyle().
				MarginLeft(1).
				MarginRight(5).
				Padding(0, 1).
				Italic(true).
				Foreground(lipgloss.Color("#FFF7DB")).
				SetString("Lip Gloss")

		descStyle = lipgloss.NewStyle().MarginTop(1)

		infoStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(true).
				BorderForeground(subtle)

		// Dialog.

		dialogBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#874BFD")).
				Padding(1, 0).
				BorderTop(true).
				BorderLeft(true).
				BorderRight(true).
				BorderBottom(true)

		buttonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#888B7E")).
				Padding(0, 3).
				MarginTop(1)

		activeButtonStyle = buttonStyle.
					Foreground(lipgloss.Color("#FFF7DB")).
					Background(lipgloss.Color("#F25D94")).
					MarginRight(2).
					Underline(true)

		// List.

		list = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(subtle).
			MarginRight(1).
			Height(8).
			Width(width / 3)

		listHeader = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(subtle).
				MarginRight(2).
				Render

		listItem = lipgloss.NewStyle().PaddingLeft(2).Render

		checkMark = lipgloss.NewStyle().SetString("✓").
				Foreground(special).
				PaddingRight(1).
				String()

		listDone = func(s string) string {
			return checkMark + lipgloss.NewStyle().
				Strikethrough(true).
				Foreground(lightDark(lipgloss.Color("#969B86"), lipgloss.Color("#696969"))).
				Render(s)
		}

		// Paragraphs/History.

		historyStyle = lipgloss.NewStyle().
				Align(lipgloss.Left).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(highlight).
				Margin(1, 3, 0, 0).
				Padding(1, 2).
				Height(19).
				Width(columnWidth)

		// Status Bar.

		statusNugget = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFDF5")).
				Padding(0, 1)

		statusBarStyle = lipgloss.NewStyle().
				Foreground(lightDark(lipgloss.Color("#343433"), lipgloss.Color("#C1C6B2"))).
				Background(lightDark(lipgloss.Color("#D9DCCF"), lipgloss.Color("#353533")))

		statusStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#FF5F87")).
				Padding(0, 1).
				MarginRight(1)

		encodingStyle = statusNugget.
				Background(lipgloss.Color("#A550DF")).
				Align(lipgloss.Right)

		statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

		fishCakeStyle = statusNugget.Background(lipgloss.Color("#6124DF"))

		// Floating thing.

		floatingStyle = lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				Padding(1, 6).
				Align(lipgloss.Center)

		// Page.

		docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	)

	physicalWidth, _, _ := term.GetSize(os.Stdout.Fd())
	doc := strings.Builder{}

	// Tabs.
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Lip Gloss"),
			tab.Render("Blush"),
			tab.Render("Eye Shadow"),
			tab.Render("Mascara"),
			tab.Render("Foundation"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	// Title.
	{
		var (
			colors = colorGrid(1, 5)
			title  strings.Builder
		)

		for i, v := range colors {
			const offset = 2
			c := lipgloss.Color(v[0])
			fmt.Fprint(&title, titleStyle.MarginLeft(i*offset).Background(c))
			if i < len(colors)-1 {
				title.WriteRune('\n')
			}
		}

		desc := lipgloss.JoinVertical(lipgloss.Left,
			descStyle.Render("Style Definitions for Nice Terminal Layouts"),
			infoStyle.Render("From Charm"+divider+url("https://github.com/charmbracelet/lipgloss")),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, title.String(), desc)
		doc.WriteString(row + "\n\n")
	}

	// Dialog.
	{
		okButton := activeButtonStyle.Render("Yes")
		cancelButton := buttonStyle.Render("Maybe")

		grad := applyGradient(
			lipgloss.NewStyle(),
			"Are you sure you want to eat marmalade?",
			lipgloss.Color("#EDFF82"),
			lipgloss.Color("#F25D94"),
		)

		question := lipgloss.NewStyle().
			Width(50).
			Align(lipgloss.Center).
			Render(grad)

		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		dialog := lipgloss.Place(width, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Foreground(subtle)),
		)

		doc.WriteString(dialog + "\n\n")
	}

	// Color grid.
	colors := func() string {
		colors := colorGrid(14, 8)

		b := strings.Builder{}
		for _, x := range colors {
			for _, y := range x {
				s := lipgloss.NewStyle().SetString("  ").Background(lipgloss.Color(y))
				b.WriteString(s.String())
			}
			b.WriteRune('\n')
		}

		return b.String()
	}()

	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		list.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Citrus Fruits to Try"),
				listDone("Grapefruit"),
				listDone("Yuzu"),
				listItem("Citron"),
				listItem("Kumquat"),
				listItem("Pomelo"),
			),
		),
		list.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Actual Lip Gloss Vendors"),
				listItem("Glossier"),
				listItem("Claire‘s Boutique"),
				listDone("Nyx"),
				listItem("Mac"),
				listDone("Milk"),
			),
		),
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists, lipgloss.NewStyle().MarginLeft(1).Render(colors)))

	// Marmalade history.
	{
		const (
			historyA = "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
			historyB = "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."
			historyC = "In 1524, Henry VIII, King of England, received a “box of marmalade” from Mr. Hull of Exeter. This was probably marmelada, a solid quince paste from Portugal, still made and sold in southern Europe today. It became a favourite treat of Anne Boleyn and her ladies in waiting."
		)

		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			historyStyle.Align(lipgloss.Right).Render(historyA),
			historyStyle.Align(lipgloss.Center).Render(historyB),
			historyStyle.MarginRight(0).Render(historyC),
		))

		doc.WriteString("\n\n")
	}

	// Status bar.
	{
		w := lipgloss.Width

		lightDarkState := "Light"
		if hasDarkBG {
			lightDarkState = "Dark"
		}

		statusKey := statusStyle.Render("STATUS")
		encoding := encodingStyle.Render("UTF-8")
		fishCake := fishCakeStyle.Render("🍥 Fish Cake")
		statusVal := statusText.
			Width(width - w(statusKey) - w(encoding) - w(fishCake)).
			Render("Ravishingly " + lightDarkState + "!")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			statusKey,
			statusVal,
			encoding,
			fishCake,
		)

		doc.WriteString(statusBarStyle.Width(width).Render(bar))
	}

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	// Render the document.
	document := docStyle.Render(doc.String())

	// Surprise! Composite some bonus content on top of the document.
	modal := floatingStyle.Render("Now with Compositing!")
	canvas := lipgloss.NewCanvas(
		lipgloss.NewLayer(document),
		lipgloss.NewLayer(modal).X(58).Y(44),
	)

	// Okay, let's print it. We use a special Lipgloss writer to downsample
	// colors to the terminal's color palette. And, if output's not a TTY, we
	// will remove color entirely.
	lipgloss.Println(canvas.Render())
}

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}

// applyGradient applies a gradient to the given string string.
func applyGradient(base lipgloss.Style, input string, from, to color.Color) string {
	// We want to get the graphemes of the input string, which is the number of
	// characters as a human would see them.
	//
	// We definitely don't want to use len(), because that returns the
	// bytes. The rune count would get us closer but there are times, like with
	// emojis, where the rune count is greater than the number of actual
	// characters.
	g := uniseg.NewGraphemes(input)
	var chars []string
	for g.Next() {
		chars = append(chars, g.Str())
	}

	// Genrate the blend.
	a, _ := colorful.MakeColor(to)
	b, _ := colorful.MakeColor(from)
	var output strings.Builder
	var hex string
	for i := 0; i < len(chars); i++ {
		hex = a.BlendLuv(b, float64(i)/float64(len(chars)-1)).Hex()
		output.WriteString(base.Foreground(lipgloss.Color(hex)).Render(chars[i]))
	}

	return output.String()
}
