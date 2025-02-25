// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"golang.org/x/exp/rand"

	"github.com/jingcheng-WU/gonum/blas"
	"github.com/jingcheng-WU/gonum/blas/blas64"
	"github.com/jingcheng-WU/gonum/lapack"
)

const (
	// dlamchE is the machine epsilon. For IEEE this is 2^{-53}.
	dlamchE = 0x1p-53
	dlamchB = 2
	dlamchP = dlamchB * dlamchE
	// dlamchS is the smallest normal number. For IEEE this is 2^{-1022}.
	dlamchS = 0x1p-1022
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// worklen describes how much workspace a test should use.
type worklen int

const (
	minimumWork worklen = iota
	mediumWork
	optimumWork
)

func (wl worklen) String() string {
	switch wl {
	case minimumWork:
		return "minimum"
	case mediumWork:
		return "medium"
	case optimumWork:
		return "optimum"
	}
	return ""
}

func normToString(norm lapack.MatrixNorm) string {
	switch norm {
	case lapack.MaxAbs:
		return "MaxAbs"
	case lapack.MaxRowSum:
		return "MaxRowSum"
	case lapack.MaxColumnSum:
		return "MaxColSum"
	case lapack.Frobenius:
		return "Frobenius"
	default:
		panic("invalid norm")
	}
}

func uploToString(uplo blas.Uplo) string {
	switch uplo {
	case blas.Lower:
		return "Lower"
	case blas.Upper:
		return "Upper"
	default:
		panic("invalid uplo")
	}
}

func diagToString(diag blas.Diag) string {
	switch diag {
	case blas.NonUnit:
		return "NonUnit"
	case blas.Unit:
		return "Unit"
	default:
		panic("invalid diag")
	}
}

func sideToString(side blas.Side) string {
	switch side {
	case blas.Left:
		return "Left"
	case blas.Right:
		return "Right"
	default:
		panic("invalid side")
	}
}

func transToString(trans blas.Transpose) string {
	switch trans {
	case blas.NoTrans:
		return "NoTrans"
	case blas.Trans:
		return "Trans"
	case blas.ConjTrans:
		return "ConjTrans"
	default:
		panic("invalid trans")
	}
}

// nanSlice allocates a new slice of length n filled with NaN.
func nanSlice(n int) []float64 {
	s := make([]float64, n)
	for i := range s {
		s[i] = math.NaN()
	}
	return s
}

// randomSlice allocates a new slice of length n filled with random values.
func randomSlice(n int, rnd *rand.Rand) []float64 {
	s := make([]float64, n)
	for i := range s {
		s[i] = rnd.NormFloat64()
	}
	return s
}

// nanGeneral allocates a new r×c general matrix filled with NaN values.
func nanGeneral(r, c, stride int) blas64.General {
	if r < 0 || c < 0 {
		panic("bad matrix size")
	}
	if r == 0 || c == 0 {
		return blas64.General{Stride: max(1, stride)}
	}
	if stride < c {
		panic("bad stride")
	}
	return blas64.General{
		Rows:   r,
		Cols:   c,
		Stride: stride,
		Data:   nanSlice((r-1)*stride + c),
	}
}

// randomGeneral allocates a new r×c general matrix filled with random
// numbers. Out-of-range elements are filled with NaN values.
func randomGeneral(r, c, stride int, rnd *rand.Rand) blas64.General {
	ans := nanGeneral(r, c, stride)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			ans.Data[i*ans.Stride+j] = rnd.NormFloat64()
		}
	}
	return ans
}

// randomHessenberg allocates a new n×n Hessenberg matrix filled with zeros
// under the first subdiagonal and with random numbers elsewhere. Out-of-range
// elements are filled with NaN values.
func randomHessenberg(n, stride int, rnd *rand.Rand) blas64.General {
	ans := nanGeneral(n, n, stride)
	for i := 0; i < n; i++ {
		for j := 0; j < i-1; j++ {
			ans.Data[i*ans.Stride+j] = 0
		}
		for j := max(0, i-1); j < n; j++ {
			ans.Data[i*ans.Stride+j] = rnd.NormFloat64()
		}
	}
	return ans
}

