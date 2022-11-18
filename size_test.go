package lipgloss

import "fmt"

func Example_size() {
	fmt.Printf("width: %d %d\n", Width("hello\nuniverse"), Width("⏩"))
	fmt.Printf("esc: %d\n", Width("hel\x1b[31mlo"))
	fmt.Printf("height: %d\n", Height("hello\nworld!"))

	w, h := Size("hello\nuniverse")
	fmt.Printf("size: %d %d\n", w, h)

	// Output:
	// width: 8 2
	// esc: 5
	// height: 2
	// size: 8 2
}
