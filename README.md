# Lip Gloss

<p>
   <img src="https://github.com/user-attachments/assets/c220c322-e8a3-42f0-a01d-3b6365935a3d" width="350"><br>
    <a href="https://github.com/charmbracelet/lipgloss/releases"><img src="https://img.shields.io/github/release/charmbracelet/lipgloss.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/charm.land/lipgloss/v2?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/charmbracelet/lipgloss/actions"><img src="https://github.com/charmbracelet/lipgloss/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Style definitions for nice terminal layouts. Built with TUIs in mind.

> [!TIP]
>
> Upgrading from v1? See the [upgrade guide](./UPGRADE_GUIDE_V2.md) for more
> details.

![Lip Gloss example](https://github.com/user-attachments/assets/92560e60-d70e-4ce0-b39e-a60bb933356b)

Lip Gloss takes an expressive, declarative approach to terminal rendering.
Users familiar with CSS will feel at home with Lip Gloss.

```go
import "charm.land/lipgloss/v2"

var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    PaddingTop(2).
    PaddingLeft(4).
    Width(22)

lipgloss.Println(style.Render("Hello, kitty"))
```

## Installation

```bash
go get charm.land/lipgloss/v2
```

> [!TIP]
> Using Lip Gloss with [Bubble Tea][tea]? Make sure you get all the latest v2s
> as they’ve been designed to work together.
>
> ```bash
> go get charm.land/bubbletea/v2@latest
> go get charm.land/bubbles/v2@latest
> go get charm.land/lipgloss/v2@latest
> ```

## Colors

Lip Gloss supports the following color profiles:

### ANSI 16 colors (4-bit)

```go
lipgloss.Color("5")  // magenta
lipgloss.Color("9")  // red
lipgloss.Color("12") // light blue
```

### ANSI 256 Colors (8-bit)

```go
lipgloss.Color("86")  // aqua
lipgloss.Color("201") // hot pink
lipgloss.Color("202") // orange
```

### True Color (16,777,216 colors; 24-bit)

```go
lipgloss.Color("#0000FF") // good ol' 100% blue
lipgloss.Color("#04B575") // a green
lipgloss.Color("#3C3C3C") // a dark gray
```

...as well as a 1-bit ASCII profile, which is black and white only.

There are also named constants for the 16 standard ANSI colors:

```go
lipgloss.Red
lipgloss.BrightCyan
lipgloss.Black
```

### Color Downsampling

One of the best things about Lip Gloss is that it can automatically downsample
colors to the best available profile, stripping colors (and ANSI) entirely when
output is not a TTY.

If you’re using Lip Gloss with Bubble Tea there’s nothing to do here:
downsampling is built into Bubble Tea v2. If you’re not using Bubble Tea, use
the Lip Gloss writer functions, which are a drop-in replacement for the `fmt`
package:

```go
s := someStyle.Render("Hello!")

// Downsample and print to stdout.
lipgloss.Println(s)

// Render to a variable.
downsampled := lipgloss.Sprint(s)

// Print to stderr.
lipgloss.Fprint(os.Stderr, s)
```

The full set: `Print`, `Println`, `Printf`, `Fprint`, `Fprintln`, `Fprintf`,
`Sprint`, `Sprintln`, `Sprintf`.

### Adaptive Colors

You can render different colors depending on whether the terminal has a light
or dark background:

```go
hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
lightDark := lipgloss.LightDark(hasDarkBG)

myColor := lightDark(lipgloss.Color("#D7FFAE"), lipgloss.Color("#D75FEE"))
```

#### With Bubble Tea

In Bubble Tea, request the background color, listen for a
`BackgroundColorMsg`, and respond accordingly:

```go
func (m model) Init() tea.Cmd {
    return tea.RequestBackgroundColor
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.BackgroundColorMsg:
        m.styles = newStyles(msg.IsDark())
        return m, nil
    }
}

func newStyles(bgIsDark bool) styles {
    lightDark := lipgloss.LightDark(bgIsDark)
    return styles{
        myHotStyle: lipgloss.NewStyle().Foreground(lightDark(
            lipgloss.Color("#f1f1f1"),
            lipgloss.Color("#333333"),
        )),
    }
}
```

#### Standalone

If you’re not using Bubble Tea you can perform the query manually:

```go
hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stderr)
lightDark := lipgloss.LightDark(hasDarkBG)

thisColor := lightDark(lipgloss.Color("#C5ADF9"), lipgloss.Color("#864EFF"))
thatColor := lightDark(lipgloss.Color("#37CD96"), lipgloss.Color("#22C78A"))

a := lipgloss.NewStyle().Foreground(thisColor).Render("this")
b := lipgloss.NewStyle().Foreground(thatColor).Render("that")

lipgloss.Fprintf(os.Stderr, "my fave colors are %s and %s", a, b)
```

### Complete Colors

For cases where you want to specify exact values for each color profile, use
the `Complete` helper:

```go
complete := lipgloss.Complete(profile)
myColor := complete(ansiColor, ansi256Color, trueColor)
```

### Compat Package

The `compat` package provides `AdaptiveColor`, `CompleteColor`, and
`CompleteAdaptiveColor` for a quicker migration from v1. These work by
looking at `stdin` and `stdout` on a global basis:

```go
import "charm.land/lipgloss/v2/compat"

color := compat.AdaptiveColor{
    Light: lipgloss.Color("#f1f1f1"),
    Dark:  lipgloss.Color("#cccccc"),
}
```

## Inline Formatting

Lip Gloss supports the usual ANSI text formatting options:

```go
var style = lipgloss.NewStyle().
    Bold(true).
    Italic(true).
    Faint(true).
    Blink(true).
    Strikethrough(true).
    Underline(true).
    Reverse(true)
```

### Underline Styles

Beyond simple on/off, underlines support multiple styles and custom colors:

```go
s := lipgloss.NewStyle().
    UnderlineStyle(lipgloss.UnderlineCurly).
    UnderlineColor(lipgloss.Color("#FF0000"))
```

Available styles: `UnderlineNone`, `UnderlineSingle`, `UnderlineDouble`,
`UnderlineCurly`, `UnderlineDotted`, `UnderlineDashed`.

### Hyperlinks

Styles can render clickable hyperlinks in supporting terminals:

```go
s := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#7B2FBE")).
    Hyperlink("https://charm.land")

lipgloss.Println(s.Render("Visit Charm"))
```

## Block-Level Formatting

Lip Gloss also supports rules for block-level formatting:

```go
// Padding
var style = lipgloss.NewStyle().
    PaddingTop(2).
    PaddingRight(4).
    PaddingBottom(2).
    PaddingLeft(4)

// Margins
var style = lipgloss.NewStyle().
    MarginTop(2).
    MarginRight(4).
    MarginBottom(2).
    MarginLeft(4)
```

There is also shorthand syntax for margins and padding, which follows the same
format as CSS:

```go
// 2 cells on all sides
lipgloss.NewStyle().Padding(2)

// 2 cells on the top and bottom, 4 cells on the left and right
lipgloss.NewStyle().Margin(2, 4)

// 1 cell on the top, 4 cells on the sides, 2 cells on the bottom
lipgloss.NewStyle().Padding(1, 4, 2)

// Clockwise, starting from the top: 2 cells on the top, 4 on the right, 3 on
// the bottom, and 1 on the left
lipgloss.NewStyle().Margin(2, 4, 3, 1)
```

You can also customize the characters used for padding and margin fill:

```go
s := lipgloss.NewStyle().
    Padding(1, 2).
    PaddingChar('·').
    Margin(1, 2).
    MarginChar('░')
```

## Aligning Text

You can align paragraphs of text to the left, right, or center.

```go
var style = lipgloss.NewStyle().
    Width(24).
    Align(lipgloss.Left).  // align it left
    Align(lipgloss.Right). // no wait, align it right
    Align(lipgloss.Center) // just kidding, align it in the center
```

## Width and Height

Setting a minimum width and height is simple and straightforward.

```go
var style = lipgloss.NewStyle().
    SetString("What’s for lunch?").
    Width(24).
    Height(32).
    Foreground(lipgloss.Color("63"))
```

## Borders

Adding borders is easy:

```go
// Add a purple, rectangular border
var style = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("63"))

// Set a rounded, yellow-on-purple border to the top and left
var anotherStyle = lipgloss.NewStyle().
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("228")).
    BorderBackground(lipgloss.Color("63")).
    BorderTop(true).
    BorderLeft(true)

// Make your own border
var myCuteBorder = lipgloss.Border{
    Top:         "._.:*:",
    Bottom:      "._.:*:",
    Left:        "|*",
    Right:       "|*",
    TopLeft:     "*",
    TopRight:    "*",
    BottomLeft:  "*",
    BottomRight: "*",
}
```

There are also shorthand functions for defining borders, which follow a similar
pattern to the margin and padding shorthand functions.

```go
// Add a thick border to the top and bottom
lipgloss.NewStyle().
    Border(lipgloss.ThickBorder(), true, false)

// Add a double border to the top and left sides. Rules are set clockwise
// from top.
lipgloss.NewStyle().
    Border(lipgloss.DoubleBorder(), true, false, false, true)
```

### Border Color Blending

Apply gradient colors to borders:

```go
s := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForegroundBlend(lipgloss.Color("#FF0000"), lipgloss.Color("#0000FF"))
```

For more on borders see [the docs][docs].

## Color Blending

Blend colors in one or two dimensions for gradient effects:

```go
// 1D gradient
colors := lipgloss.Blend1D(10, lipgloss.Color("#FF0000"), lipgloss.Color("#0000FF"))

// 2D gradient with rotation
colors := lipgloss.Blend2D(80, 24, 45.0, color1, color2, color3)
```

### Color Utilities

| Function              | Description                     |
| --------------------- | ------------------------------- |
| `Alpha(c, alpha)`     | Set a color’s alpha channel     |
| `Complementary(c)`    | Get the complementary color     |
| `Darken(c, percent)`  | Darken a color by a percentage  |
| `Lighten(c, percent)` | Lighten a color by a percentage |

## Copying Styles

Just use assignment:

```go
style := lipgloss.NewStyle().Foreground(lipgloss.Color("219"))

copiedStyle := style // this is a true copy

wildStyle := style.Blink(true) // this is also true copy, with blink added
```

Since `Style` is a pure value type, assigning a style to another effectively
creates a new copy of the style without mutating the original.

## Inheritance

Styles can inherit rules from other styles. When inheriting, only unset rules
on the receiver are inherited.

```go
var styleA = lipgloss.NewStyle().
    Foreground(lipgloss.Color("229")).
    Background(lipgloss.Color("63"))

// Only the background color will be inherited here, because the foreground
// color will have been already set:
var styleB = lipgloss.NewStyle().
    Foreground(lipgloss.Color("201")).
    Inherit(styleA)
```

## Unsetting Rules

All rules can be unset:

```go
var style = lipgloss.NewStyle().
    Bold(true).                        // make it bold
    UnsetBold().                       // jk don't make it bold
    Background(lipgloss.Color("227")). // yellow background
    UnsetBackground()                  // never mind
```

When a rule is unset, it won’t be inherited or copied.

## Enforcing Rules

Sometimes, such as when developing a component, you want to make sure style
definitions respect their intended purpose in the UI. This is where `Inline`
and `MaxWidth`, and `MaxHeight` come in:

```go
// Force rendering onto a single line, ignoring margins, padding, and borders.
someStyle.Inline(true).Render("yadda yadda")

// Also limit rendering to five cells
someStyle.Inline(true).MaxWidth(5).Render("yadda yadda")

// Limit rendering to a 5x5 cell block
someStyle.MaxWidth(5).MaxHeight(5).Render("yadda yadda")
```

## Tabs

The tab character (`\t`) is rendered differently in different terminals (often
as 8 spaces, sometimes 4). Because of this inconsistency, Lip Gloss converts
tabs to 4 spaces at render time. This behavior can be changed on a per-style
basis, however:

```go
style := lipgloss.NewStyle() // tabs will render as 4 spaces, the default
style = style.TabWidth(2)    // render tabs as 2 spaces
style = style.TabWidth(0)    // remove tabs entirely
style = style.TabWidth(lipgloss.NoTabConversion) // leave tabs intact
```

## Wrapping

The `Wrap` function wraps text while preserving ANSI styles and hyperlinks
across line boundaries:

```go
wrapped := lipgloss.Wrap(styledText, 40, " ")
```

## Rendering

Generally, you just call the `Render(string...)` method on a `lipgloss.Style`:

```go
style := lipgloss.NewStyle().Bold(true).SetString("Hello,")
lipgloss.Println(style.Render("kitty.")) // Hello, kitty.
lipgloss.Println(style.Render("puppy.")) // Hello, puppy.
```

But you could also use the Stringer interface:

```go
var style = lipgloss.NewStyle().SetString("你好，猫咪。").Bold(true)
lipgloss.Println(style) // 你好，猫咪。
```

## Compositing

Lip Gloss includes a cell-buffer based canvas system for composing layered
content:

```go
canvas := lipgloss.NewCanvas(80, 24)

layer := lipgloss.NewLayer(content).X(10).Y(5).Z(1).ID("panel")
compositor := lipgloss.NewCompositor(layer)

output := canvas.Compose(compositor).Render()
```

## Utilities

In addition to pure styling, Lip Gloss also ships with some utilities to help
assemble your layouts.

### Joining Paragraphs

Horizontally and vertically joining paragraphs is a cinch.

```go
// Horizontally join three paragraphs along their bottom edges
lipgloss.JoinHorizontal(lipgloss.Bottom, paragraphA, paragraphB, paragraphC)

// Vertically join two paragraphs along their center axes
lipgloss.JoinVertical(lipgloss.Center, paragraphA, paragraphB)

// Horizontally join three paragraphs, with the shorter ones aligning 20%
// from the top of the tallest
lipgloss.JoinHorizontal(0.2, paragraphA, paragraphB, paragraphC)
```

### Measuring Width and Height

Sometimes you’ll want to know the width and height of text blocks when building
your layouts.

```go
// Render a block of text.
var style = lipgloss.NewStyle().
    Width(40).
    Padding(2)
var block string = style.Render(someLongString)

// Get the actual, physical dimensions of the text block.
width := lipgloss.Width(block)
height := lipgloss.Height(block)

// Here's a shorthand function.
w, h := lipgloss.Size(block)
```

### Placing Text in Whitespace

Sometimes you’ll simply want to place a block of text in whitespace.

```go
// Center a paragraph horizontally in a space 80 cells wide. The height of
// the block returned will be as tall as the input paragraph.
block := lipgloss.PlaceHorizontal(80, lipgloss.Center, fancyStyledParagraph)

// Place a paragraph at the bottom of a space 30 cells tall. The width of
// the text block returned will be as wide as the input paragraph.
block := lipgloss.PlaceVertical(30, lipgloss.Bottom, fancyStyledParagraph)

// Place a paragraph in the bottom right corner of a 30x80 cell space.
block := lipgloss.Place(30, 80, lipgloss.Right, lipgloss.Bottom, fancyStyledParagraph)
```

You can also style the whitespace. For details, see [the docs][docs].

## Rendering Tables

Lip Gloss ships with a table rendering sub-package.

```go
import "charm.land/lipgloss/v2/table"
```

Define some rows of data.

```go
rows := [][]string{
    {"Chinese", "您好", "你好"},
    {"Japanese", "こんにちは", "やあ"},
    {"Arabic", "أهلين", "أهلا"},
    {"Russian", "Здравствуйте", "Привет"},
    {"Spanish", "Hola", "¿Qué tal?"},
}
```

Use the table package to style and render the table.

```go
var (
    purple    = lipgloss.Color("99")
    gray      = lipgloss.Color("245")
    lightGray = lipgloss.Color("241")

    headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
    cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14)
    oddRowStyle  = cellStyle.Foreground(gray)
    evenRowStyle = cellStyle.Foreground(lightGray)
)

t := table.New().
    Border(lipgloss.NormalBorder()).
    BorderStyle(lipgloss.NewStyle().Foreground(purple)).
    StyleFunc(func(row, col int) lipgloss.Style {
        switch {
        case row == table.HeaderRow:
            return headerStyle
        case row%2 == 0:
            return evenRowStyle
        default:
            return oddRowStyle
        }
    }).
    Headers("LANGUAGE", "FORMAL", "INFORMAL").
    Rows(rows...)

// You can also add tables row-by-row
t.Row("English", "You look absolutely fabulous.", "How's it going?")
```

Print the table.

```go
lipgloss.Println(t)
```

![Table Example](https://github.com/charmbracelet/lipgloss/assets/42545625/6e4b70c4-f494-45da-a467-bdd27df30d5d)

### Table Borders

There are helpers to generate tables in markdown or ASCII style:

#### Markdown Table

```go
table.New().Border(lipgloss.MarkdownBorder()).BorderTop(false).BorderBottom(false)
```

```
| LANGUAGE |    FORMAL    | INFORMAL  |
|----------|--------------|-----------|
| Chinese  | Nǐn hǎo      | Nǐ hǎo    |
| French   | Bonjour      | Salut     |
| Russian  | Zdravstvuyte | Privet    |
| Spanish  | Hola         | ¿Qué tal? |
```

#### ASCII Table

```go
table.New().Border(lipgloss.ASCIIBorder())
```

```
+----------+--------------+-----------+
| LANGUAGE |    FORMAL    | INFORMAL  |
+----------+--------------+-----------+
| Chinese  | Nǐn hǎo      | Nǐ hǎo    |
| French   | Bonjour      | Salut     |
| Russian  | Zdravstvuyte | Privet    |
| Spanish  | Hola         | ¿Qué tal? |
+----------+--------------+-----------+
```

For more on tables see [the docs][docs] and [examples](https://github.com/charmbracelet/lipgloss/tree/master/examples/table).

## Rendering Lists

Lip Gloss ships with a list rendering sub-package.

```go
import "charm.land/lipgloss/v2/list"
```

Define a new list.

```go
l := list.New("A", "B", "C")
```

Print the list.

```go
lipgloss.Println(l)

// • A
// • B
// • C
```

Lists have the ability to nest.

```go
l := list.New(
    "A", list.New("Artichoke"),
    "B", list.New("Baking Flour", "Bananas", "Barley", "Bean Sprouts"),
    "C", list.New("Cashew Apple", "Cashews", "Coconut Milk", "Curry Paste", "Currywurst"),
    "D", list.New("Dill", "Dragonfruit", "Dried Shrimp"),
    "E", list.New("Eggs"),
    "F", list.New("Fish Cake", "Furikake"),
    "J", list.New("Jicama"),
    "K", list.New("Kohlrabi"),
    "L", list.New("Leeks", "Lentils", "Licorice Root"),
)
```

Print the list.

```go
lipgloss.Println(l)
```

<p align="center">
<img width="600" alt="image" src="https://github.com/charmbracelet/lipgloss/assets/42545625/0dc9f440-0748-4151-a3b0-7dcf29dfcdb0">
</p>

Lists can be customized via their enumeration function as well as using
`lipgloss.Style`s.

```go
enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

l := list.New(
    "Glossier",
    "Claire's Boutique",
    "Nyx",
    "Mac",
    "Milk",
    ).
    Enumerator(list.Roman).
    EnumeratorStyle(enumeratorStyle).
    ItemStyle(itemStyle)
```

Print the list.

<p align="center">
<img width="600" alt="List example" src="https://github.com/charmbracelet/lipgloss/assets/42545625/360494f1-57fb-4e13-bc19-0006efe01561">
</p>

In addition to the predefined enumerators (`Arabic`, `Alphabet`, `Roman`, `Bullet`, `Tree`),
you may also define your own custom enumerator:

```go
l := list.New("Duck", "Duck", "Duck", "Duck", "Goose", "Duck", "Duck")

func DuckDuckGooseEnumerator(l list.Items, i int) string {
    if l.At(i).Value() == "Goose" {
        return "Honk →"
    }
    return ""
}

l = l.Enumerator(DuckDuckGooseEnumerator)
```

Print the list:

<p align="center">
<img width="600" alt="image" src="https://github.com/charmbracelet/lipgloss/assets/42545625/157aaf30-140d-4948-9bb4-dfba46e5b87e">
</p>

If you need, you can also build lists incrementally:

```go
l := list.New()

for i := 0; i < repeat; i++ {
    l.Item("Lip Gloss")
}
```

## Rendering Trees

Lip Gloss ships with a tree rendering sub-package.

```go
import "charm.land/lipgloss/v2/tree"
```

Define a new tree.

```go
t := tree.Root(".").
    Child("A", "B", "C")
```

Print the tree.

```go
lipgloss.Println(t)

// .
// ├── A
// ├── B
// └── C
```

Trees have the ability to nest.

```go
t := tree.Root(".").
    Child("macOS").
    Child(
        tree.New().
            Root("Linux").
            Child("NixOS").
            Child("Arch Linux (btw)").
            Child("Void Linux"),
        ).
    Child(
        tree.New().
            Root("BSD").
            Child("FreeBSD").
            Child("OpenBSD"),
    )
```

Print the tree.

```go
lipgloss.Println(t)
```

<p align="center">
<img width="663" alt="Tree Example (simple)" src="https://github.com/user-attachments/assets/5ef14eb8-a5d4-4f94-8834-e15d1e714f89">
</p>

Trees can be customized via their enumeration function as well as using
`lipgloss.Style`s.

```go
enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

t := tree.
    Root("⁜ Makeup").
    Child(
        "Glossier",
        "Fenty Beauty",
        tree.New().Child(
            "Gloss Bomb Universal Lip Luminizer",
            "Hot Cheeks Velour Blushlighter",
        ),
        "Nyx",
        "Mac",
        "Milk",
    ).
    Enumerator(tree.RoundedEnumerator).
    EnumeratorStyle(enumeratorStyle).
    RootStyle(rootStyle).
    ItemStyle(itemStyle)
```

Print the tree.

<p align="center">
<img width="663" alt="Tree Example (makeup)" src="https://github.com/user-attachments/assets/06d12d87-744a-4c89-bd98-45de9094a97e">
</p>

The predefined enumerators for trees are `DefaultEnumerator` and `RoundedEnumerator`.

If you need, you can also build trees incrementally:

```go
t := tree.New()

for i := 0; i < repeat; i++ {
    t.Child("Lip Gloss")
}
```

## What about [Bubble Tea][tea]?

Lip Gloss doesn’t replace Bubble Tea. Rather, it is an excellent Bubble Tea
companion. It was designed to make assembling terminal user interface views as
simple and fun as possible so that you can focus on building your application
instead of concerning yourself with low-level layout details.

In simple terms, you can use Lip Gloss to help build your Bubble Tea views.

[tea]: https://github.com/charmbracelet/bubbletea

## Rendering Markdown

For a more document-centric rendering solution with support for things like
lists, tables, and syntax-highlighted code have a look at [Glamour][glamour],
the stylesheet-based Markdown renderer.

[glamour]: https://github.com/charmbracelet/glamour

## Contributing

See [contributing][contribute].

[contribute]: https://github.com/charmbracelet/lipgloss/contribute

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

- [Discord](https://charm.land/chat)
- [Matrix](https://charm.land/matrix)

## License

[MIT](https://github.com/charmbracelet/lipgloss/raw/master/LICENSE)

---

Part of [Charm](https://charm.land).

<a href="https://charm.land/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-banner-next.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source

[docs]: https://pkg.go.dev/charm.land/lipgloss/v2?tab=doc
