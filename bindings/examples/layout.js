// TODO:
// few problems found
// 3. JoinVertical and JoinHorizontal are not decoding hidden characters.

const {
  lightDark,
  Color,
  roundedBorder,
  joinVertical,
  Center,
  Style,
  marginLeft,
  normalBorder,
  Left,
  Right,
  Top,
  Bottom,
  Width,
  joinHorizontal,
  Place,
} = require("@charmland/lipgloss");

const width = 96;
const columnWidth = 30;
const hasDarkBG = true;
let selectLightDark = lightDark(hasDarkBG);

let subtle = selectLightDark(Color("#D9DCCF"), Color("#383838"));
let highlight = selectLightDark(Color("#874BFD"), Color("#7D56F4"));
let special = selectLightDark(Color("#43BF6D"), Color("#73F59F"));

let divider = new Style()
  .setString("•")
  .padding(0, 1)
  .foreground(subtle)
  .string();

let url = new Style().foreground(special).render;

// Tabs.

// let activeTabBorder = lipgloss.Border{
// 	Top:         "─",
// 	Bottom:      " ",
// 	Left:        "│",
// 	Right:       "│",
// 	TopLeft:     "╭",
// 	TopRight:    "╮",
// 	BottomLeft:  "┘",
// 	BottomRight: "└",
// }
let activeTabBorder = normalBorder();

// tabBorder = lipgloss.Border{
// 	Top:         "─",
// 	Bottom:      "─",
// 	Left:        "│",
// 	Right:       "│",
// 	TopLeft:     "╭",
// 	TopRight:    "╮",
// 	BottomLeft:  "┴",
// 	BottomRight: "┴",
// }
let tabBorder = normalBorder();

// tab = lipgloss.NewStyle().
// 			Border(tabBorder, true).
// 			BorderForeground(highlight).
// 			Padding(0, 1)

// 		activeTab = tab.Border(activeTabBorder, true)

let tab = new Style()
  .border(tabBorder, true)
  .borderForeground(highlight)
  .padding(0, 1).render;

let activeTab = new Style()
  .border(tabBorder, true)
  .borderForeground(highlight)
  .padding(0, 1)
  .border(activeTabBorder, true).render;

let tabGap = new Style()
  .border(activeTabBorder, true)
  .borderForeground(highlight)
  .padding(0, 1)
  .borderTop(false)
  .borderLeft(false)
  .borderRight(false);

// Title.
let titleStyle = new Style()
  .marginLeft(1)
  .marginRight(5)
  .padding(0, 1)
  .italic(true)
  .foreground(Color("#FFF7DB"))
  .setString("Lip Gloss");

let descStyle = new Style().marginTop(1);

let infoStyle = new Style()
  .borderStyle(normalBorder())
  .borderTop(true)
  .borderForeground(subtle);

let dialogBoxStyle = new Style()
  .border(roundedBorder())
  .borderForeground(Color("#874BFD"))
  .padding(1, 0)
  .borderTop(true)
  .borderLeft(true)
  .borderRight(true)
  .borderBottom(true);

let buttonStyle = new Style()
  .foreground(Color("#FFF7DB"))
  .background(Color("#888B7E"))
  .padding(0, 3)
  .marginTop(1);

let activeButtonStyle = new Style()
  .padding(0, 3)
  .marginTop(1)
  .foreground(Color("#FFF7DB"))
  .background(Color("#F25D94"))
  .marginRight(2)
  .underline(true);

// List.
let list = new Style()
  .border(normalBorder(), false, true, false, false)
  .borderForeground(subtle)
  .marginRight(1)
  .height(8)
  .width(width / 3).render;

let listHeader = new Style()
  .borderStyle(normalBorder())
  .borderBottom(true)
  .borderForeground(subtle)
  .marginRight(2).render;

let listItem = new Style().paddingLeft(2).render;

let checkMark = new Style()
  .setString("✓")
  .foreground(special)
  .paddingRight(1)
  .string();

function listDone(s) {
  return (
    checkMark +
    new Style()
      .strikethrough(true)
      .foreground(selectLightDark(Color("#969B86"), Color("#696969")))
      .render(s)
  );
}

// Paragraphs/History.
let historyStyle = new Style()
  .align(Left)
  .foreground(Color("#FAFAFA"))
  .background(highlight)
  .margin(1, 3, 0, 0)
  .padding(1, 2)
  .height(19)
  .width(columnWidth);

// Status Bar.

let statusNugget = new Style().foreground(Color("#FFFDF5")).padding(0, 1);

let statusBarStyle = new Style()
  .foreground(selectLightDark(Color("#343433"), Color("#C1C6B2")))
  .background(selectLightDark(Color("#D9DCCF"), Color("#353533")));

let statusStyle = new Style()
  .inherit(statusBarStyle)
  .foreground(Color("#FFFDF5"))
  .background(Color("#FF5F87"))
  .padding(0, 1)
  .marginRight(1);

