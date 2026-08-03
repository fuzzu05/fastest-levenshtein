package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// originalJSDistance uses Node to run the original npm library to get the absolute ground-truth JS answer
func originalJSDistance(a, b string) int {
	script := `console.log(require("fastest-levenshtein").distance(process.argv[1], process.argv[2]))`

	// We use npx to temporarily download and run the original JS implementation
	cmd := exec.Command("npx", "-p", "fastest-levenshtein", "node", "-e", script, a, b)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	res := strings.TrimSpace(out.String())
	val, _ := strconv.Atoi(res)
	return val
}

func FuzzDistance(f *testing.F) {
	// Seed the fuzzer with some interesting edge cases
	f.Add("", "")
	f.Add("kitten", "sitting")
	f.Add("😎", "😭") // Surrogate pairs (UTF-16 boundary test)

	// The fuzzer automatically generates random strings to test
	f.Fuzz(func(t *testing.T, a, b string) {
		expected := originalJSDistance(a, b)
		got := Distance(a, b)

		if expected != got {
			t.Errorf("Mismatch for a=%q, b=%q\nExpected (JS): %d\nGot (Go): %d", a, b, expected, got)
		}
	})
}
