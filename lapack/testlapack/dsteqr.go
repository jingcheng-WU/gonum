// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"testing"

	"golang.org/x/exp/rand"

	"github.com/jingcheng-WU/gonum/blas"
	"github.com/jingcheng-WU/gonum/blas/blas64"
	"github.com/jingcheng-WU/gonum/floats"
	"github.com/jingcheng-WU/gonum/lapack"
)

type Dsteqrer interface {
	Dsteqr(compz lapack.EVComp, n int, d, e, z []float64, ldz int, work []float64) (ok bool)
	Dorgtrer
}

func DsteqrTest(t *testing.T, impl Dsteqrer) {
	rnd := rand.New(rand.NewSource(1))
	for _, compz := range []lapack.EVComp{lapack.EVOrig, lapack.EVTridiag} {
		for _, test := range []struct {
			n, lda int
		}{
			{1, 0},
			{4, 0},
			{8, 0},
			{10, 0},

			{2, 10},
			{8, 10},
			{10, 20},
		} {
			for cas := 0; cas < 100; cas++ {
				n := test.n
				lda := test.lda
				if lda == 0 {
					lda = n
				}
				d := make([]float64, n)
				for i := range d {
					d[i] = rnd.Float64()
				}
				e := make([]float64, n-1)
				for i := range e {
					e[i] = rnd.Float64()
				}
				a := make([]float64, n*lda)
				for i := range a {
					a[i] = rnd.Float64()
				}
				dCopy := make([]float64, len(d))
				copy(dCopy, d)
				eCopy := make([]float64, len(e))
				copy(eCopy, e)
				aCopy := make([]float64, len(a))
				copy(aCopy, a)
				if compz == lapack.EVOrig {
					uplo := blas.Upper
					tau := make([]float64, n)
					work := make([]float64, 1)
					impl.Dsytrd(blas.Upper, n, a, lda, d, e, tau, work, -1)
					work = make([]float64, int(work[0]))
					// Reduce A to symmetric tridiagonal form.
					impl.Dsytrd(uplo, n, a, lda, d, e, tau, work, len(work))
					// Compute the orthogonal matrix Q.
					impl.Dorgtr(uplo, n, a, lda, tau, work, len(work))
				} else {
					for i := 0; i < n; i++ {
						for j := 0; j < n; j++ {
							a[i*lda+j] = 0
							if i == j {
								a[i*lda+j] = 1
							}
						}
					}
				}
				work := make([]float64, 2*n)

				aDecomp := make([]float64, len(a))
				copy(aDecomp, a)
				dDecomp := make([]float64, len(d))
				copy(dDecomp, d)
				eDecomp := make([]float64, len(e))
				copy(eDecomp, e)
				impl.Dsteqr(compz, n, d, e, a, lda, work)
				dAns := make([]float64, len(d))
				copy(dAns, d)

				var truth blas64.General
				if compz == lapack.EVOrig {
					truth = blas64.General{
						Rows:   n,
						Cols:   n,
						Stride: n,
						Data:   make([]float64, n*n),
					}
					for i := 0; i < n; i++ {
						for j := i; j < n; j++ {
							v := aCopy[i*lda+j]
							truth.Data[i*truth.Stride+j] = v
							truth.Data[j*truth.Stride+i] = v
						}
					}
				} else {
					truth = blas64.General{
						Rows:   n,
						Cols:   n,
						Stride: n,
						Data:   make([]float64, n*n),
					}
					for i := 0; i < n; i++ {
						truth.Data[i*truth.Stride+i] = dCopy[i]
						if i != n-1 {
							truth.Data[(i+1)*truth.Stride+i] = eCopy[i]
							truth.Data[i*truth.Stride+i+1] = eCopy[i]
						}
					}
				}

				V := blas64.General{
					Rows:   n,
					Cols:   n,
					Stride: lda,
					Data:   a,
				}
				if !eigenDecompCorrect(d, truth, V) {
					t.Errorf("Eigen reconstruction mismatch. fromFull = %v, n = %v",
						compz == lapack.EVOrig, n)
				}

				// Compare eigenvalues when not computing eigenvectors.
				for i := range work {
					work[i] = rnd.Float64()
				}
				impl.Dsteqr(lapack.EVCompNone, n, dDecomp, eDecomp, aDecomp, lda, work)
				if !floats.EqualApprox(d, dAns, 1e-8) {
					t.Errorf("Eigenvalue mismatch when eigenvectors not computed")
				}
			}
		}
	}
}

// eigenDecompCorrect returns whether the eigen decomposition is correct.
// It checks if
//  A * v ≈ λ * v
// where the eigenvalues λ are stored in values, and the eigenvectors are stored
// in the columns of v.
func eigenDecompCorrect(values []float64, A, V blas64.General) bool {
	n := A.Rows
	for i := 0; i < n; i++ {
		lambda := values[i]
		vector := make([]float64, n)
		ans2 := make([]float64, n)
		for j := range vector {
			v := V.Data[j*V.Stride+i]
			vector[j] = v
			ans2[j] = lambda * v
		}
		v := blas64.Vector{Inc: 1, Data: vector}
		ans1 := blas64.Vector{Inc: 1, Data: make([]float64, n)}
		blas64.Gemv(blas.NoTrans, 1, A, v, 0, ans1)
		if !floats.EqualApprox(ans1.Data, ans2, 1e-8) {
			return false
		}
	}
	return true
}
