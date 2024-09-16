package lipgloss

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"
)

func revealNL(o string) string {
	// Make newlines visible.
	o = strings.ReplaceAll(o, "\n", "␤\n")
	o = strings.ReplaceAll(o, " ", "_")
	o = strings.ReplaceAll(o, "\t", ">-->")
	o = strings.ReplaceAll(o, "\x1b", "^[")
	// Add a "no newline at end" marker if there was no newline at the end.
	if len(o) == 0 || o[len(o)-1] != '\n' {
		o += "🛇"
	}
	return o
}

func out(o string) {
	fmt.Println(revealNL(o))
	fmt.Println()
}

func Example_position() {
	fmt.Println("wide:")
	out(PlaceHorizontal(3, Left, "hello"))

	fmt.Println("horizontal:")
	out(PlaceHorizontal(10, Left, "hello"))
	out(PlaceHorizontal(10, Right, "hello"))
	out(PlaceHorizontal(10, Center, "hello"))

	fmt.Println("multiline:")
	out(PlaceHorizontal(10, Left, "hello\nworld"))

	fmt.Println("tall:")
	out(PlaceVertical(1, Top, "hello\nworld"))

	fmt.Println("vertical:")
	out(PlaceVertical(3, Top, "hello"))
	out(PlaceVertical(3, Bottom, "hello"))
	out(PlaceVertical(3, Center, "hello"))

	fmt.Println("both:")
	out(Place(10, 3, Left, Top, "hello"))
	out(Place(10, 3, Center, Top, "hello"))
	out(Place(10, 3, Right, Top, "hello"))
	out(Place(10, 3, Left, Center, "hello"))
	out(Place(10, 3, Center, Center, "hello"))
	out(Place(10, 3, Right, Center, "hello"))
	out(Place(10, 3, Left, Bottom, "hello"))
	out(Place(10, 3, Center, Bottom, "hello"))
	out(Place(10, 3, Right, Bottom, "hello"))

	// Output:
	// wide:
	// hello🛇
	//
	// horizontal:
	// hello_____🛇
	//
	// _____hello🛇
	//
	// __hello___🛇
	//
	// multiline:
	// hello_____␤
	// world_____🛇
	//
	// tall:
	// hello␤
	// world🛇
	//
	// vertical:
	// hello␤
	// _____␤
	// _____🛇
	//
	// _____␤
	// _____␤
	// hello🛇
	//
	// _____␤
	// hello␤
	// _____🛇
	//
	// both:
	// hello_____␤
	// __________␤
	// __________🛇
	//
	// __hello___␤
	// __________␤
	// __________🛇
	//
	// _____hello␤
	// __________␤
	// __________🛇
	//
	// __________␤
	// hello_____␤
	// __________🛇
	//
	// __________␤
	// __hello___␤
	// __________🛇
	//
	// __________␤
	// _____hello␤
	// __________🛇
	//
	// __________␤
	// __________␤
	// hello_____🛇
	//
	// __________␤
	// __________␤
	// __hello___🛇
	//
	// __________␤
	// __________␤
	// _____hello🛇
}

func Example_ws_chars() {
	out(Place(10, 3, Center, Center, "hello",
		WithWhitespaceChars("."),
	))
	out(Place(10, 3, Center, Center, "hello",
		WithWhitespaceChars("☃"),
	))
	out(Place(10, 3, Center, Center, "hello",
		WithWhitespaceChars("12"),
	))
	out(Place(10, 3, Center, Center, "hello",
		WithWhitespaceChars("⏩"),
	))

	// Output:
	// ..........␤
	// ..hello...␤
	// ..........🛇
	//
	//☃☃☃☃☃☃☃☃☃☃␤
	//☃☃hello☃☃☃␤
	//☃☃☃☃☃☃☃☃☃☃🛇
	//
	// 1212121212␤
	// 12hello121␤
	// 1212121212🛇
	//
	// ⏩⏩⏩⏩⏩⏩␤
	// ⏩hello⏩⏩␤
	// ⏩⏩⏩⏩⏩⏩🛇
}

func Example_color() {
	// Force color output.
	curProfile := ColorProfile()
	defer SetColorProfile(curProfile)
	SetColorProfile(termenv.ANSI)

	out(PlaceHorizontal(10, Left, "hello", WithWhitespaceForeground(Color("1"))))
	out(PlaceHorizontal(10, Left, "hello", WithWhitespaceBackground(Color("2"))))

	// Output:
	// hello^[[31m_____^[[0m🛇
	//
	// hello^[[42m_____^[[0m🛇
}
