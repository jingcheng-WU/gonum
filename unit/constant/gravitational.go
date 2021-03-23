// Code generated by "go generate github.com/jingcheng-WU/gonum/unit/constant”; DO NOT EDIT.

// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constant

import (
	"fmt"

	"github.com/jingcheng-WU/gonum/unit"
)

// Gravitational is the universal gravitational constant (G), the proportionality constant connecting the gravitational force between two bodies.
// The dimensions of Gravitational are m^3 kg^-1 s^-2. The standard uncertainty of the constant is 1.5e-15 m^3 kg^-1 s^-2.
const Gravitational = gravitationalUnits(6.6743e-11)

type gravitationalUnits float64

// Unit converts the gravitationalUnits to a *unit.Unit
func (cnst gravitationalUnits) Unit() *unit.Unit {
	return unit.New(float64(cnst), unit.Dimensions{
		unit.MassDim:   -1,
		unit.LengthDim: 3,
		unit.TimeDim:   -2,
	})
}

func (cnst gravitationalUnits) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", cnst, float64(cnst))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), w, p, cnst.Unit())
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, cnst.Unit())
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), w, cnst.Unit())
		default:
			fmt.Fprintf(fs, "%"+string(c), cnst.Unit())
		}
	default:
		fmt.Fprintf(fs, "%%!"+string(c)+"(constant.gravitationalUnits=%v m^3 kg^-1 s^-2)", float64(cnst))
	}
}
