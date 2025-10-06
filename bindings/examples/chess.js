const {
  Table,
  Style,
  Color,
  normalBorder,
  joinHorizontal,
  joinVertical,
  Right,
  Center,
} = require("@charmland/lipgloss");

var s = new Style().foreground(Color("241")).render;
const board = [
  ["♜", "♞", "♝", "♛", "♚", "♝", "♞", "♜"],
  ["♟", "♟", "♟", "♟", "♟", "♟", "♟", "♟"],
  [" ", " ", " ", " ", " ", " ", " ", " "],
  [" ", " ", " ", " ", " ", " ", " ", " "],
  [" ", " ", " ", " ", " ", " ", " ", " "],
  [" ", " ", " ", " ", " ", " ", " ", " "],
  ["♙", "♙", "♙", "♙", "♙", "♙", "♙", "♙"],
  ["♖", "♘", "♗", "♕", "♔", "♗", "♘", "♖"],
];

let table = new Table()
  .border(normalBorder())
  .borderRow(true)
  .borderColumn(true)
  .rows(board)
  .styleFunc((row, col) => {
    return new Style().padding(0, 1);
  });

let ranks = s([" A", "B", "C", "D", "E", "F", "G", "H  "].join("   "));
let files = s([" 1", "2", "3", "4", "5", "6", "7", "8 "].join("\n\n "));

console.log(
  joinVertical(Right, joinHorizontal(Center, files, table.render()), ranks) +
    "\n",
);
