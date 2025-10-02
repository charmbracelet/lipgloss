# Quick Test Reference

## Running Tests

### From the main bindings directory:
```bash
cd /path/to/lipgloss/bindings

# Run all tests
npm test

# Run specific tests
npm run test:simple           # Basic table functionality
npm run test:style            # Style function tests  
npm run test:comprehensive    # Complete TableData tests
npm run test:table-data       # TableData examples
npm run test:ascii            # ASCII-only tests (no Unicode issues)
npm run test:all              # All tests including Unicode tests

# Run examples
npm run examples              # Showcase functionality
```

### From the examples directory:
```bash
cd /path/to/lipgloss/bindings/examples

# Same commands work here too
npm test
npm run test:simple
# etc...
```

### Using Task (from bindings directory):
```bash
task test                     # Runs npm test
task build                    # Build WASM
task dev                      # Build WASM with debug symbols
```

## Test Files Location

All test files are in `bindings/examples/tests/`:
- `test-simple-table.js` - Basic table tests
- `test-style-function.js` - Style function tests
- `test-comprehensive.js` - Complete TableData functionality
- `test-ascii-table.js` - ASCII-only tests
- `test-table-debug.js` - Debug utilities

## Examples Location

Example files are in `bindings/examples/`:
- `table-data.js` - Comprehensive TableData examples
- `table.js` - Original table example (may have Unicode issues)
- `size.js`, `color.js`, `list.js` - Other feature examples

## Known Issues

- Unicode characters in table content may cause WASM memory issues
- Use ASCII-only content for reliable results
- The `test:ascii` script tests functionality without Unicode