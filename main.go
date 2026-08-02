package main

import (
	"syscall/js"
)

// jsDistance wraps the Distance function for JavaScript.
func jsDistance(this js.Value, args []js.Value) any {
	a := args[0].String()
	b := args[1].String()
	return Distance(a, b)
}

// jsClosest wraps the Closest function for JavaScript.
func jsClosest(this js.Value, args []js.Value) any {
	str := args[0].String()

	// Convert JavaScript array to Go slice of strings
	jsArr := args[1]
	length := jsArr.Length()
	arr := make([]string, length)
	for i := 0; i < length; i++ {
		arr[i] = jsArr.Index(i).String()
	}

	return Closest(str, arr)
}

func main() {
	// Expose the functions to the global JS environment.
	js.Global().Set("wasmDistance", js.FuncOf(jsDistance))
	js.Global().Set("wasmClosest", js.FuncOf(jsClosest))

	// Keep the Wasm instance running indefinitely.
	select {}
}
