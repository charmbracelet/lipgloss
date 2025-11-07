export function insertString(text, module) {
  // Get the address of the writable memory.
  let addr = module.exports.getBuffer();
  let buffer = module.exports.memory.buffer;

  // Properly encode the string as UTF-8
  const encoder = new TextEncoder();
  const encodedBytes = encoder.encode(text);
  
  let mem = new Uint8Array(buffer);
  let view = mem.subarray(addr, addr + encodedBytes.length);

  // Copy the properly encoded bytes
  view.set(encodedBytes);

  // Return the address we started at.
  return addr;
}

/**
 * Clears a string from memory by zeroing out the buffer
 * @param {number} addr - The starting address of the string in memory
 * @param {number} length - The length of the string to clear
 * @param {WebAssembly.Module} module - The WebAssembly module with memory
 * @returns {boolean} - True if successfully cleared, false otherwise
 */
export function clearString(addr, length, module) {
  try {
    // Get the memory buffer
    let buffer = module.exports.memory.buffer;
    let mem = new Int8Array(buffer);

    // Zero out the memory region
    let view = mem.subarray(addr, addr + length);
    for (let i = 0; i < length; i++) {
      view[i] = 0;
    }

    return true;
  } catch (error) {
    console.error("Failed to clear string:", error);
    return false;
  }
}

/**
 * @param {WebAssembly.Instance} instance - The WebAssembly instance
 * @param {number} ptr - Pointer to the string in memory
 * @param {number} len - Length of the string
 * @returns {string} - The string read from memory
 */
export function readGoString(instance, ptr, len) {
  try {
    // Get the memory buffer
    const memory = instance.exports.memory.buffer;

    // Check if the pointer and length are within memory bounds
    if (ptr + len > memory.byteLength) {
      console.error("String access would be out of bounds:", {
        ptr,
        len,
        memorySize: memory.byteLength,
      });
      return "";
    }

    // Read the string from memory
    const bytes = new Uint8Array(memory, ptr, len);
    return new TextDecoder().decode(bytes);
  } catch (error) {
    console.error("Error in readGoString:", error);
    return "";
  }
}

// Memory management for string arrays
let memoryOffset = 4096; // Start after some safe space
const MEMORY_CHUNK_SIZE = 8192; // 8KB chunks

/**
 * Grows WASM memory if needed to accommodate the required size
 * @param {WebAssembly.Instance} instance - The WASM instance
 * @param {number} requiredBytes - Number of bytes needed
 * @returns {boolean} - True if memory is sufficient, false if growth failed
 */
export function ensureMemorySize(instance, requiredBytes) {
  const memory = instance.exports.memory;
  const currentSize = memory.buffer.byteLength;
  
  if (currentSize >= requiredBytes) {
    return true; // Already have enough memory
  }
  
  // Calculate how many pages we need to grow
  const pageSize = 65536; // WASM page size is 64KB
  const currentPages = currentSize / pageSize;
  const requiredPages = Math.ceil(requiredBytes / pageSize);
  const pagesToGrow = requiredPages - currentPages;
  
  // Check environment variables for debug output
  const debugMemory = process.env.LIPGLOSS_DEBUG_MEMORY === 'true' || 
                     process.env.LIPGLOSS_DEBUG === 'true' ||
                     process.env.DEBUG === 'lipgloss' ||
                     process.env.DEBUG === '*';
  
  if (debugMemory) {
    console.log(`Growing WASM memory: ${currentSize} -> ${requiredPages * pageSize} bytes (${pagesToGrow} pages)`);
  }
  
  try {
    const previousPages = memory.grow(pagesToGrow);
    if (debugMemory) {
      console.log(`Memory growth successful: ${previousPages} -> ${previousPages + pagesToGrow} pages`);
    }
    return true;
  } catch (error) {
    console.error(`Failed to grow WASM memory by ${pagesToGrow} pages:`, error);
    return false;
  }
}

export function passStringArray(instance, strings) {
  const memory = instance.exports.memory;
  const wasmMalloc = instance.exports.wasmMalloc;
  const wasmFree = instance.exports.wasmFree;

  // Always use manual memory management to avoid TinyGo malloc issues
  return passStringArrayWithBoundsCheck(instance, strings);
}

