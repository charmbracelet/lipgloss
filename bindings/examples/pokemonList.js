const {
  Table,
  Style,
  Color,
  normalBorder,
  Center,
} = require("@charmland/lipgloss");
const readline = require("readline");

// Create styles once
const purple = Color("99");
const gray = Color("245");
const lightGray = Color("241");
const headerStyle = new Style().foreground(purple).bold(true).align(Center);
const oddRowStyle = new Style().padding(0, 1).width(16).foreground(gray);
const evenRowStyle = new Style().padding(0, 1).width(16).foreground(lightGray);
const borderStyle = new Style().foreground(purple);

const typeColors = {
  Grass: "#78C850",
  Poison: "#A040A0",
  Fire: "#F08030",
  Flying: "#A890F0",
  Water: "#6890F0",
  Rock: "#B8A038",
  Electric: "#F8D030",
  Bug: "#A8B820",
  Normal: "#A8A878",
  Fighting: "#C03028",
  Ice: "#98D8D8",
  Psychic: "#F85888",
  Ground: "#E0C068",
  Ghost: "#705898",
  Dragon: "#7038F8",
  Dark: "#705848",
  Steel: "#B8B8D0",
  Fairy: "#EE99AC",
};

// Function to create a colored cell for types
function createColorCell(typeString) {
  const types = typeString.split("/");
  return types.map(type => {
    const colorMap = {
      'Grass': '🌿',
      'Poison': '☠️',
      'Fire': '🔥',
      'Flying': '🪶',
      'Water': '💧',
      'Rock': '🪨',
      'Electric': '⚡',
      'Bug': '🐛',
      'Normal': '⚪',
      'Fighting': '👊',
      'Ice': '❄️',
      'Psychic': '🔮',
      'Ground': '🌍',
      'Ghost': '👻',
      'Dragon': '🐉',
      'Dark': '🌑',
      'Steel': '⚙️',
      'Fairy': '🧚'
    };
    return (colorMap[type] || '●') + type;
  }).join('/');
}

// Pokemon data
const pokemonData = [
  [1, "Bulbasaur", createColorCell("Grass/Poison"), "Fire, Flying"],
  [2, "Ivysaur", createColorCell("Grass/Poison"), "Fire, Flying"],
  [3, "Venusaur", createColorCell("Grass/Poison"), "Fire, Flying"],
  [4, "Charmander", createColorCell("Fire"), "Water, Rock"],
  [5, "Charmeleon", createColorCell("Fire"), "Water, Rock"],
  [6, "Charizard", createColorCell("Fire/Flying"), "Water, Electric"],
  [7, "Squirtle", createColorCell("Water"), "Grass, Electric"],
  [8, "Wartortle", createColorCell("Water"), "Grass, Electric"],
  [9, "Blastoise", createColorCell("Water"), "Grass, Electric"],
  [10, "Caterpie", createColorCell("Bug"), "Fire, Flying"],
  [11, "Metapod", createColorCell("Bug"), "Fire, Flying"],
  [12, "Butterfree", createColorCell("Bug/Flying"), "Fire, Electric"],
  [13, "Weedle", createColorCell("Bug/Poison"), "Fire, Flying"],
  [14, "Kakuna", createColorCell("Bug/Poison"), "Fire, Flying"],
  [15, "Beedrill", createColorCell("Bug/Poison"), "Fire, Flying"],
  [16, "Pidgey", createColorCell("Normal/Flying"), "Electric, Ice"],
  [17, "Pidgeotto", createColorCell("Normal/Flying"), "Electric, Ice"],
  [18, "Pidgeot", createColorCell("Normal/Flying"), "Electric, Ice"],
  [19, "Rattata", createColorCell("Normal"), "Fighting"],
  [20, "Raticate", createColorCell("Normal"), "Fighting"],
];

// Pagination variables
const rowsPerPage = 5;
let currentPage = 0;
const totalPages = Math.ceil(pokemonData.length / rowsPerPage);

// Cache for rendered pages
const pageCache = new Map();

// Function to render a specific page on-demand
function renderPageToCache(page) {
  if (pageCache.has(page)) {
    return pageCache.get(page);
  }

  console.log(`Rendering page ${page + 1}...`);
  
  // Calculate the start and end indices for the page
  const startIndex = page * rowsPerPage;
  const endIndex = Math.min((page + 1) * rowsPerPage, pokemonData.length);
  const currentRows = pokemonData.slice(startIndex, endIndex);

  // Create the table for this specific page
  const table = new Table()
    .border(normalBorder())
    .borderStyle(borderStyle)
    .headers("#", "NAME", "TYPE", "WEAKNESS")
    .styleFunc((row, col) => {
      if (row === -1) {
        return headerStyle;
      }
      // Calculate absolute row for consistent styling
      const absoluteRow = startIndex + row;
      return absoluteRow % 2 === 0 ? evenRowStyle : oddRowStyle;
    })
    .rows(currentRows);

  // Render and cache the result
  const renderedTable = table.render();
  pageCache.set(page, renderedTable);
  
  return renderedTable;
}

// Function to display a page (from cache or render on-demand)
function renderTable(page) {
  const renderedTable = renderPageToCache(page);
  
  console.clear();
  console.log(renderedTable);
  console.log(`\nPage ${page + 1} of ${totalPages} | Use ↑ and ↓ arrow keys to navigate | Press 'q' to quit`);
}

// Pre-load the first page for instant startup
renderPageToCache(0);

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

console.log("Use ↑ and ↓ arrow keys to navigate through the Pokémon table pages. Press 'q' to quit.");
