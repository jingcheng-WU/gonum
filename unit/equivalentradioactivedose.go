// Code generated by "go generate github.com/jingcheng-WU/gonum/unit”; DO NOT EDIT.

// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
	"unicode/utf8"
)

// EquivalentRadioactiveDose is a measure of equivalent dose of ionizing radiation in sieverts.
type EquivalentRadioactiveDose float64

const Sievert EquivalentRadioactiveDose = 1

// Unit converts the EquivalentRadioactiveDose to a *Unit.
func (a EquivalentRadioactiveDose) Unit() *Unit {
	return New(float64(a), Dimensions{
		LengthDim: 2,
		TimeDim:   -2,
	})
}

// EquivalentRadioactiveDose allows EquivalentRadioactiveDose to implement a EquivalentRadioactiveDoseer interface.
func (a EquivalentRadioactiveDose) EquivalentRadioactiveDose() EquivalentRadioactiveDose {
	return a
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension.
func (a *EquivalentRadioactiveDose) From(u Uniter) error {
	if !DimensionsMatch(u, Sievert) {
		*a = EquivalentRadioactiveDose(math.NaN())
		return errors.New("unit: dimension mismatch")
	}
	*a = EquivalentRadioactiveDose(u.Unit().Value())
	return nil
}

func (a EquivalentRadioactiveDose) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", a, float64(a))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		const unit = " Sy"
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), p, float64(a))
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, float64(a))
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), float64(a))
		default:
			fmt.Fprintf(fs, "%"+string(c), float64(a))
		}
		fmt.Fprint(fs, unit)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g Sy)", c, a, float64(a))
	}
}
