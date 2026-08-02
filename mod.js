const fs = require('fs');
const path = require('path');

// Load the Go WebAssembly runtime environment
require('./wasm_exec.js');

const go = new Go();
const wasmBuffer = fs.readFileSync(path.join(__dirname, 'fastest-levenshtein.wasm'));
const wasmModule = new WebAssembly.Module(wasmBuffer);
const wasmInstance = new WebAssembly.Instance(wasmModule, go.importObject);

// Run the Go program (this registers the global functions)
go.run(wasmInstance);

// Export the WASM functions exactly like the original TypeScript library did
module.exports = {
    distance: global.wasmDistance,
    closest: global.wasmClosest
};
