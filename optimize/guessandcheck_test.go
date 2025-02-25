// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"testing"

	"github.com/jingcheng-WU/gonum/mat"
	"github.com/jingcheng-WU/gonum/optimize/functions"
	"github.com/jingcheng-WU/gonum/stat/distmv"
)

func TestGuessAndCheck(t *testing.T) {
	t.Parallel()
	dim := 30
	problem := Problem{
		Func: functions.ExtendedRosenbrock{}.Func,
	}
	mu := make([]float64, dim)
	sigma := mat.NewSymDense(dim, nil)
	for i := 0; i < dim; i++ {
		sigma.SetSym(i, i, 1)
	}
	d, ok := distmv.NewNormal(mu, sigma, nil)
	if !ok {
		panic("bad test")
	}
	initX := make([]float64, dim)
	_, err := Minimize(problem, initX, nil, &GuessAndCheck{Rander: d})
	if err != nil {
		t.Errorf("unexpected error running Minimize with nil settings: %v", err)
	}

	settings := &Settings{}
	settings.Concurrent = 5
	settings.MajorIterations = 15
	_, err = Minimize(problem, initX, settings, &GuessAndCheck{Rander: d})
	if err != nil {
		t.Errorf("unexpected error running Minimize with settings %+v: %v", settings, err)
	}
}
