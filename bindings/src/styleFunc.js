// Callback registry
const styleFunction = {};
let nextFunctionId = 1;

// Register a style callback
export function registerStyleFunction(callback) {
  const id = nextFunctionId++;
  styleFunction[id] = callback;
  return id;
}

// Call style function from javascript
export function callStyleFunction(functionId, row, col) {
  if (!styleFunction[functionId]) {
    console.error(`Style function with ID ${functionId} not found!`);
    return 0; // Return null pointer
  }
  
  try {
    const style = styleFunction[functionId](row, col);
    if (!style || !style.addr) {
      console.error(`Style function returned invalid style:`, style);
      return 0;
    }
    return style.addr;
  } catch (error) {
    console.error(`Error in style function:`, error);
    return 0;
  }
}
