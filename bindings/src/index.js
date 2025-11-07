import * as fs from "fs";
import path from "node:path";
import { fileURLToPath } from "url";
import {
  passStringArray,
  readGoString,
  insertString,
  clearString,
  ensureMemorySize,
} from "./core.js";
import { whichSidesBool, whichSidesInt, whichSidesColor } from "./helpers.js";
import { callStyleFunction, registerStyleFunction } from "./styleFunc.js";
import "./wasmExec.js";

const go = new Go();

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const wasmBuffer = fs.readFileSync(path.join(__dirname, "./lipgloss.wasm"));

// Synchronous initialization and instantiation
const wasm = new WebAssembly.Module(wasmBuffer);

go.importObject.env["callStyleFunc"] = callStyleFunction;

const instance = new WebAssembly.Instance(wasm, {
  ...go.importObject,
});

go.run(instance);

// Required since wasm runtime is not tty
process.env["TTY_FORCE"] = process.stdout.isTTY;
let envArr = Object.entries(process.env).map(([key, value]) => {
  return `${key}=${value}`;
});
let envArrPtr = passStringArray(instance, envArr);
// Configure color profile through wasm environment variables
instance.exports.DetectFromEnvVars(envArrPtr, envArr.length);

const Top = instance.exports.PositionTop();
const Bottom = instance.exports.PositionBottom();
const Right = instance.exports.PositionRight();
const Left = instance.exports.PositionLeft();
const Center = instance.exports.PositionCenter();

function Width(str) {
  let addr = insertString(str, instance);
  return instance.exports.Width(addr, str.length);
}

function Height(str) {
  let addr = insertString(str, instance);
  return instance.exports.Height(addr, str.length);
}

function join(isVertical, position, ...strings) {
  // Early return for empty cases
  if (strings.length === 0) return "";
  if (strings.length === 1) return strings[0];

  // Access the WebAssembly memory and necessary exports
  const memory = instance.exports.memory;
  const malloc = instance.exports.malloc;
  const free = instance.exports.free;
  let wasmJoin = instance.exports.JoinHorizontal;
  if (isVertical) {
    wasmJoin = instance.exports.JoinVertical;
  }

  // Initialize variables outside try block so they're available in finally
  let ptrArrayPtr = 0;
  let allocatedPtrs = [];

  try {
    // Allocate memory for the string pointers and lengths array
    const stringsCount = strings.length;
    const pointerArraySize = stringsCount * 2 * 4; // Each string needs 2 uint32 values (ptr, len) * 4 bytes each
    ptrArrayPtr = malloc(pointerArraySize);

    if (!ptrArrayPtr) {
      throw new Error("Failed to allocate memory for pointer array");
    }

    // Create a view of the pointer array in the WASM memory
    const ptrArray = new Uint32Array(
      memory.buffer,
      ptrArrayPtr,
      stringsCount * 2,
    );

    // Process each input string
    for (let i = 0; i < strings.length; i++) {
      const str = strings[i];
      if (!str.length) {
        continue;
      }

      const strBytes = new TextEncoder().encode(str);
      const strLen = strBytes.length;

      // Allocate memory for the string data
      const strPtr = malloc(strLen);
      if (!strPtr) {
        throw new Error(`Failed to allocate memory for string ${i}`);
      }
      allocatedPtrs.push(strPtr);

      // Copy string data to WASM memory
      const strBuffer = new Uint8Array(memory.buffer, strPtr, strLen);
      strBuffer.set(strBytes);

      // Store pointer and length in the pointer array
      ptrArray[i * 2] = strPtr;
      ptrArray[i * 2 + 1] = strLen;
    }

    // Call the WASM function
    // The result is a Go string, which in WASM is returned as an integer representing
    // a pointer to a structure containing both the string pointer and length
    const resultGoStrPtr = wasmJoin(position, ptrArrayPtr, stringsCount);

    // Extract the result string using Go's string representation
    // We need to read the pointer and length of the string from memory
    let result = "";

    if (resultGoStrPtr) {
      // In Go's WebAssembly ABI for strings:
      // - First 4 bytes at resultGoStrPtr: pointer to string data
      // - Next 4 bytes: length of string
      const resultPtrView = new DataView(memory.buffer, resultGoStrPtr, 8);
      const resultDataPtr = resultPtrView.getUint32(0, true); // true = little endian
      const resultLength = resultPtrView.getUint32(4, true);

      if (resultLength > 0 && resultDataPtr) {
        const resultBytes = new Uint8Array(
          memory.buffer,
          resultDataPtr,
          resultLength,
        );
        result = new TextDecoder().decode(resultBytes);
      }
    }

    return result;
  } catch (err) {
    console.error(
      "Error in join",
      isVertical ? "Vertical" : "Horizontal",
      ": ",
      err,
    );
    return "";
  } finally {
    // Clean up allocated memory
    // Note: We don't free the resultGoStrPtr as it's managed by Go's runtime
    for (const ptr of allocatedPtrs) {
      if (ptr) free(ptr);
    }
    if (ptrArrayPtr) free(ptrArrayPtr);
  }
}

