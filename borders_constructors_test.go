package lipgloss

import "testing"

// Each exported border constructor should return its corresponding predefined
// border value. These are trivial accessors, but pinning them guards against a
// constructor accidentally being wired to the wrong border.
func TestBorderConstructors(t *testing.T) {
	tests := []struct {
		name string
		got  Border
		want Border
	}{
		{"RoundedBorder", RoundedBorder(), roundedBorder},
		{"BlockBorder", BlockBorder(), blockBorder},
		{"OuterHalfBlockBorder", OuterHalfBlockBorder(), outerHalfBlockBorder},
		{"InnerHalfBlockBorder", InnerHalfBlockBorder(), innerHalfBlockBorder},
		{"ThickBorder", ThickBorder(), thickBorder},
		{"DoubleBorder", DoubleBorder(), doubleBorder},
		{"HiddenBorder", HiddenBorder(), hiddenBorder},
		{"MarkdownBorder", MarkdownBorder(), markdownBorder},
		{"ASCIIBorder", ASCIIBorder(), asciiBorder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s() = %+v, want %+v", tt.name, tt.got, tt.want)
			}
		})
	}
}
