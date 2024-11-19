package lipgloss

import (
	"image"
	"image/color"
	"testing"
)

type testCase struct {
	maxColors       int
	expectedPalette []sixelColor
}

func TestPaletteCreationRedGreen(t *testing.T) {
	redGreen := image.NewRGBA(image.Rect(0, 0, 2, 2))
	redGreen.Set(0, 0, color.RGBA{255, 0, 0, 255})
	redGreen.Set(0, 1, color.RGBA{128, 0, 0, 255})
	redGreen.Set(1, 0, color.RGBA{0, 255, 0, 255})
	redGreen.Set(1, 1, color.RGBA{0, 128, 0, 255})

	testCases := map[string]testCase{
		"way too many colors": {
			maxColors: 16,
			expectedPalette: []sixelColor{
				{100, 0, 0, 100},
				{50, 0, 0, 100},
				{0, 100, 0, 100},
				{0, 50, 0, 100},
			},
		},
		"just the right number of colors": {
			maxColors: 4,
			expectedPalette: []sixelColor{
				{100, 0, 0, 100},
				{50, 0, 0, 100},
				{0, 100, 0, 100},
				{0, 50, 0, 100},
			},
		},
		"color reduction": {
			maxColors: 2,
			expectedPalette: []sixelColor{
				{75, 0, 0, 100},
				{0, 75, 0, 100},
			},
		},
	}

	runTests(t, redGreen, testCases)
}

func TestPaletteWithSemiTransparency(t *testing.T) {
	blueAlpha := image.NewRGBA(image.Rect(0, 0, 2, 2))
	blueAlpha.Set(0, 0, color.RGBA{0, 0, 255, 255})
	blueAlpha.Set(0, 1, color.RGBA{0, 0, 128, 255})
	blueAlpha.Set(1, 0, color.RGBA{0, 0, 255, 128})
	blueAlpha.Set(1, 1, color.RGBA{0, 0, 255, 0})

	testCases := map[string]testCase{
		"just the right number of colors": {
			maxColors: 4,
			expectedPalette: []sixelColor{
				{0, 0, 100, 100},
				{0, 0, 50, 100},
				{0, 0, 100, 50},
				{0, 0, 100, 0},
			},
		},
		"color reduction": {
			maxColors: 2,
			expectedPalette: []sixelColor{
				{0, 0, 75, 100},
				{0, 0, 100, 25},
			},
		},
	}
	runTests(t, blueAlpha, testCases)
}

func runTests(t *testing.T, img image.Image, testCases map[string]testCase) {
	for testName, test := range testCases {
		t.Run(testName, func(t *testing.T) {
			palette := newSixelPalette(img, test.maxColors)
			if len(palette.PaletteColors) != len(test.expectedPalette) {
				t.Errorf("Expected colors %+v in palette, but got %+v", test.expectedPalette, palette.PaletteColors)
				return
			}

			for _, c := range test.expectedPalette {
				var foundColor bool
				for _, paletteColor := range palette.PaletteColors {
					if paletteColor == c {
						foundColor = true
						break
					}
				}

				if !foundColor {
					t.Errorf("Expected colors %+v in palette, but got %+v", test.expectedPalette, palette.PaletteColors)
					return
				}
			}

			for lookupRawColor, lookupPaletteColor := range palette.colorConvert {
				paletteIndex, inReverseLookup := palette.paletteIndexes[lookupRawColor]
				if !inReverseLookup {
					t.Errorf("Color %+v maps to color %+v in the colorConvert map, but %+v is does not have a corresponding palette index.", lookupRawColor, lookupPaletteColor, lookupPaletteColor)
					return
				}

				if paletteIndex >= len(palette.PaletteColors) {
					t.Errorf("Image color %+v maps to palette index %d, but there are only %d palette colors.", lookupRawColor, paletteIndex, len(palette.PaletteColors))
					return
				}

				colorFromPalette := palette.PaletteColors[paletteIndex]
				if colorFromPalette != lookupPaletteColor {
					t.Errorf("Image color %+v maps to palette color %+v and palette index %d, but palette index %d is actually palette color %+v", lookupRawColor, lookupPaletteColor, paletteIndex, paletteIndex, colorFromPalette)
					return
				}
			}
		})
	}
}
