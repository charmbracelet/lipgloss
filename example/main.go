package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	columnWidth = 30

	historyA = "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum—Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
	historyB = "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."
	historyC = "In 1524, Henry VIII, King of England, received a “box of marmalade” from Mr. Hull of Exeter. This was probably marmelada, a solid quince paste from Portugal, still made and sold in southern Europe today. It became a favourite treat of Anne Boleyn and her ladies in waiting."
)

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	accent = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().
		Foreground(accent).
		Render

	// Tabs

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
		BorderForegroundColor(lipgloss.Color("#874BFD")).
		Padding(0, 1)

	activeTab = tab.Copy().
			Border(activeTabBorder, true)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	// Title

	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")

	descStyle = lipgloss.NewStyle().
			MarginTop(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForegroundColor(subtle)

	// Dialog

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	// List

	list = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForegroundColor(subtle).
		Height(8)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForegroundColor(subtle).
			MarginRight(2).
			Render

	listItem = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(accent).
			PaddingRight(1).
			String()

	listDone = lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "#696969"}).
			Render

	// Paragraphs

	historyStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#7E51D6")).
			Background(lipgloss.AdaptiveColor{Light: "#888B7E", Dark: "#F6FFFD"}).
			Margin(1, 3, 0, 0).
			Padding(1, 2).
			Height(18).
			Width(columnWidth)

	helloStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForegroundColor(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 0, 2)
)

func main() {
	physicalWidth, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	width := min(94, max(120, physicalWidth))

	var doc strings.Builder

	// Tabs
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Lip Gloss"),
			tab.Render("Blush"),
			tab.Render("Eye Shadow"),
			tab.Render("Mascara"),
			tab.Render("Foundation"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row))))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	// Title
	{
		var title strings.Builder

		colors := colorGrid(1, 5)

		for i, v := range colors {
			const offset = 2
			c := lipgloss.Color(v[0])
			fmt.Fprint(&title, titleStyle.Copy().MarginLeft(i*offset).Background(c))
			if i < len(colors)-1 {
				title.WriteRune('\n')
			}
		}

		desc := lipgloss.JoinVertical(lipgloss.Left,
			descStyle.Render("Style Definitions for Terminal-Based Layouts"),
			"From Charm"+divider+url("https://github.com/charmbracelet/lipgloss"),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, title.String(), desc)
		doc.WriteString(row + "\n\n")
	}

	// Dialog
	var dialog string
	{
		okButton := activeButtonStyle.Render("Yes")
		cancelButton := buttonStyle.Render("Maybe")

		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		dialog = lipgloss.JoinVertical(lipgloss.Center, question, buttons)
		dialog = lipgloss.Place(width, 9,
			lipgloss.Center, lipgloss.Center,
			helloStyle.Render(ui),
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceForeground(subtle),
		)

	}
	doc.WriteString(dialog + "\n\n")

	// Color grid
	grid := func() string {
		x0y0, _ := colorful.Hex("#F25D94")
		x1y0, _ := colorful.Hex("#EDFF82")
		x0y1, _ := colorful.Hex("#874BFD")
		x1y1, _ := colorful.Hex("#14F9D5")

		const xSteps = 10
		const ySteps = 8

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

		b := strings.Builder{}
		for _, x := range grid {
			for _, y := range x {
				s := lipgloss.NewStyle().SetString("  ").Background(lipgloss.Color(y))
				b.WriteString(s.String())
			}
			b.WriteRune('\n')
		}

		return b.String()
	}()

	fruits := lipgloss.JoinVertical(lipgloss.Left,
		listHeader("Citrus Fruits to Try"),
		checkMark+listDone("Quince"),
		checkMark+listDone("Yuzu"),
		listItem("Citron"),
		listItem("Kumquat"),
		listItem("Pomelo"),
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, fruits, grid))

	doc.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		historyStyle.Copy().Align(lipgloss.Right).Render(historyA),
		historyStyle.Copy().Align(lipgloss.Center).Render(historyB),
		historyStyle.Copy().MarginRight(0).Render(historyC),
	))
	doc.WriteString("\n\n")

	fmt.Println(docStyle.MaxWidth(physicalWidth).Render(doc.String()))
}

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#874BFD")
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
