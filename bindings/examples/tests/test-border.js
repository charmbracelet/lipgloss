const {
  Table,
  Style,
  Color,
  normalBorder,
  Center,
} = require("@charmland/lipgloss");
const readline = require("readline");

// Test data
const testData = [
  [1, "Item A", "Type A", "Info A"],
  [2, "Item B", "Type B", "Info B"],
  [3, "Item C", "Type C", "Info C"],
  [4, "Item D", "Type D", "Info D"],
  [5, "Item E", "Type E", "Info E"],
  [6, "Item F", "Type F", "Info F"],
];

const rowsPerPage = 2;
let currentPage = 0;
const totalPages = Math.ceil(testData.length / rowsPerPage);

// Pre-render all pages at startup to avoid repeated WASM calls
const preRenderedPages = [];

function initializeAllPages() {
  console.log("Initializing all pages...");
  
  // Create styles once
  const purple = Color("99");
  const gray = Color("245");
  const lightGray = Color("241");
  
  const headerStyle = new Style().foreground(purple).bold(true).align(Center);
  const oddRowStyle = new Style().padding(0, 1).width(16).foreground(gray);
  const evenRowStyle = new Style().padding(0, 1).width(16).foreground(lightGray);
  const borderStyle = new Style().foreground(purple);

  // Pre-render each page
  for (let page = 0; page < totalPages; page++) {
    const startIndex = page * rowsPerPage;
    const endIndex = Math.min((page + 1) * rowsPerPage, testData.length);
    const currentRows = testData.slice(startIndex, endIndex);

    // Create a table for this specific page
    const table = new Table()
      .border(normalBorder())
      .borderStyle(borderStyle)
      .headers("#", "NAME", "TYPE", "INFO")
      .styleFunc((row, col) => {
        if (row === -1) {
          return headerStyle;
        }
        // Calculate absolute row for consistent styling
        const absoluteRow = startIndex + row;
        return absoluteRow % 2 === 0 ? evenRowStyle : oddRowStyle;
      })
      .rows(currentRows);

    // Render once and store the result
    const renderedTable = table.render();
    preRenderedPages[page] = renderedTable;
    
    console.log(`Page ${page + 1}/${totalPages} rendered`);
  }
  
  console.log("All pages initialized successfully!");
}

function renderTable(page) {
  // Simply display the pre-rendered page - no WASM calls during navigation
  console.clear();
  console.log(preRenderedPages[page]);
  console.log(`\nPage ${page + 1} of ${totalPages} | Use ↑ and ↓ arrow keys to navigate | Press 'q' to quit`);
}

// Initialize all pages at startup
initializeAllPages();

// Set up readline for arrow key navigation
readline.emitKeypressEvents(process.stdin);
if (process.stdin.isTTY) {
  process.stdin.setRawMode(true);
}

// Initial render
renderTable(currentPage);

// Handle keypress events
process.stdin.on("keypress", (str, key) => {
  if (key.name === "q" || (key.ctrl && key.name === "c")) {
    process.exit();
  }

  if (key.name === "down") {
    if (currentPage < totalPages - 1) {
      currentPage++;
      renderTable(currentPage);
    }
  } else if (key.name === "up") {
    if (currentPage > 0) {
      currentPage--;
      renderTable(currentPage);
    }
  }
});