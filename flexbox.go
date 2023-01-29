package lipgloss

import (
	"fmt"
	"math"
	"strings"

	"github.com/muesli/reflow/ansi"
)

// FlexDirection represent a flexbox flow direction, establishing the main
// axis inside the container, as well as the direction along this axis.
type FlexDirection int

// FlexDirection enum.
const (
	FlexDirRow FlexDirection = iota
	FlexDirRowReverse
	FlexDirColumn
	FlexDirColumnReverse
)

// FlexWrap represent a flexbox wrap setting, establishing how item within the
// container wrap or not.
type FlexWrap int

// FlexWrap enum.
const (
	FlexWrapNoWrap FlexWrap = iota
	FlexWrapWrap
	FlexWrapWrapReverse
)

// FlexJustifyContent establish how items within a flexbox container are
// distributed along the main axis.
type FlexJustifyContent int

// FlexJustifyContent enum.
const (
	// Flex items are packed toward the start of the line.
	FlexJustifyContentFlexStart FlexJustifyContent = iota
	// Flex items are packed toward the end of the line.
	FlexJustifyContentFlexEnd
	// Flex items are packed toward the center of the line.
	FlexJustifyContentCenter
	// Flex items are evenly distributed in the line.
	FlexJustifyContentSpaceBetween
	// Flex items are evenly distributed in the line, with half-size spaces on
	// either end.
	FlexJustifyContentSpaceAround
)

// FlexAlignItems establish how items are laid out along the cross axis of the
// flexbox container.
type FlexAlignItems int

// FlexAlignItems enum.
const (
	// The cross-start margin edge of the flex item is placed flush with the
	// cross-start edge of the line.
	FlexAlignItemFlexStart FlexAlignItems = iota
	// The cross-end margin edge of the flex item is placed flush with the
	// cross-end edge of the line.
	FlexAlignItemFlexEnd
	// The flex item’s margin box is centered in the cross axis within the line.
	FlexAlignItemCenter
	// If the cross size property of the flex item computes to auto, and neither
	// of the cross-axis margins are auto, the flex item is stretched.
	FlexAlignItemStretch
)

type minMax struct {
	min, max int // int for convenience, but really int32
}

type itemConstraint struct {
	minMax
	baseSize int
	hypoSize int
}

