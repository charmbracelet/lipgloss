#!/usr/bin/env node

const { execSync } = require('child_process');
const path = require('path');

const tests = [
  {
    name: 'Simple Table Test',
    file: 'tests/test-simple-table.js',
    description: 'Tests basic table functionality without styling'
  },
  {
    name: 'Style Function Test',
    file: 'tests/test-style-function.js',
    description: 'Tests table with style functions'
  },
  {
    name: 'Comprehensive TableData Test',
    file: 'tests/test-comprehensive.js',
    description: 'Tests all TableData functionality'
  },
  {
    name: 'TableData Example',
    file: 'table-data.js',
    description: 'Tests TableData examples and use cases'
  }
];

let passed = 0;
let failed = 0;

console.log('🧪 Running Lipgloss WASM Tests\n');

for (const test of tests) {
  try {
    console.log(`📋 ${test.name}`);
    console.log(`   ${test.description}`);
    
    const output = execSync(`node ${test.file}`, { 
      cwd: __dirname,
      encoding: 'utf8',
      timeout: 30000 // 30 second timeout
    });
    
    // Check if the output contains error indicators
    if (output.includes('Error:') || output.includes('RuntimeError:') || output.includes('panic:')) {
      throw new Error('Test output contains errors');
    }
    
    console.log(`   ✅ PASSED\n`);
    passed++;
  } catch (error) {
    console.log(`   ❌ FAILED`);
    console.log(`   Error: ${error.message}\n`);
    failed++;
  }
}

console.log('📊 Test Results:');
console.log(`   ✅ Passed: ${passed}`);
console.log(`   ❌ Failed: ${failed}`);
console.log(`   📈 Total:  ${passed + failed}`);

if (failed > 0) {
  console.log('\n❌ Some tests failed. Check the output above for details.');
  process.exit(1);
} else {
  console.log('\n🎉 All tests passed!');
  process.exit(0);
}