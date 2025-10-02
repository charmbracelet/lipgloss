const {
  Table,
  TableData,
  Style,
  Color,
  normalBorder,
  Center,
} = require("@charmland/lipgloss");

// Example 1: Basic TableData usage
console.log("=== Basic TableData Example ===");

const data = new TableData()
  .append(["Name", "Age", "City"])
  .append(["Alice", "25", "New York"])
  .append(["Bob", "30", "San Francisco"])
  .append(["Charlie", "35", "Chicago"]);

const basicTable = new Table()
  .data(data)
  .border(normalBorder())
  .render();

console.log(basicTable);

// Example 2: TableData with multiple rows at once
console.log("\n=== TableData with Multiple Rows ===");

const employeeData = new TableData(
  ["Employee ID", "Name", "Department", "Salary"],
  ["001", "John Doe", "Engineering", "$75,000"],
  ["002", "Jane Smith", "Marketing", "$65,000"],
  ["003", "Mike Johnson", "Sales", "$70,000"],
  ["004", "Sarah Wilson", "HR", "$60,000"]
);

const purple = Color("99");
const gray = Color("245");
const lightGray = Color("241");

let headerStyle = new Style().foreground(purple).bold(true).align(Center);
let oddRowStyle = new Style().padding(0, 1).foreground(gray);
let evenRowStyle = new Style().padding(0, 1).foreground(lightGray);

const styledTable = new Table()
  .data(employeeData)
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
  .render();

console.log(styledTable);

// Example 3: Accessing individual cells
console.log("\n=== Accessing TableData Cells ===");

console.log(`Data has ${employeeData.rowCount()} rows and ${employeeData.columnCount()} columns`);
console.log(`Cell at (1, 1): "${employeeData.at(1, 1)}"`);
console.log(`Cell at (2, 2): "${employeeData.at(2, 2)}"`);
console.log(`Cell at (3, 3): "${employeeData.at(3, 3)}"`);

// Example 4: Building data dynamically
console.log("\n=== Dynamic TableData Building ===");

const dynamicData = new TableData()
  .append(["Product", "Price", "Stock"]);

// Simulate adding products dynamically
const products = [
  ["Laptop", "$999", "15"],
  ["Mouse", "$25", "50"],
  ["Keyboard", "$75", "30"],
  ["Monitor", "$299", "8"]
];

products.forEach(product => {
  dynamicData.append(product);
});

const productTable = new Table()
  .data(dynamicData)
  .border(normalBorder())
  .borderStyle(new Style().foreground(Color("cyan")))
  .styleFunc((row, col) => {
    if (row === -1) {
      return new Style().foreground(Color("cyan")).bold(true).align(Center);
    }
    if (col === 1) { // Price column
      return new Style().foreground(Color("green")).align(Center);
    }
    if (col === 2) { // Stock column
      const stock = parseInt(dynamicData.at(row, col));
      if (stock < 10) {
        return new Style().foreground(Color("red")).align(Center);
      } else if (stock < 20) {
        return new Style().foreground(Color("yellow")).align(Center);
      } else {
        return new Style().foreground(Color("green")).align(Center);
      }
    }
    return new Style();
  })
  .render();

console.log(productTable);