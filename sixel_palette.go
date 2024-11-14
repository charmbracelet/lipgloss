package lipgloss

import (
	"cmp"
	"container/heap"
	"image"
	"image/color"
	"math"
	"slices"
)

// sixelPalette is a palette of up to 256 colors that lists the colors that will be used by
// a SixelImage. Most images, especially jpegs, have more than 256 colors, so creating a sixelPalette
// requires the use of color quantization. For this we use the Median Cut algorithm.
//
// Median cut requires all pixels in an image to be positioned in a 4D color cube, with one axis per channel.
// The cube is sliced in half along its longest axis such that half the pixels in the cube end up in one of
// the sub-cubes and half end up in the other. We continue slicing the cube with the longest axis in half
// along that axis until there are 256 sub-cubes.  Then, the average of all pixels in each subcube is used
// as that cube's color.
//
// Colors are converted to palette colors based on which they are closest to (it's not
// always their cube's color).
//
// This implementation has a few minor (but seemingly very common) differences from the Official algorithm:
//   - When determining the longest axis, the number of pixels in the cube are multiplied against axis length
//     This improves color selection by quite a bit in cases where an image has a lot of space taken up by different
//     shades of the same color.
//   - If a single color sits on a cut line, all pixels of that color are assigned to one of the subcubes
//     rather than try to split them up between the subcubes.  This allows us to use a slice of unique colors
//     and a map of pixel counts rather than try to represent each pixel individually.
type sixelPalette struct {
	// Map used to convert colors from the image to palette colors
	colorConvert map[sixelColor]sixelColor
	// Backward lookup to get index from palette color
	paletteIndexes map[sixelColor]int
	PaletteColors  []sixelColor
}

// quantizationChannel is an enum type which indicates an axis in the color cube. Used to indicate which
// axis in a cube is the longest
type quantizationChannel int

const (
	sixelMaxColors  int                 = 256
	quantizationRed quantizationChannel = iota
	quantizationGreen
	quantizationBlue
	quantizationAlpha
)

// quantizationCube represents a single cube in the median cut algorithm.
type quantizationCube struct {
	// startIndex is the index in the uniqueColors slice where this cube starts
	startIndex int
	// length is the number of elements in the uniqueColors slice this cube occupies
	length int
	// sliceChannel is the axis that will be cut if this cube is cut in half
	sliceChannel quantizationChannel
	// score is a heuristic value: higher means this cube is more likely to be cut
	score uint64
	// pixelCount is how many pixels are contained in this cube
	pixelCount uint64
}

// cubePriorityQueue is a heap used to sort quantizationCube objects in order to select the correct
// one to cut next. Pop will remove the queue with the highest score
type cubePriorityQueue []any

func (p *cubePriorityQueue) Push(x any) {
	*p = append(*p, x)
}

func (p *cubePriorityQueue) Pop() any {
	popped := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return popped
}

func (p *cubePriorityQueue) Len() int {
	return len(*p)
}

func (p *cubePriorityQueue) Less(i, j int) bool {
	left := (*p)[i].(quantizationCube)
	right := (*p)[j].(quantizationCube)

	// We want the largest channel variance
	return left.score > right.score
}

func (p *cubePriorityQueue) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