// Flexbox requires to have used Style.SetString beforehand on all items.
func Flexbox(container Style, items ...Style) string {
	// See https://drafts.csswg.org/css-flexbox/#layout-algorithm
	// Adapted for terminal and lipgloss specificities.
	// Some capabilities are not supported.

	// getMainSize returns the min/max size of a Style, in each dimension.
	// If an exact size is known, min == max.
	// If no min is known, zero is used.
	// If no max is known, MaxInt32 is used.
	getMainSize := func(s Style) (width minMax, height minMax) {
		width = minMax{0, math.MaxInt32}
		height = minMax{0, math.MaxInt32}
		// Note: definitive width/height take priority over min/max
		if mw := s.GetMaxWidth(); mw > 0 {
			width.max = mw
		}
		// TODO: support for MinWidth
		if mh := s.GetMaxHeight(); mh > 0 {
			height.max = mh
		}
		// TODO: support for MinHeight
		if w := s.GetWidth(); w > 0 {
			width.min = w
			width.max = w
		}
		if h := s.GetHeight(); h > 0 {
			height.min = h
			height.max = h
		}
		return width, height
	}

	// 2. Determine the available main and cross space for the flex items.
	widthSpace, heightSpace := getMainSize(container)

	// TODO: subtract the flex container’s margin, border, and padding from the space
	//  available to the flex container in that dimension and use that value.
	//  --> can we do that with lipgloss?

	var mainSpace, crossSpace minMax

	switch container.GetFlexDirection() {
	case FlexDirRow, FlexDirRowReverse:
		mainSpace = widthSpace
		crossSpace = heightSpace
	case FlexDirColumn, FlexDirColumnReverse:
		mainSpace = heightSpace
		crossSpace = widthSpace
	default:
		panic(fmt.Sprintf("invalid flex direction: %v", container.GetFlexDirection()))
	}

	fmt.Println("mainSpace", mainSpace)
	fmt.Println("crossSpace", crossSpace)

	// 3. Determine the flex base size and hypothetical main size of each item:
	//   Flex base size: the "natural" size of the item, not clamped
	//   Hypothetical main size: the flex previous base size, but clamped to known min/max
	// Note: as we don't support flex-basis, we use the default "auto" for that value. This defers to
	//   the item explicit width/height, then fallback to "content" value, which means the content's
	//   size itself.

	constraints := make([]itemConstraint, len(items))

	for i, item := range items {
		switch container.GetFlexDirection() {
		case FlexDirRow, FlexDirRowReverse:
			constraints[i].minMax, _ = getMainSize(item)
		case FlexDirColumn, FlexDirColumnReverse:
			_, constraints[i].minMax = getMainSize(item)
		}

		// there is other rules that don't really apply, it seems like we only really
		// care about the "natural" size of the item
		rendered := item.String()
		lines := strings.Split(rendered, "\n")
		switch container.GetFlexDirection() {
		case FlexDirRow, FlexDirRowReverse:
			naturalWidth := 0
			for _, line := range lines {
				naturalWidth = max(naturalWidth, ansi.PrintableRuneWidth(line))
			}
			constraints[i].baseSize = naturalWidth
		case FlexDirColumn, FlexDirColumnReverse:
			naturalHeight := len(lines)
			constraints[i].baseSize = naturalHeight
		}

		// apply clamping to get hypothetical main size
		constraints[i].hypoSize = min(constraints[i].baseSize, constraints[i].max)
		constraints[i].hypoSize = max(constraints[i].hypoSize, constraints[i].min)

		// TODO: and flooring the content box size at zero
	}

	fmt.Println("constraints", constraints)

	// 4. Determine the main size of the flex container using the rules of the formatting
	//   context in which it participates.

	var mainSize int64

	if mainSpace.min == mainSpace.max {
		// we already have a fixed size
		// Note: strictly speaking, that's not necessary a defined size (width/height), it could also be two matching
		// min-XX and max-XX, but I suppose it's an OK simplification.
		mainSize = int64(mainSpace.min)
	} else {
		// we take the max size of all items
		for _, constraint := range constraints {
			mainSize += int64(constraint.hypoSize)
			if constraint.hypoSize == math.MaxInt32 {
				mainSize = math.MaxInt32
				break
			}
		}
	}

	fmt.Println("mainSize", mainSize)

	// 5. Collect flex items into flex lines:
	var lines [][]Style
	var constraintPerLine [][]itemConstraint

	switch container.GetFlexWrap() {
	case FlexWrapNoWrap:
		lines = [][]Style{items}
		constraintPerLine = [][]itemConstraint{constraints}
	case FlexWrapWrap, FlexWrapWrapReverse:
		var currentLine []Style
		var currentConstraints []itemConstraint
		var currentSize int64
		for i, constraint := range constraints {
			if currentSize+int64(constraint.hypoSize) <= mainSize {
				currentLine = append(currentLine, items[i])
				currentConstraints = append(currentConstraints, constraints[i])
				currentSize += int64(constraint.hypoSize)
			} else {
				lines = append(lines, currentLine)
				constraintPerLine = append(constraintPerLine, currentConstraints)
				currentLine = []Style{items[i]}
				currentConstraints = []itemConstraint{constraints[i]}
				currentSize = int64(constraint.hypoSize)
			}
		}
		if len(currentLine) > 0 {
			lines = append(lines, currentLine)
			constraintPerLine = append(constraintPerLine, currentConstraints)
		}
	default:
		panic(fmt.Sprintf("invalid flex wrap: %v", container.GetFlexWrap()))
	}

	fmt.Println("lines:")
	for i, _ := range lines {
		fmt.Printf("%v: ", i)
		for _, constraint := range constraintPerLine[i] {
			fmt.Printf("%v, ", constraint)
		}
		fmt.Println()
	}

	// 6. Resolve the flexible lengths of all the flex items to find their used main size.
	targetMainSizesPerLine := make([][]int, len(lines))
	for i, line := range lines {
		targetMainSizesPerLine[i] = flexResolveFlexibleLength(int(mainSize), line, constraintPerLine[i])
	}

	fmt.Println("targetMainSizes:")
	for i, _ := range lines {
		fmt.Printf("%v: ", i)
		for _, mainSize := range targetMainSizesPerLine[i] {
			fmt.Printf("%v, ", mainSize)
		}
		fmt.Println()
	}

	return ""
}

