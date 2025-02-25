// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c128_test

import (
	"math"
	"math/cmplx"
	"testing"

	"github.com/jingcheng-WU/gonum/cmplxs/cscalar"
	"github.com/jingcheng-WU/gonum/floats/scalar"
)

const (
	msgVal   = "%v: unexpected value at %v Got: %v Expected: %v"
	msgGuard = "%v: Guard violated in %s vector %v %v"
)

var (
	nan = math.NaN()

	cnan = cmplx.NaN()
	cinf = cmplx.Inf()
)

// TODO(kortschak): Harmonise the situation in asm/{f32,f64} and their sinks.
const testLen = 1e5

var x = make([]complex128, testLen)

// guardVector copies the source vector (vec) into a new slice with guards.
// Guards guarded[:gdLn] and guarded[len-gdLn:] will be filled with sigil value gdVal.
func guardVector(vec []complex128, gdVal complex128, gdLn int) (guarded []complex128) {
	guarded = make([]complex128, len(vec)+gdLn*2)
	copy(guarded[gdLn:], vec)
	for i := 0; i < gdLn; i++ {
		guarded[i] = gdVal
		guarded[len(guarded)-1-i] = gdVal
	}
	return guarded
}

// isValidGuard will test for violated guards, generated by guardVector.
func isValidGuard(vec []complex128, gdVal complex128, gdLn int) bool {
	for i := 0; i < gdLn; i++ {
		if !cscalar.Same(vec[i], gdVal) || !cscalar.Same(vec[len(vec)-1-i], gdVal) {
			return false
		}
	}
	return true
}

// guardIncVector copies the source vector (vec) into a new incremented slice with guards.
// End guards will be length gdLen.
// Internal and end guards will be filled with sigil value gdVal.
func guardIncVector(vec []complex128, gdVal complex128, inc, gdLen int) (guarded []complex128) {
	if inc < 0 {
		inc = -inc
	}
	inrLen := len(vec) * inc
	guarded = make([]complex128, inrLen+gdLen*2)
	for i := range guarded {
		guarded[i] = gdVal
	}
	for i, v := range vec {
		guarded[gdLen+i*inc] = v
	}
	return guarded
}

// checkValidIncGuard will test for violated guards, generated by guardIncVector
func checkValidIncGuard(t *testing.T, vec []complex128, gdVal complex128, inc, gdLen int) {
	srcLn := len(vec) - 2*gdLen
	for i := range vec {
		switch {
		case cscalar.Same(vec[i], gdVal):
			// Correct value
		case (i-gdLen)%inc == 0 && (i-gdLen)/inc < len(vec):
			// Ignore input values
		case i < gdLen:
			t.Errorf("Front guard violated at %d %v", i, vec[:gdLen])
		case i > gdLen+srcLn:
			t.Errorf("Back guard violated at %d %v", i-gdLen-srcLn, vec[gdLen+srcLn:])
		default:
			t.Errorf("Internal guard violated at %d %v", i-gdLen, vec[gdLen:gdLen+srcLn])
		}
	}
}

// sameApprox tests for nan-aware equality within tolerance.
func sameApprox(a, b, tol float64) bool {
	return scalar.Same(a, b) || scalar.EqualWithinAbsOrRel(a, b, tol, tol)
}

// sameCmplxApprox tests for nan-aware equality within tolerance.
func sameCmplxApprox(a, b complex128, tol float64) bool {
	return cscalar.Same(a, b) || cscalar.EqualWithinAbsOrRel(a, b, tol, tol)
}

var ( // Offset sets for testing alignment handling in Unitary assembly functions.
	align1 = []int{0, 1}
)
