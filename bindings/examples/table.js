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
