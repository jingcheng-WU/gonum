// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"reflect"
	"testing"

	"golang.org/x/exp/rand"

	"github.com/jingcheng-WU/gonum/blas/blas64"
)

func TestNewVecDense(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		n      int
		data   []float64
		vector *VecDense
	}{
		{
			n:    3,
			data: []float64{4, 5, 6},
			vector: &VecDense{
				mat: blas64.Vector{
					N:    3,
					Data: []float64{4, 5, 6},
					Inc:  1,
				},
			},
		},
		{
			n:    3,
			data: nil,
			vector: &VecDense{
				mat: blas64.Vector{
					N:    3,
					Data: []float64{0, 0, 0},
					Inc:  1,
				},
			},
		},
	} {
		v := NewVecDense(test.n, test.data)
		rows, cols := v.Dims()
		if rows != test.n {
			t.Errorf("unexpected number of rows for test %d: got: %d want: %d", i, rows, test.n)
		}
		if cols != 1 {
			t.Errorf("unexpected number of cols for test %d: got: %d want: 1", i, cols)
		}
		if !reflect.DeepEqual(v, test.vector) {
			t.Errorf("unexpected data slice for test %d: got: %v want: %v", i, v, test.vector)
		}
	}
}

func TestCap(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		vector *VecDense
		want   int
	}{
		{vector: NewVecDense(3, nil), want: 3},
		{
			vector: &VecDense{
				mat: blas64.Vector{
					N:    3,
					Data: make([]float64, 7, 10),
					Inc:  3,
				},
			},
			want: 4,
		},
		{
			vector: &VecDense{
				mat: blas64.Vector{
					N:    4,
					Data: make([]float64, 10),
					Inc:  3,
				},
			},
			want: 4,
		},
		{
			vector: &VecDense{
				mat: blas64.Vector{
					N:    4,
					Data: make([]float64, 11),
					Inc:  3,
				},
			},
			want: 4,
		},
		{
			vector: &VecDense{
				mat: blas64.Vector{
					N:    4,
					Data: make([]float64, 12),
					Inc:  3,
				},
			},
			want: 4,
		},
		{
			vector: &VecDense{
				mat: blas64.Vector{
					N:    4,
					Data: make([]float64, 13),
					Inc:  3,
				},
			},
			want: 5,
		},
	} {
		got := test.vector.Cap()
		if got != test.want {
			t.Errorf("unexpected capacty for test %d: got: %d want: %d", i, got, test.want)
		}
	}
}

func TestVecDenseAtSet(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		vector *VecDense
	}{
		{
			vector: &VecDense{
				mat: blas64.Vector{
					N:    3,
					Data: []float64{0, 1, 2},
					Inc:  1,
				},
			},
		},
		{
			vector: &VecDense{
				mat: blas64.Vector{
					N:    3,
					Data: []float64{0, 10, 10, 1, 10, 10, 2},
					Inc:  3,
				},
			},
		},
	} {
		v := test.vector
		n := test.vector.mat.N

		for _, row := range []int{-1, n} {
			panicked, message := panics(func() { v.At(row, 0) })
			if !panicked || message != ErrRowAccess.Error() {
				t.Errorf("expected panic for invalid row access for test %d n=%d r=%d", i, n, row)
			}
		}
		for _, col := range []int{-1, 1} {
			panicked, message := panics(func() { v.At(0, col) })
			if !panicked || message != ErrColAccess.Error() {
				t.Errorf("expected panic for invalid column access for test %d n=%d c=%d", i, n, col)
			}
		}

		for _, row := range []int{0, 1, n - 1} {
			if e := v.At(row, 0); e != float64(row) {
				t.Errorf("unexpected value for At(%d, 0) for test %d : got: %v want: %v", row, i, e, float64(row))
			}
		}

		for _, row := range []int{-1, n} {
			panicked, message := panics(func() { v.SetVec(row, 100) })
			if !panicked || message != ErrVectorAccess.Error() {
				t.Errorf("expected panic for invalid row access for test %d n=%d r=%d", i, n, row)
			}
		}

		for inc, row := range []int{0, 2} {
			v.SetVec(row, 100+float64(inc))
			if e := v.At(row, 0); e != 100+float64(inc) {
				t.Errorf("unexpected value for At(%d, 0) after SetVec(%[1]d, %v) for test %d: got: %v want: %[2]v", row, 100+float64(inc), i, e)
			}
		}
	}
}

