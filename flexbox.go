package lipgloss

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
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

type flexItem struct {
	style             Style
	constraint        minMax
	baseSize          int // main dimension
	hypoMainSize      int
	targetMainSize    int
	hypoCrossSize     int
	crossSize         int
	mainPos, crossPos int
}

func (f flexItem) String() string {
	return fmt.Sprintf("min: %v\t| max: %v\t| baseSize: %v\t| hypoMainSize: %v\t| targetMainSize: %v\t| hypoCrossSize: %v\t| crossSize: %v\t| mainPos: %v\t| crossPos: %v\n",
		f.constraint.min, f.constraint.max, f.baseSize, f.hypoMainSize, f.targetMainSize, f.hypoCrossSize, f.crossSize, f.mainPos, f.crossPos)
}

type flexLine struct {
	items              []*flexItem
	remainingFreeSpace int
	crossSize          int
}

func (f flexLine) String() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "\nline: remainingFreeSpace: %v\t| crossSize: %v\n", f.remainingFreeSpace, f.crossSize)
	for i, item := range f.items {
		_, _ = fmt.Fprintf(&sb, "    [%v] %v", i, item)
	}
	return sb.String()
}

// Flexbox requires to have used Style.SetString beforehand on all items.
func Flexbox(container Style, items ...Style) string {
	// start by packing all item's style with their respective layout variables
	fItems := make([]*flexItem, len(items))
	for i, item := range items {
		fItems[i] = &flexItem{style: item}
	}

	// See https://drafts.csswg.org/css-flexbox/#layout-algorithm
	// Adapted for terminal and lipgloss specificities.
	// Some capabilities are not supported.

	// 2. Determine the available main and cross space for the flex items.
	widthSpace, heightSpace := getInnerMinMax(container)

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

	// fmt.Println("mainSpace", mainSpace)
	// fmt.Println("crossSpace", crossSpace)

	// 3. Determine the flex base size and hypothetical main size of each item:
	//   Flex base size: the "natural" size of the item, not clamped
	//   Hypothetical main size: the previous flex base size, but clamped to known min/max
	// Note: as we don't support flex-basis, we use the default "auto" for that value. This defers to
	//   the item explicit width/height, then fallback to "content" value, which means the content's
	//   size itself.
	computeFlexBaseSizeAndHypoSize(container, fItems)

	// fmt.Println("items\n", fItems)

	// 4. Determine the main size of the flex container using the rules of the formatting
	//   context in which it participates.

	var mainSize int64

	if mainSpace.min == mainSpace.max {
		// we already have a fixed size
		// Note: strictly speaking, that's not necessary a defined size (width/height), it could also be
		// two matching min-XX and max-XX, but I suppose it's an OK simplification.
		mainSize = int64(mainSpace.min)
	} else {
		// we take the max size of all items
		for _, item := range fItems {
			mainSize += int64(item.hypoMainSize)
			if item.hypoMainSize == math.MaxInt32 {
				mainSize = math.MaxInt32
				break
			}
		}
	}

	// fmt.Println("mainSize", mainSize)

	// 5. Collect flex items into flex lines:
	lines := collectIntoFlexLines(container, mainSize, fItems)

	// fmt.Println("organized in lines\n", lines)

	// 6. Resolve the flexible lengths of all the flex items to find their used main size.
	for _, line := range lines {
		flexResolveFlexibleLength(int(mainSize), line)
	}

	// fmt.Println("target main sizes\n", lines)

	// 7. Determine the hypothetical cross size of each item
	computeItemsHypotheticalCrossSize(container, lines)

	// fmt.Println("hypoCrossSize\n", lines)

	// 8. Calculate the cross size of each flex line.
	computeLineCrossSize(container, crossSpace, lines)

	// fmt.Println("lineCrossSize\n", lines)

	// 11. Determine the used cross size of each flex item.
	for _, line := range lines {
		computeItemCrossSize(line)
	}

	// 12. Distribute any remaining free space.
	for _, line := range lines {
		distributeMainSpace(container, line)
	}

	// fmt.Println("mainSpace\n", lines)

	// 15. Determine the flex container’s used cross size
	var crossSize int

	if crossSpace.min == crossSpace.max {
		// we already have a fixed size
		// Note: strictly speaking, that's not necessary a defined size (width/height), it could also be
		// two matching min-XX and max-XX, but I suppose it's an OK simplification.
		crossSize = crossSpace.min
	} else {
		// we take the max size of each line, which are the max size of each items
		for _, line := range lines {
			crossSize += line.crossSize
		}
	}

	// 16. Align all flex lines per align-content.
	distributeCrossSpace(container, crossSize, lines)

	fmt.Println("crossSpace\n", lines)

	return renderFlexbox(container, int(mainSize), crossSize, lines)
}