function passStringArrayWithTinyGoOptimization(instance, strings, wasmMalloc, wasmFree) {
  const memory = instance.exports.memory;
  const wasmGC = instance.exports.wasmGC;

  // Calculate total size needed
  const encoder = new TextEncoder();
  let totalStringSize = 0;
  const encodedStrings = [];
  
  for (const str of strings) {
    // Don't truncate strings - let memory grow instead
    const encoded = encoder.encode(str);
    encodedStrings.push(encoded);
    totalStringSize += encoded.length;
  }

  // Calculate total memory needed (pointer array + string data + safety buffer)
  const pointerArraySize = strings.length * 8;
  const totalMemoryNeeded = pointerArraySize + totalStringSize + 4096; // 4KB safety buffer
  
  // Try to ensure we have enough memory
  if (!ensureMemorySize(instance, totalMemoryNeeded)) {
    console.warn(`Failed to grow memory for ${totalMemoryNeeded} bytes, using fallback`);
    return passStringArrayWithBoundsCheck(instance, strings);
  }

  // Allocate memory for the pointer/length array (8 bytes per string: 4 for ptr, 4 for length)
  const pointerArrayOffset = wasmMalloc(pointerArraySize);
  
  if (!pointerArrayOffset) {
    console.error("TinyGo malloc failed for pointer array");
    return passStringArrayWithBoundsCheck(instance, strings);
  }

  const allocatedPtrs = [];

  try {
    // Create a view for writing the pointer/length pairs
    const pointerArray = new Uint32Array(
      memory.buffer,
      pointerArrayOffset,
      strings.length * 2,
    );

    // Encode and store each string with careful memory management
    for (let i = 0; i < encodedStrings.length; i++) {
      const stringBytes = encodedStrings[i];
      
      // Allocate memory for this string
      const stringDataOffset = wasmMalloc(stringBytes.length);
      if (!stringDataOffset) {
        console.warn(`TinyGo malloc failed for string ${i}, using fallback`);
        // Clean up and fall back
        for (const ptr of allocatedPtrs) {
          if (ptr) wasmFree(ptr);
        }
        wasmFree(pointerArrayOffset);
        return passStringArrayWithBoundsCheck(instance, strings);
      }
      allocatedPtrs.push(stringDataOffset);

      // Copy string data to memory with bounds checking
      try {
        const memoryView = new Uint8Array(memory.buffer);
        if (stringDataOffset + stringBytes.length > memory.buffer.byteLength) {
          throw new Error("Memory bounds exceeded");
        }
        memoryView.set(stringBytes, stringDataOffset);
      } catch (e) {
        console.warn(`Memory copy failed for string ${i}:`, e.message);
        // Clean up and fall back
        for (const ptr of allocatedPtrs) {
          if (ptr) wasmFree(ptr);
        }
        wasmFree(pointerArrayOffset);
        return passStringArrayWithBoundsCheck(instance, strings);
      }

      // Store pointer and length in the array
      pointerArray[i * 2] = stringDataOffset; // String pointer
      pointerArray[i * 2 + 1] = stringBytes.length; // String length
    }

    // Trigger periodic cleanup for TinyGo
    if (wasmGC) {
      wasmGC();
    }

    return pointerArrayOffset;
  } catch (error) {
    // Clean up on error and fall back
    for (const ptr of allocatedPtrs) {
      if (ptr) wasmFree(ptr);
    }
    wasmFree(pointerArrayOffset);
    console.warn("TinyGo memory allocation failed, using fallback:", error.message);
    return passStringArrayWithBoundsCheck(instance, strings);
  }
}

function passStringArrayWithMalloc(instance, strings, malloc, free) {
  const memory = instance.exports.memory;

  // Allocate memory for the pointer/length array (8 bytes per string: 4 for ptr, 4 for length)
  const pointerArraySize = strings.length * 8;
  const pointerArrayOffset = malloc(pointerArraySize);

  if (!pointerArrayOffset) {
    console.error("Failed to allocate memory for pointer array");
    return 0;
  }

  // Create a view for writing the pointer/length pairs
  const pointerArray = new Uint32Array(
    memory.buffer,
    pointerArrayOffset,
    strings.length * 2,
  );

  const allocatedPtrs = [];

  try {
    // Encode and store each string
    for (let i = 0; i < strings.length; i++) {
      const encoder = new TextEncoder();
      const stringBytes = encoder.encode(strings[i]);

      // Allocate memory for this string
      const stringDataOffset = malloc(stringBytes.length);
      if (!stringDataOffset) {
        throw new Error(`Failed to allocate memory for string ${i}`);
      }
      allocatedPtrs.push(stringDataOffset);

      // Copy string data to memory
      const memoryView = new Uint8Array(memory.buffer);
      memoryView.set(stringBytes, stringDataOffset);

      // Store pointer and length in the array
      pointerArray[i * 2] = stringDataOffset; // String pointer
      pointerArray[i * 2 + 1] = stringBytes.length; // String length
    }

    return pointerArrayOffset;
  } catch (error) {
    // Clean up on error
    for (const ptr of allocatedPtrs) {
      if (ptr) free(ptr);
    }
    free(pointerArrayOffset);
    console.error("Error in passStringArray:", error);
    return 0;
  }
}