// randomSchurCanonical returns a random, general matrix in Schur canonical
// form, that is, block upper triangular with 1×1 and 2×2 diagonal blocks where
// each 2×2 diagonal block has its diagonal elements equal and its off-diagonal
// elements of opposite sign. bad controls whether the returned matrix will have
// zero or tiny eigenvalues.
func randomSchurCanonical(n, stride int, bad bool, rnd *rand.Rand) (t blas64.General, wr, wi []float64) {
	t = randomGeneral(n, n, stride, rnd)
	// Zero out the lower triangle including the diagonal which will be set later.
	for i := 0; i < t.Rows; i++ {
		for j := 0; j <= i; j++ {
			t.Data[i*t.Stride+j] = 0
		}
	}
	// Randomly create 2×2 diagonal blocks.
	for i := 0; i < t.Rows; {
		a := rnd.NormFloat64()
		if bad && rnd.Float64() < 0.5 {
			if rnd.Float64() < 0.5 {
				// A quarter of real parts of eigenvalues will be tiny.
				a = dlamchS
			} else {
				// A quarter of them will be zero.
				a = 0
			}
		}

		// A half of eigenvalues will be real.
		if rnd.Float64() < 0.5 || i == t.Rows-1 {
			// Store 1×1 block at the diagonal of T.
			t.Data[i*t.Stride+i] = a
			wr = append(wr, a)
			wi = append(wi, 0)
			i++
			continue
		}

		// Diagonal elements are equal.
		d := a
		// Element under the diagonal is "normal".
		c := rnd.NormFloat64()
		// Element above the diagonal cannot be zero.
		var b float64
		if bad && rnd.Float64() < 0.5 {
			b = dlamchS
		} else {
			b = rnd.NormFloat64()
		}
		// Make sure off-diagonal elements are of opposite sign.
		if math.Signbit(b) == math.Signbit(c) {
			c *= -1
		}

		// Store 2×2 block at the diagonal of T.
		t.Data[i*t.Stride+i], t.Data[i*t.Stride+i+1] = a, b
		t.Data[(i+1)*t.Stride+i], t.Data[(i+1)*t.Stride+i+1] = c, d

		wr = append(wr, a, a)
		im := math.Sqrt(math.Abs(b)) * math.Sqrt(math.Abs(c))
		wi = append(wi, im, -im)
		i += 2
	}
	return t, wr, wi
}

// blockedUpperTriGeneral returns a normal random, general matrix in the form
//
//            c-k-l  k    l
//  A =    k [  0   A12  A13 ] if r-k-l >= 0;
//         l [  0    0   A23 ]
//     r-k-l [  0    0    0  ]
//
//          c-k-l  k    l
//  A =  k [  0   A12  A13 ] if r-k-l < 0;
//     r-k [  0    0   A23 ]
//
// where the k×k matrix A12 and l×l matrix is non-singular
// upper triangular. A23 is l×l upper triangular if r-k-l >= 0,
// otherwise A23 is (r-k)×l upper trapezoidal.
func blockedUpperTriGeneral(r, c, k, l, stride int, kblock bool, rnd *rand.Rand) blas64.General {
	t := l
	if kblock {
		t += k
	}
	ans := zeros(r, c, stride)
	for i := 0; i < min(r, t); i++ {
		var v float64
		for v == 0 {
			v = rnd.NormFloat64()
		}
		ans.Data[i*ans.Stride+i+(c-t)] = v
	}
	for i := 0; i < min(r, t); i++ {
		for j := i + (c - t) + 1; j < c; j++ {
			ans.Data[i*ans.Stride+j] = rnd.NormFloat64()
		}
	}
	return ans
}

// nanTriangular allocates a new r×c triangular matrix filled with NaN values.
func nanTriangular(uplo blas.Uplo, n, stride int) blas64.Triangular {
	if n < 0 {
		panic("bad matrix size")
	}
	if n == 0 {
		return blas64.Triangular{
			Stride: max(1, stride),
			Uplo:   uplo,
			Diag:   blas.NonUnit,
		}
	}
	if stride < n {
		panic("bad stride")
	}
	return blas64.Triangular{
		N:      n,
		Stride: stride,
		Data:   nanSlice((n-1)*stride + n),
		Uplo:   uplo,
		Diag:   blas.NonUnit,
	}
}

// generalOutsideAllNaN returns whether all out-of-range elements have NaN
// values.
func generalOutsideAllNaN(a blas64.General) bool {
	// Check after last column.
	for i := 0; i < a.Rows-1; i++ {
		for _, v := range a.Data[i*a.Stride+a.Cols : i*a.Stride+a.Stride] {
			if !math.IsNaN(v) {
				return false
			}
		}
	}
	// Check after last element.
	last := (a.Rows-1)*a.Stride + a.Cols
	if a.Rows == 0 || a.Cols == 0 {
		last = 0
	}
	for _, v := range a.Data[last:] {
		if !math.IsNaN(v) {
			return false
		}
	}
	return true
}

// triangularOutsideAllNaN returns whether all out-of-triangle elements have NaN
// values.
func triangularOutsideAllNaN(a blas64.Triangular) bool {
	if a.Uplo == blas.Upper {
		// Check below diagonal.
		for i := 0; i < a.N; i++ {
			for _, v := range a.Data[i*a.Stride : i*a.Stride+i] {
				if !math.IsNaN(v) {
					return false
				}
			}
		}
		// Check after last column.
		for i := 0; i < a.N-1; i++ {
			for _, v := range a.Data[i*a.Stride+a.N : i*a.Stride+a.Stride] {
				if !math.IsNaN(v) {
					return false
				}
			}
		}
	} else {
		// Check above diagonal.
		for i := 0; i < a.N-1; i++ {
			for _, v := range a.Data[i*a.Stride+i+1 : i*a.Stride+a.Stride] {
				if !math.IsNaN(v) {
					return false
				}
			}
		}
	}
	// Check after last element.
	for _, v := range a.Data[max(0, a.N-1)*a.Stride+a.N:] {
		if !math.IsNaN(v) {
			return false
		}
	}
	return true
}

