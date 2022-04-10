package lipgloss

import (
	"fmt"
	"strings"
	"testing"
)

const (
	topRight uint = iota
	topLeft
	bottomRight
	bottomLeft
)

func TestPlace(t *testing.T) {
	// TODO: make tests with word len that changes
	word := "Hello"

	tests := []struct {
		name               string
		size               int
		horizontalPosition Position
		verticalPosition   Position
		want               uint
	}{
		{
			name:               "top left, no padding",
			size:               1,
			horizontalPosition: Top,
			verticalPosition:   Top,
			want:               topLeft,
		},
		{
			name:               "top left, some padding",
			size:               3,
			horizontalPosition: Left,
			verticalPosition:   Top,
			want:               topLeft,
		},
		{
			name:               "top right, some padding",
			size:               3,
			horizontalPosition: Right,
			verticalPosition:   Top,
			want:               topRight,
		},
		{
			name:               "bottom left, some padding",
			size:               3,
			horizontalPosition: Left,
			verticalPosition:   Bottom,
			want:               bottomLeft,
		},
		{
			name:               "bottom right, some padding",
			size:               3,
			horizontalPosition: Right,
			verticalPosition:   Bottom,
			want:               bottomRight,
		},
		{
			name:               "top left, padding gt word",
			size:               7,
			horizontalPosition: Left,
			verticalPosition:   Top,
			want:               topLeft,
		},
		{
			name:               "top right, padding gt word",
			size:               7,
			horizontalPosition: Right,
			verticalPosition:   Top,
			want:               topRight,
		},
		{
			name:               "bottom left, padding gt word",
			size:               7,
			horizontalPosition: Left,
			verticalPosition:   Bottom,
			want:               bottomLeft,
		},
		{
			name:               "bottom right, padding gt word",
			size:               7,
			horizontalPosition: Right,
			verticalPosition:   Bottom,
			want:               bottomRight,
		},
	}
	for _, tc := range tests {
		got := Place(tc.size, tc.size, tc.horizontalPosition, tc.verticalPosition, word, WithWhitespaceForeground(NoColor{}))
		want := manualFormatString(tc.want, tc.size, word)
		if got != want {
			fmt.Println([]byte(got))
			fmt.Println([]byte(want))
			t.Errorf("%s\ngot: %s\nwant: %s", tc.name, []byte(got), []byte(want))
		}
	}
}

func manualFormatString(choice uint, size int, word string) string {
	var cols string
	var linewithword string

	if len(word) < size {
		cols = strings.Repeat(" ", size)
		linewithword = strings.Repeat(" ", size-len(word))
	} else {
		cols = strings.Repeat(" ", len(word))
		linewithword = ""
	}
	leftFiller := strings.Repeat(cols+"\n", size-1)
	rightFiller := strings.Repeat("\n"+cols, size-1)
	switch choice {
	case topRight:
		return linewithword + word + rightFiller
	case bottomRight:
		return leftFiller + linewithword + word
	case bottomLeft:
		return leftFiller + word + linewithword
	default:
		return word + linewithword + rightFiller
	}
}
