// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scalar_test

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/jingcheng-WU/gonum/floats/scalar"
	"github.com/jingcheng-WU/gonum/stat"
)

func ExampleParseWithNA() {
	// Calculate the mean of a list of numbers
	// ignoring missing values.
	const data = `6
missing
4
`

	var vals, weights []float64
	sc := bufio.NewScanner(strings.NewReader(data))
	for sc.Scan() {
		v, w, err := scalar.ParseWithNA(sc.Text(), "missing")
		if err != nil {
			log.Fatal(err)
		}
		vals = append(vals, v)
		weights = append(weights, w)
	}
	err := sc.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stat.Mean(vals, weights))

	// Output:
	// 5
}
