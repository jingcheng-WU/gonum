// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import (
	"testing"

	"github.com/jingcheng-WU/gonum/blas"
)

type Dsyrer interface {
	Dsyr(ul blas.Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int)
}

func DsyrTest(t *testing.T, blasser Dsyrer) {
	for i, test := range []struct {
		ul    blas.Uplo
		n     int
		a     [][]float64
		x     []float64
		alpha float64
		ans   [][]float64
	}{
		{
			ul: blas.Upper,
			n:  4,
			a: [][]float64{
				{10, 2, 0, 1},
				{0, 1, 2, 3},
				{0, 0, 9, 15},
				{0, 0, 0, -6},
			},
			x:     []float64{1, 2, 0, 5},
			alpha: 8,
			ans: [][]float64{
				{18, 18, 0, 41},
				{0, 33, 2, 83},
				{0, 0, 9, 15},
				{0, 0, 0, 194},
			},
		},
		{
			ul: blas.Lower,
			n:  3,
			a: [][]float64{
				{10, 2, 0},
				{4, 1, 2},
				{2, 7, 9},
			},
			x:     []float64{3, 0, 5},
			alpha: 8,
			ans: [][]float64{
				{82, 2, 0},
				{4, 1, 2},
				{122, 7, 209},
			},
		},
	} {
		incTest := func(incX, extra int) {
			xnew := makeIncremented(test.x, incX, extra)
			aFlat := flatten(test.a)
			ans := flatten(test.ans)
			lda := test.n
			blasser.Dsyr(test.ul, test.n, test.alpha, xnew, incX, aFlat, lda)
			if !dSliceTolEqual(aFlat, ans) {
				t.Errorf("Case %v, idx %v: Want %v, got %v.", i, incX, ans, aFlat)
			}
		}
		incTest(1, 3)
		incTest(1, 0)
		incTest(3, 2)
		incTest(-2, 2)
	}
}
