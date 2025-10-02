const { ensureMemorySize } = require('./src/core.js');

// Create a simple WASM module with memory export
const wasmCode = new Uint8Array([
  0x00, 0x61, 0x73, 0x6d, // WASM magic number
  0x01, 0x00, 0x00, 0x00, // version
  0x05, 0x03, 0x01, 0x00, 0x01, // memory section: initial=1 page, no maximum
  0x07, 0x0a, 0x01, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x02, 0x00, // export memory
]);

try {
  const wasmModule = new WebAssembly.Module(wasmCode);
  const instance = new WebAssembly.Instance(wasmModule);

  console.log('Initial memory:', instance.exports.memory.buffer.byteLength / 1024 / 1024, 'MB');

  // Test growing to different sizes
  const testSizes = [
    50 * 1024 * 1024,   // 50MB
    100 * 1024 * 1024,  // 100MB
    500 * 1024 * 1024,  // 500MB
    1024 * 1024 * 1024, // 1GB
  ];

  for (const size of testSizes) {
    try {
      const success = ensureMemorySize(instance, size);
      const actualSize = instance.exports.memory.buffer.byteLength;
      console.log(`Requested: ${size / 1024 / 1024}MB, Success: ${success}, Actual: ${actualSize / 1024 / 1024}MB`);
      if (!success) break;
    } catch (error) {
      console.log(`Failed at ${size / 1024 / 1024}MB:`, error.message);
      break;
    }
  }
} catch (error) {
  console.log('WASM module creation failed:', error.message);
  console.log('Testing with a simpler approach...');
  
  // Test with the actual lipgloss WASM module
  const fs = require('fs');
  const wasmBytes = fs.readFileSync('./src/lipgloss.wasm');
  const wasmModule = new WebAssembly.Module(wasmBytes);
  
  // Create minimal imports for Go runtime
  const imports = {
    go: {
      'runtime.wasmExit': () => {},
      'runtime.wasmWrite': () => {},
      'runtime.nanotime1': () => BigInt(Date.now() * 1000000),
      'runtime.walltime': () => {
        const now = Date.now();
        return [Math.floor(now / 1000), (now % 1000) * 1000000];
      },
      'runtime.scheduleTimeoutEvent': () => {},
      'runtime.clearTimeoutEvent': () => {},
      'runtime.getRandomData': () => {},
      'syscall/js.finalizeRef': () => {},
      'syscall/js.stringVal': () => {},
      'syscall/js.valueGet': () => {},
      'syscall/js.valueSet': () => {},
      'syscall/js.valueDelete': () => {},
      'syscall/js.valueIndex': () => {},
      'syscall/js.valueSetIndex': () => {},
      'syscall/js.valueCall': () => {},
      'syscall/js.valueInvoke': () => {},
      'syscall/js.valueNew': () => {},
      'syscall/js.valueLength': () => {},
      'syscall/js.valuePrepareString': () => {},
      'syscall/js.valueLoadString': () => {},
      'syscall/js.valueInstanceOf': () => {},
      'syscall/js.copyBytesToGo': () => {},
      'syscall/js.copyBytesToJS': () => {},
    },
    env: {}
  };
  
  try {
    const instance = new WebAssembly.Instance(wasmModule, imports);
    console.log('Lipgloss WASM initial memory:', instance.exports.memory.buffer.byteLength / 1024 / 1024, 'MB');
    
    // Test memory growth
    const testSize = 100 * 1024 * 1024; // 100MB
    const success = ensureMemorySize(instance, testSize);
    const actualSize = instance.exports.memory.buffer.byteLength;
    console.log(`Requested: ${testSize / 1024 / 1024}MB, Success: ${success}, Actual: ${actualSize / 1024 / 1024}MB`);
  } catch (error) {
    console.log('Lipgloss WASM test failed:', error.message);
  }
}

console.log('\nNode.js memory info:');
const memUsage = process.memoryUsage();
Object.entries(memUsage).forEach(([key, value]) => {
  console.log(`${key}: ${Math.round(value / 1024 / 1024)}MB`);
});