func TestVecDenseZero(t *testing.T) {
	t.Parallel()
	// Elements that equal 1 should be set to zero, elements that equal -1
	// should remain unchanged.
	for _, test := range []*VecDense{
		{
			mat: blas64.Vector{
				N:   5,
				Inc: 2,
				Data: []float64{
					1, -1,
					1, -1,
					1, -1,
					1, -1,
					1,
				},
			},
		},
	} {
		dataCopy := make([]float64, len(test.mat.Data))
		copy(dataCopy, test.mat.Data)
		test.Zero()
		for i, v := range test.mat.Data {
			if dataCopy[i] != -1 && v != 0 {
				t.Errorf("Matrix not zeroed in bounds")
			}
			if dataCopy[i] == -1 && v != -1 {
				t.Errorf("Matrix zeroed out of bounds")
			}
		}
	}
}

func TestVecDenseMul(t *testing.T) {
	t.Parallel()
	method := func(receiver, a, b Matrix) {
		type mulVecer interface {
			MulVec(a Matrix, b Vector)
		}
		rd := receiver.(mulVecer)
		rd.MulVec(a, b.(Vector))
	}
	denseComparison := func(receiver, a, b *Dense) {
		receiver.Mul(a, b)
	}
	legalSizeMulVec := func(ar, ac, br, bc int) bool {
		var legal bool
		if bc != 1 {
			legal = false
		} else {
			legal = ac == br
		}
		return legal
	}
	testTwoInput(t, "MulVec", &VecDense{}, method, denseComparison, legalTypesMatrixVector, legalSizeMulVec, 1e-14)
}

func TestVecDenseScale(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		a     Vector
		alpha float64
		want  *VecDense
	}{
		{
			a:     NewVecDense(3, []float64{0, 1, 2}),
			alpha: 0,
			want:  NewVecDense(3, []float64{0, 0, 0}),
		},
		{
			a:     NewVecDense(3, []float64{0, 1, 2}),
			alpha: 1,
			want:  NewVecDense(3, []float64{0, 1, 2}),
		},
		{
			a:     NewVecDense(3, []float64{0, 1, 2}),
			alpha: -2,
			want:  NewVecDense(3, []float64{0, -2, -4}),
		},
		{
			a:     NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			alpha: 0,
			want:  NewVecDense(3, []float64{0, 0, 0}),
		},
		{
			a:     NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			alpha: 1,
			want:  NewVecDense(3, []float64{0, 1, 2}),
		},
		{
			a:     NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			alpha: -2,
			want:  NewVecDense(3, []float64{0, -2, -4}),
		},
		{
			a: NewDense(3, 3, []float64{
				0, 1, 2,
				3, 4, 5,
				6, 7, 8,
			}).ColView(1),
			alpha: -2,
			want:  NewVecDense(3, []float64{-2, -8, -14}),
		},
	} {
		var v VecDense
		v.ScaleVec(test.alpha, test.a.(*VecDense))
		if !reflect.DeepEqual(v.RawVector(), test.want.RawVector()) {
			t.Errorf("test %d: unexpected result for v = alpha * a: got: %v want: %v", i, v.RawVector(), test.want.RawVector())
		}

		v.CopyVec(test.a.(*VecDense))
		v.ScaleVec(test.alpha, &v)
		if !reflect.DeepEqual(v.RawVector(), test.want.RawVector()) {
			t.Errorf("test %d: unexpected result for v = alpha * v: got: %v want: %v", i, v.RawVector(), test.want.RawVector())
		}
	}

	for _, alpha := range []float64{0, 1, -1, 2.3, -2.3} {
		method := func(receiver, a Matrix) {
			type scaleVecer interface {
				ScaleVec(float64, Vector)
			}
			v := receiver.(scaleVecer)
			v.ScaleVec(alpha, a.(Vector))
		}
		denseComparison := func(receiver, a *Dense) {
			receiver.Scale(alpha, a)
		}
		testOneInput(t, "ScaleVec", &VecDense{}, method, denseComparison, legalTypeVector, isAnyColumnVector, 0)
	}
}

