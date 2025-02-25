// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"github.com/jingcheng-WU/gonum/floats/scalar"
)

type Dlapy2er interface {
	Dlapy2(float64, float64) float64
}

func Dlapy2Test(t *testing.T, impl Dlapy2er) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 10; i++ {
		x := math.Abs(1e200 * rnd.NormFloat64())
		y := math.Abs(1e200 * rnd.NormFloat64())
		got := impl.Dlapy2(x, y)
		want := math.Hypot(x, y)
		if !scalar.EqualWithinRel(got, want, 1e-16) {
			t.Errorf("Dlapy2(%g, %g) = %g, want %g", x, y, got, want)
		}
	}
}
