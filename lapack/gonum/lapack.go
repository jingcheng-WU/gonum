// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "github.com/jingcheng-WU/gonum/lapack"

// Implementation is the native Go implementation of LAPACK routines. It
// is built on top of calls to the return of blas64.Implementation(), so while
// this code is in pure Go, the underlying BLAS implementation may not be.
type Implementation struct{}

var _ lapack.Float64 = Implementation{}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

const (
	// dlamchE is the machine epsilon. For IEEE this is 2^{-53}.
	dlamchE = 1.1102230246251565e-16

	// dlamchB is the radix of the machine (the base of the number system).
	dlamchB = 2

	// dlamchP is base * eps.
	dlamchP = dlamchB * dlamchE

	// dlamchS is the "safe minimum", that is, the lowest number such that
	// 1/dlamchS does not overflow, or also the smallest normal number.
	// For IEEE this is 2^{-1022}.
	dlamchS = 2.2250738585072014e-308

	// (rtmin,rtmax) is a range of well-scaled numbers whose square
	// or sum of squares is also safe.
	// drtmin is sqrt(dlamchS/dlamchP)
	drtmin = 1.0010415475915505e-146
	drtmax = 1 / drtmin
)
