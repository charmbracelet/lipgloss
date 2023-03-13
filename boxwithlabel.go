package lipgloss

import "strings"

// BoxWithLabel defines box with a label on the border.
//
//	box := BoxWithLabel{
//		BoxStyle: NewStyle().
//			Border(RoundedBorder()).
//			BorderForeground(Color("63")).
//			Padding(1),
//
//		LabelStyle: NewStyle().
//			PaddingTop(0).
//			PaddingBottom(0).
//			PaddingLeft(1).
//			PaddingRight(1),
//	}
//
//	box.Render("Label", "Box content...", 20)
type BoxWithLabel struct {
	BoxStyle   Style
	LabelStyle Style
}

// NewDefaultBoxWithLabel creates a new BoxWithLabel with default styling. You
// can continue to customize the BoxWithLabel returned.
func NewDefaultBoxWithLabel() BoxWithLabel {
	return BoxWithLabel{
		BoxStyle: NewStyle().
			Border(RoundedBorder()).
			BorderForeground(Color("63")).
			Padding(1),

		// You could, of course, also set background and foreground colors here
		// as well.
		LabelStyle: NewStyle().
			PaddingTop(0).
			PaddingBottom(0).
			PaddingLeft(1).
			PaddingRight(1),
	}
}

// Render renders the labeled box.
func (b BoxWithLabel) Render(label, content string, width int) string {
	var (
		// Query the box style for some of its border properties so we can
		// essentially take the top border apart and put it around the label.
		border          Border                 = b.BoxStyle.GetBorderStyle()
		topBorderStyler func(...string) string = NewStyle().Foreground(b.BoxStyle.GetBorderTopForeground()).Render
		topLeft         string                 = topBorderStyler(border.TopLeft)
		topRight        string                 = topBorderStyler(border.TopRight)

		renderedLabel string = b.LabelStyle.Render(label)
	)

	// Render top row with the label
	borderWidth := b.BoxStyle.GetHorizontalBorderSize()
	cellsShort := max(0, width+borderWidth-Width(topLeft+topRight+renderedLabel))
	gap := strings.Repeat(border.Top, cellsShort)
	top := topLeft + renderedLabel + topBorderStyler(gap) + topRight

	// Render the rest of the box
	bottom := b.BoxStyle.Copy().
		BorderTop(false).
		Width(width).
		Render(content)

	// Stack the pieces
	return top + "\n" + bottom
}
