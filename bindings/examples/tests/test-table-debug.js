const {
  Table,
  Style,
  Color,
  normalBorder,
  Center,
} = require("@charmland/lipgloss");

console.log("=== First Table from table.js ===");

var s = new Style().foreground(Color("240")).render;
console.log(
  new Table()
    .wrap(false)
    .headers("Drink", "Description")
    .row("Bubble Tea", s("Milky"))
    .row("Milk Tea", s("Also milky"))
    .row("Actual milk", s("Milky as well"))
    .render(),
);

console.log("\n=== Second Table Setup ===");

const purple = Color("99");
const gray = Color("245");
const lightGray = Color("241");

let headerStyle = new Style().foreground(purple).bold(true).align(Center);
let oddRowStyle = new Style().padding(0, 1).width(14).foreground(gray);
let evenRowStyle = new Style().padding(0, 1).width(14).foreground(lightGray);

let rows = [
  ["Chinese", "您好", "你好"],
  ["Japanese", "こんにちは", "やあ"],
  ["Arabic", "أهلين", "أهلا"],
  ["Russian", "Здравствуйте", "Привет"],
  ["Spanish", "Hola", "¿Qué tal?"],
];

console.log("Creating table...");

const t = new Table()
  .border(normalBorder())
  .borderStyle(new Style().foreground(purple))
  .headers("LANGUAGE", "FORMAL", "INFORMAL")
  .rows(rows);

console.log("Table created, adding style function...");

t.styleFunc((row, col) => {
  console.log(`Style function called with row=${row}, col=${col}`);
  if (row === -1) {
    return headerStyle;
  } else if (row % 2 === 0) {
    return evenRowStyle;
  } else {
    return oddRowStyle;
  }
});

console.log("Style function added, rendering...");

console.log(t.render());