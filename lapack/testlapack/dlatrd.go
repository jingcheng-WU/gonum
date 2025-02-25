// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"github.com/jingcheng-WU/gonum/blas"
	"github.com/jingcheng-WU/gonum/blas/blas64"
)

type Dlatrder interface {
	Dlatrd(uplo blas.Uplo, n, nb int, a []float64, lda int, e, tau, w []float64, ldw int)
}

func DlatrdTest(t *testing.T, impl Dlatrder) {
	const tol = 1e-14

	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		for _, test := range []struct {
			n, nb, lda, ldw int
		}{
			{5, 2, 0, 0},
			{5, 5, 0, 0},

			{5, 3, 10, 11},
			{5, 5, 10, 11},
		} {
			n := test.n
			nb := test.nb
			lda := test.lda
			if lda == 0 {
				lda = n
			}
			ldw := test.ldw
			if ldw == 0 {
				ldw = nb
			}

			// Allocate n×n matrix A and fill it with random numbers.
			a := make([]float64, n*lda)
			for i := range a {
				a[i] = rnd.NormFloat64()
			}

			// Allocate output slices and matrix W and fill them
			// with NaN. All their elements should be overwritten by
			// Dlatrd.
			e := make([]float64, n-1)
			for i := range e {
				e[i] = math.NaN()
			}
			tau := make([]float64, n-1)
			for i := range tau {
				tau[i] = math.NaN()
			}
			w := make([]float64, n*ldw)
			for i := range w {
				w[i] = math.NaN()
			}

			aCopy := make([]float64, len(a))
			copy(aCopy, a)

			// Reduce nb rows and columns of the symmetric matrix A
			// defined by uplo triangle to symmetric tridiagonal
			// form.
			impl.Dlatrd(uplo, n, nb, a, lda, e, tau, w, ldw)

			// Construct Q from elementary reflectors stored in
			// columns of A.
			q := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: n,
				Data:   make([]float64, n*n),
			}
			// Initialize Q to the identity matrix.
			for i := 0; i < n; i++ {
				q.Data[i*q.Stride+i] = 1
			}
			if uplo == blas.Upper {
				for i := n - 1; i >= n-nb; i-- {
					if i == 0 {
						continue
					}

					// Extract the elementary reflector v from A.
					v := blas64.Vector{
						Inc:  1,
						Data: make([]float64, n),
					}
					for j := 0; j < i-1; j++ {
						v.Data[j] = a[j*lda+i]
					}
					v.Data[i-1] = 1

					// Compute H = I - tau[i-1] * v * vᵀ.
					h := blas64.General{
						Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n),
					}
					for j := 0; j < n; j++ {
						h.Data[j*n+j] = 1
					}
					blas64.Ger(-tau[i-1], v, v, h)

					// Update Q <- Q * H.
					qTmp := blas64.General{
						Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n),
					}
					copy(qTmp.Data, q.Data)
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qTmp, h, 0, q)
				}
			} else {
				for i := 0; i < nb; i++ {
					if i == n-1 {
						continue
					}

					// Extract the elementary reflector v from A.
					v := blas64.Vector{
						Inc:  1,
						Data: make([]float64, n),
					}
					v.Data[i+1] = 1
					for j := i + 2; j < n; j++ {
						v.Data[j] = a[j*lda+i]
					}

					// Compute H = I - tau[i] * v * vᵀ.
					h := blas64.General{
						Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n),
					}
					for j := 0; j < n; j++ {
						h.Data[j*n+j] = 1
					}
					blas64.Ger(-tau[i], v, v, h)

					// Update Q <- Q * H.
					qTmp := blas64.General{
						Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n),
					}
					copy(qTmp.Data, q.Data)
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qTmp, h, 0, q)
				}
			}
			name := fmt.Sprintf("uplo=%c,n=%v,nb=%v", uplo, n, nb)
			if resid := residualOrthogonal(q, false); resid > tol {
				t.Errorf("Case %v: Q is not orthogonal; resid=%v, want<=%v", name, resid, tol)
			}
			aGen := genFromSym(blas64.Symmetric{N: n, Stride: lda, Uplo: uplo, Data: aCopy})
			if !dlatrdCheckDecomposition(t, uplo, n, nb, e, a, lda, aGen, q, tol) {
				t.Errorf("Case %v: Decomposition mismatch", name)
			}
		}
	}
}

// dlatrdCheckDecomposition checks that the first nb rows have been successfully
// reduced.
func dlatrdCheckDecomposition(t *testing.T, uplo blas.Uplo, n, nb int, e, a []float64, lda int, aGen, q blas64.General, tol float64) bool {
	// Compute ans = Qᵀ * A * Q.
	// ans should be a tridiagonal matrix in the first or last nb rows and
	// columns, depending on uplo.
	tmp := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	ans := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, aGen, 0, tmp)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, q, 0, ans)

	// Compare the output of Dlatrd (stored in a and e) with the explicit
	// reduction to tridiagonal matrix Qᵀ * A * Q (stored in ans).
	if uplo == blas.Upper {
		for i := n - nb; i < n; i++ {
			for j := 0; j < n; j++ {
				v := ans.Data[i*ans.Stride+j]
				switch {
				case i == j:
					// Diagonal elements of a and ans should match.
					if math.Abs(v-a[i*lda+j]) > tol {
						return false
					}
				case i == j-1:
					// Superdiagonal elements in a should be 1.
					if math.Abs(a[i*lda+j]-1) > tol {
						return false
					}
					// Superdiagonal elements of ans should match e.
					if math.Abs(v-e[i]) > tol {
						return false
					}
				case i == j+1:
				default:
					// All other elements should be 0.
					if math.Abs(v) > tol {
						return false
					}
				}
			}
		}
	} else {
		for i := 0; i < nb; i++ {
			for j := 0; j < n; j++ {
				v := ans.Data[i*ans.Stride+j]
				switch {
				case i == j:
					// Diagonal elements of a and ans should match.
					if math.Abs(v-a[i*lda+j]) > tol {
						return false
					}
				case i == j-1:
				case i == j+1:
					// Subdiagonal elements in a should be 1.
					if math.Abs(a[i*lda+j]-1) > tol {
						return false
					}
					// Subdiagonal elements of ans should match e.
					if math.Abs(v-e[i-1]) > tol {
						return false
					}
				default:
					// All other elements should be 0.
					if math.Abs(v) > tol {
						return false
					}
				}
			}
		}
	}
	return true
}

// genFromSym constructs a (symmetric) general matrix from the data in the
// symmetric.
// TODO(btracey): Replace other constructions of this with a call to this function.
func genFromSym(a blas64.Symmetric) blas64.General {
	n := a.N
	lda := a.Stride
	uplo := a.Uplo
	b := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			v := a.Data[i*lda+j]
			if uplo == blas.Lower {
				v = a.Data[j*lda+i]
			}
			b.Data[i*n+j] = v
			b.Data[j*n+i] = v
		}
	}
	return b
}
