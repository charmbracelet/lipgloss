const {
  Table,
  Style,
  Color,
  normalBorder,
} = require("@charmland/lipgloss");

console.log("=== Simple Table Test ===");

// Test 1: Basic table without style function
const simpleTable = new Table()
  .headers("Name", "Age")
  .row("Alice", "25")
  .row("Bob", "30")
  .render();

console.log(simpleTable);

console.log("\n=== Table with Border ===");

// Test 2: Table with border but no style function
const borderTable = new Table()
  .border(normalBorder())
  .headers("Name", "Age")
  .row("Alice", "25")
  .row("Bob", "30")
  .render();

console.log(borderTable);