func TestCopyVec(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		src   *VecDense
		dst   *VecDense
		want  *VecDense
		wantN int
	}{
		{src: NewVecDense(1, nil), dst: NewVecDense(1, nil), want: NewVecDense(1, nil), wantN: 1},
		{src: NewVecDense(3, []float64{1, 2, 3}), dst: NewVecDense(2, []float64{-1, -2}), want: NewVecDense(2, []float64{1, 2}), wantN: 2},
		{src: NewVecDense(2, []float64{1, 2}), dst: NewVecDense(3, []float64{-1, -2, -3}), want: NewVecDense(3, []float64{1, 2, -3}), wantN: 2},
	} {
		got := test.dst
		var n int
		panicked, message := panics(func() { n = got.CopyVec(test.src) })
		if panicked {
			t.Errorf("unexpected panic during vector copy for test %d: %s", i, message)
		}
		if !Equal(got, test.want) {
			t.Errorf("test %d: unexpected result CopyVec:\ngot: %v\nwant:%v", i, got, test.want)
		}
		if n != test.wantN {
			t.Errorf("test %d: unexpected result number of elements copied: got:%d want:%d", i, n, test.wantN)
		}
	}
}

func TestVecDenseAddScaled(t *testing.T) {
	t.Parallel()
	for _, alpha := range []float64{0, 1, -1, 2.3, -2.3} {
		method := func(receiver, a, b Matrix) {
			type addScaledVecer interface {
				AddScaledVec(Vector, float64, Vector)
			}
			v := receiver.(addScaledVecer)
			v.AddScaledVec(a.(Vector), alpha, b.(Vector))
		}
		denseComparison := func(receiver, a, b *Dense) {
			var sb Dense
			sb.Scale(alpha, b)
			receiver.Add(a, &sb)
		}
		testTwoInput(t, "AddScaledVec", &VecDense{}, method, denseComparison, legalTypesVectorVector, legalSizeSameVec, 1e-14)
	}
}

func TestVecDenseAdd(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		a, b Vector
		want *VecDense
	}{
		{
			a:    NewVecDense(3, []float64{0, 1, 2}),
			b:    NewVecDense(3, []float64{0, 2, 3}),
			want: NewVecDense(3, []float64{0, 3, 5}),
		},
		{
			a:    NewVecDense(3, []float64{0, 1, 2}),
			b:    NewDense(3, 1, []float64{0, 2, 3}).ColView(0),
			want: NewVecDense(3, []float64{0, 3, 5}),
		},
		{
			a:    NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			b:    NewDense(3, 1, []float64{0, 2, 3}).ColView(0),
			want: NewVecDense(3, []float64{0, 3, 5}),
		},
	} {
		var v VecDense
		v.AddVec(test.a.(*VecDense), test.b.(*VecDense))
		if !reflect.DeepEqual(v.RawVector(), test.want.RawVector()) {
			t.Errorf("unexpected result for test %d: got: %v want: %v", i, v.RawVector(), test.want.RawVector())
		}
	}
}

