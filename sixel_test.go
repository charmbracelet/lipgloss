package lipgloss

import (
	"image"
	"image/color"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestFullImage(t *testing.T) {
	colorDrawRegex, err := regexp.Compile("^#(\\d+)(.+)$")
	if err != nil {
		t.Errorf("failed to compile regex: %+v", err)
		return
	}

	testCases := map[string]struct {
		imageWidth  int
		imageHeight int
		bandCount   int
		// When filling the image, we'll use a map of indices to colors and change colors every
		// time the current index is in the map- this will prevent dozens of lines with the same color
		// in a row and make this slightly more legible
		colors map[int]color.RGBA
		// Two bands, with a map of color to pixel strings, since we don't have any control over what
		// order the colors will appear
		pixels []map[sixelColor]string
	}{
		"3x12 single color filled": {
			3, 12, 2,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
			},
			[]map[sixelColor]string{
				{
					// ~ means all 6 pixels filled
					{100, 0, 0, 100}: "~~~",
				},
				{
					{100, 0, 0, 100}: "~~~",
				},
			},
		},
		"3x12 two color filled": {
			3, 12, 2,
			map[int]color.RGBA{
				// 3-pixel high alternating bands
				0:  {0, 0, 255, 255},
				9:  {0, 255, 0, 255},
				18: {0, 0, 255, 255},
				27: {0, 255, 0, 255},
			},
			[]map[sixelColor]string{
				{
					// F means top 3 pixels filled, w means bottom 3 pixels filled
					{0, 0, 100, 100}: "FFF",
					{0, 100, 0, 100}: "www",
				},
				{
					{0, 0, 100, 100}: "FFF",
					{0, 100, 0, 100}: "www",
				},
			},
		},
		"3x12 8 color with right gutter": {
			3, 12, 2,
			map[int]color.RGBA{
				0:  {255, 0, 0, 255},
				2:  {0, 255, 0, 255},
				3:  {255, 0, 0, 255},
				5:  {0, 255, 0, 255},
				6:  {255, 0, 0, 255},
				8:  {0, 255, 0, 255},
				9:  {0, 0, 255, 255},
				11: {128, 128, 0, 255},
				12: {0, 0, 255, 255},
				14: {128, 128, 0, 255},
				15: {0, 0, 255, 255},
				17: {128, 128, 0, 255},
				18: {0, 128, 128, 255},
				20: {128, 0, 128, 255},
				21: {0, 128, 128, 255},
				23: {128, 0, 128, 255},
				24: {0, 128, 128, 255},
				26: {128, 0, 128, 255},
				27: {64, 0, 0, 255},
				29: {0, 64, 0, 255},
				30: {64, 0, 0, 255},
				32: {0, 64, 0, 255},
				33: {64, 0, 0, 255},
				35: {0, 64, 0, 255},
			},
			[]map[sixelColor]string{
				{
					// F means top 3 pixels filled, w means bottom 3 pixels filled
					// ? means no pixels filled
					{100, 0, 0, 100}: "FF?",
					{0, 100, 0, 100}: "??F",
					{0, 0, 100, 100}: "ww?",
					{50, 50, 0, 100}: "??w",
				},
				{
					{0, 50, 50, 100}: "FF?",
					{50, 0, 50, 100}: "??F",
					{25, 0, 0, 100}:  "ww?",
					{0, 25, 0, 100}:  "??w",
				},
			},
		},
		"3x12 single color with transparent band in the middle": {
			3, 12, 2,
			map[int]color.RGBA{
				0:  {255, 0, 0, 255},
				15: {0, 0, 0, 0},
				21: {255, 0, 0, 255},
			},
			[]map[sixelColor]string{
				{
					// ^ means all pixels except the bottom one filled
					{100, 0, 0, 100}: "^^^",
				},
				{
					// } means all pixels except the top one filled
					{100, 0, 0, 100}: "}}}",
				},
			},
		},
		"3x5 single color": {
			3, 5, 1,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
			},
			[]map[sixelColor]string{
				{
					// ^ means all pixels except the bottom one filled
					{100, 0, 0, 100}: "^^^",
				},
			},
		},
		"12x4 single color use RLE": {
			12, 4, 1,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
			},
			[]map[sixelColor]string{
				{
					// N means the top 4 pixels filled,
					{100, 0, 0, 100}: "!12N",
				},
			},
		},
		"12x1 two color use RLE": {
			12, 1, 1,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
				6: {0, 255, 0, 255},
			},
			[]map[sixelColor]string{
				{
					// @ means just the top pixel filled, ? means no pixels filled
					{100, 0, 0, 100}: "!6@!6?",
					{0, 100, 0, 100}: "!6?!6@",
				},
			},
		},
		"12x12 single color use RLE": {
			12, 12, 2,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
			},
			[]map[sixelColor]string{
				{
					// ~ means all six pixels filled
					{100, 0, 0, 100}: "!12~",
				},
				{
					{100, 0, 0, 100}: "!12~",
				},
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			img := image.NewRGBA(image.Rect(0, 0, testCase.imageWidth, testCase.imageHeight))

			currentColor := color.RGBA{0, 0, 0, 0}
			for y := 0; y < testCase.imageHeight; y++ {
				for x := 0; x < testCase.imageWidth; x++ {
					index := y*testCase.imageWidth + x
					newColor, changingColor := testCase.colors[index]
					if changingColor {
						currentColor = newColor
					}

					img.Set(x, y, currentColor)
				}
			}

			sixelImage := Sixel(img)
			if sixelImage.PixelHeight() != testCase.imageHeight {
				t.Errorf("SixelImage had a height of %d, but a height of %d was expected", sixelImage.PixelHeight(), testCase.imageHeight)
				return
			}
			if sixelImage.PixelWidth() != testCase.imageWidth {
				t.Errorf("SixelImage had a width of %d, but a width of %d was expected", sixelImage.PixelWidth(), testCase.imageWidth)
				return
			}

			if !strings.HasSuffix(sixelImage.pixels, string(sixelLineBreak)) {
				t.Errorf("SixelImage pixels were expected to end with a linebreak character '-' but instead ended with %q", sixelImage.pixels[len(sixelImage.pixels)-1])
				return
			}

			pixelLines := strings.Split(sixelImage.pixels, string(sixelLineBreak))
			// We expect some bands of pixels, followed by a linebreak, so splitting on linebreak should produce
			// one more entry than the number of bands
			lineCount := len(pixelLines) - 1

			if lineCount != testCase.bandCount {
				t.Errorf("The SixelImage pixels are supposed to have %d bands of pixels but have %d instead", testCase.bandCount, lineCount)
				return
			}

			// Check each band of pixels in the image
			for bandY := 0; bandY < lineCount; bandY++ {
				expectedBandColors := testCase.pixels[bandY]
				colorDraws := strings.Split(pixelLines[bandY], string(sixelCarriageReturn))

				if len(colorDraws) != len(expectedBandColors) {
					t.Errorf("The SixelImage had %d colors in band %d, but was expected to have %d.", len(colorDraws), bandY, len(expectedBandColors))
					return
				}

				// Individually check each color in the band
				for _, colorPixels := range colorDraws {
					// #<palette index><pixels>
					matches := colorDrawRegex.FindStringSubmatch(colorPixels)
					if matches == nil {
						t.Errorf("Could not locate color palette change in color draw substring %q", colorPixels)
						return
					}

					paletteIndex, err := strconv.Atoi(matches[1])
					if err != nil {
						t.Errorf("failed to convert %q to a palette index", matches[1])
						return
					}

					if paletteIndex >= len(sixelImage.palette.PaletteColors) {
						t.Errorf("Found palette index %d in the pixel buffer, but the palette only has %d colors.", paletteIndex, len(sixelImage.palette.PaletteColors))
						return
					}

					foundColor := sixelImage.palette.PaletteColors[paletteIndex]
					expectedPixelString, hasColor := expectedBandColors[foundColor]
					if !hasColor {
						t.Errorf("Found palette index %d (color %+v) with pixels %q in the pixel buffer, but it was not expected on band %d.", paletteIndex, foundColor, matches[2], bandY)
						return
					}

					if expectedPixelString != matches[2] {
						t.Errorf("Palette index %d (color %+v) had pixel buffer %q, but expected %q", paletteIndex, foundColor, matches[2], expectedPixelString)
						return
					}
				}
			}
		})
	}
}
