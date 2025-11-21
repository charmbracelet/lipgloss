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

// X sets the x-coordinate of the Layer relative to its parent.
func (l *Layer) X(x int) *Layer {
	l.x = x
	return l
}

// Y sets the y-coordinate of the Layer relative to its parent.
func (l *Layer) Y(y int) *Layer {
	l.y = y
	return l
}

// Z sets the z-index of the Layer relative to its parent.
func (l *Layer) Z(z int) *Layer {
	l.z = z
	return l
}

// GetX returns the x-coordinate of the Layer relative to its parent.
func (l *Layer) GetX() int {
	return l.x
}

// GetY returns the y-coordinate of the Layer relative to its parent.
func (l *Layer) GetY() int {
	return l.y
}

// GetZ returns the z-index of the Layer relative to its parent.
func (l *Layer) GetZ() int {
	return l.z
}

// absolutePosition calculates the absolute position by adding parent offsets.
func (l *Layer) absolutePosition(parentX, parentY, parentZ int) (x, y, z int) {
	return l.x + parentX, l.y + parentY, l.z + parentZ
}

// GetLayer returns a child [Layer] by its ID, or nil if not found.
func (l *Layer) GetLayer(id string) *Layer {
	if l.id == id {
		return l
	}
	for _, layer := range l.layers {
		found := layer.getLayerRecursive(id)
		if found != nil {
			return found
		}
	}
	return nil
}

// getLayerRecursive recursively searches for a layer by ID.
func (l *Layer) getLayerRecursive(id string) *Layer {
	if l.id == id {
		return l
	}
	for _, layer := range l.layers {
		found := layer.getLayerRecursive(id)
		if found != nil {
			return found
		}
	}
	return nil
}

// Bounds returns the bounds of the Layer as a [image.Rectangle].
func (l *Layer) Bounds() image.Rectangle {
	return l.boundsWithOffset(0, 0, 0)
}

// boundsWithOffset calculates bounds with parent offset applied.
func (l *Layer) boundsWithOffset(parentX, parentY, parentZ int) image.Rectangle {
	absX, absY, _ := l.absolutePosition(parentX, parentY, parentZ)

	width, height := Width(l.content), Height(l.content)
	this := image.Rectangle{
		Min: image.Pt(absX, absY),
		Max: image.Pt(absX+width, absY+height),
	}

	for _, layer := range l.layers {
		area := layer.boundsWithOffset(absX, absY, 0)
		this = this.Union(area)
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
	return l.hitWithOffset(x, y, 0, 0, 0)
}

// hitWithOffset recursively checks for hits with parent offset applied.
func (l *Layer) hitWithOffset(x, y, parentX, parentY, parentZ int) string {
	absX, absY, absZ := l.absolutePosition(parentX, parentY, parentZ)

	// Sort children by z-index for hit testing (highest z first)
	sortedChildren := make([]*Layer, len(l.layers))
	copy(sortedChildren, l.layers)
	slices.SortFunc(sortedChildren, func(a, b *Layer) int {
		aZ := a.z + absZ
		bZ := b.z + absZ
		return bZ - aZ
	})

	// Check children first (top-most first)
	for _, child := range sortedChildren {
		if hit := child.hitWithOffset(x, y, absX, absY, absZ); hit != "" {
			return hit
		}
	}

	// Check this layer
	bounds := image.Rectangle{
		Min: image.Pt(absX, absY),
		Max: image.Pt(absX+Width(l.content), absY+Height(l.content)),
	}
	if image.Pt(x, y).In(bounds) {
		return l.id
	}

	return ""
}

// AddLayers adds child layers to the Layer.
func (l *Layer) AddLayers(layers ...*Layer) *Layer {
	for i, layer := range layers {
		if layer == nil {
			panic(fmt.Sprintf("layer at index %d is nil", i))
		}
		l.layers = append(l.layers, layer)
	}
	return l
}

// Draw draws the [Layer] and its children onto the given [uv.Screen]. This can
// be a [Canvas].
func (l *Layer) Draw(scr uv.Screen, area image.Rectangle) {
	l.drawWithOffset(scr, area, 0, 0, 0)
}

// drawWithOffset recursively draws the layer and its children with parent offset applied.
func (l *Layer) drawWithOffset(scr uv.Screen, area image.Rectangle, parentX, parentY, parentZ int) {
	absX, absY, absZ := l.absolutePosition(parentX, parentY, parentZ)

	// Draw this layer's content at absolute position
	width, height := Width(l.content), Height(l.content)
	layerBounds := image.Rectangle{
		Min: image.Pt(absX, absY),
		Max: image.Pt(absX+width, absY+height),
	}

	// Only draw if the layer intersects with the area
	if layerBounds.Overlaps(area) {
		content := uv.NewStyledString(l.content)
		content.Draw(scr, layerBounds)
	}

	// Sort and draw children by z-index (lowest to highest)
	sortedChildren := make([]*Layer, len(l.layers))
	copy(sortedChildren, l.layers)
	slices.SortFunc(sortedChildren, func(a, b *Layer) int {
		aZ := a.z + absZ
		bZ := b.z + absZ
		return aZ - bZ
	})

	for _, child := range sortedChildren {
		child.drawWithOffset(scr, area, absX, absY, absZ)
	}
}
