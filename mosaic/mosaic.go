package mosaic

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"math"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// mosaic.New().Scale(mosaic.Fit).Width(100).Render(img)

// Blocks definition.
var (
	HalfBlocks = []Block{
		{Char: '▀', Coverage: [4]bool{true, true, false, false}, CoverageMap: "██\n  "},   // Upper half block.
		{Char: '▄', Coverage: [4]bool{false, false, true, true}, CoverageMap: "  \n██"},   // Lower half block.
		{Char: ' ', Coverage: [4]bool{false, false, false, false}, CoverageMap: "  \n  "}, // Space.
		{Char: '█', Coverage: [4]bool{true, true, true, true}, CoverageMap: "██\n██"},     // Full block.
	}
	QuarterBlocks = []Block{
		{Char: '▘', Coverage: [4]bool{true, false, false, false}, CoverageMap: "█ \n  "}, // Quadrant upper left.
		{Char: '▝', Coverage: [4]bool{false, true, false, false}, CoverageMap: " █\n  "}, // Quadrant upper right.
		{Char: '▖', Coverage: [4]bool{false, false, true, false}, CoverageMap: "  \n█ "}, // Quadrant lower left.
		{Char: '▗', Coverage: [4]bool{false, false, false, true}, CoverageMap: "  \n █"}, // Quadrant lower right.
		{Char: '▌', Coverage: [4]bool{true, false, true, false}, CoverageMap: "█ \n█ "},  // Left half block.
		{Char: '▐', Coverage: [4]bool{false, true, false, true}, CoverageMap: " █\n █"},  // Right half block.
		{Char: '▀', Coverage: [4]bool{true, true, false, false}, CoverageMap: "██\n  "},  // Upper half block (already added).
		{Char: '▄', Coverage: [4]bool{false, false, true, true}, CoverageMap: "  \n██"},  // Lower half block (already added).
	}
	ComplexBlocks = []Block{
		{Char: '▙', Coverage: [4]bool{true, false, true, true}, CoverageMap: "█ \n██"},  // Quadrant upper left and lower half.
		{Char: '▟', Coverage: [4]bool{false, true, true, true}, CoverageMap: " █\n██"},  // Quadrant upper right and lower half.
		{Char: '▛', Coverage: [4]bool{true, true, true, false}, CoverageMap: "██\n█ "},  // Quadrant upper half and lower left.
		{Char: '▜', Coverage: [4]bool{true, true, false, true}, CoverageMap: "██\n █"},  // Quadrant upper half and lower right.
		{Char: '▚', Coverage: [4]bool{true, false, false, true}, CoverageMap: "█ \n █"}, // Quadrant upper left and lower right.
		{Char: '▞', Coverage: [4]bool{false, true, true, false}, CoverageMap: " █\n█ "}, // Quadrant upper right and lower left.
	}
)

// Block represents different Unicode block characters.
type Block struct {
	Char        rune
	Coverage    [4]bool // Which parts of the block are filled (true = filled).
	CoverageMap string  // Visual representation of coverage for debugging.
}

// Symbol represents the symbol type to use when rendering the image.
type Symbol string

// Symbol types.
const (
	AllSymbols     Symbol = "all"
	HalfSymbols    Symbol = "half"
	QuarterSymbols Symbol = "quarter"
)

// Scale represents scale mode that will be used when rendering the image.
type Scale uint8

// Symbol types.
const (
	None Scale = iota
	Fit
	Stretch
	Center
)

// Render mosaic with default values
func Render(img image.Image) string {
	m := New()
	return m.Render(img)
}

// Options contains all configurable settings.
type Mosaic struct {
	blocks         []Block
	outputWidth    int    // Output width.
	outputHeight   int    // Output height (0 for auto).
	thresholdLevel uint8  // Threshold for considering a pixel as set (0-255).
	dither         bool   // Enable Dithering (false as default).
	useFgBgOnly    bool   // Use only foreground/background colors (no block symbols).
	invertColors   bool   // Invert colors.
	scaleMode      Scale  // fit, stretch, center.
	symbols        Symbol // Which symbols to use: "half", "quarter", "all".
}