// transposeGeneral returns a new general matrix that is the transpose of the
// input. Nothing is done with data outside the {rows, cols} limit of the general.
func transposeGeneral(a blas64.General) blas64.General {
	ans := blas64.General{
		Rows:   a.Cols,
		Cols:   a.Rows,
		Stride: a.Rows,
		Data:   make([]float64, a.Cols*a.Rows),
	}
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			ans.Data[j*ans.Stride+i] = a.Data[i*a.Stride+j]
		}
	}
	return ans
}

// columnNorms returns the column norms of a.
func columnNorms(m, n int, a []float64, lda int) []float64 {
	bi := blas64.Implementation()
	norms := make([]float64, n)
	for j := 0; j < n; j++ {
		norms[j] = bi.Dnrm2(m, a[j:], lda)
	}
	return norms
}

// extractVMat collects the single reflectors from a into a matrix.
func extractVMat(m, n int, a []float64, lda int, direct lapack.Direct, store lapack.StoreV) blas64.General {
	k := min(m, n)
	switch {
	default:
		panic("not implemented")
	case direct == lapack.Forward && store == lapack.ColumnWise:
		v := blas64.General{
			Rows:   m,
			Cols:   k,
			Stride: k,
			Data:   make([]float64, m*k),
		}
		for i := 0; i < k; i++ {
			for j := 0; j < i; j++ {
				v.Data[j*v.Stride+i] = 0
			}
			v.Data[i*v.Stride+i] = 1
			for j := i + 1; j < m; j++ {
				v.Data[j*v.Stride+i] = a[j*lda+i]
			}
		}
		return v
	case direct == lapack.Forward && store == lapack.RowWise:
		v := blas64.General{
			Rows:   k,
			Cols:   n,
			Stride: n,
			Data:   make([]float64, k*n),
		}
		for i := 0; i < k; i++ {
			for j := 0; j < i; j++ {
				v.Data[i*v.Stride+j] = 0
			}
			v.Data[i*v.Stride+i] = 1
			for j := i + 1; j < n; j++ {
				v.Data[i*v.Stride+j] = a[i*lda+j]
			}
		}
		return v
	}
}

// constructBidiagonal constructs a bidiagonal matrix with the given diagonal
// and off-diagonal elements.
func constructBidiagonal(uplo blas.Uplo, n int, d, e []float64) blas64.General {
	bMat := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}

	for i := 0; i < n-1; i++ {
		bMat.Data[i*bMat.Stride+i] = d[i]
		if uplo == blas.Upper {
			bMat.Data[i*bMat.Stride+i+1] = e[i]
		} else {
			bMat.Data[(i+1)*bMat.Stride+i] = e[i]
		}
	}
	bMat.Data[(n-1)*bMat.Stride+n-1] = d[n-1]
	return bMat
}

// constructVMat transforms the v matrix based on the storage.
func constructVMat(vMat blas64.General, store lapack.StoreV, direct lapack.Direct) blas64.General {
	m := vMat.Rows
	k := vMat.Cols
	switch {
	default:
		panic("not implemented")
	case store == lapack.ColumnWise && direct == lapack.Forward:
		ldv := k
		v := make([]float64, m*k)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				if j > i {
					v[i*ldv+j] = 0
				} else if j == i {
					v[i*ldv+i] = 1
				} else {
					v[i*ldv+j] = vMat.Data[i*vMat.Stride+j]
				}
			}
		}
		return blas64.General{
			Rows:   m,
			Cols:   k,
			Stride: k,
			Data:   v,
		}
	case store == lapack.RowWise && direct == lapack.Forward:
		ldv := m
		v := make([]float64, m*k)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				if j > i {
					v[j*ldv+i] = 0
				} else if j == i {
					v[j*ldv+i] = 1
				} else {
					v[j*ldv+i] = vMat.Data[i*vMat.Stride+j]
				}
			}
		}
		return blas64.General{
			Rows:   k,
			Cols:   m,
			Stride: m,
			Data:   v,
		}
	case store == lapack.ColumnWise && direct == lapack.Backward:
		rowsv := m
		ldv := k
		v := make([]float64, m*k)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				vrow := rowsv - i - 1
				vcol := k - j - 1
				if j > i {
					v[vrow*ldv+vcol] = 0
				} else if j == i {
					v[vrow*ldv+vcol] = 1
				} else {
					v[vrow*ldv+vcol] = vMat.Data[i*vMat.Stride+j]
				}
			}
		}
		return blas64.General{
			Rows:   rowsv,
			Cols:   ldv,
			Stride: ldv,
			Data:   v,
		}
	case store == lapack.RowWise && direct == lapack.Backward:
		rowsv := k
		ldv := m
		v := make([]float64, m*k)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				vcol := ldv - i - 1
				vrow := k - j - 1
				if j > i {
					v[vrow*ldv+vcol] = 0
				} else if j == i {
					v[vrow*ldv+vcol] = 1
				} else {
					v[vrow*ldv+vcol] = vMat.Data[i*vMat.Stride+j]
				}
			}
		}
		return blas64.General{
			Rows:   rowsv,
			Cols:   ldv,
			Stride: ldv,
			Data:   v,
		}
	}
}

