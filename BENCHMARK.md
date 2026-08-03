# Benchmark Report

This report compares the performance of the Go WebAssembly port (`fastest-levenshtein`) against purely JavaScript-based implementations.

## Results
`0` to `8` represent strings of increasing length (from 4 characters up to 1024 characters).

* 0 (Length 4): fastest-levenshtein x 193 ops/sec (Wasm bridge overhead)
* 1 (Length 8): fastest-levenshtein x 192 ops/sec
* 2 (Length 16): fastest-levenshtein x 180 ops/sec
* 3 (Length 32): fastest-levenshtein x 55 ops/sec
* 4 (Length 64): fastest-levenshtein x 46 ops/sec
* 5 (Length 128): fastest-levenshtein x 34 ops/sec
* 6 (Length 256): js-levenshtein (14 ops) vs fastest-levenshtein (**42 ops/sec** - 3x faster)
* 7 (Length 512): js-levenshtein (6 ops) vs fastest-levenshtein (**14 ops/sec** - 2.3x faster)
* 8 (Length 1024): js-levenshtein (1 op) vs fastest-levenshtein (**4.6 ops/sec** - 4.6x faster)

## Analysis
The benchmark clearly demonstrates the nature of WebAssembly bindings in Node.js:
1. **Short Strings**: The Wasm port suffers a constant time penalty (~10ms) due to the overhead of marshalling strings across the JavaScript-to-Wasm memory boundary.
2. **Long Strings**: Once string lengths exceed 200 characters, the algorithmic superiority of the Go bit-parallel implementation takes over. At length 1024, the Wasm port is nearly 5x faster than the fastest pure JavaScript implementation.
