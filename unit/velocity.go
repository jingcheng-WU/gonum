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

// Velocity represents a velocity in metres per second.
type Velocity float64

// Unit converts the Velocity to a *Unit.
func (v Velocity) Unit() *Unit {
	return New(float64(v), Dimensions{
		LengthDim: 1,
		TimeDim:   -1,
	})
}

// Velocity allows Velocity to implement a Velocityer interface.
func (v Velocity) Velocity() Velocity {
	return v
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension.
func (v *Velocity) From(u Uniter) error {
	if !DimensionsMatch(u, Velocity(0)) {
		*v = Velocity(math.NaN())
		return errors.New("unit: dimension mismatch")
	}
	*v = Velocity(u.Unit().Value())
	return nil
}

func (v Velocity) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", v, float64(v))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		const unit = " m s^-1"
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), p, float64(v))
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, float64(v))
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), float64(v))
		default:
			fmt.Fprintf(fs, "%"+string(c), float64(v))
		}
		fmt.Fprint(fs, unit)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g m s^-1)", c, v, float64(v))
	}
}