func constructH(tau []float64, v blas64.General, store lapack.StoreV, direct lapack.Direct) blas64.General {
	m := v.Rows
	k := v.Cols
	if store == lapack.RowWise {
		m, k = k, m
	}
	h := blas64.General{
		Rows:   m,
		Cols:   m,
		Stride: m,
		Data:   make([]float64, m*m),
	}
	for i := 0; i < m; i++ {
		h.Data[i*m+i] = 1
	}
	for i := 0; i < k; i++ {
		vecData := make([]float64, m)
		if store == lapack.ColumnWise {
			for j := 0; j < m; j++ {
				vecData[j] = v.Data[j*v.Cols+i]
			}
		} else {
			for j := 0; j < m; j++ {
				vecData[j] = v.Data[i*v.Cols+j]
			}
		}
		vec := blas64.Vector{
			Inc:  1,
			Data: vecData,
		}

		hi := blas64.General{
			Rows:   m,
			Cols:   m,
			Stride: m,
			Data:   make([]float64, m*m),
		}
		for i := 0; i < m; i++ {
			hi.Data[i*m+i] = 1
		}
		// hi = I - tau * v * vᵀ
		blas64.Ger(-tau[i], vec, vec, hi)

		hcopy := blas64.General{
			Rows:   m,
			Cols:   m,
			Stride: m,
			Data:   make([]float64, m*m),
		}
		copy(hcopy.Data, h.Data)
		if direct == lapack.Forward {
			// H = H * H_I in forward mode
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, hcopy, hi, 0, h)
		} else {
			// H = H_I * H in backward mode
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, hi, hcopy, 0, h)
		}
	}
	return h
}

// constructQ constructs the Q matrix from the result of dgeqrf and dgeqr2.
func constructQ(kind string, m, n int, a []float64, lda int, tau []float64) blas64.General {
	k := min(m, n)
	return constructQK(kind, m, n, k, a, lda, tau)
}

// constructQK constructs the Q matrix from the result of dgeqrf and dgeqr2 using
// the first k reflectors.
func constructQK(kind string, m, n, k int, a []float64, lda int, tau []float64) blas64.General {
	var sz int
	switch kind {
	case "QR":
		sz = m
	case "LQ", "RQ":
		sz = n
	}

	q := blas64.General{
		Rows:   sz,
		Cols:   sz,
		Stride: max(1, sz),
		Data:   make([]float64, sz*sz),
	}
	for i := 0; i < sz; i++ {
		q.Data[i*sz+i] = 1
	}
	qCopy := blas64.General{
		Rows:   q.Rows,
		Cols:   q.Cols,
		Stride: q.Stride,
		Data:   make([]float64, len(q.Data)),
	}
	for i := 0; i < k; i++ {
		h := blas64.General{
			Rows:   sz,
			Cols:   sz,
			Stride: max(1, sz),
			Data:   make([]float64, sz*sz),
		}
		for j := 0; j < sz; j++ {
			h.Data[j*sz+j] = 1
		}
		vVec := blas64.Vector{
			Inc:  1,
			Data: make([]float64, sz),
		}
		switch kind {
		case "QR":
			vVec.Data[i] = 1
			for j := i + 1; j < sz; j++ {
				vVec.Data[j] = a[lda*j+i]
			}
		case "LQ":
			vVec.Data[i] = 1
			for j := i + 1; j < sz; j++ {
				vVec.Data[j] = a[i*lda+j]
			}
		case "RQ":
			for j := 0; j < n-k+i; j++ {
				vVec.Data[j] = a[(m-k+i)*lda+j]
			}
			vVec.Data[n-k+i] = 1
		}
		blas64.Ger(-tau[i], vVec, vVec, h)
		copy(qCopy.Data, q.Data)
		// Multiply q by the new h.
		switch kind {
		case "QR", "RQ":
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qCopy, h, 0, q)
		case "LQ":
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, h, qCopy, 0, q)
		}
	}
	return q
}