function Size(str) {
  return {
    height: Height(str),
    width: Width(str),
  };
}

function normalBorder() {
  return instance.exports.BorderNormalBorder();
}

function roundedBorder() {
  return instance.exports.BorderRoundedBorder();
}

function blockBorder() {
  return instance.exports.BorderBlockBorder();
}

function outerHalfBlockBorder() {
  return instance.exports.BorderOuterHalfBlockBorder();
}

function innerHalfBlockBorder() {
  return instance.exports.BorderInnerHalfBlockBorder();
}

function thickBorder() {
  return instance.exports.BorderThickBorder();
}

function doubleBorder() {
  return instance.exports.BorderDoubleBorder();
}

function hiddenBorder() {
  return instance.exports.BorderHiddenBorder();
}

function markdownBorder() {
  return instance.exports.BorderMarkdownBorder();
}

function ASCIIBorder() {
  return instance.exports.BorderASCIIBorder();
}

function joinHorizontal(position, ...strings) {
  return join(false, position, ...strings);
}

function joinVertical(position, ...strings) {
  return join(true, position, ...strings);
}

function lightDark(isDark) {
  return (light, dark) => {
    if (isDark) {
      return dark;
    }
    return light;
  };
}

// function Color(str) {
//   let addr = insertString(str, instance);
//   let result = instance.exports.Color(addr, str.length);
//   clearString(addr, str.length, instance);
//   return result;
// }

// wasmPlace(width, height int32, hPos, vPos *Position, str string, opts ...WhitespaceOption) *string {
function Place(width, height, hPos, vPos, str) {
  const memory = instance.exports.memory;

  let addr = insertString(str, instance);
  let resultGoStrPtr = instance.exports.PositionPlace(
    width,
    height,
    hPos,
    vPos,
    addr,
    str.length,
  );

  let result = "";

  if (resultGoStrPtr) {
    // In Go's WebAssembly ABI for strings:
    // - First 4 bytes at resultGoStrPtr: pointer to string data
    // - Next 4 bytes: length of string
    const resultPtrView = new DataView(memory.buffer, resultGoStrPtr, 8);
    const resultDataPtr = resultPtrView.getUint32(0, true); // true = little endian
    const resultLength = resultPtrView.getUint32(4, true);

    if (resultLength > 0 && resultDataPtr) {
      const resultBytes = new Uint8Array(
        memory.buffer,
        resultDataPtr,
        resultLength,
      );
      result = new TextDecoder().decode(resultBytes);
    }
  }

  return result;
}

function Color(str) {
  let addr = insertString(str, instance);
  let result = instance.exports.Color(addr, str.length);
  clearString(addr, str.length, instance);
  return result;
}

const NoColor = instance.exports.NoColor();

// List enumerators
const Alphabet = instance.exports.ListEnumeratorAlphabet();
const Arabic = instance.exports.ListEnumeratorArabic();
const Bullet = instance.exports.ListEnumeratorBullet();
const Dash = instance.exports.ListEnumeratorDash();
const Roman = instance.exports.ListEnumeratorRoman();
const Asterisk = instance.exports.ListEnumeratorAsterisk();

