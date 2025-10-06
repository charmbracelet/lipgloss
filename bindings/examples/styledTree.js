const {
  Tree,
  Style,
  Color,
} = require("@charmland/lipgloss");

// Create styles matching the Go example
const purple = new Style()
  .foreground(Color("99"))
  .marginRight(1);

const pink = new Style()
  .foreground(Color("212"))
  .marginRight(1);

// Create the tree structure
const t = new Tree()
  .child(
    "Glossier",
    "Claire's Boutique",
    new Tree()
      .root("Nyx")
      .child("Lip Gloss", "Foundation")
      .enumeratorStyle(pink),
    "Mac",
    "Milk"
  )
  .enumeratorStyle(purple);

console.log(t.render());