// checkBidiagonal checks the bidiagonal decomposition from dlabrd and dgebd2.
// The input to this function is the answer returned from the routines, stored
// in a, d, e, tauP, and tauQ. The data of original A matrix (before
// decomposition) is input in aCopy.
//
// checkBidiagonal constructs the V and U matrices, and from them constructs Q
// and P. Using these constructions, it checks that Qᵀ * A * P and checks that
// the result is bidiagonal.
func checkBidiagonal(t *testing.T, m, n, nb int, a []float64, lda int, d, e, tauP, tauQ, aCopy []float64) {
	// Check the answer.
	// Construct V and U.
	qMat := constructQPBidiagonal(lapack.ApplyQ, m, n, nb, a, lda, tauQ)
	pMat := constructQPBidiagonal(lapack.ApplyP, m, n, nb, a, lda, tauP)

	// Compute Qᵀ * A * P.
	aMat := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: lda,
		Data:   make([]float64, len(aCopy)),
	}
	copy(aMat.Data, aCopy)

	tmp1 := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, m*n),
	}
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, qMat, aMat, 0, tmp1)
	tmp2 := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, m*n),
	}
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp1, pMat, 0, tmp2)

	// Check that the first nb rows and cols of tm2 are upper bidiagonal
	// if m >= n, and lower bidiagonal otherwise.
	correctDiag := true
	matchD := true
	matchE := true
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if i >= nb && j >= nb {
				continue
			}
			v := tmp2.Data[i*tmp2.Stride+j]
			if i == j {
				if math.Abs(d[i]-v) > 1e-12 {
					matchD = false
				}
				continue
			}
			if m >= n && i == j-1 {
				if math.Abs(e[j-1]-v) > 1e-12 {
					matchE = false
				}
				continue
			}
			if m < n && i-1 == j {
				if math.Abs(e[i-1]-v) > 1e-12 {
					matchE = false
				}
				continue
			}
			if math.Abs(v) > 1e-12 {
				correctDiag = false
			}
		}
	}
	if !correctDiag {
		t.Errorf("Updated A not bi-diagonal")
	}
	if !matchD {
		fmt.Println("d = ", d)
		t.Errorf("D Mismatch")
	}
	if !matchE {
		t.Errorf("E mismatch")
	}
}

// constructQPBidiagonal constructs Q or P from the Bidiagonal decomposition
// computed by dlabrd and bgebd2.
func constructQPBidiagonal(vect lapack.ApplyOrtho, m, n, nb int, a []float64, lda int, tau []float64) blas64.General {
	sz := n
	if vect == lapack.ApplyQ {
		sz = m
	}

	var ldv int
	var v blas64.General
	if vect == lapack.ApplyQ {
		ldv = nb
		v = blas64.General{
			Rows:   m,
			Cols:   nb,
			Stride: ldv,
			Data:   make([]float64, m*ldv),
		}
	} else {
		ldv = n
		v = blas64.General{
			Rows:   nb,
			Cols:   n,
			Stride: ldv,
			Data:   make([]float64, m*ldv),
		}
	}

	if vect == lapack.ApplyQ {
		if m >= n {
			for i := 0; i < m; i++ {
				for j := 0; j <= min(nb-1, i); j++ {
					if i == j {
						v.Data[i*ldv+j] = 1
						continue
					}
					v.Data[i*ldv+j] = a[i*lda+j]
				}
			}
		} else {
			for i := 1; i < m; i++ {
				for j := 0; j <= min(nb-1, i-1); j++ {
					if i-1 == j {
						v.Data[i*ldv+j] = 1
						continue
					}
					v.Data[i*ldv+j] = a[i*lda+j]
				}
			}
		}
	} else {
		if m < n {
			for i := 0; i < nb; i++ {
				for j := i; j < n; j++ {
					if i == j {
						v.Data[i*ldv+j] = 1
						continue
					}
					v.Data[i*ldv+j] = a[i*lda+j]
				}
			}
		} else {
			for i := 0; i < nb; i++ {
				for j := i + 1; j < n; j++ {
					if j-1 == i {
						v.Data[i*ldv+j] = 1
						continue
					}
					v.Data[i*ldv+j] = a[i*lda+j]
				}
			}
		}
	}

	// The variable name is a computation of Q, but the algorithm is mostly the
	// same for computing P (just with different data).
	qMat := blas64.General{
		Rows:   sz,
		Cols:   sz,
		Stride: sz,
		Data:   make([]float64, sz*sz),
	}
	hMat := blas64.General{
		Rows:   sz,
		Cols:   sz,
		Stride: sz,
		Data:   make([]float64, sz*sz),
	}
	// set Q to I
	for i := 0; i < sz; i++ {
		qMat.Data[i*qMat.Stride+i] = 1
	}
	for i := 0; i < nb; i++ {
		qCopy := blas64.General{Rows: qMat.Rows, Cols: qMat.Cols, Stride: qMat.Stride, Data: make([]float64, len(qMat.Data))}
		copy(qCopy.Data, qMat.Data)

		// Set g and h to I
		for i := 0; i < sz; i++ {
			for j := 0; j < sz; j++ {
				if i == j {
					hMat.Data[i*sz+j] = 1
				} else {
					hMat.Data[i*sz+j] = 0
				}
			}
		}
		var vi blas64.Vector
		// H -= tauQ[i] * v[i] * v[i]^t
		if vect == lapack.ApplyQ {
			vi = blas64.Vector{
				Inc:  v.Stride,
				Data: v.Data[i:],
			}
		} else {
			vi = blas64.Vector{
				Inc:  1,
				Data: v.Data[i*v.Stride:],
			}
		}
		blas64.Ger(-tau[i], vi, vi, hMat)
		// Q = Q * G[1]
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qCopy, hMat, 0, qMat)
	}
	return qMat
}