func TestVecDenseSub(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		a, b Vector
		want *VecDense
	}{
		{
			a:    NewVecDense(3, []float64{0, 1, 2}),
			b:    NewVecDense(3, []float64{0, 0.5, 1}),
			want: NewVecDense(3, []float64{0, 0.5, 1}),
		},
		{
			a:    NewVecDense(3, []float64{0, 1, 2}),
			b:    NewDense(3, 1, []float64{0, 0.5, 1}).ColView(0),
			want: NewVecDense(3, []float64{0, 0.5, 1}),
		},
		{
			a:    NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			b:    NewDense(3, 1, []float64{0, 0.5, 1}).ColView(0),
			want: NewVecDense(3, []float64{0, 0.5, 1}),
		},
	} {
		var v VecDense
		v.SubVec(test.a.(*VecDense), test.b.(*VecDense))
		if !reflect.DeepEqual(v.RawVector(), test.want.RawVector()) {
			t.Errorf("unexpected result for test %d: got: %v want: %v", i, v.RawVector(), test.want.RawVector())
		}
	}
}

func TestVecDenseMulElem(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		a, b Vector
		want *VecDense
	}{
		{
			a:    NewVecDense(3, []float64{0, 1, 2}),
			b:    NewVecDense(3, []float64{0, 2, 3}),
			want: NewVecDense(3, []float64{0, 2, 6}),
		},
		{
			a:    NewVecDense(3, []float64{0, 1, 2}),
			b:    NewDense(3, 1, []float64{0, 2, 3}).ColView(0),
			want: NewVecDense(3, []float64{0, 2, 6}),
		},
		{
			a:    NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			b:    NewDense(3, 1, []float64{0, 2, 3}).ColView(0),
			want: NewVecDense(3, []float64{0, 2, 6}),
		},
	} {
		var v VecDense
		v.MulElemVec(test.a.(*VecDense), test.b.(*VecDense))
		if !reflect.DeepEqual(v.RawVector(), test.want.RawVector()) {
			t.Errorf("unexpected result for test %d: got: %v want: %v", i, v.RawVector(), test.want.RawVector())
		}
	}
}

func TestVecDenseDivElem(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		a, b Vector
		want *VecDense
	}{
		{
			a:    NewVecDense(3, []float64{0.5, 1, 2}),
			b:    NewVecDense(3, []float64{0.5, 0.5, 1}),
			want: NewVecDense(3, []float64{1, 2, 2}),
		},
		{
			a:    NewVecDense(3, []float64{0.5, 1, 2}),
			b:    NewDense(3, 1, []float64{0.5, 0.5, 1}).ColView(0),
			want: NewVecDense(3, []float64{1, 2, 2}),
		},
		{
			a:    NewDense(3, 1, []float64{0.5, 1, 2}).ColView(0),
			b:    NewDense(3, 1, []float64{0.5, 0.5, 1}).ColView(0),
			want: NewVecDense(3, []float64{1, 2, 2}),
		},
	} {
		var v VecDense
		v.DivElemVec(test.a.(*VecDense), test.b.(*VecDense))
		if !reflect.DeepEqual(v.RawVector(), test.want.RawVector()) {
			t.Errorf("unexpected result for test %d: got: %v want: %v", i, v.RawVector(), test.want.RawVector())
		}
	}
}

