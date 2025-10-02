const { Style, Color } = require('./src/index.js');

function hexToRgb(hex) {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return { r, g, b };
}

function rgbToHex(rgb) {
  const toHex = (c) => {
    const hex = Math.round(c).toString(16);
    return hex.length === 1 ? "0" + hex : hex;
  };
  return `#${toHex(rgb.r)}${toHex(rgb.g)}${toHex(rgb.b)}`;
}

function blendLuv(color1, color2, t) {
  return {
    r: color1.r + (color2.r - color1.r) * t,
    g: color1.g + (color2.g - color1.g) * t,
    b: color1.b + (color2.b - color1.b) * t,
  };
}

function colorGrid(xSteps, ySteps) {
  const x0y0 = hexToRgb("#F25D94");
  const x1y0 = hexToRgb("#EDFF82");
  const x0y1 = hexToRgb("#643AFF");
  const x1y1 = hexToRgb("#14F9D5");

  const x0 = [];
  for (let i = 0; i < ySteps; i++) {
    x0.push(blendLuv(x0y0, x0y1, i / ySteps));
  }

  const x1 = [];
  for (let i = 0; i < ySteps; i++) {
    x1.push(blendLuv(x1y0, x1y1, i / ySteps));
  }

  const grid = [];
  for (let x = 0; x < ySteps; x++) {
    const y0 = x0[x];
    const row = [];
    for (let y = 0; y < xSteps; y++) {
      row.push(rgbToHex(blendLuv(y0, x1[x], y / xSteps)));
    }
    grid.push(row);
  }

  return grid;
}

console.log('Testing colorGrid function...');
const colors = colorGrid(14, 8);
console.log('Generated color grid, now creating styles...');

let b = "";
for (let i = 0; i < colors.length; i++) {
  const x = colors[i];
  for (let j = 0; j < x.length; j++) {
    const y = x[j];
    console.log(`Creating style ${i},${j} with color ${y}`);
    const s = new Style().setString("  ").background(Color(y));
    b += s.render();
  }
  b += "\n";
}

console.log('Success!');