// Create Mosaic
func New() Mosaic {
	return Mosaic{
		blocks:         HalfBlocks,
		outputWidth:    80,     // Default width.
		outputHeight:   0,      // Auto height.
		thresholdLevel: 128,    // Middle threshold.
		dither:         false,  // Enable dithering.
		useFgBgOnly:    false,  // Use block symbols.
		invertColors:   false,  // Don't invert.
		scaleMode:      None,   // Fit to terminal.
		symbols:        "half", // Use half blocks.
	}
}

// PixelBlock represents a 2x2 pixel block from the image.
type pixelBlock struct {
	Pixels      [2][2]color.Color // 2x2 pixel grid.
	AvgFg       color.Color       // Average foreground color.
	AvgBg       color.Color       // Average background color.
	BestSymbol  rune              // Best matching character.
	BestFgColor color.Color       // Best foreground color.
	BestBgColor color.Color       // Best background color.
}

type shiftable interface {
	~uint | ~uint16 | ~uint32 | ~uint64
}

func shift[T shiftable](x T) T {
	if x > 0xff {
		x >>= 8
	}
	return x
}

// Set ScaleMode on Mosaic
func (m *Mosaic) Scale(scale Scale) *Mosaic {
	m.scaleMode = scale
	return m
}

// Set UseFgBgOnly on Mosaic
func (m *Mosaic) OnlyForeground(fgOnly bool) *Mosaic {
	m.useFgBgOnly = fgOnly
	return m
}

// Set DitherLevel on Mosaic
func (m *Mosaic) Dither(dither bool) *Mosaic {
	m.dither = dither
	return m
}

// Set ThresholdLevel on Mosaic
func (m *Mosaic) Threshold(threshold uint8) *Mosaic {
	m.thresholdLevel = threshold
	return m
}

// Set InvertColors on Mosaic
func (m *Mosaic) WithInvertColors(invertColors bool) *Mosaic {
	m.invertColors = invertColors
	return m
}

// Set OutputWidth on Mosaic
func (m *Mosaic) Width(width int) *Mosaic {
	m.outputWidth = width
	return m
}

// Set OutputHeight on Mosaic
func (m *Mosaic) Height(height int) *Mosaic {
	m.outputHeight = height
	return m
}

// Set Symbols on Mosaic
func (m *Mosaic) Symbol(symbol Symbol) *Mosaic {
	m.symbols = symbol
	return m
}

// Render creates a new renderer with the given options.
func (m *Mosaic) Render(img image.Image) string {
	// Quarter blocks.
	if m.symbols == "quarter" || m.symbols == "all" {
		m.blocks = append(m.blocks, QuarterBlocks...)
	}

	// All block elements (including complex combinations).
	if m.symbols == "all" {
		m.blocks = append(m.blocks, ComplexBlocks...)
	}

	// Calculate dimensions.
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	// Determine output dimensions.
	outWidth := m.outputWidth
	outHeight := m.outputHeight

	if outHeight <= 0 {
		// Calculate height based on aspect ratio and character cell proportions.
		// Terminal characters are roughly twice as tall as wide, so we divide by 2.
		outHeight = int(float64(outWidth) * float64(srcHeight) / float64(srcWidth) / 2)
		if outHeight < 1 {
			outHeight = 1
		}
	}

	// Scale image according to the selected mode.
	var scaledImg image.Image
	switch m.scaleMode {
	case Stretch:
		scaledImg = m.scaleImage(img, outWidth*2, outHeight*2) // *2 for sub-character precision.
	case Center:
		// Center the image, maintaining original size.
		scaledImg = m.centerImage(img, outWidth*2, outHeight*2)
	case Fit:
		// Scale while preserving aspect ratio.
		scaledImg = m.fitImage(img, outWidth*2, outHeight*2)
	default:
		// Do nothing.
		scaledImg = m.scaleImageWithoutDistortion(img, outWidth, outHeight)
	}

	// Apply dithering if enabled.
	if m.dither {
		scaledImg = m.applyDithering(scaledImg)
	}

	// Invert colors if needed.
	if m.invertColors {
		scaledImg = m.invertImage(scaledImg)
	}

	// Generate terminal outpum.
	var output strings.Builder

	// Process the image by 2x2 blocks (representing one character cell).

	imageBounds := scaledImg.Bounds()

	for y := 0; y < imageBounds.Max.Y; y += 2 {
		for x := 0; x < imageBounds.Max.X; x++ {
			// Create and analyze the 2x2 pixel block.
			block := m.createPixelBlock(scaledImg, x, y)

			// Determine best symbol and colors.
			m.findBestRepresentation(block)

			// Append to output.
			output.WriteString(
				ansi.Style{}.ForegroundColor(block.BestFgColor).BackgroundColor(block.BestBgColor).Styled(string(block.BestSymbol)),
			)
		}
		output.WriteString("\n")
	}

	return output.String()
}

