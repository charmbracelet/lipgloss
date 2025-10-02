const {
  List,
  Style,
  Color,
  Table,
  Bullet,
  Dash,
  Roman,
  Center,
  Right,
  Left,
} = require("@charmland/lipgloss");

// Helper function to create color grid (simplified version)
function colorGrid(xSteps, ySteps) {
  const colors = [
    ["#F25D94", "#EDFF82"],
    ["#E85A4F", "#D4E157"],
    ["#D32F2F", "#C0CA33"],
    ["#7B1FA2", "#8BC34A"],
    ["#643AFF", "#14F9D5"]
  ];
  return colors;
}

// Styles
const purple = new Style()
  .foreground(Color("99"))
  .marginRight(1);

const pink = new Style()
  .foreground(Color("212"))
  .marginRight(1);

const base = new Style()
  .marginBottom(1)
  .marginLeft(1);

const faint = new Style().faint(true);

const dim = Color("250");
const highlight = Color("#EE6FF8");
const special = Color("#73F59F");

// Checklist functions
const checklistEnumStyle = (items, index) => {
  const itemIndex = items;
  if ([1, 2, 4].includes(itemIndex)) {
    return new Style()
      .foreground(special)
      .paddingRight(1);
  }
  return new Style().paddingRight(1);
};

const checklistEnum = (items, index) => {
  if ([1, 2, 4].includes(index)) {
    return "✓";
  }
  return "•";
};

const checklistStyle = (items, index) => {
  // The 'items' parameter is actually the item index, 'index' is always 0
  const itemIndex = items;
  
  // Items with ✓ should be strikethrough: indices 1, 2, 4 (Yuzu, Citron, Pomelo) and 1, 2, 4 (Claire's, Nyx, Milk)
  if ([1, 2, 4].includes(itemIndex)) {
    return new Style()
      .strikethrough(true)
      .foreground(Color("#696969"));
  }
  return new Style();
};

// Helper function to create styled checklist items
const createChecklistItem = (text, isChecked) => {
  if (isChecked) {
    const checkmark = new Style().foreground(special).render("✓");
    const itemText = new Style().strikethrough(true).foreground(Color("#696969")).render(" " + text);
    return checkmark + itemText;
  } else {
    return "• " + text;
  }
};

const colors = colorGrid(1, 5);

const titleStyle = new Style()
  .italic(true)
  .foreground(Color("#FFF7DB"));

const lipglossStyleFunc = (items, index) => {
  // The 'items' parameter is actually the item index, 'index' is always 0
  const itemIndex = items;
  const itemsLength = 5; // We know there are 5 items in this list
  
  if (itemIndex === itemsLength - 1) {
    return titleStyle
      .padding(1, 2)
      .margin(0, 0, 1, 0)
      .maxWidth(20)
      .background(Color(colors[itemIndex][0]));
  }
  return titleStyle
    .padding(0, 5 - itemIndex, 0, itemIndex + 2)
    .maxWidth(20)
    .background(Color(colors[itemIndex][0]));
};

const history = "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac.";

// Create the complex nested list
const l = new List()
  .enumeratorStyle(purple)
  .item("Lip Gloss")
  .item("Blush")
  .item("Eye Shadow")
  .item("Mascara")
  .item("Foundation")
  .item(
    new List()
      .enumeratorStyle(pink)
      .item("Citrus Fruits to Try")
      .item(
        new List()
          .enumerator(Dash)
          .enumeratorStyle(new Style().foreground(Color("0")).width(0))
          .item(createChecklistItem("Grapefruit", false))
          .item(createChecklistItem("Yuzu", true))
          .item(createChecklistItem("Citron", true))
          .item(createChecklistItem("Kumquat", false))
          .item(createChecklistItem("Pomelo", true))
      )
      .item("Actual Lip Gloss Vendors")
      .item(
        new List()
          .enumerator(Dash)
          .enumeratorStyle(new Style().foreground(Color("0")).width(0))
          .item(createChecklistItem("Glossier", false))
          .item(createChecklistItem("Claire's Boutique", true))
          .item(createChecklistItem("Nyx", true))
          .item(createChecklistItem("Mac", false))
          .item(createChecklistItem("Milk", true))
          .item(
            new List()
              .enumeratorStyle(purple)
              .enumerator(Dash)
              .itemStyleFunc(lipglossStyleFunc)
              .item("Lip Gloss")
              .item("Lip Gloss")
              .item("Lip Gloss")
              .item("Lip Gloss")
              .item(
                new List()
                  .enumeratorStyle(new Style().foreground(Color(colors[4][0])).marginRight(1))
                  .item("\nStyle Definitions for Nice Terminal Layouts\n─────")
                  .item("From Charm")
                  .item("https://github.com/charmbracelet/lipgloss")
                  .item(
                    new List()
                      .enumeratorStyle(new Style().foreground(Color(colors[3][0])).marginRight(1))
                      .item("Emperors: Julio-Claudian dynasty")
                      .item(
                        new Style().padding(1).render(
                          new List(
                            "Augustus",
                            "Tiberius",
                            "Caligula",
                            "Claudius",
                            "Nero"
                          ).enumerator(Roman).render()
                        )
                      )
                      .item(
                        new Style()
                          .bold(true)
                          .foreground(Color("#FAFAFA"))
                          .background(Color("#7D56F4"))
                          .align(Center, Center)
                          .padding(1, 3)
                          .margin(0, 1, 1, 1)
                          .width(40)
                          .render(history)
                      )
                      .item("Simple table placeholder")
                      .item("Documents")
                      .item(
                        new List()
                          .enumerator((items, i) => {
                            if (i === 1) {
                              return "│\n│";
                            }
                            return " ";
                          })
                          .itemStyleFunc((items, i) => {
                            if (i === 1) {
                              return base.foreground(highlight);
                            }
                            return base.foreground(dim);
                          })
                          .enumeratorStyleFunc((items, i) => {
                            if (i === 1) {
                              return new Style().foreground(highlight);
                            }
                            return new Style().foreground(dim);
                          })
                          .item("Foo Document\n" + faint.render("1 day ago"))
                          .item("Bar Document\n" + faint.render("2 days ago"))
                          .item("Baz Document\n" + faint.render("10 minutes ago"))
                          .item("Qux Document\n" + faint.render("1 month ago"))
                      )
                      .item("EOF")
                  )
                  .item("go get github.com/charmbracelet/lipgloss/list\n")
              )
              .item("See ya later")
          )
      )
      .item("List")
  )
  .item("xoxo, Charm_™");

console.log(l.render());