// printRowise prints the matrix with one row per line. This is useful for debugging.
// If beyond is true, it prints beyond the final column to lda. If false, only
// the columns are printed.
//nolint:deadcode,unused
func printRowise(a []float64, m, n, lda int, beyond bool) {
	for i := 0; i < m; i++ {
		end := n
		if beyond {
			end = lda
		}
		fmt.Println(a[i*lda : i*lda+end])
	}
}

func copyGeneral(dst, src blas64.General) {
	r := min(dst.Rows, src.Rows)
	c := min(dst.Cols, src.Cols)
	for i := 0; i < r; i++ {
		copy(dst.Data[i*dst.Stride:i*dst.Stride+c], src.Data[i*src.Stride:i*src.Stride+c])
	}
}

// cloneGeneral allocates and returns an exact copy of the given general matrix.
func cloneGeneral(a blas64.General) blas64.General {
	c := a
	c.Data = make([]float64, len(a.Data))
	copy(c.Data, a.Data)
	return c
}

// equalGeneral returns whether the general matrices a and b are equal.
func equalGeneral(a, b blas64.General) bool {
	if a.Rows != b.Rows || a.Cols != b.Cols {
		panic("bad input")
	}
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			if a.Data[i*a.Stride+j] != b.Data[i*b.Stride+j] {
				return false
			}
		}
	}
	return true
}

// equalApproxGeneral returns whether the general matrices a and b are
// approximately equal within given tolerance.
func equalApproxGeneral(a, b blas64.General, tol float64) bool {
	if a.Rows != b.Rows || a.Cols != b.Cols {
		panic("bad input")
	}
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			diff := a.Data[i*a.Stride+j] - b.Data[i*b.Stride+j]
			if math.IsNaN(diff) || math.Abs(diff) > tol {
				return false
			}
		}
	}
	return true
}

// randSymBand returns an n×n random symmetric positive definite band matrix
// with kd diagonals.
func randSymBand(uplo blas.Uplo, n, kd, ldab int, rnd *rand.Rand) []float64 {
	// Allocate a triangular band matrix U or L and fill it with random numbers.
	var ab []float64
	if n > 0 {
		ab = make([]float64, (n-1)*ldab+kd+1)
	}
	for i := range ab {
		ab[i] = rnd.NormFloat64()
	}
	// Make sure that the matrix U or L has a sufficiently positive diagonal.
	switch uplo {
	case blas.Upper:
		for i := 0; i < n; i++ {
			ab[i*ldab] = float64(n) + rnd.Float64()
		}
	case blas.Lower:
		for i := 0; i < n; i++ {
			ab[i*ldab+kd] = float64(n) + rnd.Float64()
		}
	}
	// Compute Uᵀ*U or L*Lᵀ. The resulting (symmetric) matrix A will be
	// positive definite and well-conditioned.
	dsbmm(uplo, n, kd, ab, ldab)
	return ab
}

// distSymBand returns the max-norm distance between the symmetric band matrices
// A and B.
func distSymBand(uplo blas.Uplo, n, kd int, a []float64, lda int, b []float64, ldb int) float64 {
	var dist float64
	switch uplo {
	case blas.Upper:
		for i := 0; i < n; i++ {
			for j := 0; j < min(kd+1, n-i); j++ {
				dist = math.Max(dist, math.Abs(a[i*lda+j]-b[i*ldb+j]))
			}
		}
	case blas.Lower:
		for i := 0; i < n; i++ {
			for j := max(0, kd-i); j < kd+1; j++ {
				dist = math.Max(dist, math.Abs(a[i*lda+j]-b[i*ldb+j]))
			}
		}
	}
	return dist
}

// eye returns an identity matrix of given order and stride.
func eye(n, stride int) blas64.General {
	ans := nanGeneral(n, n, stride)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			ans.Data[i*ans.Stride+j] = 0
		}
		ans.Data[i*ans.Stride+i] = 1
	}
	return ans
}

// zeros returns an m×n matrix with given stride filled with zeros.
func zeros(m, n, stride int) blas64.General {
	a := nanGeneral(m, n, stride)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			a.Data[i*a.Stride+j] = 0
		}
	}
	return a
}

// extract2x2Block returns the elements of T at [0,0], [0,1], [1,0], and [1,1].
func extract2x2Block(t []float64, ldt int) (a, b, c, d float64) {
	return t[0], t[1], t[ldt], t[ldt+1]
}

// isSchurCanonical returns whether the 2×2 matrix [a b; c d] is in Schur
// canonical form.
func isSchurCanonical(a, b, c, d float64) bool {
	return c == 0 || (b != 0 && a == d && math.Signbit(b) != math.Signbit(c))
}