// createPixelBlock extracts a 2x2 block of pixels from the image.
func (m *Mosaic) createPixelBlock(img image.Image, x, y int) *pixelBlock {
	block := &pixelBlock{}

	// Extract the 2x2 pixel grid.
	for dy := 0; dy < 2; dy++ {
		for dx := 0; dx < 2; dx++ {
			block.Pixels[dy][dx] = m.getPixelSafe(img, x+dx, y+dy)
		}
	}

	return block
}

// findBestRepresentation finds the best block character and colors for a 2x2 pixel block.
func (m *Mosaic) findBestRepresentation(block *pixelBlock) {
	// Simple case: use only foreground/background colors.
	if m.useFgBgOnly {
		// Just use the upper half block with top pixels as background and bottom as foreground.
		block.BestSymbol = '▀'
		block.BestBgColor = m.averageColors([]color.Color{block.Pixels[0][0], block.Pixels[0][1]})
		block.BestFgColor = m.averageColors([]color.Color{block.Pixels[1][0], block.Pixels[1][1]})
		return
	}

	// Determine which pixels are "set" based on threshold.
	pixelMask := [2][2]bool{}
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			// Calculate luminance.
			luma := rgbaToLuminance(block.Pixels[y][x])
			pixelMask[y][x] = luma >= m.thresholdLevel
		}
	}

	// Find the best matching block character.
	bestChar := ' '
	bestScore := math.MaxFloat64

	for _, blockChar := range m.blocks {
		score := 0.0
		for i := 0; i < 4; i++ {
			y, x := i/2, i%2
			if blockChar.Coverage[i] != pixelMask[y][x] {
				score += 1.0
			}
		}

		if score < bestScore {
			bestScore = score
			bestChar = blockChar.Char
		}
	}

	// Determine foreground and background colors based on the best character.
	var fgPixels, bgPixels []color.Color

	// Get the coverage pattern for the selected character.
	var coverage [4]bool
	for _, b := range m.blocks {
		if b.Char == bestChar {
			coverage = b.Coverage
			break
		}
	}

	// Assign pixels to foreground or background based on the character's coverage.
	for i := 0; i < 4; i++ {
		y, x := i/2, i%2
		if coverage[i] {
			fgPixels = append(fgPixels, block.Pixels[y][x])
		} else {
			bgPixels = append(bgPixels, block.Pixels[y][x])
		}
	}

	// Calculate average colors.
	if len(fgPixels) > 0 {
		block.BestFgColor = m.averageColors(fgPixels)
	} else {
		// Default to black if no foreground pixels.
		block.BestFgColor = color.Black
	}

	if len(bgPixels) > 0 {
		block.BestBgColor = m.averageColors(bgPixels)
	} else {
		// Default to black if no background pixels.
		block.BestBgColor = color.Black
	}

	block.BestSymbol = bestChar
}

// averageColors calculates the average color from a slice of colors.
func (m *Mosaic) averageColors(colors []color.Color) color.Color {
	if len(colors) == 0 {
		return color.Black
	}

	var sumR, sumG, sumB, sumA uint32

	for _, c := range colors {
		r, g, b, a := c.RGBA()
		r, g, b, a = shift(r), shift(g), shift(b), shift(a)
		sumR += r
		sumG += g
		sumB += b
		sumA += a
	}

	count := uint32(len(colors))
	return color.RGBA{
		R: uint8(sumR / count),
		G: uint8(sumG / count),
		B: uint8(sumB / count),
		A: uint8(sumA / count),
	}
}

