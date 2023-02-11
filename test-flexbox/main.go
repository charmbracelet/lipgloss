package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	// render, adjustContainer := makeSimpleItems()
	render, adjustContainer := makePageLayout()

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	adjustContainer(width, height)
	render()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH, syscall.SIGTERM, syscall.SIGKILL)

	for {
		switch <-sig {
		case syscall.SIGWINCH:
			width, height, _ = term.GetSize(int(os.Stdout.Fd()))
			adjustContainer(width, height)
			render()

		case syscall.SIGTERM, syscall.SIGKILL:
			return
		}
	}
}

func makeSimpleItems() (func(), func(width int, height int)) {
	container := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		FlexWrap(lipgloss.FlexWrapWrap).
		// FlexDirection(lipgloss.FlexDirRowReverse).
		FlexDirection(lipgloss.FlexDirColumn)

	var items []lipgloss.Style

	for i := 0; i < 10; i++ {
		items = append(items, lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			Padding(2).
			SetString(fmt.Sprintf("content %d", i)))
	}

	render := func() {
		fmt.Print(lipgloss.Flexbox(container, items...))
	}

	adjustContainer := func(width, height int) {
		if width > 2 {
			container = container.Width(width - 2)
		}
		if height > 2 {
			container = container.Height(height - 2)
		}
	}

	return render, adjustContainer
}

func makePageLayout() (func(), func(width int, height int)) {
	container := lipgloss.NewStyle().
		FlexWrap(lipgloss.FlexWrapNoWrap).
		FlexAlignItems(lipgloss.FlexAlignItemStretch)

	// header := lipgloss.NewStyle().
	// 	Border(lipgloss.DoubleBorder()).
	// 	FlexGrow(1).
	// 	Padding(2).
	// 	Align(lipgloss.Center).
	// 	SetString("Page Header")

	nav := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1).
		FlexAlignItems(lipgloss.FlexAlignItemStretch).
		SetString(strings.Repeat("• foobar\n", 6))

	main := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		FlexShrink(2).
		FlexGrow(2).
		SetString(strings.Repeat("lorem ipsum ", 60))

	render := func() {
		fmt.Print(lipgloss.Flexbox(container, nav, main))
	}

	adjustContainer := func(width, height int) {
		if width > 0 {
			container = container.Width(width - 0)
		}
		if height > 0 {
			container = container.Height(height - 0)
		}
	}

	return render, adjustContainer
}
