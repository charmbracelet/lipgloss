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

func TestCompositorRenderOffsetLayers(t *testing.T) {
	tests := []struct {
		name     string
		layers   []*Layer
		expected string
	}{
		{
			name:     "single layer away from the origin",
			layers:   []*Layer{NewLayer("AAAA").X(2).Y(1)},
			expected: "\n  AAAA",
		},
		{
			name: "no layer at the origin",
			layers: []*Layer{
				NewLayer("Box\nBox").X(3).Y(2),
				NewLayer("x").X(1).Y(1).Z(1),
			},
			expected: strings.Join([]string{
				"",
				" x",
				"   Box",
				"   Box",
			}, "\n"),
		},
		{
			name: "layer at the origin",
			layers: []*Layer{
				NewLayer("BB"),
				NewLayer("CC").X(4).Y(2),
			},
			expected: "BB\n\n    CC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rendered := NewCompositor(tt.layers...).Render(); rendered != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, rendered)
			}
		})
	}
}