// Tree enumerators
const DefaultEnumerator = instance.exports.TreeEnumeratorDefault();
const RoundedEnumerator = instance.exports.TreeEnumeratorRounded();

// Tree indenters
const DefaultIndenter = instance.exports.TreeIndenterDefault();

class Style {
  constructor() {
    this.addr = instance.exports.StyleNewStyle();
    this.render = this.render.bind(this);
  }
  render(...strs) {
    let addr = this.addr;
    const memory = instance.exports.memory;

    if (strs.length > 0) {
      let ptr = passStringArray(instance, strs);
      instance.exports.StyleJoinString(addr, ptr, strs.length);
    }

    // Estimate memory needed for rendering and ensure we have enough
    const totalInputSize = strs.reduce((sum, str) => sum + str.length, 0);
    const estimatedOutputSize = Math.max(totalInputSize * 2, 32768); // At least 32KB for complex rendering
    const currentMemory = memory.buffer.byteLength;
    const memoryNeeded = currentMemory + estimatedOutputSize;
    
    ensureMemorySize(instance, memoryNeeded);

    let resultGoStrPtr = instance.exports.StyleRender(addr);
    let result = "";

    if (resultGoStrPtr) {
      // In Go's WebAssembly ABI for strings:
      // - First 4 bytes at resultGoStrPtr: pointer to string data
      // - Next 4 bytes: length of string
      const resultPtrView = new DataView(memory.buffer, resultGoStrPtr, 8);
      const resultDataPtr = resultPtrView.getUint32(0, true); // true = little endian
      const resultLength = resultPtrView.getUint32(4, true);

      if (resultLength > 0 && resultDataPtr) {
        const resultBytes = new Uint8Array(
          memory.buffer,
          resultDataPtr,
          resultLength,
        );
        result = new TextDecoder().decode(resultBytes);
      }
    }

    instance.exports.StyleClearValue(addr);
    return result;
  }
  setString(...strs) {
    let addr = this.addr;
    let ptr = passStringArray(instance, strs);
    instance.exports.StyleSetString(addr, ptr, strs.length);
    return this;
  }
  string() {
    return this.render();
  }
  bold(b) {
    instance.exports.StyleBold(this.addr, b);
    return this;
  }
  italic(b) {
    instance.exports.StyleItalic(this.addr, b);
    return this;
  }
  inherit(s) {
    instance.exports.StyleInherit(this.addr, s.addr);
    return this;
  }
  strikethrough(b) {
    instance.exports.StyleStrikethrough(this.addr, b);
    return this;
  }
  underline(b) {
    instance.exports.StyleUnderline(this.addr, b);
    return this;
  }
  blink(b) {
    instance.exports.StyleBlink(this.addr, b);
    return this;
  }
  reverse(b) {
    instance.exports.StyleReverse(this.addr, b);
    return this;
  }
  faint(b) {
    instance.exports.StyleFaint(this.addr, b);
    return this;
  }
  foreground(color) {
    instance.exports.StyleForeground(this.addr, color);
    return this;
  }
  background(color) {
    instance.exports.StyleBackground(this.addr, color);
    return this;
  }
  width(val) {
    instance.exports.StyleWidth(this.addr, val);
    return this;
  }
  height(val) {
    instance.exports.StyleHeight(this.addr, val);
    return this;
  }
  inline(b) {
    instance.exports.StyleInline(this.addr, b);
    return this;
  }
  underlineSpaces(b) {
    instance.exports.StyleUnderlineSpaces(this.addr, b);
    return this;
  }
  strikethroughSpaces(b) {
    instance.exports.StyleStrikethroughSpaces(this.addr, b);
    return this;
  }
  align(...positions) {
    if (positions.length > 0) {
      instance.exports.StyleAlignHorizontal(this.addr, positions[0]);
    }
    if (positions.length > 1) {
      instance.exports.StyleAlignVertical(this.addr, positions[1]);
    }
    return this;
  }
  alignHorizontal(position) {
    instance.exports.StyleAlignHorizontal(position);
    return this;
  }
  alignVertical(position) {
    instance.exports.StyleAlignVertical(position);
    return this;
  }
  border(b, ...sides) {
    instance.exports.StyleBorder(this.addr, b);

    let [top, right, bottom, left, ok] = whichSidesBool(...sides);
    if (!ok) {
      top = true;
      right = true;
      bottom = true;
      left = true;
    }

    instance.exports.SetBorderTop(this.addr, top);
    instance.exports.SetBorderRight(this.addr, right);
    instance.exports.SetBorderBottom(this.addr, bottom);
    instance.exports.SetBorderLeft(this.addr, left);

    return this;
  }
  borderStyle(b) {
    instance.exports.StyleBorder(this.addr, b);
    return this;
  }
  padding(...paddings) {
    if (paddings.length > 0 && paddings.length <= 4) {
      const [top, right, bottom, left, ok] = whichSidesInt(...paddings);
      if (ok) {
        instance.exports.StylePaddingTop(this.addr, top);
        instance.exports.StylePaddingRight(this.addr, right);
        instance.exports.StylePaddingBottom(this.addr, bottom);
        instance.exports.StylePaddingLeft(this.addr, left);
      }
    }

    return this;
  }
  paddingLeft(val) {
    instance.exports.StylePaddingLeft(this.addr, val);
    return this;
  }
  paddingRight(val) {
    instance.exports.StylePaddingRight(this.addr, val);
    return this;
  }
  paddingTop(val) {
    instance.exports.StylePaddingTop(this.addr, val);
    return this;
  }
  paddingBottom(val) {
    instance.exports.StylePaddingBottom(this.addr, val);
    return this;
  }
  colorWhitespace(val) {
    instance.exports.StyleColorWhitespace(this.addr, val);
    return this;
  }
  margin(...margins) {
    if (margins.length > 0) {
      let [top, right, bottom, left, ok] = whichSidesInt(...margins);
      if (ok) {
        instance.exports.StyleMarginTop(this.addr, top);
        instance.exports.StyleMarginRight(this.addr, right);
        instance.exports.StyleMarginBottom(this.addr, bottom);
        instance.exports.StyleMarginLeft(this.addr, left);
      }
    }
    return this;
  }
  marginTop(val) {
    instance.exports.StyleMarginTop(this.addr, val);
    return this;
  }
  marginRight(val) {
    instance.exports.StyleMarginRight(this.addr, val);
    return this;
  }
  marginLeft(val) {
    instance.exports.StyleMarginLeft(this.addr, val);
    return this;
  }
  marginBottom(val) {
    instance.exports.StyleMarginBottom(this.addr, val);
    return this;
  }
  marginBackground(val) {
    instance.exports.StyleMarginBackground(this.addr, val);
    return this;
  }
  borderForeground(...colors) {
    if (colors.length > 0) {
      let [top, right, bottom, left, ok] = whichSidesColor(...colors);
      if (ok) {
        instance.exports.StyleBorderTopForeground(this.addr, top);
        instance.exports.StyleBorderRightForeground(this.addr, right);
        instance.exports.StyleBorderBottomForeground(this.addr, bottom);
        instance.exports.StyleBorderLeftForeground(this.addr, left);
      }
    }

    return this;
  }
  borderTopForeground(color) {
    instance.exports.StyleBorderTopForeground(this.addr, color);
    return this;
  }
  borderRightForeground(color) {
    instance.exports.StyleBorderRightForeground(this.addr, color);
    return this;
  }
  borderLeftForeground(color) {
    instance.exports.StyleBorderLeftForeground(this.addr, color);
    return this;
  }
  borderBottomForeground(color) {
    instance.exports.StyleBorderBottomForeground(this.addr, color);
    return this;
  }
  borderBackground(...colors) {
    if (colors.length > 0) {
      let [top, right, bottom, left, ok] = whichSidesColor(...colors);
      if (ok) {
        instance.exports.StyleBorderTopBackground(this.addr, top);
        instance.exports.StyleBorderRightBackground(this.addr, right);
        instance.exports.StyleBorderBottomBackground(this.addr, bottom);
        instance.exports.StyleBorderLeftBackground(this.addr, left);
      }
    }

    return this;
  }
  borderTopBackground(color) {
    instance.exports.StyleBorderTopBackground(this.addr, color);
    return this;
  }
  borderRightBackground(color) {
    instance.exports.StyleBorderRightBackground(this.addr, color);
    return this;
  }
  borderLeftBackground(color) {
    instance.exports.StyleBorderLeftBackground(this.addr, color);
    return this;
  }
  borderBottomBackground(color) {
    instance.exports.StyleBorderBottomBackground(this.addr, color);
    return this;
  }
  borderRight(val) {
    instance.exports.SetBorderRight(this.addr, val);
    return this;
  }
  borderLeft(val) {
    instance.exports.SetBorderLeft(this.addr, val);
    return this;
  }
  borderTop(val) {
    instance.exports.SetBorderTop(this.addr, val);
    return this;
  }
  borderBottom(val) {
    instance.exports.SetBorderBottom(this.addr, val);
    return this;
  }
  maxWidth(val) {
    instance.exports.StyleMaxWidth(this.addr, val);
    return this;
  }
  maxHeight(val) {
    instance.exports.StyleMaxHeight(this.addr, val);
    return this;
  }
  tabWidth(val) {
    instance.exports.StyleTabWidth(this.addr, val);
    return this;
  }
}

