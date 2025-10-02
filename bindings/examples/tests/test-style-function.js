const {
  Table,
  Style,
  Color,
  normalBorder,
  Center,
} = require("@charmland/lipgloss");

console.log("=== Table with Style Function Test ===");

const purple = Color("99");
const gray = Color("245");

let headerStyle = new Style().foreground(purple).bold(true).align(Center);
let rowStyle = new Style().padding(0, 1).foreground(gray);

// Test with style function
const styledTable = new Table()
  .border(normalBorder())
  .headers("Name", "Age")
  .row("Alice", "25")
  .row("Bob", "30")
  .styleFunc((row, col) => {
    console.log(`Style function called with row=${row}, col=${col}`);
    if (row === -1) {
      return headerStyle;
    } else {
      return rowStyle;
    }
  })
  .render();

console.log(styledTable);