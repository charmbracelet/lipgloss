const {
  Table,
  Style,
  Color,
  normalBorder,
  Center,
} = require("@charmland/lipgloss");

console.log("=== ASCII-only Table Test ===");

const purple = Color("99");
const gray = Color("245");
const lightGray = Color("241");

let headerStyle = new Style().foreground(purple).bold(true).align(Center);
let oddRowStyle = new Style().padding(0, 1).width(14).foreground(gray);
let evenRowStyle = new Style().padding(0, 1).width(14).foreground(lightGray);

// ASCII-only data
let rows = [
  ["Chinese", "Hello", "Hi"],
  ["Japanese", "Hello", "Hi"],
  ["Arabic", "Hello", "Hi"],
  ["Russian", "Hello", "Hi"],
  ["Spanish", "Hola", "Que tal?"],
];

const t = new Table()
  .border(normalBorder())
  .borderStyle(new Style().foreground(purple))
  .styleFunc((row, col) => {
    if (row === -1) {
      return headerStyle;
    } else if (row % 2 === 0) {
      return evenRowStyle;
    } else {
      return oddRowStyle;
    }
  })
  .headers("LANGUAGE", "FORMAL", "INFORMAL")
  .rows(rows);

console.log(t.render());