const {
  List,
  Style,
  Color,
} = require("@charmland/lipgloss");

const checklistEnum = (items, index) => {
  console.log(`checklistEnum called with index ${index}`);
  if ([1, 2].includes(index)) {
    return "✓";
  }
  return "•";
};

const checklistStyle = (items, index) => {
  console.log(`checklistStyle called with index ${index}`);
  if ([1, 2].includes(index)) {
    return new Style()
      .strikethrough(true)
      .foreground(Color("#696969"));
  }
  return new Style();
};

const testList = new List()
  .enumerator(checklistEnum)
  .itemStyleFunc(checklistStyle)
  .item("First item")
  .item("Second item (should be checked)")
  .item("Third item (should be checked)")
  .item("Fourth item");

console.log(testList.render());