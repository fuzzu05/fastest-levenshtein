# fastest-levenshtein (Go Port) ⚡

A pure Go port of the [fastest-levenshtein](https://github.com/ka-weihe/fastest-levenshtein) library, originally written in TypeScript.

Built for the **🔥 PORT MORTEM 2026** Hackathon, this project aims to preserve the original behavior while bringing the implementation into the Go ecosystem.

## ✨ Highlights

- **Behavioral Equivalence**: Mirrors the original library’s behavior and edge cases, including the hashed test suite.
- **Memory Safe**: Uses no `unsafe` imports, keeping the implementation straightforward and safe.
- **Wasm-Friendly**: Designed to work cleanly with Node.js through WebAssembly bindings.

## 🐛 Bug Catcher

While porting and reviewing the original TypeScript implementation, four issues were identified and addressed:

1. **Empty-array handling in `closest`**: Calling `closest("a", [])` previously returned `undefined` at runtime, which violated the expected string return contract.
2. **Reentrancy bug from global state**: The shared `peq` state could be corrupted when a custom `charCodeAt` getter triggered nested calls.
3. **Unicode distance inflation**: The original logic measured UTF-16 code units rather than true Unicode code points, which could overcount characters like emoji.
4. **Unnecessary `Math.min` work**: The `myers_x` loop performed a redundant minimum calculation that was mathematically constant.

## 🧪 Building and Testing

This project follows a simple workflow:

```bash
make test
```

Run the original unmodified test suite and build the Wasm binary.

```bash
make bench
```

Run the benchmark suite.

```