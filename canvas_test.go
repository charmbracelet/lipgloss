package lipgloss

import (
	"strings"
	"testing"
)

func TestCanvasRender(t *testing.T) {
	c := NewCanvas(5, 3)

	// Fill the canvas with dots
	for y := 0; y < c.Height(); y++ {
		for x := 0; x < c.Width(); x++ {
			cell := c.CellAt(x, y)
			cell.Content = "."
		}
	}

	// Draw a rectangle
	for y := 1; y < 2; y++ {
		for x := 1; x < 4; x++ {
			cell := c.CellAt(x, y)
			cell.Content = "#"
		}
	}

	expected := strings.Join([]string{
		".....",
		".###.",
		".....",
	}, "\n")

	if rendered := c.Render(); rendered != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, rendered)
	}
}

func TestCanvasRenderWithTrailingSpaces(t *testing.T) {
	c := NewCanvas(5, 2)

	// Fill the canvas with spaces and some trailing spaces
	for y := 0; y < c.Height(); y++ {
		for x := 0; x < c.Width(); x++ {
			cell := c.CellAt(x, y)
			if x < 3 {
				cell.Content = "A"
			} else {
				cell.Content = " "
			}
		}
	}

	expected := strings.Join([]string{
		"AAA",
		"AAA",
	}, "\n")

	if rendered := c.Render(); rendered != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, rendered)
	}
}
