Lip Gloss
=========

<p>
    <a href="https://github.com/charmbracelet/lipgloss/releases"><img src="https://img.shields.io/github/release/charmbracelet/lipgloss.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/lipgloss?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/charmbracelet/lipgloss/actions"><img src="https://github.com/charmbracelet/lipgloss/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Style definitions for building command line applications.

```go

import "github.com/charmbracelet/lipgloss"

var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#04B575")).
    Background(lipgloss.Color("#3C3C3C")).
    PaddingTop(2).
    PaddingLeft(4).
    Width(24)

fmt.Println(style.Render("Hello, kitty."))
```


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

### True Color (24-bit)

```go
lipgloss.Color("#0000FF") // good ol' 100% blue
lipgloss.Color("#04B575") // a green
lipgloss.Color("#3C3C3C") // a dark gray
```

The terminal's color profile will be automatically detected, and colors outside
the gamut of the current palette will be automatically coerced to their closest
available value.


### Adaptive Colors

You can also specify color options for light and dark backgrounds:

```go
lipgloss.AdaptiveColor{Light: "236", Dark: "248"}
```

The terminal's background color will automatically be detected and the
appropriate color will be chosen at runtime.


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
    RightMarginRight(4).
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


## Aligning Text

You can align paragraphs of text to the left, right, or center.

```go
var style = lipgloss.NewStyle().
    Width(24).
    Align(lipgloss.AlignLeft).  // align it left
    Align(lipgloss.AlignRight). // no wait, align it right
    Align(lipgloss.AlignCenter) // just kidding, align it in the center
```


## Setting the height

In addition to the width, you can also set a minimum height.

```go
var str = lipgloss.NewStyle().
    Width(24).
    Height(32).
    Foreground(lipgloss.Color("63")).
    Render("What’s for lunch?")
```


## Copying Styles

Just use `Copy()`:

```go
var style = lipgloss.NewStyle().Foreground(lipgloss.Color("219"))

var wildStyle = style.Copy().Blink(true)
```

`Copy()` performs a copy on the underlying data structure ensuring that you get
a true, dereferenced copy of a style. Without copying it's possible to mutate
styles.


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


## Unsetting rules

All rules can be unset:

```go
var style = lipgloss.NewStyle().
    Bold(true).                        // make it bold
    UnsetBold().                       // jk don't make it bold
    Background(lipgloss.Color("227")). // yellow background
    UnsetBackground()                  // never mind
```

Note that in the case of boolean values, unset values are different from falsey
ones in that values specifically set to `false` will be copied and inherited,
whereas unset values will not.


## Enforcing Rules

Sometimes, such as when developing a component, you want to make sure style
definitions respect their indended purpose in the UI. This is where `Inline`
and `MaxWidth`, and `MaxHeight` come in:

```go
// Force rendering onto a single line
someStyle.Inline().Render("yadda yadda")

// Also limit rendering to five cells
someStyle.Inline().MaxWidth(5).Render("yadda yadda")

// Limit rendering to a 5x5 cell block
someStyle.MaxWidth(5).MaxHeight(5).Render("yadda yadda")
```

## Rendering

Generally, you just call the `Render(string)` method on a `lipgloss.Style`:

```go
fmt.Println(lipgloss.NewStyle().Bold(true).Render("Hello, kitty."))
```

But you could also use the Stringer interface:

```go
var style = lipgloss.NewStyle().String("你好，猫咪。").Bold(true)

fmt.Printf("%s\n", style)
```


## Joining Paragraphs

There are also some utility functions for horizontally and vertically joining
paragraphs of text.

```go
// Horizontally join three paragraphs along their bottom edges
lipgloss.HorizontalJoin(lipgloss.JoinTop, paragraphA, paragraphB, paragraphC)

// Vertically join two paragraphs along their center axes
lipgloss.VerticalJoin(lipgloss.JoinCenter, paragraphA, paragraphB)

// Horizontally join three paragraphs, with the shorter ones aligning 20%
// from the top of the tallest.
lipgloss.HorizontalJoin(0.2, paragraphA, paragraphB, paragraphC)
```


## Under the Hood

Lip Gloss is built on the excellent [Termenv][termenv] and [Reflow][reflow]
libraries which deal with color and ANSI-aware text operations, respectively.
For many use cases Termenv and Reflow will be sufficient for your needs.

[termenv]: https://github.com/muesli/termenv
[reflow]: https://github.com/muesli/reflow


## Rendering Markdown

Lip Gloss is intended for use in building terminal user interfaces. For a more
document-centric rendering solution with support for things like lists, tables,
and syntax-highlighted code have a look at [Glamour][glamour], the
stylesheet-based Markdown renderer.

[glamour]: https://github.com/charmbracelet/glamour


## License

[MIT](https://github.com/charmbracelet/lipgloss/raw/master/LICENSE)


***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge-unrounded.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
