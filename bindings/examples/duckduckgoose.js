const {
  List,
  Style,
  Color,
} = require("../src/index.js");

// Create styles matching the original Go example
const enumStyle = new Style()
  .foreground(Color("#00d787"))
  .marginRight(1);

const itemStyle = new Style()
  .foreground(Color("255"));

// Recreate the exact original Go example behavior
const items = ["Duck", "Duck", "Duck", "Goose", "Duck"];
items.forEach((item) => {
  // Custom enumerator: "Honk →" for Goose, 7 spaces for Duck to align properly
  const enumerator = item === "Goose" ? "Honk →" : "      ";
  const styledEnum = enumStyle.render(enumerator);
  const styledItem = itemStyle.render(item);
  console.log(styledEnum + " " + styledItem);
});