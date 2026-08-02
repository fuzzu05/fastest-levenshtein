package main

import (
	"testing"
)

// referenceDistance is a simple, clearly-correct Dynamic Programming
// implementation of Levenshtein distance. It explicitly calculates distance
// on UTF-16 code units to match the original JavaScript semantics.
func referenceDistance(a, b string) int {
	aUnits := utf16Units(a)
	bUnits := utf16Units(b)

	if len(aUnits) == 0 {
		return len(bUnits)
	}
	if len(bUnits) == 0 {
		return len(aUnits)
	}

	row := make([]int, len(aUnits)+1)
	for i := 0; i <= len(aUnits); i++ {
		row[i] = i
	}

	for j := 1; j <= len(bUnits); j++ {
		prev := j
		for i := 1; i <= len(aUnits); i++ {
			var val int
			if aUnits[i-1] == bUnits[j-1] {
				val = row[i-1]
			} else {
				// Min of (insertion, deletion, substitution)
				min := row[i-1] + 1
				if prev+1 < min {
					min = prev + 1
				}
				if row[i]+1 < min {
					min = row[i] + 1
				}
				val = min
			}
			row[i-1] = prev
			prev = val
		}
		row[len(aUnits)] = prev
	}
	return row[len(aUnits)]
}

func FuzzDistance(f *testing.F) {
	// Seed the fuzzer with some interesting edge cases and boundaries
	f.Add("", "")
	f.Add("a", "")
	f.Add("", "a")
	f.Add("abc", "abc")
	f.Add("kitten", "sitting")
	f.Add("😎", "😭") // Surrogate pairs (UTF-16 boundary test)
	f.Add("a😎b", "ab😎")

	// Add strings exactly length 32 to stress-test the myers32 -> myersX transition
	f.Add("12345678901234567890123456789012", "12345678901234567890123456789012")
	f.Add("123456789012345678901234567890123", "1234567890")

	// The fuzzer automatically generates thousands of random strings to test
	f.Fuzz(func(t *testing.T, a, b string) {
		expected := referenceDistance(a, b)
		got := Distance(a, b)

		if expected != got {
			t.Errorf("Mismatch for a=%q, b=%q\nExpected: %d\nGot: %d", a, b, expected, got)
		}
	})
}