class TableData {
  constructor(...rows) {
    this.addr = instance.exports.TableDataNew();
    if (rows.length > 0) {
      this.rows(...rows);
    }
  }
  
  append(row) {
    if (Array.isArray(row)) {
      let ptr = passStringArray(instance, row);
      instance.exports.TableDataAppend(this.addr, ptr, row.length);
    }
    return this;
  }
  
  rows(...rows) {
    for (const row of rows) {
      this.append(row);
    }
    return this;
  }
  
  at(row, col) {
    let ptr = instance.exports.TableDataAtPtr(this.addr, row, col);
    let ptrLen = instance.exports.TableDataAtLength(this.addr, row, col);
    return readGoString(instance, ptr, ptrLen);
  }
  
  rowCount() {
    return instance.exports.TableDataRows(this.addr);
  }
  
  columnCount() {
    return instance.exports.TableDataColumns(this.addr);
  }
}

class Table {
  constructor() {
    this.addr = instance.exports.TableNew();
  }
  row(...rows) {
    let ptr = passStringArray(instance, rows);
    instance.exports.TableRow(this.addr, ptr, rows.length);
    return this;
  }
  rows(rows) {
    for (const rowIndex in rows) {
      const row = rows[rowIndex];
      let ptr = passStringArray(instance, row);
      instance.exports.TableRow(this.addr, ptr, row.length);
    }
    return this;
  }
  clearRows() {
    instance.exports.TableClearRows(this.addr);
    return this;
  }
  headers(...headers) {
    let ptr = passStringArray(instance, headers);
    instance.exports.TableHeaders(this.addr, ptr, headers.length);
    return this;
  }
  data(data) {
    if (data instanceof TableData) {
      instance.exports.TableSetData(this.addr, data.addr);
    }
    return this;
  }
  border(border) {
    instance.exports.TableBorder(this.addr, border);
    return this;
  }
  borderStyle(style) {
    if (style.addr) {
      instance.exports.TableBorderStyle(this.addr, style.addr);
    }
    return this;
  }
  styleFunc(fn) {
    let id = registerStyleFunction(fn);
    instance.exports.TableStyleFunc(this.addr, id);
    return this;
  }
  // TODO: These JS functions can be made by a compositor function
  wrap(b) {
    instance.exports.TableWrap(this.addr, b);
    return this;
  }
  borderLeft(b) {
    instance.exports.TableBorderLeft(this.addr, b);
    return this;
  }
  borderRight(b) {
    instance.exports.TableBorderRight(this.addr, b);
    return this;
  }
  borderHeader(b) {
    instance.exports.TableBorderHeader(this.addr, b);
    return this;
  }
  borderColumn(b) {
    instance.exports.TableBorderColumn(this.addr, b);
    return this;
  }
  borderRow(b) {
    instance.exports.TableBorderRow(this.addr, b);
    return this;
  }
  borderBottom(b) {
    instance.exports.TableBorderBottom(this.addr, b);
    return this;
  }
  borderLeft(b) {
    instance.exports.TableBorderLeft(this.addr, b);
    return this;
  }
  string() {
    return this.render();
  }
  render() {
    let ptr = instance.exports.TableRenderPtr(this.addr);
    let ptrLen = instance.exports.TableRenderLength(this.addr);
    return readGoString(instance, ptr, ptrLen);
  }
}