func flexResolveFlexibleLength(containerSize int, items []Style, constraints []itemConstraint) []int {
	var flexFactor int
	for _, constraint := range constraints {
		flexFactor += constraint.hypoSize
	}

	isGrowFactor := flexFactor < containerSize

	targetMainSize := make([]int, len(items))
	frozen := make([]bool, len(items))

	// Size and freeze inflexible items
	for i, constraint := range constraints {
		if isGrowFactor && items[i].GetFlexGrow() == 0 {
			frozen[i] = true
			targetMainSize[i] = constraint.hypoSize
		}
		if !isGrowFactor && items[i].GetFlexShrink() == 0 {
			frozen[i] = true
			targetMainSize[i] = constraint.hypoSize
		}
		if isGrowFactor && constraint.baseSize > constraint.hypoSize {
			frozen[i] = true
			targetMainSize[i] = constraint.hypoSize
		}
		if !isGrowFactor && constraint.baseSize < constraint.hypoSize {
			frozen[i] = true
			targetMainSize[i] = constraint.hypoSize
		}
	}

	getFreeSpace := func() int {
		freeSpace := containerSize
		for i, constraint := range constraints {
			if frozen[i] {
				freeSpace -= targetMainSize[i]
			} else {
				freeSpace -= constraint.baseSize
			}
		}
		return freeSpace
	}

	// Calculate initial free space
	initialFreeSpace := getFreeSpace()

	for {
		allFrozen := true
		for _, f := range frozen {
			allFrozen = allFrozen && f
		}
		if allFrozen {
			break
		}

		// Calculate the remaining free space
		remainingFreeSpace := getFreeSpace()
		var unfrozenFlexFactor float32
		for i, item := range items {
			if frozen[i] {
				continue
			}
			if isGrowFactor {
				unfrozenFlexFactor += item.GetFlexGrow()
			} else {
				unfrozenFlexFactor += item.GetFlexShrink()
			}
		}
		if unfrozenFlexFactor < 1 {
			val := float32(initialFreeSpace) * unfrozenFlexFactor
			if math.Abs(float64(val)) < math.Abs(float64(remainingFreeSpace)) {
				remainingFreeSpace = int(val + .5)
			}
		}

		// distribute the remaining free space
		var sumScaledShrinkFactor float32
		if !isGrowFactor {
			for i, constraint := range constraints {
				sumScaledShrinkFactor += float32(constraint.baseSize) * items[i].GetFlexShrink()
			}
		}
		for i, constraint := range constraints {
			if frozen[i] {
				continue
			}
			if isGrowFactor {
				ratio := items[i].GetFlexGrow() / unfrozenFlexFactor
				// TODO: do we need adjustment due to int?
				targetMainSize[i] = constraint.baseSize + int(ratio*float32(remainingFreeSpace)+.5)
			} else {
				ratio := (float32(constraint.baseSize) * items[i].GetFlexShrink()) / sumScaledShrinkFactor
				// TODO: do we need adjustment due to int?
				targetMainSize[i] = constraint.baseSize + int(ratio*float32(abs(remainingFreeSpace))+.5)
			}
		}

		// Fix min/max violations
		var totalViolation int
		minViolation := make([]bool, len(items))
		maxViolation := make([]bool, len(items))
		for i, constraint := range constraints {
			if frozen[i] {
				continue
			}
			if targetMainSize[i] < constraint.min {
				totalViolation += constraint.min - targetMainSize[i]
				targetMainSize[i] = constraint.min
				minViolation[i] = true
			}
			if targetMainSize[i] > constraint.max {
				totalViolation += constraint.max - targetMainSize[i]
				targetMainSize[i] = constraint.max
				maxViolation[i] = true
			}
		}

		// Freeze over-flexed items
		switch {
		case totalViolation == 0:
			for i, _ := range frozen {
				frozen[i] = true
			}
		case totalViolation > 0:
			for i, violation := range minViolation {
				frozen[i] = frozen[i] || violation
			}
		case totalViolation < 0:
			for i, violation := range maxViolation {
				frozen[i] = frozen[i] || violation
			}
		}
	}

	return targetMainSize
}
