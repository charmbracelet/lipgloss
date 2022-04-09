package lipgloss

import (
	"fmt"
	"strings"
	"testing"
)

func TestPlace(t *testing.T) {
	// use bytes to compare strings
	defaultSize := 3
	word := "Hello"
	var cols string
	var stringLineWidth string
	if len(word) < defaultSize {
		cols = strings.Repeat(" ", defaultSize)
		stringLineWidth = strings.Repeat(" ", defaultSize-len(word))
	} else {
		cols = strings.Repeat(" ", len(word))
		stringLineWidth = ""
	}
	leftFiller := strings.Repeat(cols+"\n", defaultSize-1)
	rightFiller := strings.Repeat("\n"+cols, defaultSize-1)

	tests := []struct {
		name               string
		size               int
		horizontalPosition Position
		verticalPosition   Position
		want               string
	}{
		{
			name:               "top left, no padding",
			size:               1,
			horizontalPosition: Top,
			verticalPosition:   Top,
			want:               word,
		},
		{
			name:               "top right, some padding",
			size:               defaultSize,
			horizontalPosition: Right,
			verticalPosition:   Top,
			want:               stringLineWidth + word + rightFiller,
		},
		{
			name:               "bottom left, some padding",
			size:               defaultSize,
			horizontalPosition: Left,
			verticalPosition:   Bottom,
			want:               leftFiller + word + stringLineWidth,
		}, {
			name:               "bottom right, some padding",
			size:               defaultSize,
			horizontalPosition: Right,
			verticalPosition:   Bottom,
			want:               leftFiller + stringLineWidth + word,
		},
	}

	for _, tc := range tests {
		got := Place(tc.size, tc.size, tc.horizontalPosition, tc.verticalPosition, word, WithWhitespaceForeground(NoColor{}))
		if got != tc.want {
			fmt.Println([]byte(got))
			fmt.Println([]byte(tc.want))
			t.Errorf("%s\ngot: %s\nwant: %s", tc.name, []byte(got), []byte(tc.want))
		}
	}
}