// createCube is used to initialize a new quantizationCube containing a region of the uniqueColors slice
func (p *sixelPalette) createCube(uniqueColors []sixelColor, pixelCounts map[sixelColor]uint64, startIndex, bucketLength int) quantizationCube {
	minRed, minGreen, minBlue, minAlpha := uint32(0xffff), uint32(0xffff), uint32(0xffff), uint32(0xffff)
	maxRed, maxGreen, maxBlue, maxAlpha := uint32(0), uint32(0), uint32(0), uint32(0)
	totalWeight := uint64(0)

	// Figure out which channel has the greatest variance
	for i := startIndex; i < startIndex+bucketLength; i++ {
		r, g, b, a := uniqueColors[i].Red, uniqueColors[i].Green, uniqueColors[i].Blue, uniqueColors[i].Alpha
		totalWeight += pixelCounts[uniqueColors[i]]

		if r < minRed {
			minRed = r
		}
		if r > maxRed {
			maxRed = r
		}
		if g < minGreen {
			minGreen = g
		}
		if g > maxGreen {
			maxGreen = g
		}
		if b < minBlue {
			minBlue = b
		}
		if b > maxBlue {
			maxBlue = b
		}
		if a < minAlpha {
			minAlpha = a
		}
		if a > maxAlpha {
			maxAlpha = a
		}
	}

	dRed := maxRed - minRed
	dGreen := maxGreen - minGreen
	dBlue := maxBlue - minBlue
	dAlpha := maxAlpha - minAlpha

	cube := quantizationCube{
		startIndex: startIndex,
		length:     bucketLength,
		pixelCount: totalWeight,
	}

	if dRed >= dGreen && dRed >= dBlue && dRed >= dAlpha {
		cube.sliceChannel = quantizationRed
		cube.score = uint64(dRed)
	} else if dGreen >= dBlue && dGreen >= dAlpha {
		cube.sliceChannel = quantizationGreen
		cube.score = uint64(dGreen)
	} else if dBlue >= dAlpha {
		cube.sliceChannel = quantizationBlue
		cube.score = uint64(dBlue)
	} else {
		cube.sliceChannel = quantizationAlpha
		cube.score = uint64(dAlpha)
	}

	// Boost the score of cubes with more pixels in them
	cube.score *= totalWeight

	return cube
}

// quantize is a method that will initialize the palette's colors and lookups, provided a set
// of unique colors and a map containing pixel counts for those colors
func (p *sixelPalette) quantize(uniqueColors []sixelColor, pixelCounts map[sixelColor]uint64, maxColors int) {
	p.colorConvert = make(map[sixelColor]sixelColor)
	p.paletteIndexes = make(map[sixelColor]int)

	// We don't need to quantize if we don't even have more than the maximum colors, and in fact, this code will explode
	// if we have fewer than maximum colors
	if len(uniqueColors) <= maxColors {
		p.PaletteColors = uniqueColors
		return
	}

	cubeHeap := make(cubePriorityQueue, 0, maxColors)

	// Start with a cube that contains all colors
	heap.Init(&cubeHeap)
	heap.Push(&cubeHeap, p.createCube(uniqueColors, pixelCounts, 0, len(uniqueColors)))

	// Slice the best cube into two cubes until we have max colors, then we have our palette
	for cubeHeap.Len() < maxColors {
		cubeToSplit := heap.Pop(&cubeHeap).(quantizationCube)

		// Sort the colors in the bucket's range along the cube's longest color axis
		slices.SortFunc(uniqueColors[cubeToSplit.startIndex:cubeToSplit.startIndex+cubeToSplit.length],
			func(left sixelColor, right sixelColor) int {
				switch cubeToSplit.sliceChannel {
				case quantizationRed:
					return cmp.Compare(left.Red, right.Red)
				case quantizationGreen:
					return cmp.Compare(left.Green, right.Green)
				case quantizationBlue:
					return cmp.Compare(left.Blue, right.Blue)
				default:
					return cmp.Compare(left.Alpha, right.Alpha)
				}
			})

		// We need to split up the colors in this cube so that the pixels are evenly split between the two,
		// or at least as close as we can reasonably get. What we do is count up the pixels as we go through
		// and place the cut point where around half of the pixels are on the left side
		countSoFar := pixelCounts[uniqueColors[cubeToSplit.startIndex]]
		targetCount := cubeToSplit.pixelCount / 2
		leftLength := 1

		for i := cubeToSplit.startIndex + 1; i < cubeToSplit.startIndex+cubeToSplit.length; i++ {
			c := uniqueColors[i]
			weight := pixelCounts[c]
			if countSoFar+weight > targetCount {
				break
			}
			leftLength++
			countSoFar += weight
		}

		rightLength := cubeToSplit.length - leftLength
		rightIndex := cubeToSplit.startIndex + leftLength
		heap.Push(&cubeHeap, p.createCube(uniqueColors, pixelCounts, cubeToSplit.startIndex, leftLength))
		heap.Push(&cubeHeap, p.createCube(uniqueColors, pixelCounts, rightIndex, rightLength))
	}

	// Once we've got max cubes in the heap, pull them all out and load them into the palette
	for cubeHeap.Len() > 0 {
		bucketToLoad := heap.Pop(&cubeHeap).(quantizationCube)
		p.loadColor(uniqueColors, pixelCounts, bucketToLoad.startIndex, bucketToLoad.length)
	}
}

