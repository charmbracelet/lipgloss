const {
  Table,
  TableData,
  Style,
  Color,
  normalBorder,
  Center,
} = require("@charmland/lipgloss");

console.log("=== Comprehensive TableData Test ===\n");

// Test 1: Basic TableData creation and usage
console.log("1. Basic TableData creation:");
const basicData = new TableData()
  .append(["Product", "Price", "Stock"])
  .append(["Laptop", "$999", "15"])
  .append(["Mouse", "$25", "50"]);

console.log(`   Rows: ${basicData.rowCount()}, Columns: ${basicData.columnCount()}`);
console.log(`   Cell (1,1): "${basicData.at(1, 1)}"`);

// Test 2: TableData with constructor
console.log("\n2. TableData with constructor:");
const constructorData = new TableData(
  ["Name", "Role", "Department"],
  ["Alice", "Engineer", "Tech"],
  ["Bob", "Designer", "Creative"]
);

console.log(`   Rows: ${constructorData.rowCount()}, Columns: ${constructorData.columnCount()}`);

// Test 3: Table with TableData
console.log("\n3. Table with TableData:");
const table1 = new Table()
  .data(basicData)
  .border(normalBorder())
  .render();

console.log(table1);

// Test 4: Table with TableData and styling
console.log("\n4. Styled Table with TableData:");
const styledTable = new Table()
  .data(constructorData)
  .border(normalBorder())
  .borderStyle(new Style().foreground(Color("blue")))
  .styleFunc((row, col) => {
    if (row === -1) {
      return new Style().bold(true).foreground(Color("cyan")).align(Center);
    }
    if (row % 2 === 0) {
      return new Style().foreground(Color("green"));
    }
    return new Style().foreground(Color("yellow"));
  })
  .render();

console.log(styledTable);

// Test 5: Dynamic data building
console.log("\n5. Dynamic TableData building:");
const dynamicData = new TableData(["ID", "Status"]);

for (let i = 1; i <= 3; i++) {
  dynamicData.append([`${i}`, i % 2 === 0 ? "Active" : "Inactive"]);
}

const dynamicTable = new Table()
  .data(dynamicData)
  .border(normalBorder())
  .styleFunc((row, col) => {
    if (row === -1) return new Style().bold(true);
    if (col === 1) {
      const status = dynamicData.at(row, col);
      return new Style().foreground(status === "Active" ? Color("green") : Color("red"));
    }
    return new Style();
  })
  .render();

console.log(dynamicTable);

console.log("\n✅ All TableData tests completed successfully!");