let encodingStyle = statusNugget.background(Color("#A550DF")).align(Right);

let statusText = new Style().inherit(statusBarStyle);

let fishCakeStyle = new Style()
  .foreground(Color("#FFFDF5"))
  .padding(0, 1)
  .background(Color("#6124DF"));

// Page.
let docStyle = new Style().padding(1, 2, 1, 2);

let physicalWidth = process.stdout.columns;
let doc = "";

// Tabs.
{
  let row = joinHorizontal(
    Top,
    activeTab("Lip Gloss"),
    tab("Blush"),
    tab("Eye Shadow"),
    tab("Mascara"),
    tab("Foundation"),
  );
  // console.log("oi", row.length);
  // let repeat = Math.max(0, width - Width(row) - 2);
  let repeat = 27;
  let gap = tabGap.render(" ".repeat(repeat));
  row = joinHorizontal(Bottom, row, gap);
  doc += row + "\n\n";
}

// row := lipgloss.JoinHorizontal(
// 			lipgloss.Top,
// 			activeTab.Render("Lip Gloss"),
// 			tab.Render("Blush"),
// 			tab.Render("Eye Shadow"),
// 			tab.Render("Mascara"),
// 			tab.Render("Foundation"),
// 		)
// 		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
// 		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
// 		doc.WriteString(row + "\n\n")

// Title.
{
  let colors = colorGrid(1, 5);
  let title = "";

  for (let i = 0; i < colors.length; i++) {
    const offset = 2;
    const v = colors[i];
    let c = Color(v[0]);
    title += titleStyle
      .marginLeft(i * offset)
      .background(c)
      .render();
    if (i < colors.length - 1) {
      title += "\n";
    }
  }

  let desc = joinVertical(
    Left,
    descStyle.render("Style Definitions for Nice Terminal Layouts"),
    infoStyle.render(
      "From Charm" + divider + url("https://github.com/charmbracelet/lipgloss"),
    ),
  );

  let row = joinHorizontal(Top, title, desc);
  doc += row + "\n\n";
}

// Dialog.
{
  let okButton = new Style().inherit(activeButtonStyle).render("Yes");
  let cancelButton = new Style().inherit(buttonStyle).render("Maybe");

  let grad = applyGradient(
    "Are you sure you want to eat marmalade?",
    "#EDFF82",
    "#F25D94",
  );

  let question = new Style().width(50).align(Center).render(grad);

  let buttons = joinHorizontal(Top, okButton, cancelButton);
  let ui = joinVertical(Center, question, buttons);

  let dialog = Place(
    width,
    9,
    Center,
    Center,
    dialogBoxStyle.render(ui),
    // withWhitespaceChars("猫咪"),
    // withWhitespaceStyle(new Style().foreground(subtle)),
  );

  doc += dialog + "\n\n";
}

// // Color grid.
let colors = () => {
  const colors = colorGrid(14, 8);

  let b = "";
  for (let i = 0; i < colors.length; i++) {
    const x = colors[i];
    for (let j = 0; j < x.length; j++) {
      const y = x[j];
      const s = new Style().setString("  ").background(Color(y));
      b += s.render();
    }
    b += "\n";
  }

  return b;
};

// doc += colors();

// let lists = joinHorizontal(
//   Top,
//   list(
//     joinVertical(
//       Left,
//       listHeader("Citrus Fruits to Try"),
//       listDone("Grapefruit"),
//       listDone("Yuzu"),
//       listItem("Citron"),
//       listItem("Kumquat"),
//       listItem("Pomelo"),
//     ),
//   ),
//   list(
//     joinVertical(
//       Left,
//       listHeader("Actual Lip Gloss Vendors"),
//       listItem("Glossier"),
//       listItem("Claire‘s Boutique"),
//       listDone("Nyx"),
//       listItem("Mac"),
//       listDone("Milk"),
//     ),
//   ),
// );

doc += joinHorizontal(Top, new Style().marginLeft(1).render(colors()));

// 	// Marmalade history.
// 	{
// 		const (
// 			historyA = "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
// 			historyB = "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."
// 			historyC = "In 1524, Henry VIII, King of England, received a “box of marmalade” from Mr. Hull of Exeter. This was probably marmelada, a solid quince paste from Portugal, still made and sold in southern Europe today. It became a favourite treat of Anne Boleyn and her ladies in waiting."
// 		)

// 		doc.WriteString(lipgloss.JoinHorizontal(
// 			lipgloss.Top,
// 			historyStyle.Align(lipgloss.Right).Render(historyA),
// 			historyStyle.Align(lipgloss.Center).Render(historyB),
// 			historyStyle.MarginRight(0).Render(historyC),
// 		))

// 		doc.WriteString("\n\n")
// 	}

// Status bar.
// {
//   let w = lipgloss.Width;

