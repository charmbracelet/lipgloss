const {
  lightDark,
  Color,
  roundedBorder,
  joinVertical,
  Center,
  Style,
  List,
} = require("@charmland/lipgloss");

// TODO: use hasDarkBackground()
let hasDarkBG = true;
let selectLightDark = lightDark(hasDarkBG);

let frameStyle = new Style()
  .border(roundedBorder())
  .borderForeground(selectLightDark(Color("#C5ADF9"), Color("#864EFF")))
  .padding(1, 3)
  .margin(1, 3);

let paragraphStyle = new Style().width(40).marginBottom(1).align(Center);
let textStyle = new Style().foreground(
  selectLightDark(Color("#696969"), Color("#bdbdbd")),
);
let keywordStyle = new Style()
  .foreground(selectLightDark(Color("#37CD96"), Color("#22C78A")))
  .bold(true);
let activeButton = new Style()
  .padding(0, 3)
  .background(Color("#FF6AD2"))
  .foreground(Color("#FFFCC2"));
let inactiveButton = new Style()
  .padding(0, 3)
  .background(selectLightDark(Color("#988F95"), Color("#978692")))
  .foreground(selectLightDark(Color("#FDFCE3"), Color("#FBFAE7")));

let text = paragraphStyle.render(
  textStyle.render("Are you sure you want to eat that ") +
    keywordStyle.render("moderatly ripe") +
    textStyle.render(" banana?"),
);
let buttons = activeButton.render("Yes") + "  " + inactiveButton.render("No");
let block = frameStyle.render(joinVertical(Center, text, buttons));
console.log(block);