class List {
  constructor(...items) {
    this.addr = instance.exports.ListNew();
    if (items.length > 0) {
      this.items(...items);
    }
  }
  item(item) {
    if (typeof item === 'string') {
      let addr = insertString(item, instance);
      const encoder = new TextEncoder();
      const byteLength = encoder.encode(item).length;
      instance.exports.ListItem(this.addr, addr, byteLength);
      clearString(addr, byteLength, instance);
    } else if (item instanceof List) {
      instance.exports.ListItemList(this.addr, item.addr);
    }
    return this;
  }
  items(...items) {
    for (const item of items) {
      this.item(item);
    }
    return this;
  }
  hidden() {
    return instance.exports.ListHidden(this.addr);
  }
  hide(hide) {
    instance.exports.ListHide(this.addr, hide);
    return this;
  }
  offset(start, end) {
    instance.exports.ListOffset(this.addr, start, end);
    return this;
  }
  enumeratorStyle(style) {
    if (style.addr) {
      instance.exports.ListEnumeratorStyle(this.addr, style.addr);
    }
    return this;
  }
  enumeratorStyleFunc(fn) {
    let id = registerStyleFunction(fn);
    instance.exports.ListEnumeratorStyleFunc(this.addr, id);
    return this;
  }
  itemStyle(style) {
    if (style.addr) {
      instance.exports.ListItemStyle(this.addr, style.addr);
    }
    return this;
  }
  itemStyleFunc(fn) {
    let id = registerStyleFunction(fn);
    instance.exports.ListStyleFunc(this.addr, id);
    return this;
  }
  enumerator(type) {
    instance.exports.ListEnumerator(this.addr, type);
    return this;
  }
  string() {
    return this.render();
  }
  render() {
    let ptr = instance.exports.ListRenderPtr(this.addr);
    let ptrLen = instance.exports.ListRenderLength(this.addr);
    return readGoString(instance, ptr, ptrLen);
  }
}

