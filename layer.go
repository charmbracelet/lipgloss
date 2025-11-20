package lipgloss

import (
	"fmt"
	"image"
	"slices"

	uv "github.com/charmbracelet/ultraviolet"
)

// Layer holds metadata about a layer in the canvas.
type Layer struct {
	id      string
	content string
	x, y, z int
	layers  []*Layer
}

// NewLayer creates a new [Layer] with the given id and styled content.
func NewLayer(id string, content string, layers ...*Layer) *Layer {
	l := &Layer{
		id:      id,
		content: content,
	}
	l.AddLayers(layers...)
	return l
}

// GetID returns the ID of the Layer.
func (l *Layer) GetID() string {
	return l.id
}

// X sets the x-coordinate of the Layer.
func (l *Layer) X(x int) *Layer {
	l.x = x
	for _, layer := range l.layers {
		layer.x += x
	}
	return l
}

// Y sets the y-coordinate of the Layer.
func (l *Layer) Y(y int) *Layer {
	l.y = y
	for _, layer := range l.layers {
		layer.y += y
	}
	return l
}

// Z sets the z-index of the Layer.
func (l *Layer) Z(z int) *Layer {
	l.z = z
	for _, layer := range l.layers {
		layer.z += z
	}
	// Re-sort layers based on new z-index
	sort(l.layers)
	return l
}

// GetX returns the x-coordinate of the Layer.
func (l *Layer) GetX() int {
	return l.x
}

// GetY returns the y-coordinate of the Layer.
func (l *Layer) GetY() int {
	return l.y
}

// GetZ returns the z-index of the Layer.
func (l *Layer) GetZ() int {
	return l.z
}

// GetLayer returns a child [Layer] by its ID, or nil if not found.
func (l *Layer) GetLayer(id string) *Layer {
	if l.id == id {
		return l
	}
	for _, layer := range l.layers {
		if layer.id == id {
			return layer
		}
	}
	return nil
}

// Bounds returns the bounds of the Layer as a [image.Rectangle].
func (l *Layer) Bounds() image.Rectangle {
	// Calculate bounds based on child layers
	width, height := Width(l.content), Height(l.content)
	this := image.Rectangle{
		Min: image.Pt(l.x, l.y),
		Max: image.Pt(l.x+width, l.y+height),
	}
	for _, layer := range l.layers {
		area := layer.Bounds()
		this = this.Union(area)
	}

	// Adjust the size of the layer if it's negative
	if this.Min.X < 0 {
		this = this.Add(image.Pt(-this.Min.X, 0))
	}
	if this.Min.Y < 0 {
		this = this.Add(image.Pt(0, -this.Min.Y))
	}

	return this
}

// InBounds checks if the given point is within the [Layer]'s bounds.
func (l *Layer) InBounds(x, y int) bool {
	return image.Pt(x, y).In(l.Bounds())
}

// Hit checks if the given point hits the Layer or any of its child layers. If
// a hit is detected, it returns the ID of the top-most Layer that was hit. If
// no hit is detected, it returns an empty string.
func (l *Layer) Hit(x, y int) string {
	// Reverse the order of the layers so that the top-most layer is checked
	// first.
	for i := len(l.layers) - 1; i >= 0; i-- {
		if l.layers[i].InBounds(x, y) {
			return l.layers[i].Hit(x, y)
		}
	}

	if image.Pt(x, y).In(l.Bounds()) {
		return l.id
	}

	return ""
}

// AddLayers adds child layers to the Layer.
func (l *Layer) AddLayers(layers ...*Layer) *Layer {
	// Make children layers relative to parent
	for i, layer := range layers {
		if layer == nil {
			panic(fmt.Sprintf("layer at index %d is nil", i))
		}
		layer.x += l.x
		layer.y += l.y
		layer.z += l.z
		l.layers = append(l.layers, layer)
	}
	sort(l.layers)
	return l
}

// Draw draws the [Layer] and its children onto the given [uv.Screen]. This can
// be a [Canvas].
func (l *Layer) Draw(scr uv.Screen, area image.Rectangle) {
	content := uv.NewStyledString(l.content)
	content.Draw(scr, area.Intersect(l.Bounds()))
	for _, child := range l.layers {
		child.Draw(scr, area.Intersect(child.Bounds()))
	}
}

// sort sorts layers by their z-index.
func sort(layers []*Layer) {
	slices.SortFunc(layers, func(a, b *Layer) int {
		return a.z - b.z
	})
}
