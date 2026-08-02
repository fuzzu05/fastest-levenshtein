# Default target
all: build

# Build the WebAssembly binary and copy the required runtime
build:
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" . || cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" .
	GOOS=js GOARCH=wasm go build -o fastest-levenshtein.wasm main.go levenshtein.go

# Run the original hashed tests against our new Wasm port
test: build
	npm install --ignore-scripts
	cp tests/original/test.ts .
	-npx tsc test.ts --esModuleInterop --skipLibCheck
	npx jest test.js
	rm test.ts test.js

# Run the performance benchmarks
bench: build
	npm install --ignore-scripts
	cp tests/original/bench.ts .
	-npx tsc bench.ts --esModuleInterop --skipLibCheck
	node bench.js
	rm bench.ts bench.js