class Tree {
  constructor() {
    this.addr = instance.exports.TreeNew();
  }
  root(value) {
    if (typeof value === 'string') {
      let addr = insertString(value, instance);
      const encoder = new TextEncoder();
      const byteLength = encoder.encode(value).length;
      instance.exports.TreeRoot(this.addr, addr, byteLength);
      clearString(addr, byteLength, instance);
    }
    return this;
  }
  child(...children) {
    for (const child of children) {
      if (typeof child === 'string') {
        let addr = insertString(child, instance);
        const encoder = new TextEncoder();
        const byteLength = encoder.encode(child).length;
        instance.exports.TreeChild(this.addr, addr, byteLength);
        clearString(addr, byteLength, instance);
      } else if (child instanceof Tree) {
        instance.exports.TreeChildTree(this.addr, child.addr);
      } else if (child instanceof Leaf) {
        instance.exports.TreeChildLeaf(this.addr, child.addr);
      }
    }
    return this;
  }
  hidden() {
    return instance.exports.TreeHidden(this.addr);
  }
  hide(hide) {
    instance.exports.TreeHide(this.addr, hide);
    return this;
  }
  offset(start, end) {
    instance.exports.TreeOffset(this.addr, start, end);
    return this;
  }
  enumeratorStyle(style) {
    if (style.addr) {
      instance.exports.TreeEnumeratorStyle(this.addr, style.addr);
    }
    return this;
  }
  enumeratorStyleFunc(fn) {
    let id = registerStyleFunction(fn);
    instance.exports.TreeEnumeratorStyleFunc(this.addr, id);
    return this;
  }
  itemStyle(style) {
    if (style.addr) {
      instance.exports.TreeItemStyle(this.addr, style.addr);
    }
    return this;
  }
  itemStyleFunc(fn) {
    let id = registerStyleFunction(fn);
    instance.exports.TreeStyleFunc(this.addr, id);
    return this;
  }
  rootStyle(style) {
    if (style.addr) {
      instance.exports.TreeRootStyle(this.addr, style.addr);
    }
    return this;
  }
  enumerator(type) {
    instance.exports.TreeEnumerator(this.addr, type);
    return this;
  }
  indenter(type) {
    instance.exports.TreeIndenter(this.addr, type);
    return this;
  }
  string() {
    return this.render();
  }
  render() {
    let ptr = instance.exports.TreeRenderPtr(this.addr);
    let ptrLen = instance.exports.TreeRenderLength(this.addr);
    return readGoString(instance, ptr, ptrLen);
  }
}

