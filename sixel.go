package lipgloss

import (
	"github.com/bits-and-blooms/bitset"
	"image"
	"image/color"
	"strconv"
	"strings"
)

const (
	sixelLineBreak      = '-'
	sixelCarriageReturn = '$'
	sixelRepeat         = '!'
	sixelUseColor       = '#'
)

// SixelImage is a processed and ready-to-render image. A sixel escape string can be
// obtained with Style.RenderSixelImage.  A SixelImage can be obtained with the Sixel
// method.
type SixelImage struct {
	pixelWidth  int
	pixelHeight int
	palette     sixelPalette
	pixels      string
}

// Width gets the width of the image in pixels
func (i SixelImage) Width() int {
	return i.pixelWidth
}

// Height gets the height of the image in pixels
func (i SixelImage) Height() int {
	return i.pixelHeight
}

// sixelBuilder is a temporary structure used to create a SixelImage. It handles
// breaking pixels out into bits, and then encoding them into a sixel data string. RLE
// handling is included.
//
// Making use of a sixelBuilder is done in two phases.  First, SetColor is used to write all
// pixels to the internal BitSet data.  Then, GeneratePixels is called to retrieve a string
// representing the pixel data encoded in the sixel format.
type sixelBuilder struct {
	SixelPalette sixelPalette

	imageHeight int
	imageWidth  int

	pixelBands []bitset.BitSet

	imageData   strings.Builder
	repeatRune  rune
	repeatCount int
}

// newSixelBuilder creates a sixelBuilder and prepares it for writing
func newSixelBuilder(width, height int, palette sixelPalette) sixelBuilder {
	scratch := sixelBuilder{
		imageWidth:   width,
		imageHeight:  height,
		SixelPalette: palette,
	}

	bandHeight := scratch.BandHeight()
	scratch.pixelBands = make([]bitset.BitSet, bandHeight)

	return scratch
}

// flushRepeats is used to actually write the current repeatRune to the imageData when
// it is about to change. This buffering is used to manage RLE in the sixelBuilder
func (s *sixelBuilder) flushRepeats() {
	if s.repeatCount == 0 {
		return
	}

	// Only write using the RLE form if it's actually providing space savings
	if s.repeatCount > 3 {
		countStr := strconv.Itoa(s.repeatCount)
		s.imageData.WriteRune(sixelRepeat)
		s.imageData.WriteString(countStr)
		s.imageData.WriteRune(s.repeatRune)
		return
	}

	for i := 0; i < s.repeatCount; i++ {
		s.imageData.WriteRune(s.repeatRune)
	}
}

// writeImageRune will write a single line of six pixels to pixel data.  The data
// doesn't get written to the imageData, it gets buffered for the purposes of RLE
func (s *sixelBuilder) writeImageRune(r rune) {
	if r == s.repeatRune {
		s.repeatCount++
		return
	}

	s.flushRepeats()
	s.repeatRune = r
	s.repeatCount = 1
}

// writeControlRune will write a special rune such as a new line or carriage return
// rune. It will call flushRepeats first, if necessary.
func (s *sixelBuilder) writeControlRune(r rune) {
	if s.repeatCount > 0 {
		s.flushRepeats()
		s.repeatCount = 0
		s.repeatRune = 0
	}

	s.imageData.WriteRune(r)
}

// BandHeight returns the number of six-pixel bands this image consists of
func (s *sixelBuilder) BandHeight() int {
	bandHeight := s.imageHeight / 6
	if s.imageHeight%6 != 0 {
		bandHeight++
	}

	return bandHeight
}

// SetColor will write a single pixel to sixelBuilder's internal bitset data to be used by
// GeneratePixels
func (s *sixelBuilder) SetColor(x int, y int, color color.Color) {
	bandY := y / 6
	paletteIndex := s.SixelPalette.ColorIndex(sixelConvertColor(color))

	bit := s.imageWidth*6*paletteIndex + (x * 6) + (y % 6)
	s.pixelBands[bandY].Set(uint(bit))
}

// GeneratePixels is used to write the pixel data to the internal imageData string builder.
// All pixels in the image must be written to the sixelBuilder using SetColor before this method is
// called. This method returns a string that represents the pixel data.  Sixel strings consist of five parts:
// ISC <header> <palette> <pixels> ST
// The header contains some arbitrary options indicating how the sixel image is to be drwan.
// The palette maps palette indices to RGB colors
// The pixels indicates which pixels are to be drawn with which palette colors.
//
// GeneratePixels only produces the <pixels> part of the string.  The rest is written by
// Style.RenderSixelImage.
func (s *sixelBuilder) GeneratePixels() string {
	s.imageData = strings.Builder{}
	bandHeight := s.BandHeight()

	for bandY := 0; bandY < bandHeight; bandY++ {
		if bandY > 0 {
			s.writeControlRune(sixelLineBreak)
		}

		hasWrittenAColor := false

		for paletteIndex := 0; paletteIndex < len(s.SixelPalette.PaletteColors); paletteIndex++ {
			if s.SixelPalette.PaletteColors[paletteIndex].Alpha < 1 {
				// Don't draw anything for purely transparent pixels
				continue
			}

			firstColorBit := uint(s.imageWidth * 6 * paletteIndex)
			nextColorBit := firstColorBit + uint(s.imageWidth*6)

			firstSetBitInBand, anySet := s.pixelBands[bandY].NextSet(firstColorBit)
			if !anySet || firstSetBitInBand >= nextColorBit {
				// Color not appearing in this row
				continue
			}

			if hasWrittenAColor {
				s.writeControlRune(sixelCarriageReturn)
			}
			hasWrittenAColor = true

			s.writeControlRune(sixelUseColor)
			s.imageData.WriteString(strconv.Itoa(paletteIndex))
			for x := 0; x < s.imageWidth; x += 4 {
				bit := firstColorBit + uint(x*6)
				word := s.pixelBands[bandY].Word(bit)

				pixel1 := rune((word & 63) + '?')
				pixel2 := rune(((word >> 6) & 63) + '?')
				pixel3 := rune(((word >> 12) & 63) + '?')
				pixel4 := rune(((word >> 18) & 63) + '?')

				s.writeImageRune(pixel1)

				if x+1 >= s.imageWidth {
					continue
				}
				s.writeImageRune(pixel2)

				if x+2 >= s.imageWidth {
					continue
				}
				s.writeImageRune(pixel3)

				if x+3 >= s.imageWidth {
					continue
				}
				s.writeImageRune(pixel4)
			}
		}
	}

	s.writeControlRune('-')
	return s.imageData.String()
}

// Sixel accepts a Go image and returns a SixelImage that can be rendered via
// Style.RenderSixelImage
func Sixel(image image.Image) SixelImage {
	imageBounds := image.Bounds()
	palette := newSixelPalette(image, sixelMaxColors)
	scratch := newSixelBuilder(imageBounds.Dx(), imageBounds.Dy(), palette)

	for y := 0; y < imageBounds.Dy(); y++ {
		for x := 0; x < imageBounds.Dx(); x++ {
			scratch.SetColor(x, y, image.At(x, y))
		}
	}

	pixels := scratch.GeneratePixels()

	return SixelImage{
		pixelWidth:  imageBounds.Dx(),
		pixelHeight: imageBounds.Dy(),
		palette:     palette,
		pixels:      pixels,
	}
}