// getPixelSafe returns the color at (x,y) or black if out of bounds.
func (m *Mosaic) getPixelSafe(img image.Image, x, y int) color.RGBA {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return color.RGBA{0, 0, 0, 255}
	}

	r8, g8, b8, a8 := img.At(x, y).RGBA()
	return color.RGBA{
		R: uint8(r8 >> 8),
		G: uint8(g8 >> 8),
		B: uint8(b8 >> 8),
		A: uint8(a8 >> 8),
	}
}

// scaleImage resizes an image to the specified dimensions.
func (m *Mosaic) scaleImage(img image.Image, width, height int) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, width, height))
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	for y := 0; y < height; y++ {
		srcY := bounds.Min.Y + y*srcHeight/height
		for x := 0; x < width; x++ {
			srcX := bounds.Min.X + x*srcWidth/width
			result.Set(x, y, img.At(srcX, srcY))
		}
	}

	return result
}

// fitImage scales image while preserving aspect ratio.
func (m *Mosaic) fitImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	// Calculate scaling factor to fit within maxWidth/maxHeight.
	widthRatio := float64(maxWidth) / float64(srcWidth)
	heightRatio := float64(maxHeight) / float64(srcHeight)

	// Use the smaller ratio to ensure image fits.
	ratio := math.Min(widthRatio, heightRatio)

	newWidth := int(float64(srcWidth) * ratio)
	newHeight := int(float64(srcHeight) * ratio)

	// Scale the image.
	scaledImg := m.scaleImage(img, newWidth, newHeight)

	return scaledImg
}

// fitImage scales image while preserving aspect ratio.
func (m *Mosaic) scaleImageWithoutDistortion(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	// Calculate scaling factor to fit within maxWidth/maxHeight.
	widthRatio := float64(maxWidth) / float64(srcWidth)
	heightRatio := float64(maxHeight) / float64(srcHeight)

	// Use the smaller ratio to ensure image fits.
	ratio := math.Min(widthRatio, heightRatio)

	newWidth := int(float64(srcWidth) * ratio)
	newHeight := int(float64(srcHeight) * ratio)

	// Scale the image.
	scaledImg := m.scaleImage(img, newWidth, newHeight)

	return scaledImg
}

// centerImage centers the original image without scaling.
func (m *Mosaic) centerImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	// If image is larger than max dimensions, scale it down.
	if srcWidth > maxWidth || srcHeight > maxHeight {
		return m.fitImage(img, maxWidth, maxHeight)
	}

	// Otherwise center it without scaling.
	return m.centerScaledImage(img, maxWidth, maxHeight)
}

// centerScaledImage places a scaled image in the center of a larger canvas.
func (m *Mosaic) centerScaledImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	// Create a new black canvas.
	result := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))

	// Calculate offsets to center the image.
	xOffset := (maxWidth - srcWidth) / 2
	yOffset := (maxHeight - srcHeight) / 2

	// Copy the image to the center of the canvas.
	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			if x+xOffset >= 0 && x+xOffset < maxWidth && y+yOffset >= 0 && y+yOffset < maxHeight {
				result.Set(x+xOffset, y+yOffset, img.At(x+bounds.Min.X, y+bounds.Min.Y))
			}
		}
	}

	return result
}

// applyDithering applies Floyd-Steinberg dithering.
func (m *Mosaic) applyDithering(img image.Image) image.Image {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.Point{})
	return pm
}

// invertImage inverts the colors of an image.
func (m *Mosaic) invertImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	result := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r8, g8, b8, a8 := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			result.Set(x, y, color.RGBA{
				R: uint8(255 - (r8 >> 8)),
				G: uint8(255 - (g8 >> 8)),
				B: uint8(255 - (b8 >> 8)),
				A: uint8(a8 >> 8),
			})
		}
	}

	return result
}

// rgbaToLuminance converts RGBA color to luminance (brightness).
func rgbaToLuminance(c color.Color) uint8 {
	r, g, b, a := c.RGBA()
	r, g, b, a = shift(r), shift(g), shift(b), shift(a)
	// Weighted RGB to account for human perception.
	return uint8(float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114)
}

// clamp ensures value is between 0-255.
func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