class Leaf {
  constructor(value, hidden = false) {
    if (typeof value === 'string') {
      let addr = insertString(value, instance);
      const encoder = new TextEncoder();
      const byteLength = encoder.encode(value).length;
      this.addr = instance.exports.TreeNewLeaf(addr, byteLength, hidden);
      clearString(addr, byteLength, instance);
    } else {
      // Convert non-string values to string
      const strValue = String(value);
      let addr = insertString(strValue, instance);
      const encoder = new TextEncoder();
      const byteLength = encoder.encode(strValue).length;
      this.addr = instance.exports.TreeNewLeaf(addr, byteLength, hidden);
      clearString(addr, byteLength, instance);
    }
  }
  value() {
    let ptr = instance.exports.TreeLeafValue(this.addr);
    let ptrLen = instance.exports.TreeLeafValueLength(this.addr);
    return readGoString(instance, ptr, ptrLen);
  }
  hidden() {
    return instance.exports.TreeLeafHidden(this.addr);
  }
  setHidden(hidden) {
    instance.exports.TreeLeafSetHidden(this.addr, hidden);
    return this;
  }
  setValue(value) {
    const strValue = String(value);
    let addr = insertString(strValue, instance);
    const encoder = new TextEncoder();
    const byteLength = encoder.encode(strValue).length;
    instance.exports.TreeLeafSetValue(this.addr, addr, byteLength);
    clearString(addr, byteLength, instance);
    return this;
  }
  string() {
    return this.value();
  }
}

function joinStyled(strings, bgColor, fgColor) {
  const memory = instance.exports.memory;
  let ptr = passStringArray(instance, strings);
  
  let resultGoStrPtr = instance.exports.StyleJoinStyled(ptr, strings.length, bgColor || 0, fgColor || 0);
  let result = "";

  if (resultGoStrPtr) {
    // In Go's WebAssembly ABI for strings:
    // - First 4 bytes at resultGoStrPtr: pointer to string data
    // - Next 4 bytes: length of string
    const resultPtrView = new DataView(memory.buffer, resultGoStrPtr, 8);
    const resultDataPtr = resultPtrView.getUint32(0, true); // true = little endian
    const resultLength = resultPtrView.getUint32(4, true);

    if (resultLength > 0 && resultDataPtr) {
      const resultBytes = new Uint8Array(
        memory.buffer,
        resultDataPtr,
        resultLength,
      );
      result = new TextDecoder().decode(resultBytes);
    }
  }

  return result;
}

function triggerGC() {
  try {
    if (instance && instance.exports.wasmGC) {
      instance.exports.wasmGC();
    }
  } catch (e) {
    // Ignore GC errors
  }
}

export {
  Place,
  Table,
  TableData,
  List,
  Tree,
  Leaf,
  Style,
  Color,
  NoColor,
  Size,
  Width,
  Height,
  Left,
  Right,
  Center,
  Top,
  Bottom,
  joinVertical,
  joinHorizontal,
  normalBorder,
  roundedBorder,
  blockBorder,
  outerHalfBlockBorder,
  innerHalfBlockBorder,
  thickBorder,
  doubleBorder,
  hiddenBorder,
  markdownBorder,
  ASCIIBorder,
  lightDark,
  joinStyled,
  triggerGC,
  Alphabet,
  Arabic,
  Bullet,
  Dash,
  Roman,
  Asterisk,
  DefaultEnumerator,
  RoundedEnumerator,
  DefaultIndenter,
};