// isSchurCanonicalGeneral returns whether T is block upper triangular with 1×1
// and 2×2 diagonal blocks, each 2×2 block in Schur canonical form. The function
// checks only along the diagonal and the first subdiagonal, otherwise the lower
// triangle is not accessed.
func isSchurCanonicalGeneral(t blas64.General) bool {
	n := t.Cols
	if t.Rows != n {
		panic("invalid matrix")
	}
	for j := 0; j < n-1; {
		if t.Data[(j+1)*t.Stride+j] == 0 {
			// 1×1 block.
			for i := j + 1; i < n; i++ {
				if t.Data[i*t.Stride+j] != 0 {
					return false
				}
			}
			j++
			continue
		}
		// 2×2 block.
		a, b, c, d := extract2x2Block(t.Data[j*t.Stride+j:], t.Stride)
		if !isSchurCanonical(a, b, c, d) {
			return false
		}
		for i := j + 2; i < n; i++ {
			if t.Data[i*t.Stride+j] != 0 {
				return false
			}
		}
		for i := j + 2; i < n; i++ {
			if t.Data[i*t.Stride+j+1] != 0 {
				return false
			}
		}
		j += 2
	}
	return true
}

// schurBlockEigenvalues returns the two eigenvalues of the 2×2 matrix [a b; c d]
// that must be in Schur canonical form.
func schurBlockEigenvalues(a, b, c, d float64) (ev1, ev2 complex128) {
	if !isSchurCanonical(a, b, c, d) {
		panic("block not in Schur canonical form")
	}
	if c == 0 {
		return complex(a, 0), complex(d, 0)
	}
	im := math.Sqrt(math.Abs(b)) * math.Sqrt(math.Abs(c))
	return complex(a, im), complex(a, -im)
}

// schurBlockSize returns the size of the diagonal block at i-th row in the
// upper quasi-triangular matrix t in Schur canonical form, and whether i points
// to the first row of the block. For zero-sized matrices the function returns 0
// and true.
func schurBlockSize(t blas64.General, i int) (size int, first bool) {
	if t.Rows != t.Cols {
		panic("matrix not square")
	}
	if t.Rows == 0 {
		return 0, true
	}
	if i < 0 || t.Rows <= i {
		panic("index out of range")
	}

	first = true
	if i > 0 && t.Data[i*t.Stride+i-1] != 0 {
		// There is a non-zero element to the left, therefore i must
		// point to the second row in a 2×2 diagonal block.
		first = false
		i--
	}
	size = 1
	if i+1 < t.Rows && t.Data[(i+1)*t.Stride+i] != 0 {
		// There is a non-zero element below, this must be a 2×2
		// diagonal block.
		size = 2
	}
	return size, first
}

// containsComplex returns whether z is approximately equal to one of the complex
// numbers in v. If z is found, its index in v will be also returned.
func containsComplex(v []complex128, z complex128, tol float64) (found bool, index int) {
	for i := range v {
		if cmplx.Abs(v[i]-z) < tol {
			return true, i
		}
	}
	return false, -1
}

// isAllNaN returns whether x contains only NaN values.
func isAllNaN(x []float64) bool {
	for _, v := range x {
		if !math.IsNaN(v) {
			return false
		}
	}
	return true
}

// isUpperHessenberg returns whether h contains only zeros below the
// subdiagonal.
func isUpperHessenberg(h blas64.General) bool {
	if h.Rows != h.Cols {
		panic("matrix not square")
	}
	n := h.Rows
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i > j+1 && h.Data[i*h.Stride+j] != 0 {
				return false
			}
		}
	}
	return true
}

// isUpperTriangular returns whether a contains only zeros below the diagonal.
func isUpperTriangular(a blas64.General) bool {
	n := a.Rows
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if a.Data[i*a.Stride+j] != 0 {
				return false
			}
		}
	}
	return true
}

// unbalancedSparseGeneral returns an m×n dense matrix with a random sparse
// structure consisting of nz nonzero elements. The matrix will be unbalanced by
// multiplying each element randomly by its row or column index.
func unbalancedSparseGeneral(m, n, stride int, nonzeros int, rnd *rand.Rand) blas64.General {
	a := zeros(m, n, stride)
	for k := 0; k < nonzeros; k++ {
		i := rnd.Intn(n)
		j := rnd.Intn(n)
		if rnd.Float64() < 0.5 {
			a.Data[i*stride+j] = float64(i+1) * rnd.NormFloat64()
		} else {
			a.Data[i*stride+j] = float64(j+1) * rnd.NormFloat64()
		}
	}
	return a
}

// rootsOfUnity returns the n complex numbers whose n-th power is equal to 1.
func rootsOfUnity(n int) []complex128 {
	w := make([]complex128, n)
	for i := 0; i < n; i++ {
		angle := math.Pi * float64(2*i) / float64(n)
		w[i] = complex(math.Cos(angle), math.Sin(angle))
	}
	return w
}

