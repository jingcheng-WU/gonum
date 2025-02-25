// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import (
	"testing"

	"github.com/jingcheng-WU/gonum/blas"
)

type Dgbmver interface {
	Dgbmv(tA blas.Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
}

func DgbmvTest(t *testing.T, blasser Dgbmver) {
	for i, test := range []struct {
		tA     blas.Transpose
		m, n   int
		kL, kU int
		alpha  float64
		a      [][]float64
		lda    int
		x      []float64
		beta   float64
		y      []float64
		ans    []float64
	}{
		{
			tA:    blas.NoTrans,
			m:     9,
			n:     6,
			lda:   4,
			kL:    2,
			kU:    1,
			alpha: 3.0,
			beta:  2.0,
			a: [][]float64{
				{5, 3, 0, 0, 0, 0},
				{-1, 2, 9, 0, 0, 0},
				{4, 8, 3, 6, 0, 0},
				{0, -1, 8, 2, 1, 0},
				{0, 0, 9, 9, 9, 5},
				{0, 0, 0, 2, -3, 2},
				{0, 0, 0, 0, 1, 5},
				{0, 0, 0, 0, 0, 6},
			},
			x:   []float64{1, 2, 3, 4, 5, 6},
			y:   []float64{-1, -2, -3, -4, -5, -6, -7, -8, -9},
			ans: []float64{31, 86, 153, 97, 404, 3, 91, 92, -18},
		},
		{
			tA:    blas.Trans,
			m:     9,
			n:     6,
			lda:   4,
			kL:    2,
			kU:    1,
			alpha: 3.0,
			beta:  2.0,
			a: [][]float64{
				{5, 3, 0, 0, 0, 0},
				{-1, 2, 9, 0, 0, 0},
				{4, 8, 3, 6, 0, 0},
				{0, -1, 8, 2, 1, 0},
				{0, 0, 9, 9, 9, 5},
				{0, 0, 0, 2, -3, 2},
				{0, 0, 0, 0, 1, 5},
				{0, 0, 0, 0, 0, 6},
			},
			x:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			y:   []float64{-1, -2, -3, -4, -5, -6},
			ans: []float64{43, 77, 306, 241, 104, 348},
		},
		{
			tA:    blas.NoTrans,
			m:     6,
			n:     3,
			lda:   1,
			kL:    0,
			kU:    0,
			alpha: 2.0,
			beta:  1.0,
			a: [][]float64{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			x:   []float64{1, 2, 3},
			y:   []float64{-1, -2, -3, -4, -5, -6},
			ans: []float64{1, 2, 3, -4, -5, -6},
		},
		{
			tA:    blas.Trans,
			m:     6,
			n:     3,
			lda:   1,
			kL:    0,
			kU:    0,
			alpha: 2.0,
			beta:  1.0,
			a: [][]float64{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			x:   []float64{1, 2, 3, 4, 5, 6},
			y:   []float64{-1, -2, -3},
			ans: []float64{1, 2, 3},
		},
		{
			tA:    blas.NoTrans,
			m:     3,
			n:     5,
			lda:   4,
			kL:    1,
			kU:    2,
			alpha: 2.0,
			beta:  1.0,
			a: [][]float64{
				{1, 2, 3, 0, 0},
				{1, 3, 6, 9, 0},
				{0, 1, 1, 1, 1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			y:   []float64{-1, -2, -3},
			ans: []float64{-1 + 2*(1+4+9), -2 + 2*(1+6+18+36), -3 + 2*(2+3+4+5)},
		},
	} {
		extra := 3
		aFlat := flattenBanded(test.a, test.kU, test.kL)
		incTest := func(incX, incY, extra int) {
			xnew := makeIncremented(test.x, incX, extra)
			ynew := makeIncremented(test.y, incY, extra)
			ans := makeIncremented(test.ans, incY, extra)
			blasser.Dgbmv(test.tA, test.m, test.n, test.kL, test.kU, test.alpha, aFlat, test.lda, xnew, incX, test.beta, ynew, incY)
			if !dSliceTolEqual(ans, ynew) {
				t.Errorf("Case %v: Want %v, got %v", i, ans, ynew)
			}
		}
		incTest(1, 1, extra)
		incTest(1, 3, extra)
		incTest(1, -3, extra)
		incTest(2, 3, extra)
		incTest(2, -3, extra)
		incTest(3, 2, extra)
		incTest(-3, 2, extra)
	}
}
