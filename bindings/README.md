# LipGloss.js

The lipgloss that you know and love, now in JavaScript.

> [!WARNING]
> LipGloss.js is stll experimental.

## Installation

```bash
npm i @charmland/lipgloss
```

## Example

```javascript
const { Table, TableData, Style, Color, List, Tree, Leaf, Bullet, RoundedEnumerator } = require("@charmland/lipgloss");

var s = new Style().foreground(Color("240")).render;
console.log(
  new Table()
    .wrap(false)
    .headers("Drink", "Description")
    .row("Bubble Tea", s("Milky"))
    .row("Milk Tea", s("Also milky"))
    .row("Actual milk", s("Milky as well"))
    .render(),
);

// TableData example - for more complex data management
const employeeData = new TableData(
  ["Employee ID", "Name", "Department"],
  ["001", "John Doe", "Engineering"],
  ["002", "Jane Smith", "Marketing"],
  ["003", "Mike Johnson", "Sales"]
);

console.log(
  new Table()
    .data(employeeData)
    .styleFunc((row, col) => {
      if (row === -1) {
        return new Style().foreground(Color("99")).bold(true);
      }
      return new Style().padding(0, 1);
    })
    .render()
);

// List example
const groceries = new List("Bananas", "Barley", "Cashews", "Milk")
  .enumerator(Bullet)
  .itemStyle(new Style().foreground(Color("255")));

console.log(groceries.render());

// Tree example with Leaf nodes
const makeupTree = new Tree()
  .root("⁜ Makeup")
  .child(
    "Glossier",
    "Fenty Beauty",
    new Tree().child(
      new Leaf("Gloss Bomb Universal Lip Luminizer"),
      new Leaf("Hot Cheeks Velour Blushlighter")
    ),
    new Leaf("Nyx"),
    new Leaf("Mac"),
    "Milk"
  )
  .enumerator(RoundedEnumerator)
  .enumeratorStyle(new Style().foreground(Color("63")).marginRight(1))
  .rootStyle(new Style().foreground(Color("35")))
  .itemStyle(new Style().foreground(Color("212")));

console.log(makeupTree.render());
```

## TableData

The `TableData` class provides a more structured way to manage table data compared to using individual rows. It's particularly useful when you need to:

- Build tables dynamically
- Access individual cells programmatically
- Manage large datasets
- Separate data logic from presentation

### Basic Usage

```javascript
const { Table, TableData, Style, Color } = require("@charmland/lipgloss");

// Create TableData with initial rows
const data = new TableData(
  ["Name", "Age", "City"],
  ["Alice", "25", "New York"],
  ["Bob", "30", "San Francisco"]
);

// Or build it incrementally
const data2 = new TableData()
  .append(["Product", "Price"])
  .append(["Laptop", "$999"])
  .append(["Mouse", "$25"]);

// Use with Table
const table = new Table()
  .data(data)
  .styleFunc((row, col) => {
    if (row === -1) return new Style().bold(true); // Header
    return new Style().padding(0, 1);
  })
  .render();
```

### TableData Methods

- `new TableData(...rows)` - Create with optional initial rows
- `.append(row)` - Add a single row (array of strings)
- `.rows(...rows)` - Add multiple rows at once
- `.at(row, col)` - Get value at specific position
- `.rowCount()` - Get number of rows
- `.columnCount()` - Get number of columns

See `examples/table-data.js` for more comprehensive examples.

## Environment Variables

You can control debug output using environment variables:

- `LIPGLOSS_DEBUG=true` - Enable all lipgloss debug output
- `DEBUG=lipgloss` - Enable lipgloss debug output (standard debug pattern)
- `DEBUG=*` - Enable all debug output

Example:
```bash
node your-script.js
```

## Compability

Lipgloss in JavaScript it's experimental and a lot of existing functionalities are not still ported to JavaScript.

### Size (100%)

| Function | Status |
| --- | --- |
| `Width` | ✅ |
| `Height` | ✅ |
| `Size` | ✅ |

### Color (100%)

| Function | Status |
| --- | --- |
| `Color` | ✅ |
| `NoColor` | ✅ |
| `Complete` | ✅ |
| `LightDark` | ✅ |
| `RGBA` | ✅ |

### Borders (100%)

| Function | Status |
| --- | --- |
| `NormalBorder` | ✅ |
| `RoundedBorder` | ✅ |
| `BlockBorder` | ✅ |
| `OuterHalfBlockBorder` | ✅ |
| `InnerHalfBlockBorder` | ✅ |
| `ThickBorder` | ✅ |
| `DoubleBorder` | ✅ |
| `HiddenBorder` | ✅ |
| `MarkdownBorder` | ✅ |
| `ASCIIBorder` | ✅ |

### Style (97.62%)