// constructGSVDresults returns the matrices [ 0 R ], D1 and D2 described
// in the documentation of Dtgsja and Dggsvd3, and the result matrix in
// the documentation for Dggsvp3.
func constructGSVDresults(n, p, m, k, l int, a, b blas64.General, alpha, beta []float64) (zeroR, d1, d2 blas64.General) {
	// [ 0 R ]
	zeroR = zeros(k+l, n, n)
	dst := zeroR
	dst.Rows = min(m, k+l)
	dst.Cols = k + l
	dst.Data = zeroR.Data[n-k-l:]
	src := a
	src.Rows = min(m, k+l)
	src.Cols = k + l
	src.Data = a.Data[n-k-l:]
	copyGeneral(dst, src)
	if m < k+l {
		// [ 0 R ]
		dst.Rows = k + l - m
		dst.Cols = k + l - m
		dst.Data = zeroR.Data[m*zeroR.Stride+n-(k+l-m):]
		src = b
		src.Rows = k + l - m
		src.Cols = k + l - m
		src.Data = b.Data[(m-k)*b.Stride+n+m-k-l:]
		copyGeneral(dst, src)
	}

	// D1
	d1 = zeros(m, k+l, k+l)
	for i := 0; i < k; i++ {
		d1.Data[i*d1.Stride+i] = 1
	}
	for i := k; i < min(m, k+l); i++ {
		d1.Data[i*d1.Stride+i] = alpha[i]
	}

	// D2
	d2 = zeros(p, k+l, k+l)
	for i := 0; i < min(l, m-k); i++ {
		d2.Data[i*d2.Stride+i+k] = beta[k+i]
	}
	for i := m - k; i < l; i++ {
		d2.Data[i*d2.Stride+i+k] = 1
	}

	return zeroR, d1, d2
}

func constructGSVPresults(n, p, m, k, l int, a, b blas64.General) (zeroA, zeroB blas64.General) {
	zeroA = zeros(m, n, n)
	dst := zeroA
	dst.Rows = min(m, k+l)
	dst.Cols = k + l
	dst.Data = zeroA.Data[n-k-l:]
	src := a
	dst.Rows = min(m, k+l)
	src.Cols = k + l
	src.Data = a.Data[n-k-l:]
	copyGeneral(dst, src)

	zeroB = zeros(p, n, n)
	dst = zeroB
	dst.Rows = l
	dst.Cols = l
	dst.Data = zeroB.Data[n-l:]
	src = b
	dst.Rows = l
	src.Cols = l
	src.Data = b.Data[n-l:]
	copyGeneral(dst, src)

	return zeroA, zeroB
}

// distFromIdentity returns the L-infinity distance of an n×n matrix A from the
// identity. If A contains NaN elements, distFromIdentity will return +inf.
func distFromIdentity(n int, a []float64, lda int) float64 {
	var dist float64
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			aij := a[i*lda+j]
			if math.IsNaN(aij) {
				return math.Inf(1)
			}
			if i == j {
				dist = math.Max(dist, math.Abs(aij-1))
			} else {
				dist = math.Max(dist, math.Abs(aij))
			}
		}
	}
	return dist
}

func sameFloat64(a, b float64) bool {
	return a == b || math.IsNaN(a) && math.IsNaN(b)
}

// sameLowerTri returns whether n×n matrices A and B are same under the diagonal.
func sameLowerTri(n int, a []float64, lda int, b []float64, ldb int) bool {
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			aij := a[i*lda+j]
			bij := b[i*ldb+j]
			if !sameFloat64(aij, bij) {
				return false
			}
		}
	}
	return true
}

// sameUpperTri returns whether n×n matrices A and B are same above the diagonal.
func sameUpperTri(n int, a []float64, lda int, b []float64, ldb int) bool {
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			aij := a[i*lda+j]
			bij := b[i*ldb+j]
			if !sameFloat64(aij, bij) {
				return false
			}
		}
	}
	return true
}

// svdJobString returns a string representation of job.
func svdJobString(job lapack.SVDJob) string {
	switch job {
	case lapack.SVDAll:
		return "All"
	case lapack.SVDStore:
		return "Store"
	case lapack.SVDOverwrite:
		return "Overwrite"
	case lapack.SVDNone:
		return "None"
	}
	return "unknown SVD job"
}

// residualOrthogonal returns the residual
//  |I - Q * Qᵀ|  if m < n or (m == n and rowwise == true),
//  |I - Qᵀ * Q|  otherwise.
// It can be used to check that the matrix Q is orthogonal.
func residualOrthogonal(q blas64.General, rowwise bool) float64 {
	m, n := q.Rows, q.Cols
	if m == 0 || n == 0 {
		return 0
	}
	var transq blas.Transpose
	if m < n || (m == n && rowwise) {
		transq = blas.NoTrans
	} else {
		transq = blas.Trans
	}
	minmn := min(m, n)

	// Set work = I.
	work := blas64.Symmetric{
		Uplo:   blas.Upper,
		N:      minmn,
		Data:   make([]float64, minmn*minmn),
		Stride: minmn,
	}
	for i := 0; i < minmn; i++ {
		work.Data[i*work.Stride+i] = 1
	}

	// Compute
	//  work = work - Q * Qᵀ = I - Q * Qᵀ
	// or
	//  work = work - Qᵀ * Q = I - Qᵀ * Q
	blas64.Syrk(transq, -1, q, 1, work)
	return dlansy(lapack.MaxColumnSum, blas.Upper, work.N, work.Data, work.Stride)
}
