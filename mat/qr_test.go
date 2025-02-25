// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"github.com/jingcheng-WU/gonum/blas/blas64"
)

func TestQR(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n int
	}{
		{5, 5},
		{10, 5},
	} {
		m := test.m
		n := test.n
		a := NewDense(m, n, nil)
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				a.Set(i, j, rnd.NormFloat64())
			}
		}
		var want Dense
		want.CloneFrom(a)

		var qr QR
		qr.Factorize(a)
		var q, r Dense
		qr.QTo(&q)

		if !isOrthonormal(&q, 1e-10) {
			t.Errorf("Q is not orthonormal: m = %v, n = %v", m, n)
		}

		qr.RTo(&r)

		var got Dense
		got.Mul(&q, &r)
		if !EqualApprox(&got, &want, 1e-12) {
			t.Errorf("QR does not equal original matrix. \nWant: %v\nGot: %v", want, got)
		}
	}
}

func isOrthonormal(q *Dense, tol float64) bool {
	m, n := q.Dims()
	if m != n {
		return false
	}
	for i := 0; i < m; i++ {
		for j := i; j < m; j++ {
			dot := blas64.Dot(blas64.Vector{N: m, Inc: 1, Data: q.mat.Data[i*q.mat.Stride:]},
				blas64.Vector{N: m, Inc: 1, Data: q.mat.Data[j*q.mat.Stride:]})
			// Dot product should be 1 if i == j and 0 otherwise.
			if i == j && math.Abs(dot-1) > tol {
				return false
			}
			if i != j && math.Abs(dot) > tol {
				return false
			}
		}
	}
	return true
}

func TestQRSolveTo(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewSource(1))
	for _, trans := range []bool{false, true} {
		for _, test := range []struct {
			m, n, bc int
		}{
			{5, 5, 1},
			{10, 5, 1},
			{5, 5, 3},
			{10, 5, 3},
		} {
			m := test.m
			n := test.n
			bc := test.bc
			a := NewDense(m, n, nil)
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					a.Set(i, j, rnd.Float64())
				}
			}
			br := m
			if trans {
				br = n
			}
			b := NewDense(br, bc, nil)
			for i := 0; i < br; i++ {
				for j := 0; j < bc; j++ {
					b.Set(i, j, rnd.Float64())
				}
			}
			var x Dense
			var qr QR
			qr.Factorize(a)
			err := qr.SolveTo(&x, trans, b)
			if err != nil {
				t.Errorf("unexpected error from QR solve: %v", err)
			}

			// Test that the normal equations hold.
			// Aᵀ * A * x = Aᵀ * b if !trans
			// A * Aᵀ * x = A * b if trans
			var lhs Dense
			var rhs Dense
			if trans {
				var tmp Dense
				tmp.Mul(a, a.T())
				lhs.Mul(&tmp, &x)
				rhs.Mul(a, b)
			} else {
				var tmp Dense
				tmp.Mul(a.T(), a)
				lhs.Mul(&tmp, &x)
				rhs.Mul(a.T(), b)
			}
			if !EqualApprox(&lhs, &rhs, 1e-10) {
				t.Errorf("Normal equations do not hold.\nLHS: %v\n, RHS: %v\n", lhs, rhs)
			}
		}
	}
	// TODO(btracey): Add in testOneInput when it exists.
}

func TestQRSolveVecTo(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewSource(1))
	for _, trans := range []bool{false, true} {
		for _, test := range []struct {
			m, n int
		}{
			{5, 5},
			{10, 5},
		} {
			m := test.m
			n := test.n
			a := NewDense(m, n, nil)
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					a.Set(i, j, rnd.Float64())
				}
			}
			br := m
			if trans {
				br = n
			}
			b := NewVecDense(br, nil)
			for i := 0; i < br; i++ {
				b.SetVec(i, rnd.Float64())
			}
			var x VecDense
			var qr QR
			qr.Factorize(a)
			err := qr.SolveVecTo(&x, trans, b)
			if err != nil {
				t.Errorf("unexpected error from QR solve: %v", err)
			}

			// Test that the normal equations hold.
			// Aᵀ * A * x = Aᵀ * b if !trans
			// A * Aᵀ * x = A * b if trans
			var lhs Dense
			var rhs Dense
			if trans {
				var tmp Dense
				tmp.Mul(a, a.T())
				lhs.Mul(&tmp, &x)
				rhs.Mul(a, b)
			} else {
				var tmp Dense
				tmp.Mul(a.T(), a)
				lhs.Mul(&tmp, &x)
				rhs.Mul(a.T(), b)
			}
			if !EqualApprox(&lhs, &rhs, 1e-10) {
				t.Errorf("Normal equations do not hold.\nLHS: %v\n, RHS: %v\n", lhs, rhs)
			}
		}
	}
	// TODO(btracey): Add in testOneInput when it exists.
}

func TestQRSolveCondTo(t *testing.T) {
	t.Parallel()
	for _, test := range []*Dense{
		NewDense(2, 2, []float64{1, 0, 0, 1e-20}),
		NewDense(3, 2, []float64{1, 0, 0, 1e-20, 0, 0}),
	} {
		m, _ := test.Dims()
		var qr QR
		qr.Factorize(test)
		b := NewDense(m, 2, nil)
		var x Dense
		if err := qr.SolveTo(&x, false, b); err == nil {
			t.Error("No error for near-singular matrix in matrix solve.")
		}

		bvec := NewVecDense(m, nil)
		var xvec VecDense
		if err := qr.SolveVecTo(&xvec, false, bvec); err == nil {
			t.Error("No error for near-singular matrix in matrix solve.")
		}
	}
}
