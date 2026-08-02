// Polyfill for globalThis.crypto required by Go's wasm_exec.js in Node/Jest environments
if (typeof globalThis.crypto === 'undefined') {
    globalThis.crypto = {
        getRandomValues: function (b) {
            require('crypto').randomFillSync(b);
        }
    };
}
if (typeof globalThis.performance === 'undefined') {
    globalThis.performance = {
        now: function () {
            return Date.now();
        }
    };
}
if (typeof globalThis.TextEncoder === 'undefined') {
    const util = require('util');
    globalThis.TextEncoder = util.TextEncoder;
    globalThis.TextDecoder = util.TextDecoder;
}

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