// ColorIndex accepts a raw image color (NOT a palette color) and provides the palette index of that color
func (p *sixelPalette) ColorIndex(c sixelColor) int {
	return p.paletteIndexes[c]
}

// loadColor accepts a range of colors representing a single median cut cube. It calculates the
// average color in the cube and adds it to the palette.
func (p *sixelPalette) loadColor(uniqueColors []sixelColor, pixelCounts map[sixelColor]uint64, startIndex, cubeLen int) {
	totalRed, totalGreen, totalBlue, totalAlpha := uint64(0), uint64(0), uint64(0), uint64(0)
	totalCount := uint64(0)
	for i := startIndex; i < startIndex+cubeLen; i++ {
		count := pixelCounts[uniqueColors[i]]
		totalRed += uint64(uniqueColors[i].Red) * count
		totalGreen += uint64(uniqueColors[i].Green) * count
		totalBlue += uint64(uniqueColors[i].Blue) * count
		totalAlpha += uint64(uniqueColors[i].Alpha) * count
		totalCount += count
	}

	averageColor := sixelColor{
		Red:   uint32(totalRed / totalCount),
		Green: uint32(totalGreen / totalCount),
		Blue:  uint32(totalBlue / totalCount),
		Alpha: uint32(totalAlpha / totalCount),
	}

	p.PaletteColors = append(p.PaletteColors, averageColor)
}

// sixelColor is a flat struct that contains a single color: all channels are 0-100
// instead of anything sensible
type sixelColor struct {
	Red   uint32
	Green uint32
	Blue  uint32
	Alpha uint32
}

// sixelConvertColor accepts an ordinary Go color and converts it to a sixelColor, which
// has channels ranging from 0-100
func sixelConvertColor(c color.Color) sixelColor {
	r, g, b, a := c.RGBA()
	return sixelColor{
		Red:   sixelConvertChannel(r),
		Green: sixelConvertChannel(g),
		Blue:  sixelConvertChannel(b),
		Alpha: sixelConvertChannel(a),
	}
}

// sixelConvertChannel converts a single color channel from go's standard 0-0xffff range to
// sixel's 0-100 range
func sixelConvertChannel(channel uint32) uint32 {
	// We add 327 because that is about 0.5 in the sixel 0-100 color range, we're trying to
	// round to the nearest value
	return (channel + 328) * 100 / 0xffff
}

// newSixelPalette accepts an image and produces an N-color quantized color palette using the median cut
// algorithm. The produced sixelPalette can convert colors from the image to the quantized palette
// in O(1) time.
func newSixelPalette(image image.Image, maxColors int) sixelPalette {
	pixelCounts := make(map[sixelColor]uint64)

	height := image.Bounds().Dy()
	width := image.Bounds().Dx()

	// Record pixel counts for every color while also getting a set of all unique colors in the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := sixelConvertColor(image.At(x, y))
			count, _ := pixelCounts[c]
			count++

			pixelCounts[c] = count
		}
	}

	p := sixelPalette{}
	uniqueColors := make([]sixelColor, 0, len(pixelCounts))
	for c := range pixelCounts {
		uniqueColors = append(uniqueColors, c)
	}

	// Build up p.PaletteColors using the median cut algorithm
	p.quantize(uniqueColors, pixelCounts, maxColors)

	// The average color for a cube a color occupies is not always the closest palette color.  As a result,
	// we need to use this very upsetting double loop to find the lookup palette color for each
	// unique color in the image.
	for _, c := range uniqueColors {
		var bestColor sixelColor
		var bestColorIndex int
		bestScore := uint32(math.MaxUint32)

		for paletteIndex, paletteColor := range p.PaletteColors {
			redDiff := c.Red - paletteColor.Red
			greenDiff := c.Green - paletteColor.Green
			blueDiff := c.Blue - paletteColor.Blue
			alphaDiff := c.Alpha - paletteColor.Alpha

			score := (redDiff * redDiff) + (greenDiff * greenDiff) + (blueDiff * blueDiff) + (alphaDiff * alphaDiff)
			if score < bestScore {
				bestColor = paletteColor
				bestColorIndex = paletteIndex
				bestScore = score
			}
		}

		p.paletteIndexes[c] = bestColorIndex
		p.colorConvert[c] = bestColor
	}

	return p
}
