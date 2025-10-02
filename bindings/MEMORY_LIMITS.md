# WASM Memory Limits and Configuration

## Current Status

The 15MB you saw in the layout example is **not a fixed limit** - it was just what was needed for that particular rendering operation. The memory can grow much larger.

## Memory Limits Hierarchy

### 1. WASM Specification Limits
- **Maximum**: 4GB (4,294,967,296 bytes)
- **Page size**: 64KB (65,536 bytes)
- **Maximum pages**: 65,536 pages
- **This is a hard limit** defined by the WebAssembly specification

### 2. JavaScript Engine Limits
- **ArrayBuffer maximum**: ~2GB on most 32-bit systems, up to available memory on 64-bit
- **Node.js default heap**: ~1.4GB (can be increased with `--max-old-space-size`)
- **Browser limits**: Vary by browser and available system memory

### 3. System Memory Limits
- **Available RAM**: Physical memory minus what's used by other processes
- **Virtual memory**: Can extend beyond physical RAM using swap

## How to Increase Limits

### For Node.js Applications
```bash
# Increase Node.js heap to 8GB
node --max-old-space-size=8192 your-app.js

# Increase to 16GB
node --max-old-space-size=16384 your-app.js
```

### For Browser Applications
- No direct configuration - limited by browser and system memory
- Modern browsers on 64-bit systems can typically handle several GB

## Practical Limits for Lipgloss

### Current Implementation
- ✅ **Up to 4GB**: Fully supported by our dynamic memory growth
- ✅ **Automatic growth**: Memory grows as needed for complex layouts
- ✅ **Efficient allocation**: Uses manual memory management to avoid TinyGo malloc issues

### Real-World Usage
- **Simple layouts**: < 1MB
- **Complex layouts with colors**: 10-50MB  
- **Very large documents**: 100MB-1GB
- **Extreme cases**: Up to 4GB (WASM limit)

## Configuration Examples

### High Memory Node.js App
```bash
# For very large layout rendering
node --max-old-space-size=8192 examples/layout.js
```

### Memory-Efficient Usage
```javascript
// Process large layouts in chunks to stay under memory limits
const chunks = splitLargeLayout(content);
const results = chunks.map(chunk => renderChunk(chunk));
const final = joinResults(results);
```

## Monitoring Memory Usage

The current implementation logs memory growth:
```
Growing WASM memory: 393216 -> 458752 bytes (1 pages)
Memory growth successful: 6 -> 7 pages
```

You can monitor this to understand your application's memory needs.

## Summary

- **No fixed 15MB limit** - memory grows dynamically as needed
- **Maximum possible**: 4GB (WASM specification limit)
- **Practical limit**: Usually constrained by available system memory
- **Node.js**: Can be increased with `--max-old-space-size` flag
- **Browser**: Limited by browser and system capabilities

The dynamic memory growth ensures your layouts will work regardless of complexity, up to the 4GB WASM limit.