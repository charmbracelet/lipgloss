# Lipgloss WASM Tests

This directory contains test files for the Lipgloss WebAssembly bindings, specifically focusing on Table and TableData functionality.

## Running Tests

From the examples directory, run:

```bash
# Run all tests
npm test

# Run individual tests
npm run test:simple
npm run test:style
npm run test:comprehensive
npm run test:table-data

# Run all tests including Unicode tests (may fail due to WASM Unicode limitations)
npm run test:all
```

## Test Files

### `test-simple-table.js`
Tests basic table functionality without styling:
- Simple table creation
- Headers and rows
- Basic borders
- No style functions

### `test-style-function.js`
Tests table styling capabilities:
- Style function registration
- Header styling
- Row styling
- Color application

### `test-comprehensive.js`
Comprehensive test of TableData functionality:
- TableData creation and manipulation
- Cell access methods
- Row and column counting
- Dynamic data building
- Integration with Table styling

### `test-ascii-table.js`
Tests table functionality with ASCII-only content:
- Verifies that styling works without Unicode issues
- Useful for debugging Unicode-related problems

### `test-table-debug.js`
Debug version of the original table example:
- Helps isolate Unicode-related issues
- Provides detailed logging

## Test Runner

The `test-runner.js` script provides:
- Organized test execution
- Clear pass/fail reporting
- Error detection and reporting
- Timeout handling (30 seconds per test)

## Known Issues

- Unicode characters in table content may cause memory access issues in the WASM implementation
- This is a limitation of the current WASM bindings, not the TableData implementation
- ASCII-only content works reliably

## Adding New Tests

To add a new test:

1. Create a new test file in the `tests/` directory
2. Follow the naming convention: `test-[description].js`
3. Add the test to the `tests` array in `test-runner.js`
4. Optionally add a specific npm script in `package.json`

Test files should:
- Use descriptive console output
- Exit cleanly on success
- Throw errors or output error messages on failure
- Complete within 30 seconds