func BenchmarkAddScaledVec10Inc1(b *testing.B)      { addScaledVecBench(b, 10, 1) }
func BenchmarkAddScaledVec100Inc1(b *testing.B)     { addScaledVecBench(b, 100, 1) }
func BenchmarkAddScaledVec1000Inc1(b *testing.B)    { addScaledVecBench(b, 1000, 1) }
func BenchmarkAddScaledVec10000Inc1(b *testing.B)   { addScaledVecBench(b, 10000, 1) }
func BenchmarkAddScaledVec100000Inc1(b *testing.B)  { addScaledVecBench(b, 100000, 1) }
func BenchmarkAddScaledVec10Inc2(b *testing.B)      { addScaledVecBench(b, 10, 2) }
func BenchmarkAddScaledVec100Inc2(b *testing.B)     { addScaledVecBench(b, 100, 2) }
func BenchmarkAddScaledVec1000Inc2(b *testing.B)    { addScaledVecBench(b, 1000, 2) }
func BenchmarkAddScaledVec10000Inc2(b *testing.B)   { addScaledVecBench(b, 10000, 2) }
func BenchmarkAddScaledVec100000Inc2(b *testing.B)  { addScaledVecBench(b, 100000, 2) }
func BenchmarkAddScaledVec10Inc20(b *testing.B)     { addScaledVecBench(b, 10, 20) }
func BenchmarkAddScaledVec100Inc20(b *testing.B)    { addScaledVecBench(b, 100, 20) }
func BenchmarkAddScaledVec1000Inc20(b *testing.B)   { addScaledVecBench(b, 1000, 20) }
func BenchmarkAddScaledVec10000Inc20(b *testing.B)  { addScaledVecBench(b, 10000, 20) }
func BenchmarkAddScaledVec100000Inc20(b *testing.B) { addScaledVecBench(b, 100000, 20) }
func addScaledVecBench(b *testing.B, size, inc int) {
	src := rand.NewSource(1)
	x := randVecDense(size, inc, 1, src)
	y := randVecDense(size, inc, 1, src)
	b.ResetTimer()
	var v VecDense
	for i := 0; i < b.N; i++ {
		v.AddScaledVec(y, 2, x)
	}
}

func BenchmarkScaleVec10Inc1(b *testing.B)      { scaleVecBench(b, 10, 1) }
func BenchmarkScaleVec100Inc1(b *testing.B)     { scaleVecBench(b, 100, 1) }
func BenchmarkScaleVec1000Inc1(b *testing.B)    { scaleVecBench(b, 1000, 1) }
func BenchmarkScaleVec10000Inc1(b *testing.B)   { scaleVecBench(b, 10000, 1) }
func BenchmarkScaleVec100000Inc1(b *testing.B)  { scaleVecBench(b, 100000, 1) }
func BenchmarkScaleVec10Inc2(b *testing.B)      { scaleVecBench(b, 10, 2) }
func BenchmarkScaleVec100Inc2(b *testing.B)     { scaleVecBench(b, 100, 2) }
func BenchmarkScaleVec1000Inc2(b *testing.B)    { scaleVecBench(b, 1000, 2) }
func BenchmarkScaleVec10000Inc2(b *testing.B)   { scaleVecBench(b, 10000, 2) }
func BenchmarkScaleVec100000Inc2(b *testing.B)  { scaleVecBench(b, 100000, 2) }
func BenchmarkScaleVec10Inc20(b *testing.B)     { scaleVecBench(b, 10, 20) }
func BenchmarkScaleVec100Inc20(b *testing.B)    { scaleVecBench(b, 100, 20) }
func BenchmarkScaleVec1000Inc20(b *testing.B)   { scaleVecBench(b, 1000, 20) }
func BenchmarkScaleVec10000Inc20(b *testing.B)  { scaleVecBench(b, 10000, 20) }
func BenchmarkScaleVec100000Inc20(b *testing.B) { scaleVecBench(b, 100000, 20) }
func scaleVecBench(b *testing.B, size, inc int) {
	src := rand.NewSource(1)
	x := randVecDense(size, inc, 1, src)
	b.ResetTimer()
	var v VecDense
	for i := 0; i < b.N; i++ {
		v.ScaleVec(2, x)
	}
}

