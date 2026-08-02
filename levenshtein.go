package main

import (
	"sync"
	"unicode/utf16"
)

const alphabetSize = 0x10000

// Reusing the 65,536-entry table avoids allocating roughly 256 KB
// every time Distance is called.
var peqPool = sync.Pool{
	New: func() any {
		return make([]uint32, alphabetSize)
	},
}

// utf16Units converts a Go UTF-8 string into JavaScript-style UTF-16
// code units, equivalent to indexing a JS string with charCodeAt().
func utf16Units(s string) []uint16 {
	return utf16.Encode([]rune(s))
}

// JavaScript masks bit-shift counts to 0-31.
// Go does not, so this helper reproduces JavaScript's behaviour.
func bit(position int) uint32 {
	return uint32(1) << uint(position&31)
}

func myers32(a, b []uint16, peq []uint32) int {
	n := len(a)
	m := len(b)

	last := bit(n - 1)

	pv := ^uint32(0)
	var mv uint32

	score := n

	for i := n - 1; i >= 0; i-- {
		peq[a[i]] |= bit(i)
	}

	for i := 0; i < m; i++ {
		eq := peq[b[i]]
		xv := eq | mv

		eq |= ((eq & pv) + pv) ^ pv
		mv |= ^(eq | pv)
		pv &= eq

		if mv&last != 0 {
			score++
		}

		if pv&last != 0 {
			score--
		}

		mv = (mv << 1) | 1
		pv = (pv << 1) | ^(xv | mv)
		mv &= xv
	}

	// Reset only the entries that were modified.
	for i := n - 1; i >= 0; i-- {
		peq[a[i]] = 0
	}

	return score
}

func myersX(longer, shorter []uint16, peq []uint32) int {
	n := len(shorter)
	m := len(longer)

	hsize := (n + 31) / 32
	vsize := (m + 31) / 32

	phc := make([]uint32, hsize)
	mhc := make([]uint32, hsize)

	for i := range phc {
		phc[i] = ^uint32(0)
	}

	j := 0

	// Process all complete 32-code-unit blocks except the final block.
	for ; j < vsize-1; j++ {
		var mv uint32
		pv := ^uint32(0)

		start := j * 32
		end := start + 32

		for k := start; k < end; k++ {
			peq[longer[k]] |= bit(k)
		}

		for i := 0; i < n; i++ {
			eq := peq[shorter[i]]

			block := i >> 5
			shift := uint(i & 31)

			pb := (phc[block] >> shift) & 1
			mb := (mhc[block] >> shift) & 1

			xv := eq | mv
			xh := ((((eq | mb) & pv) + pv) ^ pv) | eq | mb

			ph := mv | ^(xh | pv)
			mh := pv & xh

			if ((ph >> 31) ^ pb) != 0 {
				phc[block] ^= uint32(1) << shift
			}

			if ((mh >> 31) ^ mb) != 0 {
				mhc[block] ^= uint32(1) << shift
			}

			ph = (ph << 1) | pb
			mh = (mh << 1) | mb

			pv = mh | ^(xv | ph)
			mv = ph & xv
		}

		for k := start; k < end; k++ {
			peq[longer[k]] = 0
		}
	}

	// Process the final block.
	var mv uint32
	pv := ^uint32(0)

	start := j * 32
	end := m

	for k := start; k < end; k++ {
		peq[longer[k]] |= bit(k)
	}

	score := m
	lastShift := uint((m - 1) & 31)

	for i := 0; i < n; i++ {
		eq := peq[shorter[i]]

		block := i >> 5
		shift := uint(i & 31)

		pb := (phc[block] >> shift) & 1
		mb := (mhc[block] >> shift) & 1

		xv := eq | mv
		xh := ((((eq | mb) & pv) + pv) ^ pv) | eq | mb

		ph := mv | ^(xh | pv)
		mh := pv & xh

		score += int((ph >> lastShift) & 1)
		score -= int((mh >> lastShift) & 1)

		if ((ph >> 31) ^ pb) != 0 {
			phc[block] ^= uint32(1) << shift
		}

		if ((mh >> 31) ^ mb) != 0 {
			mhc[block] ^= uint32(1) << shift
		}

		ph = (ph << 1) | pb
		mh = (mh << 1) | mb

		pv = mh | ^(xv | ph)
		mv = ph & xv
	}

	for k := start; k < end; k++ {
		peq[longer[k]] = 0
	}

	return score
}

func distanceUnits(a, b []uint16, peq []uint32) int {
	if len(a) < len(b) {
		a, b = b, a
	}

	if len(b) == 0 {
		return len(a)
	}

	if len(a) <= 32 {
		return myers32(a, b, peq)
	}

	return myersX(a, b, peq)
}

// Distance returns the Levenshtein edit distance between a and b.
func Distance(a, b string) int {
	aUnits := utf16Units(a)
	bUnits := utf16Units(b)

	peq := peqPool.Get().([]uint32)
	defer peqPool.Put(peq)

	return distanceUnits(aUnits, bUnits, peq)
}

// Closest returns the string in arr having the smallest edit distance
// from str. In the event of a tie, the earliest string is returned.
//
// Closest panics when arr is empty.
func Closest(str string, arr []string) string {
	if len(arr) == 0 {
		panic("levenshtein: Closest called with an empty slice")
	}

	query := utf16Units(str)

	peq := peqPool.Get().([]uint32)
	defer peqPool.Put(peq)

	minDistance := int(^uint(0) >> 1)
	minIndex := 0

	for i, candidate := range arr {
		dist := distanceUnits(query, utf16Units(candidate), peq)

		if dist < minDistance {
			minDistance = dist
			minIndex = i
		}
	}

	return arr[minIndex]
}
