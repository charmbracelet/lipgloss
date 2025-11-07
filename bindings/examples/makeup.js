const {
  Tree,
  Style,
  Color,
  RoundedEnumerator,
} = require("@charmland/lipgloss");

// Create styles matching the Go example
const enumeratorStyle = new Style()
  .foreground(Color("63"))
  .marginRight(1);

const rootStyle = new Style()
  .foreground(Color("35"));

const itemStyle = new Style()
  .foreground(Color("212"));

// Create the tree structure
const makeupTree = new Tree()
  .root("⁜ Makeup")
  .child(
    "Glossier",
    "Fenty Beauty",
    new Tree().child(
      "Gloss Bomb Universal Lip Luminizer",
      "Hot Cheeks Velour Blushlighter"
    ),
    "Nyx",
    "Mac",
    "Milk"
  )
  .enumerator(RoundedEnumerator)
  .enumeratorStyle(enumeratorStyle)
  .rootStyle(rootStyle)
  .itemStyle(itemStyle);

console.log(makeupTree.render());