| Function | Status |
| --- | --- |
| `Foreground` | ✅ |
| `Background` | ✅ |
| `Width` | ✅ |
| `Height` | ✅ |
| `Align` | ✅ |
| `AlignHorizontal` | ✅ |
| `AlignVertical` | ✅ |
| `Padding` | ✅ |
| `PaddingLeft` | ✅ |
| `PaddingRight` | ✅ |
| `PaddingTop` | ✅ |
| `PaddingBottom` | ✅ |
| `ColorWhitespace` | ✅ |
| `Margin` | ✅ |
| `MarginLeft` | ✅ |
| `MarginRight` | ✅ |
| `MarginTop` | ✅ |
| `MarginBottom` | ✅ |
| `MarginBackground` | ✅ |
| `Border` | ✅ |
| `BorderStyle` | ✅ |
| `SetBorderRight` | ✅ |
| `SetBorderLeft` | ✅ |
| `SetBorderTop` | ✅ |
| `SetBorderBottom` | ✅ |
| `BorderForeground` | ✅ |
| `BorderTopForeground` | ✅ |
| `BorderRightForeground` | ✅ |
| `BorderBottomForeground` | ✅ |
| `BorderLeftForeground` | ✅ |
| `BorderBackground` | ✅ |
| `BorderTopBackground` | ✅ |
| `BorderRightBackground` | ✅ |
| `BorderBottomBackground` | ✅ |
| `BorderLeftBackground` | ✅ |
| `Inline` | ✅ |
| `MaxWidth` | ✅ |
| `MaxHeight` | ✅ |
| `TabWidth` | ✅ |
| `UnderlineSpaces` | ✅ |
| `Underline` | ✅ |
| `Reverse` | ✅ |
| `SetString` | ✅ |
| `Inherit` | ✅ |
| `Faint` | ✅ |
| `Italic` | ✅ |
| `Strikethrough` | ✅ |
| `StrikethroughSpaces` | ✅ |
| `Transform` | ⏳ |

### Table (100%)

| Function | Status |
| --- | --- |
| `table.New` | ✅ |
| `table.Rows` | ✅ |
| `table.Headers` | ✅ |
| `table.Render` | ✅ |
| `table.ClearRows` | ✅ |
| `table.BorderTop` | ✅ |
| `table.BorderBottom` | ✅ |
| `table.BorderLeft` | ✅ |
| `table.BorderRight` | ✅ |
| `table.BorderHeader` | ✅ |
| `table.BorderColumn` | ✅ |
| `table.BorderRow` | ✅ |
| `table.Width` | ✅ |
| `table.Height` | ✅ |
| `table.Offset` | ✅ |
| `table.String` | ✅ |
| `table.StyleFunc` | ✅ |
| `table.Data` | ✅ |
| `table.Border` | ✅ |
| `table.BorderStyle` | ✅ |

### List (100%)

| Function | Status |
| --- | --- |
| `Hidden` | ✅ |
| `Hide` | ✅ |
| `Offset` | ✅ |
| `Value` | ✅ |
| `String` | ✅ |
| `Indenter` | ✅ |
| `ItemStyle` | ✅ |
| `ItemStyleFunc` | ✅ |
| `EnumeratorStyle` | ✅ |
| `EnumeratorStyleFunc` | ✅ |
| `Item` | ✅ |
| `Items` | ✅ |
| `Enumerator` | ✅ |

### Tree (95%)

| Function | Status |
| --- | --- |
| `tree.New` | ✅ |
| `tree.Root` | ✅ |
| `tree.Child` | ✅ |
| `tree.Hidden` | ✅ |
| `tree.Hide` | ✅ |
| `tree.SetHidden` | ✅ |
| `tree.Offset` | ✅ |
| `tree.Value` | ✅ |
| `tree.SetValue` | ✅ |
| `tree.String` | ✅ |
| `tree.Render` | ✅ |
| `tree.EnumeratorStyle` | ✅ |
| `tree.EnumeratorStyleFunc` | ✅ |
| `tree.ItemStyle` | ✅ |
| `tree.ItemStyleFunc` | ✅ |
| `tree.RootStyle` | ✅ |
| `tree.Enumerator` | ✅ |
| `tree.Indenter` | ✅ |
| `tree.DefaultEnumerator` | ✅ |
| `tree.RoundedEnumerator` | ✅ |
| `tree.DefaultIndenter` | ✅ |
| `NewLeaf` | ✅ |
| `Leaf.Value` | ✅ |
| `Leaf.SetValue` | ✅ |
| `Leaf.Hidden` | ✅ |
| `Leaf.SetHidden` | ✅ |
| `Leaf.String` | ✅ |
| Custom Enumerators | ⏳ |
| Custom Indenters | ⏳ |

### Join (100%)

| Function | Status |
| --- | --- |
| `JoinHorizontal` | ✅ |
| `JoinVertical` | ✅ |

### Position (100%)

| Function | Status |
| --- | --- |
| `Center` | ✅ |
| `Right` | ✅ |
| `Bottom` | ✅ |
| `Top` | ✅ |
| `Left` | ✅ |
| `Place` | ✅ |

### Query (50%)

| Function | Status |
| --- | --- |
| `BackgroundColor` | ⏳ |
| `HasDarkBackground` | ✅ |

### Align (0%)

| Function | Status |
| --- | --- |
| `EnableLegacyWindowsANSI` | ⏳ |

## Contributing

We'd love to have you contribute! Please see the [contributing guidelines](https://github.com/charmbracelet/lipgloss/contribute) for more information.

### Testing

The JavaScript bindings include a comprehensive test suite:

```bash
# From bindings directory
npm test                      # Run all tests
npm run test:simple           # Basic table functionality
npm run test:comprehensive    # Complete TableData tests
npm run examples              # Showcase functionality

# Or from examples directory
cd examples && npm test
```

Tests are located in `examples/tests/` and cover:
- Basic table functionality
- TableData operations
- Style functions
- Unicode handling (with known limitations)

See [contributing][contribute].

[contribute]: https://github.com/charmbracelet/lipgloss/contribute

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/lipgloss/raw/master/LICENSE)

---

Part of [Charm](https://charm.land).

<a href="https://charm.land/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-banner-next.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source

[docs]: https://pkg.go.dev/github.com/charmbracelet/lipgloss?tab=doc
[wish]: https://github.com/charmbracelet/wish
[ssh-example]: examples/ssh