function passStringArrayWithWasmMalloc(instance, strings, wasmMalloc, wasmFree) {
  const memory = instance.exports.memory;
  const wasmGC = instance.exports.wasmGC;

  // Calculate total size needed
  const encoder = new TextEncoder();
  let totalStringSize = 0;
  const encodedStrings = [];
  
  for (const str of strings) {
    const encoded = encoder.encode(str);
    encodedStrings.push(encoded);
    totalStringSize += encoded.length;
  }

  // For TinyGo, use chunked processing for large strings
  const maxChunkSize = 64 * 1024; // 64KB chunks
  if (totalStringSize > maxChunkSize) {
    return passStringArrayChunked(instance, encodedStrings, wasmMalloc, wasmFree);
  }

  // Allocate memory for the pointer/length array (8 bytes per string: 4 for ptr, 4 for length)
  const pointerArraySize = strings.length * 8;
  const pointerArrayOffset = wasmMalloc(pointerArraySize);
  
  if (!pointerArrayOffset) {
    console.error("Failed to allocate memory for pointer array");
    return 0;
  }

  const allocatedPtrs = [];

  try {
    // Create a view for writing the pointer/length pairs
    const pointerArray = new Uint32Array(
      memory.buffer,
      pointerArrayOffset,
      strings.length * 2,
    );

    // Encode and store each string
    for (let i = 0; i < encodedStrings.length; i++) {
      const stringBytes = encodedStrings[i];
      
      // Allocate memory for this string
      const stringDataOffset = wasmMalloc(stringBytes.length);
      if (!stringDataOffset) {
        throw new Error(`Failed to allocate memory for string ${i}`);
      }
      allocatedPtrs.push(stringDataOffset);

      // Copy string data to memory
      const memoryView = new Uint8Array(memory.buffer);
      memoryView.set(stringBytes, stringDataOffset);

      // Store pointer and length in the array
      pointerArray[i * 2] = stringDataOffset; // String pointer
      pointerArray[i * 2 + 1] = stringBytes.length; // String length
    }

    // Trigger periodic cleanup for TinyGo
    if (wasmGC) {
      wasmGC();
    }

    return pointerArrayOffset;
  } catch (error) {
    // Clean up on error
    for (const ptr of allocatedPtrs) {
      if (ptr) wasmFree(ptr);
    }
    wasmFree(pointerArrayOffset);
    console.error("Error in passStringArrayWithWasmMalloc:", error);
    return 0;
  }
}

// Chunked processing for large string arrays
function passStringArrayChunked(instance, encodedStrings, wasmMalloc, wasmFree) {
  const memory = instance.exports.memory;
  
  // Process strings in smaller chunks to avoid TinyGo string length limits
  const chunkSize = 8; // Process 8 strings at a time
  const chunks = [];
  
  for (let i = 0; i < encodedStrings.length; i += chunkSize) {
    chunks.push(encodedStrings.slice(i, i + chunkSize));
  }
  
  // Allocate memory for the main pointer array
  const pointerArraySize = encodedStrings.length * 8;
  const pointerArrayOffset = wasmMalloc(pointerArraySize);
  
  if (!pointerArrayOffset) {
    console.error("Failed to allocate memory for chunked pointer array");
    return 0;
  }
  
  try {
    const pointerArray = new Uint32Array(
      memory.buffer,
      pointerArrayOffset,
      encodedStrings.length * 2,
    );
    
    let stringIndex = 0;
    
    // Process each chunk
    for (const chunk of chunks) {
      for (const stringBytes of chunk) {
        const stringDataOffset = wasmMalloc(stringBytes.length);
        if (!stringDataOffset) {
          throw new Error(`Failed to allocate memory for chunked string ${stringIndex}`);
        }
        
        // Copy string data to memory
        const memoryView = new Uint8Array(memory.buffer);
        memoryView.set(stringBytes, stringDataOffset);
        
        // Store pointer and length in the array
        pointerArray[stringIndex * 2] = stringDataOffset;
        pointerArray[stringIndex * 2 + 1] = stringBytes.length;
        
        stringIndex++;
      }
      
      // Small delay between chunks to let TinyGo process
      // This is a no-op but helps with memory management
      if (chunks.length > 1) {
        const dummy = new Uint8Array(1);
      }
    }
    
    return pointerArrayOffset;
  } catch (error) {
    wasmFree(pointerArrayOffset);
    console.error("Error in passStringArrayChunked:", error);
    return 0;
  }
}

