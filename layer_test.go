package lipgloss

import (
	"fmt"
	"testing"
)

func TestCompositorEqualZKeepsInsertionOrder(t *testing.T) {
	// Layers that share a z-index are drawn in the order they were added, so
	// the last one ends up on top no matter how many layers there are.
	for _, n := range []int{2, 11, 12, 20, 100} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			// A raised layer away from the stack, followed by n overlapping
			// layers at the default z-index.
			layers := []*Layer{NewLayer("top").ID("top").X(10).Z(1)}
			for i := range n {
				layers = append(layers, NewLayer(fmt.Sprintf("%03d", i)).ID(fmt.Sprint(i)))
			}
			comp := NewCompositor(layers...)

			want := fmt.Sprintf("%03d       top", n-1)
			if got := comp.Render(); got != want {
				t.Errorf("Render() = %q, want %q", got, want)
			}
			if got, want := comp.Hit(0, 0).ID(), fmt.Sprint(n-1); got != want {
				t.Errorf("Hit(0, 0).ID() = %q, want %q", got, want)
			}
		})
	}
}
