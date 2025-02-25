// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import (
	"math"
	"testing"

	"github.com/jingcheng-WU/gonum/blas"
	"github.com/jingcheng-WU/gonum/floats"
)

type Dsymmer interface {
	Dsymm(s blas.Side, ul blas.Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
}

func DsymmTest(t *testing.T, blasser Dsymmer) {
	for i, test := range []struct {
		m     int
		n     int
		side  blas.Side
		ul    blas.Uplo
		a     [][]float64
		b     [][]float64
		c     [][]float64
		alpha float64
		beta  float64
		ans   [][]float64
	}{
		{
			side: blas.Left,
			ul:   blas.Upper,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 3, 4},
				{0, 6, 7},
				{0, 0, 10},
			},
			b: [][]float64{
				{2, 3, 4, 8},
				{5, 6, 7, 15},
				{8, 9, 10, 20},
			},
			c: [][]float64{
				{8, 12, 2, 1},
				{9, 12, 9, 9},
				{12, 1, -1, 5},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{126, 156, 144, 285},
				{211, 252, 275, 535},
				{282, 291, 327, 689},
			},
		},
		{
			side: blas.Left,
			ul:   blas.Upper,
			m:    4,
			n:    3,
			a: [][]float64{
				{2, 3, 4, 8},
				{0, 6, 7, 9},
				{0, 0, 10, 10},
				{0, 0, 0, 11},
			},
			b: [][]float64{
				{2, 3, 4},
				{5, 6, 7},
				{8, 9, 10},
				{2, 1, 1},
			},
			c: [][]float64{
				{8, 12, 2},
				{9, 12, 9},
				{12, 1, -1},
				{1, 9, 5},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{158, 172, 160},
				{247, 270, 293},
				{322, 311, 347},
				{329, 385, 427},
			},
		},
		{
			side: blas.Left,
			ul:   blas.Lower,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 0, 0},
				{3, 6, 0},
				{4, 7, 10},
			},
			b: [][]float64{
				{2, 3, 4, 8},
				{5, 6, 7, 15},
				{8, 9, 10, 20},
			},
			c: [][]float64{
				{8, 12, 2, 1},
				{9, 12, 9, 9},
				{12, 1, -1, 5},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{126, 156, 144, 285},
				{211, 252, 275, 535},
				{282, 291, 327, 689},
			},
		},
		{
			side: blas.Left,
			ul:   blas.Lower,
			m:    4,
			n:    3,
			a: [][]float64{
				{2, 0, 0, 0},
				{3, 6, 0, 0},
				{4, 7, 10, 0},
				{8, 9, 10, 11},
			},
			b: [][]float64{
				{2, 3, 4},
				{5, 6, 7},
				{8, 9, 10},
				{2, 1, 1},
			},
			c: [][]float64{
				{8, 12, 2},
				{9, 12, 9},
				{12, 1, -1},
				{1, 9, 5},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{158, 172, 160},
				{247, 270, 293},
				{322, 311, 347},
				{329, 385, 427},
			},
		},
		{
			side: blas.Right,
			ul:   blas.Upper,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 0, 0, 0},
				{3, 6, 0, 0},
				{4, 7, 10, 0},
				{3, 4, 5, 6},
			},
			b: [][]float64{
				{2, 3, 4, 9},
				{5, 6, 7, -3},
				{8, 9, 10, -2},
			},
			c: [][]float64{
				{8, 12, 2, 10},
				{9, 12, 9, 10},
				{12, 1, -1, 10},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{32, 72, 86, 138},
				{47, 108, 167, -6},
				{68, 111, 197, 6},
			},
		},
		{
			side: blas.Right,
			ul:   blas.Upper,
			m:    4,
			n:    3,
			a: [][]float64{
				{2, 0, 0},
				{3, 6, 0},
				{4, 7, 10},
			},
			b: [][]float64{
				{2, 3, 4},
				{5, 6, 7},
				{8, 9, 10},
				{2, 1, 1},
			},
			c: [][]float64{
				{8, 12, 2},
				{9, 12, 9},
				{12, 1, -1},
				{1, 9, 5},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{32, 72, 86},
				{47, 108, 167},
				{68, 111, 197},
				{11, 39, 35},
			},
		},
		{
			side: blas.Right,
			ul:   blas.Lower,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 0, 0, 0},
				{3, 6, 0, 0},
				{4, 7, 10, 0},
				{3, 4, 5, 6},
			},
			b: [][]float64{
				{2, 3, 4, 2},
				{5, 6, 7, 1},
				{8, 9, 10, 1},
			},
			c: [][]float64{
				{8, 12, 2, 1},
				{9, 12, 9, 9},
				{12, 1, -1, 5},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{94, 156, 164, 103},
				{145, 244, 301, 187},
				{208, 307, 397, 247},
			},
		},
		{
			side: blas.Right,
			ul:   blas.Lower,
			m:    4,
			n:    3,
			a: [][]float64{
				{2, 0, 0},
				{3, 6, 0},
				{4, 7, 10},
			},
			b: [][]float64{
				{2, 3, 4},
				{5, 6, 7},
				{8, 9, 10},
				{2, 1, 1},
			},
			c: [][]float64{
				{8, 12, 2},
				{9, 12, 9},
				{12, 1, -1},
				{1, 9, 5},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{82, 140, 144},
				{139, 236, 291},
				{202, 299, 387},
				{25, 65, 65},
			},
		},

		{
			side: blas.Left,
			ul:   blas.Upper,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 3, 4},
				{0, 6, 7},
				{0, 0, 10},
			},
			b: [][]float64{
				{2, 3, 4, 8},
				{5, 6, 7, 15},
				{8, 9, 10, 20},
			},
			c: [][]float64{
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
			},
			alpha: 2,
			ans: [][]float64{
				{102, 120, 138, 282},
				{184, 216, 248, 508},
				{246, 288, 330, 674},
			},
		},
		{
			side: blas.Left,
			ul:   blas.Upper,
			m:    4,
			n:    3,
			a: [][]float64{
				{2, 3, 4, 8},
				{0, 6, 7, 9},
				{0, 0, 10, 10},
				{0, 0, 0, 11},
			},
			b: [][]float64{
				{2, 3, 4},
				{5, 6, 7},
				{8, 9, 10},
				{2, 1, 1},
			},
			c: [][]float64{
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
			},
			alpha: 2,
			ans: [][]float64{
				{134, 136, 154},
				{220, 234, 266},
				{286, 308, 350},
				{326, 358, 412},
			},
		},
		{
			side: blas.Left,
			ul:   blas.Lower,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 0, 0},
				{3, 6, 0},
				{4, 7, 10},
			},
			b: [][]float64{
				{2, 3, 4, 8},
				{5, 6, 7, 15},
				{8, 9, 10, 20},
			},
			c: [][]float64{
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
			},
			alpha: 2,
			ans: [][]float64{
				{102, 120, 138, 282},
				{184, 216, 248, 508},
				{246, 288, 330, 674},
			},
		},
		{
			side: blas.Left,
			ul:   blas.Lower,
			m:    4,
			n:    3,
			a: [][]float64{
				{2, 0, 0, 0},
				{3, 6, 0, 0},
				{4, 7, 10, 0},
				{8, 9, 10, 11},
			},
			b: [][]float64{
				{2, 3, 4},
				{5, 6, 7},
				{8, 9, 10},
				{2, 1, 1},
			},
			c: [][]float64{
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
			},
			alpha: 2,
			ans: [][]float64{
				{134, 136, 154},
				{220, 234, 266},
				{286, 308, 350},
				{326, 358, 412},
			},
		},
		{
			side: blas.Right,
			ul:   blas.Upper,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 0, 0, 0},
				{3, 6, 0, 0},
				{4, 7, 10, 0},
				{3, 4, 5, 6},
			},
			b: [][]float64{
				{2, 3, 4, 9},
				{5, 6, 7, -3},
				{8, 9, 10, -2},
			},
			c: [][]float64{
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
			},
			alpha: 2,
			ans: [][]float64{
				{8, 36, 80, 108},
				{20, 72, 140, -36},
				{32, 108, 200, -24},
			},
		},
		{
			side: blas.Right,
			ul:   blas.Upper,
			m:    4,
			n:    3,
			a: [][]float64{
				{2, 0, 0},
				{3, 6, 0},
				{4, 7, 10},
			},
			b: [][]float64{
				{2, 3, 4},
				{5, 6, 7},
				{8, 9, 10},
				{2, 1, 1},
			},
			c: [][]float64{
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
			},
			alpha: 2,
			ans: [][]float64{
				{8, 36, 80, 20},
				{72, 140, 32, 108},
				{200, 8, 12, 20},
			},
		},
		{
			side: blas.Right,
			ul:   blas.Lower,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 0, 0, 0},
				{3, 6, 0, 0},
				{4, 7, 10, 0},
				{3, 4, 5, 6},
			},
			b: [][]float64{
				{2, 3, 4, 2},
				{5, 6, 7, 1},
				{8, 9, 10, 1},
			},
			c: [][]float64{
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
			},
			alpha: 2,
			ans: [][]float64{
				{70, 120, 158, 100},
				{118, 208, 274, 160},
				{172, 304, 400, 232},
			},
		},
		{
			side: blas.Right,
			ul:   blas.Lower,
			m:    4,
			n:    3,
			a: [][]float64{
				{2, 0, 0},
				{3, 6, 0},
				{4, 7, 10},
			},
			b: [][]float64{
				{2, 3, 4},
				{5, 6, 7},
				{8, 9, 10},
				{2, 1, 1},
			},
			c: [][]float64{
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
				{math.NaN(), math.NaN(), math.NaN()},
			},
			alpha: 2,
			ans: [][]float64{
				{58, 104, 138},
				{112, 200, 264},
				{166, 296, 390},
				{22, 38, 50},
			},
		},
	} {
		aFlat := flatten(test.a)
		bFlat := flatten(test.b)
		cFlat := flatten(test.c)
		ansFlat := flatten(test.ans)
		blasser.Dsymm(test.side, test.ul, test.m, test.n, test.alpha, aFlat, len(test.a[0]), bFlat, test.n, test.beta, cFlat, test.n)
		if !floats.EqualApprox(cFlat, ansFlat, 1e-14) {
			t.Errorf("Case %v: Want %v, got %v.", i, ansFlat, cFlat)
		}
	}
}
