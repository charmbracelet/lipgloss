const fs = require('fs');
const path = require('path');
const {
  Tree,
  Style,
  Color,
} = require("@charmland/lipgloss");

function addBranches(root, dirPath) {
  try {
    const items = fs.readdirSync(dirPath, { withFileTypes: true });
    
    for (const item of items) {
      // Skip items that start with a dot
      if (item.name.startsWith('.')) {
        continue;
      }
      
      if (item.isDirectory()) {
        // It's a directory
        const treeBranch = new Tree().root(item.name);
        root.child(treeBranch);
        
        // Recurse into subdirectory
        const branchPath = path.join(dirPath, item.name);
        addBranches(treeBranch, branchPath);
      } else {
        // It's a file
        root.child(item.name);
      }
    }
  } catch (err) {
    console.error(`Error reading directory ${dirPath}:`, err.message);
  }
}

function main() {
  const enumeratorStyle = new Style()
    .foreground(Color("240"))
    .paddingRight(1);
    
  const itemStyle = new Style()
    .foreground(Color("99"))
    .bold(true)
    .paddingRight(1);
  
  // Get current working directory
  const pwd = process.cwd();
  
  const t = new Tree()
    .root(pwd)
    .enumeratorStyle(enumeratorStyle)
    .rootStyle(itemStyle)
    .itemStyle(itemStyle);
  
  // Build the tree starting from current directory
  addBranches(t, ".");
  
  console.log(t.render());
}

main();