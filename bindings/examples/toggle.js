const {
  Tree,
  Style,
  Color,
  RoundedEnumerator,
} = require("@charmland/lipgloss");

function main() {
  // Create styles that match the Go version exactly
  const base = new Style().background(Color("57")).foreground(Color("225"));

  const block = base.padding(1, 3).margin(1, 3).width(40);

  const enumerator = new Style()
    .background(Color("57"))
    .foreground(Color("212"))
    .paddingRight(1);

  const toggle = new Style()
    .background(Color("57"))
    .foreground(Color("207"))
    .paddingRight(1);

  const dir = new Style().background(Color("57")).foreground(Color("225"));

  const file = new Style().background(Color("57")).foreground(Color("225"));

  // Create the tree structure exactly like the Go version
  const t = new Tree()
    .root(toggle.render("▼") + dir.render("~/charm"))
    .enumerator(RoundedEnumerator)
    .enumeratorStyle(enumerator)
    .child(
      toggle.render("▶") + dir.render("ayman"),
      new Tree()
        .root(toggle.render("▼") + dir.render("bash"))
        .child(
          new Tree()
            .root(toggle.render("▼") + dir.render("tools"))
            .child(file.render("zsh"), file.render("doom-emacs")),
        ),
      new Tree()
        .root(toggle.render("▼") + dir.render("carlos"))
        .child(
          new Tree()
            .root(toggle.render("▼") + dir.render("emotes"))
            .child(file.render("chefkiss.png"), file.render("kekw.png")),
        ),
      toggle.render("▶") + dir.render("maas"),
    );

  // Apply the block styling to the entire tree output
  console.log(block.render(t.render()));
}

main();