func computeFlexBaseSizeAndHypoSize(container Style, items []*flexItem) {
	for i, item := range items {
		switch container.GetFlexDirection() {
		case FlexDirRow, FlexDirRowReverse:
			items[i].constraint, _ = getOuterMinMax(item.style)
		case FlexDirColumn, FlexDirColumnReverse:
			_, items[i].constraint = getOuterMinMax(item.style)
		}

		offsetWidth, offsetHeight := getOuterOffset(item.style)

		// there is other rules that don't really apply, it seems like we only really
		// care about the "natural" size of the item, possibly augmented by borders and margins.
		rendered := item.style.String()
		lines := strings.Split(rendered, "\n")
		switch container.GetFlexDirection() {
		case FlexDirRow, FlexDirRowReverse:
			naturalWidth := 0
			for _, line := range lines {
				naturalWidth = max(naturalWidth, ansi.PrintableRuneWidth(line))
			}
			items[i].baseSize = naturalWidth + offsetWidth
		case FlexDirColumn, FlexDirColumnReverse:
			naturalHeight := len(lines)
			items[i].baseSize = naturalHeight + offsetHeight
		}

		// apply clamping to get hypothetical main size
		items[i].hypoMainSize = min(items[i].baseSize, items[i].constraint.max)
		items[i].hypoMainSize = max(items[i].hypoMainSize, items[i].constraint.min)

		// TODO: and flooring the content box size at zero
	}
}

func collectIntoFlexLines(container Style, mainSize int64, items []*flexItem) []*flexLine {
	var lines []*flexLine

	switch container.GetFlexWrap() {
	case FlexWrapNoWrap:
		lines = []*flexLine{{items: items}}
	case FlexWrapWrap, FlexWrapWrapReverse:
		currentLine := &flexLine{}
		var currentSize int64
		for _, item := range items {
			if currentSize+int64(item.hypoMainSize) <= mainSize {
				currentLine.items = append(currentLine.items, item)
				currentSize += int64(item.hypoMainSize)
			} else {
				lines = append(lines, currentLine)
				currentLine = &flexLine{items: []*flexItem{item}}
				currentSize = int64(item.hypoMainSize)
			}
		}
		if len(currentLine.items) > 0 {
			lines = append(lines, currentLine)
		}
	default:
		panic(fmt.Sprintf("invalid flex wrap: %v", container.GetFlexWrap()))
	}
	return lines
}

