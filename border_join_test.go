package lipgloss

import "testing"

func TestJoinBordersLeft(t *testing.T) {
	b := JoinBordersLeft(NormalBorder())
	if b.TopRight != "┬" {
		t.Errorf("expected TopRight to be ┬, got %s", b.TopRight)
	}
	if b.BottomRight != "┴" {
		t.Errorf("expected BottomRight to be ┴, got %s", b.BottomRight)
	}
	if b.Right != "┼" {
		t.Errorf("expected Right to be ┼, got %s", b.Right)
	}
	// Left side should be unchanged.
	if b.TopLeft != "┌" {
		t.Errorf("expected TopLeft to be ┌, got %s", b.TopLeft)
	}
}

func TestJoinBordersRight(t *testing.T) {
	b := JoinBordersRight(NormalBorder())
	if b.TopLeft != "" {
		t.Errorf("expected TopLeft to be empty, got %s", b.TopLeft)
	}
	if b.BottomLeft != "" {
		t.Errorf("expected BottomLeft to be empty, got %s", b.BottomLeft)
	}
	if b.Left != "" {
		t.Errorf("expected Left to be empty, got %s", b.Left)
	}
	// Right side should be unchanged.
	if b.TopRight != "┐" {
		t.Errorf("expected TopRight to be ┐, got %s", b.TopRight)
	}
}

func TestJoinBordersTop(t *testing.T) {
	b := JoinBordersTop(NormalBorder())
	if b.BottomLeft != "├" {
		t.Errorf("expected BottomLeft to be ├, got %s", b.BottomLeft)
	}
	if b.BottomRight != "┤" {
		t.Errorf("expected BottomRight to be ┤, got %s", b.BottomRight)
	}
	if b.Bottom != "┼" {
		t.Errorf("expected Bottom to be ┼, got %s", b.Bottom)
	}
}

func TestJoinBordersBottom(t *testing.T) {
	b := JoinBordersBottom(NormalBorder())
	if b.TopLeft != "" {
		t.Errorf("expected TopLeft to be empty, got %s", b.TopLeft)
	}
	if b.TopRight != "" {
		t.Errorf("expected TopRight to be empty, got %s", b.TopRight)
	}
	if b.Top != "" {
		t.Errorf("expected Top to be empty, got %s", b.Top)
	}
}

func TestJoinBordersHorizontalRendering(t *testing.T) {
	left := JoinBordersLeft(NormalBorder())
	right := NormalBorder()

	leftStyle := NewStyle().Border(left).Width(6)
	rightStyle := NewStyle().Border(right).Width(6)

	leftBox := leftStyle.Render("A")
	rightBox := rightStyle.Render("B")

	result := JoinHorizontal(Top, leftBox, rightBox)

	// The result should not have doubled borders (i.e., no "┐┌" pattern).
	if containsSubstring(result, "┐┌") {
		t.Errorf("joined borders should not have doubled corners, got:\n%s", result)
	}
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
