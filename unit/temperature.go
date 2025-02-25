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

// Temperature represents a temperature in Kelvin.
type Temperature float64

const Kelvin Temperature = 1

// Unit converts the Temperature to a *Unit.
func (t Temperature) Unit() *Unit {
	return New(float64(t), Dimensions{
		TemperatureDim: 1,
	})
}

// Temperature allows Temperature to implement a Temperaturer interface.
func (t Temperature) Temperature() Temperature {
	return t
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension.
func (t *Temperature) From(u Uniter) error {
	if !DimensionsMatch(u, Kelvin) {
		*t = Temperature(math.NaN())
		return errors.New("unit: dimension mismatch")
	}
	*t = Temperature(u.Unit().Value())
	return nil
}

func (t Temperature) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", t, float64(t))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		const unit = " K"
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), p, float64(t))
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, float64(t))
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), pos(w-utf8.RuneCount([]byte(unit))), float64(t))
		default:
			fmt.Fprintf(fs, "%"+string(c), float64(t))
		}
		fmt.Fprint(fs, unit)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g K)", c, t, float64(t))
	}
}