//   let lightDarkState = "Light";
// 	if (hasDarkBG) {
//     lightDarkState = "Dark";
// 	}

//   let statusKey = statusStyle.render("STATUS");
//   let encoding = encodingStyle.render("UTF-8");
//   let fishCake = fishCakeStyle.render("🍥 Fish Cake");
// 	let statusVal = statusText.
// 		width(width - w(statusKey) - w(encoding) - w(fishCake)).
// 		render("Ravishingly " + lightDarkState + "!")

// 	let bar = lipgloss.joinHorizontal(Top,
// 		statusKey,
// 		statusVal,
// 		encoding,
// 		fishCake,
// 	)

// 	doc += statusBarStyle.width(with)
// 	doc.WriteString(statusBarStyle.Width(width).Render(bar))
// }

if (physicalWidth > 0) {
  docStyle = docStyle.maxWidth(physicalWidth);
}

// Okay, let's print it. We use a special Lipgloss writer to downsample
// colors to the terminal's color palette. And, if output's not a TTY, we
// will remove color entirely.
console.log("%s", docStyle.render(doc));

/**
 * Creates a grid of hex color values by blending between four corner colors
 * @param {number} xSteps - Number of horizontal color steps
 * @param {number} ySteps - Number of vertical color steps
 * @returns {string[][]} 2D array of hex color strings
 */
function colorGrid(xSteps, ySteps) {
  // Define the four corner colors in hex
  const x0y0 = hexToRgb("#F25D94");
  const x1y0 = hexToRgb("#EDFF82");
  const x0y1 = hexToRgb("#643AFF");
  const x1y1 = hexToRgb("#14F9D5");

  // Create the left edge colors (x0)
  const x0 = [];
  for (let i = 0; i < ySteps; i++) {
    x0.push(blendLuv(x0y0, x0y1, i / ySteps));
  }

  // Create the right edge colors (x1)
  const x1 = [];
  for (let i = 0; i < ySteps; i++) {
    x1.push(blendLuv(x1y0, x1y1, i / ySteps));
  }

  // Create the color grid
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

/**
 * Converts hex color string to RGB object
 * @param {string} hex - Color in hex format (e.g. "#F25D94")
 * @returns {object} RGB color object with r, g, b properties
 */
function hexToRgb(hex) {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return { r, g, b };
}

// /**
//  * Converts RGB object to hex color string
//  * @param {object} rgb - RGB color object with r, g, b properties
//  * @returns {string} Color in hex format (e.g. "#F25D94")
//  */
function rgbToHex(rgb) {
  const toHex = (c) => {
    const hex = Math.round(c).toString(16);
    return hex.length === 1 ? "0" + hex : hex;
  };

  return `#${toHex(rgb.r)}${toHex(rgb.g)}${toHex(rgb.b)}`;
}

/**
 * Blend two colors in LUV color space (approximation)
 * Note: This is a simplified approximation as JavaScript doesn't have built-in LUV color space.
 * For accurate LUV blending, you would need to implement full color space conversion.
 * @param {object} color1 - First RGB color
 * @param {object} color2 - Second RGB color
 * @param {number} t - Blend factor (0 to 1)
 * @returns {object} Resulting blended RGB color
 */
function blendLuv(color1, color2, t) {
  // For simplicity, we'll do a simple linear interpolation in RGB space
  // For production, consider using a color library with proper LUV support
  return {
    r: color1.r + (color2.r - color1.r) * t,
    g: color1.g + (color2.g - color1.g) * t,
    b: color1.b + (color2.b - color1.b) * t,
  };
}

/**
 * Applies a color gradient to a string
 * @param {string} input - The input string to apply gradient to
 * @param {string|object} from - Starting color (hex string or RGB object)
 * @param {string|object} to - Ending color (hex string or RGB object)
 * @returns {string} String with gradient applied
 */
function applyGradient(input, from, to) {
  // Convert colors to RGB if they're hex strings
  const fromRgb = typeof from === "string" ? hexToRgb(from) : from;
  const toRgb = typeof to === "string" ? hexToRgb(to) : to;

  // Split the input string into graphemes (user-perceived characters)
  // Since we don't have the uniseg library, we'll use a simple approach
  // Note: This won't handle complex Unicode correctly like the original Go code
  // For production use, consider using a library like "grapheme-splitter"
  const chars = Array.from(input);

  let output = "";

  // Apply the gradient
  for (let i = 0; i < chars.length; i++) {
    // Calculate the blend ratio
    const t = chars.length > 1 ? i / (chars.length - 1) : 0;

    // Blend the colors
    const blendedColor = blendLuv(fromRgb, toRgb, t);

    // Convert back to hex
    const hex = rgbToHex(blendedColor);

    // Apply the style and add to output
    // Assuming baseStyle has a method similar to lipgloss.Style's Foreground() and Render()
    output += new Style().foreground(Color(hex)).render(chars[i]);
  }

  return output;
}