function passStringArrayWithBoundsCheck(instance, strings) {
  const memory = instance.exports.memory;

  // Calculate total size needed
  const encoder = new TextEncoder();
  let totalStringSize = 0;
  const encodedStrings = [];

  for (const str of strings) {
    const encoded = encoder.encode(str);
    encodedStrings.push(encoded);
    totalStringSize += encoded.length;
  }

  // Allocate memory for the pointer/length array (8 bytes per string: 4 for ptr, 4 for length)
  const pointerArraySize = strings.length * 8;
  const totalSize = pointerArraySize + totalStringSize;
  const totalMemoryNeeded = totalSize + 4096; // 4KB safety buffer

  // Try to ensure we have enough memory
  if (!ensureMemorySize(instance, totalMemoryNeeded)) {
    console.error(`Failed to grow memory for ${totalMemoryNeeded} bytes in fallback`);
    return 0;
  }

  const memorySize = memory.buffer.byteLength; // Get updated size after potential growth

  // Use a much larger safe memory area for large ANSI strings
  const safeMemoryStart = 1024;
  const safeMemoryEnd = Math.min(memorySize - totalSize - 2048, memorySize * 0.8); // Use up to 80% of memory, leave 2KB buffer
  
  if (safeMemoryEnd <= safeMemoryStart || totalSize > (memorySize * 0.8)) {
    console.error(`Not enough memory: need ${totalSize}, available ${safeMemoryEnd - safeMemoryStart}`);
    // Try to reset memory offset and retry once
    memoryOffset = safeMemoryStart;
    if (totalSize > (memorySize * 0.8)) {
      console.error("String data too large for available memory");
      return 0;
    }
  }

  // Use a rotating buffer within safe bounds, but with larger chunks for ANSI strings
  const chunkSize = Math.max(totalSize, 16384); // At least 16KB chunks
  if (memoryOffset < safeMemoryStart || memoryOffset + totalSize > safeMemoryEnd) {
    memoryOffset = safeMemoryStart;
  }

  const pointerArrayOffset = memoryOffset;

  // Check bounds before proceeding
  if (pointerArrayOffset + totalSize > memorySize) {
    console.error(
      `Memory bounds check failed: ${pointerArrayOffset + totalSize} > ${memorySize}`,
    );
    return 0;
  }

  try {
    // Allocate memory for string data (after the pointer array)
    let stringDataOffset = pointerArrayOffset + pointerArraySize;

    // Create a view for writing the pointer/length pairs
    const pointerArray = new Uint32Array(
      memory.buffer,
      pointerArrayOffset,
      strings.length * 2,
    );

    const memoryView = new Uint8Array(memory.buffer);

    // Store each string
    for (let i = 0; i < encodedStrings.length; i++) {
      const stringBytes = encodedStrings[i];

      // Double-check bounds for each string
      if (stringDataOffset + stringBytes.length > memorySize) {
        throw new Error(`String ${i} would exceed memory bounds`);
      }

      // Copy string data to memory
      memoryView.set(stringBytes, stringDataOffset);

      // Store pointer and length in the array
      pointerArray[i * 2] = stringDataOffset; // String pointer
      pointerArray[i * 2 + 1] = stringBytes.length; // String length

      // Update offset for next string
      stringDataOffset += stringBytes.length;
    }

    // Update global offset for next call with larger increment to avoid conflicts
    memoryOffset = Math.max(stringDataOffset, memoryOffset + chunkSize);
    
    return pointerArrayOffset;
  } catch (error) {
    console.error("Error in passStringArrayWithBoundsCheck:", error);
    // Reset memory offset on error
    memoryOffset = safeMemoryStart;
    return 0;
  }
}
