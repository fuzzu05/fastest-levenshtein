
# 🧠 Architectural Decision Log: fastest-levenshtein (Go Port)

This document captures the key design choices behind the Go port and the reasoning behind each one.

## 1. Zero Unsafe

**Decision:** Implement the Myers bit-parallel algorithm in pure Go without using the `unsafe` package.

**Why:** Go’s built-in bounds checks and type safety are already fast enough for this use case. Using only safe operations preserves memory safety while keeping the implementation straightforward and reliable.

## 2. UTF-16 Semantic Equivalence (Quirk Preservation)

**Decision:** Convert Go UTF-8 strings into UTF-16 code units with `utf16.Encode` before processing.

**Why:** The original JavaScript implementation calculates distance using UTF-16 code units rather than true Unicode code points. To preserve behavioral equivalence and satisfy the hashed test suite, the Go port intentionally mirrors that behavior instead of silently changing it.

## 3. WebAssembly FFI Adapter for Test Parity

**Decision:** Expose the implementation to Node.js through WebAssembly using `syscall/js`, rather than CGO or native binaries.

**Why:** WebAssembly allows the original JavaScript test suite to run with minimal friction across platforms. The accompanying `mod.js` shim preserves the original module interface, making the port behave like a drop-in replacement.

## 4. sync.Pool for the Peq Array

**Decision:** Use a `sync.Pool` for the 65,536-entry bitmask array used by the algorithm.

**Why:** The original JavaScript version relied on a global array to avoid repeated allocations during the inner loop. In Go, a global array would introduce concurrency issues. A `sync.Pool` provides the same performance benefit while remaining safe under concurrent use.