func BenchmarkAddVec10Inc1(b *testing.B)      { addVecBench(b, 10, 1) }
func BenchmarkAddVec100Inc1(b *testing.B)     { addVecBench(b, 100, 1) }
func BenchmarkAddVec1000Inc1(b *testing.B)    { addVecBench(b, 1000, 1) }
func BenchmarkAddVec10000Inc1(b *testing.B)   { addVecBench(b, 10000, 1) }
func BenchmarkAddVec100000Inc1(b *testing.B)  { addVecBench(b, 100000, 1) }
func BenchmarkAddVec10Inc2(b *testing.B)      { addVecBench(b, 10, 2) }
func BenchmarkAddVec100Inc2(b *testing.B)     { addVecBench(b, 100, 2) }
func BenchmarkAddVec1000Inc2(b *testing.B)    { addVecBench(b, 1000, 2) }
func BenchmarkAddVec10000Inc2(b *testing.B)   { addVecBench(b, 10000, 2) }
func BenchmarkAddVec100000Inc2(b *testing.B)  { addVecBench(b, 100000, 2) }
func BenchmarkAddVec10Inc20(b *testing.B)     { addVecBench(b, 10, 20) }
func BenchmarkAddVec100Inc20(b *testing.B)    { addVecBench(b, 100, 20) }
func BenchmarkAddVec1000Inc20(b *testing.B)   { addVecBench(b, 1000, 20) }
func BenchmarkAddVec10000Inc20(b *testing.B)  { addVecBench(b, 10000, 20) }
func BenchmarkAddVec100000Inc20(b *testing.B) { addVecBench(b, 100000, 20) }
func addVecBench(b *testing.B, size, inc int) {
	src := rand.NewSource(1)
	x := randVecDense(size, inc, 1, src)
	y := randVecDense(size, inc, 1, src)
	b.ResetTimer()
	var v VecDense
	for i := 0; i < b.N; i++ {
		v.AddVec(x, y)
	}
}

func BenchmarkSubVec10Inc1(b *testing.B)      { subVecBench(b, 10, 1) }
func BenchmarkSubVec100Inc1(b *testing.B)     { subVecBench(b, 100, 1) }
func BenchmarkSubVec1000Inc1(b *testing.B)    { subVecBench(b, 1000, 1) }
func BenchmarkSubVec10000Inc1(b *testing.B)   { subVecBench(b, 10000, 1) }
func BenchmarkSubVec100000Inc1(b *testing.B)  { subVecBench(b, 100000, 1) }
func BenchmarkSubVec10Inc2(b *testing.B)      { subVecBench(b, 10, 2) }
func BenchmarkSubVec100Inc2(b *testing.B)     { subVecBench(b, 100, 2) }
func BenchmarkSubVec1000Inc2(b *testing.B)    { subVecBench(b, 1000, 2) }
func BenchmarkSubVec10000Inc2(b *testing.B)   { subVecBench(b, 10000, 2) }
func BenchmarkSubVec100000Inc2(b *testing.B)  { subVecBench(b, 100000, 2) }
func BenchmarkSubVec10Inc20(b *testing.B)     { subVecBench(b, 10, 20) }
func BenchmarkSubVec100Inc20(b *testing.B)    { subVecBench(b, 100, 20) }
func BenchmarkSubVec1000Inc20(b *testing.B)   { subVecBench(b, 1000, 20) }
func BenchmarkSubVec10000Inc20(b *testing.B)  { subVecBench(b, 10000, 20) }
func BenchmarkSubVec100000Inc20(b *testing.B) { subVecBench(b, 100000, 20) }
func subVecBench(b *testing.B, size, inc int) {
	src := rand.NewSource(1)
	x := randVecDense(size, inc, 1, src)
	y := randVecDense(size, inc, 1, src)
	b.ResetTimer()
	var v VecDense
	for i := 0; i < b.N; i++ {
		v.SubVec(x, y)
	}
}

func randVecDense(size, inc int, rho float64, src rand.Source) *VecDense {
	if size <= 0 {
		panic("bad vector size")
	}
	rnd := rand.New(src)
	data := make([]float64, size*inc)
	for i := range data {
		if rnd.Float64() < rho {
			data[i] = rnd.NormFloat64()
		}
	}
	return &VecDense{
		mat: blas64.Vector{
			N:    size,
			Inc:  inc,
			Data: data,
		},
	}
}

func BenchmarkVectorSum100000(b *testing.B) { vectorSumBench(b, 100000) }

var vectorSumForBench float64

func vectorSumBench(b *testing.B, size int) {
	src := rand.NewSource(1)
	a := randVecDense(size, 1, 1.0, src)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vectorSumForBench = Sum(a)
	}
}
