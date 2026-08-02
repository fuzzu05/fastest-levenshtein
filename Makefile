# Default target
all: build

# Build the WebAssembly binary and copy the required runtime
build:
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" .
	GOOS=js GOARCH=wasm go build -o fastest-levenshtein.wasm main.go levenshtein.go

# Run the original hashed tests against our new Wasm port
test: build
	npm install
	cp tests/original/test.ts .
	npx tsc test.ts --esModuleInterop
	npx jest test.js
	rm test.ts test.js

# Run the performance benchmarks
bench: build
	npm install
	cp tests/original/bench.ts .
	npx tsc bench.ts --esModuleInterop
	node bench.js
	rm bench.ts bench.js
