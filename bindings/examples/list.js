const {
  List,
  Style,
  Color,
  Bullet,
  Arabic,
  Alphabet,
  Dash,
} = require("@charmland/lipgloss");

// Simple bullet list
console.log("Simple bullet list:");
const groceries = new List(
  "Bananas",
  "Barley", 
  "Cashews",
  "Milk",
  "Eggs"
).enumerator(Bullet);

console.log(groceries.render());

// Numbered list with styling
console.log("\nNumbered list with styling:");
const purple = Color("99");
const gray = Color("245");

let itemStyle = new Style().foreground(gray).padding(0, 1);
let enumeratorStyle = new Style().foreground(purple).bold(true);

const tasks = new List()
  .item("Write documentation")
  .item("Add tests")
  .item("Review code")
  .item("Deploy to production")
  .enumerator(Arabic)
  .itemStyle(itemStyle)
  .enumeratorStyle(enumeratorStyle);

console.log(tasks.render());

// Nested list
console.log("\nNested list:");
const subList = new List(
  "Almond Milk",
  "Coconut Milk", 
  "Full Fat Milk"
).enumerator(Dash);

const nestedGroceries = new List(
  "Bananas",
  "Barley",
  "Cashews", 
  subList,
  "Eggs",
  "Fish Cake"
).enumerator(Bullet);

console.log(nestedGroceries.render());

// Alphabetical list
console.log("\nAlphabetical list:");
const alphabet = new List(
  "Apple",
  "Banana", 
  "Cherry",
  "Date",
  "Elderberry"
).enumerator(Alphabet);

console.log(alphabet.render());