// getInnerMinMax returns the min/max size of the content of a Style, in each dimension.
// Those sizes correspond to height/width of the "border-box" CSS box model, which include padding, but not borders and margins.
// If an exact size is known, min == max.
// If no min is known, zero is used.
// If no max is known, MaxInt32 is used.
func getInnerMinMax(s Style) (width minMax, height minMax) {
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

// getOuterOffset returns the total offset in width/height between the inner content box + padding, and the
// outer box (content+padding+border+margin).
func getOuterOffset(s Style) (width, height int) {
	{
		_, top, right, bottom, left := s.GetBorder()
		if top {
			height += 1
		}
		if bottom {
			height += 1
		}
		if right {
			width += 1
		}
		if left {
			width += 1
		}
	}
	{
		top, right, bottom, left := s.GetMargin()
		height += top + bottom
		width += right + left
	}
	return width, height
}

// getOuterMinMax returns the min/max size of the outer box of a Style, in each dimension.
// Those sizes include padding, borders and margins.
// If an exact size is known, min == max.
// If no min is known, zero is used.
// If no max is known, MaxInt32 is used.
func getOuterMinMax(s Style) (width minMax, height minMax) {
	width, height = getInnerMinMax(s)
	offsetWidth, offsetHeight := getOuterOffset(s)

	if width.min != math.MaxInt32 {
		width.min += offsetWidth
	}
	if width.max != math.MaxInt32 {
		width.max += offsetWidth
	}
	if height.min != math.MaxInt32 {
		height.min += offsetHeight
	}
	if height.max != math.MaxInt32 {
		height.max += offsetHeight
	}

	return width, height
}

func flexResolveFlexibleLength(containerSize int, line *flexLine) {
	var flexFactor int
	for _, item := range line.items {
		flexFactor += item.hypoMainSize
	}

	isGrowFactor := flexFactor < containerSize

	frozen := make([]bool, len(line.items))

	// Size and freeze inflexible items
	for i, item := range line.items {
		if isGrowFactor && item.style.GetFlexGrow() == 0 {
			frozen[i] = true
			item.targetMainSize = item.hypoMainSize
		}
		if !isGrowFactor && item.style.GetFlexShrink() == 0 {
			frozen[i] = true
			item.targetMainSize = item.hypoMainSize
		}
		if isGrowFactor && item.baseSize > item.hypoMainSize {
			frozen[i] = true
			item.targetMainSize = item.hypoMainSize
		}
		if !isGrowFactor && item.baseSize < item.hypoMainSize {
			frozen[i] = true
			item.targetMainSize = item.hypoMainSize
		}
	}

	getFreeSpace := func() int {
		freeSpace := containerSize
		for i, item := range line.items {
			if frozen[i] {
				freeSpace -= item.targetMainSize
			} else {
				freeSpace -= item.baseSize
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
			line.remainingFreeSpace = getFreeSpace()
			break
		}

		// Calculate the remaining free space
		line.remainingFreeSpace = getFreeSpace()
		var unfrozenFlexFactor float32
		for i, item := range line.items {
			if frozen[i] {
				continue
			}
			if isGrowFactor {
				unfrozenFlexFactor += item.style.GetFlexGrow()
			} else {
				unfrozenFlexFactor += item.style.GetFlexShrink()
			}
		}
		if unfrozenFlexFactor < 1 {
			val := float32(initialFreeSpace) * unfrozenFlexFactor
			if math.Abs(float64(val)) < math.Abs(float64(line.remainingFreeSpace)) {
				line.remainingFreeSpace = int(val + .5)
			}
		}

		// distribute the remaining free space
		var sumScaledShrinkFactor float32
		if !isGrowFactor {
			for _, item := range line.items {
				sumScaledShrinkFactor += float32(item.baseSize) * item.style.GetFlexShrink()
			}
		}
		for i, item := range line.items {
			if frozen[i] {
				continue
			}
			if isGrowFactor {
				ratio := item.style.GetFlexGrow() / unfrozenFlexFactor
				// TODO: do we need adjustment due to int?
				item.targetMainSize = item.baseSize + int(ratio*float32(line.remainingFreeSpace)+.5)
			} else {
				ratio := (float32(item.baseSize) * item.style.GetFlexShrink()) / sumScaledShrinkFactor
				// TODO: do we need adjustment due to int?
				item.targetMainSize = item.baseSize - int(ratio*float32(abs(line.remainingFreeSpace))+.5)
			}
		}

		// Fix min/max violations
		var totalViolation int
		minViolation := make([]bool, len(line.items))
		maxViolation := make([]bool, len(line.items))
		for i, item := range line.items {
			if frozen[i] {
				continue
			}
			if item.targetMainSize < item.constraint.min {
				totalViolation += item.constraint.min - item.targetMainSize
				item.targetMainSize = item.constraint.min
				minViolation[i] = true
			}
			if item.targetMainSize > item.constraint.max {
				totalViolation += item.constraint.max - item.targetMainSize
				item.targetMainSize = item.constraint.max
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
}

func computeItemsHypotheticalCrossSize(container Style, lines []*flexLine) {
	for _, line := range lines {
		for _, item := range line.items {
			// make a copy to not change the original
			style := item.style.Copy()
			offsetWidth, offsetHeight := getOuterOffset(style)
			switch container.GetFlexDirection() {
			case FlexDirRow, FlexDirRowReverse:
				style = style.Width(item.targetMainSize - offsetWidth)
				rendered := style.String()
				item.hypoCrossSize = len(strings.Split(rendered, "\n"))
			case FlexDirColumn, FlexDirColumnReverse:
				style = style.Height(item.targetMainSize - offsetHeight)
				rendered := style.String()
				for _, s := range strings.Split(rendered, "\n") {
					item.hypoCrossSize = max(item.hypoCrossSize, ansi.PrintableRuneWidth(s))
				}
			}
		}
	}
}

func computeLineCrossSize(container Style, crossSpace minMax, lines []*flexLine) {
	if len(lines) == 1 {
		containerCrossSize := math.MinInt32
		flexDir := container.GetFlexDirection()
		if (flexDir == FlexDirRow || flexDir == FlexDirRowReverse) && container.GetHeight() > 0 {
			containerCrossSize = container.GetHeight()
		}
		if (flexDir == FlexDirColumn || flexDir == FlexDirColumnReverse) && container.GetWidth() > 0 {
			containerCrossSize = container.GetWidth()
		}
		if containerCrossSize != math.MinInt32 {
			lines[0].crossSize = containerCrossSize
			return
		}
	}

	for _, line := range lines {
		for _, item := range line.items {
			line.crossSize = max(line.crossSize, item.hypoCrossSize)
		}
	}

	if len(lines) == 1 {
		lines[0].crossSize = min(lines[0].crossSize, crossSpace.max)
		lines[0].crossSize = max(lines[0].crossSize, crossSpace.min)
	}
}

func computeItemCrossSize(line *flexLine) {
	for _, item := range line.items {
		// as far as I understand, that's really all there is, once what we don't support is removed
		item.crossSize = item.hypoCrossSize
	}
}

func distributeMainSpace(container Style, line *flexLine) {
	justifyContent := container.GetFlexJustifyContent()

	// Support for reversed direction, we flip start and end alignment, and reverse the item order
	switch container.GetFlexDirection() {
	case FlexDirRowReverse, FlexDirColumnReverse:
		switch justifyContent {
		case FlexJustifyContentFlexStart:
			justifyContent = FlexJustifyContentFlexEnd
		case FlexJustifyContentFlexEnd:
			justifyContent = FlexJustifyContentFlexStart
		}
		for i, j := 0, len(line.items)-1; i < j; i, j = i+1, j-1 {
			line.items[i], line.items[j] = line.items[j], line.items[i]
		}
	}

	// if we don't have space, items are simply in order
	if line.remainingFreeSpace <= 0 {
		var pos int
		for _, item := range line.items {
			item.mainPos = pos
			pos += item.targetMainSize
		}
		return
	}

	switch container.GetFlexJustifyContent() {
	case FlexJustifyContentFlexStart:
		var pos int
		for _, item := range line.items {
			item.mainPos = pos
			pos += item.targetMainSize
		}
	case FlexJustifyContentFlexEnd:
		pos := line.remainingFreeSpace
		for _, item := range line.items {
			item.mainPos = pos
			pos += item.targetMainSize
		}
	case FlexJustifyContentCenter:
		pos := line.remainingFreeSpace / 2
		for _, item := range line.items {
			item.mainPos = pos
			pos += item.targetMainSize
		}
	case FlexJustifyContentSpaceBetween:
		var pos int
		freeSpace := line.remainingFreeSpace
		for i, item := range line.items {
			item.mainPos = pos
			spacing := freeSpace / (len(line.items) - i)
			freeSpace -= spacing
			pos += item.targetMainSize + spacing
		}
	case FlexJustifyContentSpaceAround:
		var pos int
		freeSpace := line.remainingFreeSpace
		for i, item := range line.items {
			spacing := freeSpace / (len(line.items) - i + 1)
			freeSpace -= spacing
			pos += spacing
			item.mainPos = pos
			pos += item.targetMainSize
		}
	}
}

func distributeCrossSpace(container Style, crossSize int, lines []*flexLine) {
	switch container.GetFlexAlignItems() {
	case FlexAlignItemFlexStart:
		var offset int
		for _, line := range lines {
			for _, item := range line.items {
				item.crossPos = offset
			}
			offset += line.crossSize
		}
	case FlexAlignItemFlexEnd:
		var offset int
		for _, line := range lines {
			for _, item := range line.items {
				item.crossSize = min(item.crossSize, line.crossSize)
				item.crossPos = offset + (line.crossSize - item.crossSize)
			}
			offset += line.crossSize
		}
	case FlexAlignItemCenter:
		var offset int
		for _, line := range lines {
			for _, item := range line.items {
				item.crossSize = min(item.crossSize, line.crossSize)
				item.crossPos = offset + (line.crossSize - item.crossSize/2)
			}
			offset += line.crossSize
		}
	case FlexAlignItemStretch:
		var offset int
		freeSpace := crossSize
		for _, line := range lines {
			freeSpace -= line.crossSize
		}
		if freeSpace > 0 {
			for i, line := range lines {
				extraSpace := freeSpace / (len(lines) - i)
				line.crossSize += extraSpace
				freeSpace -= extraSpace
			}
		}
		for _, line := range lines {
			for _, item := range line.items {
				item.crossPos = offset
			}
			offset += line.crossSize
		}
	}
}

func renderFlexbox(container Style, mainSize int, crossSize int, lines []*flexLine) string {
	if len(lines) == 0 {
		return ""
	}

	// First step:
	// - normalize container main/cross inner size into regular X/Y
	// - normalize main/cross axis into regular X/Y
	// - render each item into a "fragment" to be assembled

	var xSize, ySize int

	type fragment struct {
		xPos, yPos   int
		xSize, ySize int
		lines        []string
	}

	var fragments []fragment

	switch container.GetFlexDirection() {
	case FlexDirRow, FlexDirRowReverse:
		xSize = mainSize
		ySize = crossSize
		for _, line := range lines {
			for _, item := range line.items {
				// make a copy to not change the original
				style := item.style.Copy()
				offsetWidth, offsetHeight := getOuterOffset(style)
				style = style.Width(item.targetMainSize - offsetWidth).Height(item.crossSize - offsetHeight)
				renderedLines, _ := getLines(style.String())
				fragments = append(fragments, fragment{
					xPos:  item.mainPos,
					yPos:  item.crossPos,
					xSize: item.targetMainSize,
					ySize: item.crossSize,
					lines: renderedLines,
				})
			}
		}
	case FlexDirColumn, FlexDirColumnReverse:
		xSize = crossSize
		ySize = mainSize
		for _, line := range lines {
			for _, item := range line.items {
				// make a copy to not change the original
				style := item.style.Copy()
				offsetWidth, offsetHeight := getOuterOffset(style)
				style = style.Width(item.crossSize - offsetWidth).Height(item.targetMainSize - offsetHeight)
				renderedLines, _ := getLines(style.String())
				fragments = append(fragments, fragment{
					xPos:  item.crossPos,
					yPos:  item.mainPos,
					xSize: item.crossSize,
					ySize: item.targetMainSize,
					lines: renderedLines,
				})
			}
		}
	}

	// Second step: order fragments per line (Y), then per position on that line (X)
	sort.Slice(fragments, func(i, j int) bool {
		if fragments[i].yPos < fragments[j].yPos {
			return true
		}
		if fragments[i].yPos > fragments[j].yPos {
			return false
		}
		return fragments[i].xPos < fragments[j].xPos
	})

	// Third step: assemble fragments, line by line
	var b strings.Builder

	for y := 0; y < ySize; y++ {
		written := 0
		for _, frag := range fragments {
			if y < frag.yPos {
				continue
			}
			if y >= frag.yPos+frag.ySize {
				continue
			}
			b.WriteString(strings.Repeat(" ", frag.xPos-written))
			written += frag.xPos - written
			toWrite := frag.lines[y-frag.yPos]
			toWriteSize := frag.xSize
			if written+frag.xSize > xSize {
				// end overflow
				toWrite = truncate.String(toWrite, uint(xSize-written))
				toWriteSize = xSize - written
			}
			if written > frag.xPos {
				// start overflow

			}
			written += toWriteSize
			b.WriteString(toWrite)
		}
		// Make lines the same length
		b.WriteString(strings.Repeat(" ", xSize-written))
		if y < ySize-1 {
			b.WriteString("\n")
		}
	}

	// Last step: decorate with container
	return container.Render(b.String())
}
