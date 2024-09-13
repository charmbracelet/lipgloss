package lipgloss

import "testing"

func Example_special_case() {
	out(JoinVertical(Left))
	out(JoinVertical(Left, ""))
	out(JoinVertical(Left, "hello"))
	out(JoinHorizontal(Top))
	out(JoinHorizontal(Top, ""))
	out(JoinHorizontal(Top, "hello"))

	// Output:
	// 🛇
	//
	// 🛇
	//
	// hello🛇
	//
	// 🛇
	//
	// 🛇
	//
	// hello🛇
}

func Example_vertical() {
	blockA := "AAA\nAA"
	blockB := "BBB\nBBBBB\nBB"
	blockC := "C"

	out(JoinVertical(Left, blockA, blockB, blockC))
	out(JoinVertical(Center, blockA, blockB, blockC))
	out(JoinVertical(Right, blockA, blockB, blockC))
	out(JoinVertical(0.25, blockA, blockB, blockC))

	// Output:
	// AAA__␤
	// AA___␤
	// BBB__␤
	// BBBBB␤
	// BB___␤
	// C____🛇
	//
	// _AAA_␤
	// __AA_␤
	// _BBB_␤
	// BBBBB␤
	// __BB_␤
	// __C__🛇
	//
	// __AAA␤
	// ___AA␤
	// __BBB␤
	// BBBBB␤
	// ___BB␤
	// ____C🛇
	//
	// _AAA_␤
	// _AA__␤
	// _BBB_␤
	// BBBBB␤
	// _BB__␤
	// _C___🛇
}

func Example_horizontal() {
	blockA := "AAA\nAA\n\n\n"
	blockB := "BBB\nBBBBB\nBB"
	blockC := "C"

	out(JoinHorizontal(Top, blockA, blockB, blockC))
	out(JoinHorizontal(Center, blockA, blockB, blockC))
	out(JoinHorizontal(Bottom, blockA, blockB, blockC))
	out(JoinHorizontal(0.25, blockA, blockB, blockC))

	// Output:
	// AAABBB__C␤
	// AA_BBBBB_␤
	// ___BB____␤
	// _________␤
	// _________🛇
	//
	// AAA______␤
	// AA_BBB___␤
	// ___BBBBBC␤
	// ___BB____␤
	// _________🛇
	//
	// AAA______␤
	// AA_______␤
	// ___BBB___␤
	// ___BBBBB_␤
	// ___BB___C🛇
	//
	// AAA______␤
	// AA_BBB__C␤
	// ___BBBBB_␤
	// ___BB____␤
	// _________🛇
}

func TestJoinVertical(t *testing.T) {
	type test struct {
		name     string
		result   string
		expected string
	}
	tests := []test{
		{"pos0", JoinVertical(Left, "A", "BBBB"), "A   \nBBBB"},
		{"pos1", JoinVertical(Right, "A", "BBBB"), "   A\nBBBB"},
		{"pos0.25", JoinVertical(0.25, "A", "BBBB"), " A  \nBBBB"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.result != test.expected {
				t.Errorf("Got \n%s\n, expected \n%s\n", test.result, test.expected)
			}
		})
	}
}

func TestJoinHorizontal(t *testing.T) {
	type test struct {
		name     string
		result   string
		expected string
	}
	tests := []test{
		{"pos0", JoinHorizontal(Top, "A", "B\nB\nB\nB"), "AB\n B\n B\n B"},
		{"pos1", JoinHorizontal(Bottom, "A", "B\nB\nB\nB"), " B\n B\n B\nAB"},
		{"pos0.25", JoinHorizontal(0.25, "A", "B\nB\nB\nB"), " B\nAB\n B\n B"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.result != test.expected {
				t.Errorf("Got \n%s\n, expected \n%s\n", test.result, test.expected)
			}
		})
	}
}
