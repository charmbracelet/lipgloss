package lipgloss

import "testing"

func TestAlignTextVertical(t *testing.T) {
	tests := []struct {
		str    string
		pos    Position
		height int
		want   string
	}{
		{str: "Foo", pos: Top, height: 2, want: "Foo\n"},
		{str: "Foo", pos: Center, height: 5, want: "\n\nFoo\n\n"},
		{str: "Foo", pos: Bottom, height: 5, want: "\n\n\n\nFoo"},

		{str: "Foo\nBar", pos: Bottom, height: 5, want: "\n\n\nFoo\nBar"},
		{str: "Foo\nBar", pos: Center, height: 5, want: "\nFoo\nBar\n\n"},
		{str: "Foo\nBar", pos: Top, height: 5, want: "Foo\nBar\n\n\n"},

		{str: "Foo\nBar\nBaz", pos: Bottom, height: 5, want: "\n\nFoo\nBar\nBaz"},
		{str: "Foo\nBar\nBaz", pos: Center, height: 5, want: "\nFoo\nBar\nBaz\n"},

		{str: "Foo\nBar\nBaz", pos: Bottom, height: 3, want: "Foo\nBar\nBaz"},
		{str: "Foo\nBar\nBaz", pos: Center, height: 3, want: "Foo\nBar\nBaz"},
		{str: "Foo\nBar\nBaz", pos: Top, height: 3, want: "Foo\nBar\nBaz"},

		{str: "Foo\n\n\n\nBar", pos: Bottom, height: 5, want: "Foo\n\n\n\nBar"},
		{str: "Foo\n\n\n\nBar", pos: Center, height: 5, want: "Foo\n\n\n\nBar"},
		{str: "Foo\n\n\n\nBar", pos: Top, height: 5, want: "Foo\n\n\n\nBar"},

		{str: "Foo\nBar\nBaz", pos: Center, height: 9, want: "\n\n\nFoo\nBar\nBaz\n\n\n"},
		{str: "Foo\nBar\nBaz", pos: Center, height: 10, want: "\n\n\nFoo\nBar\nBaz\n\n\n\n"},
	}

	for _, test := range tests {
		got := alignTextVertical(test.str, test.pos, test.height, nil)
		if got != test.want {
			t.Errorf("alignTextVertical(%q, %v, %d) = %q, want %q", test.str, test.pos, test.height, got, test.want)
		}
	}
}
