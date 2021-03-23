// Code generated by "go generate github.com/jingcheng-WU/gonum/unit/constant”; DO NOT EDIT.

// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constant

import (
	"fmt"
	"testing"
)

func TestPlanckFormat(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		format string
		want   string
	}{
		{"%v", "6.62607015e-34 kg m^2 s^-1"},
		{"%.1v", "7e-34 kg m^2 s^-1"},
		{"%50.1v", "                                 7e-34 kg m^2 s^-1"},
		{"%50v", "                        6.62607015e-34 kg m^2 s^-1"},
		{"%1v", "6.62607015e-34 kg m^2 s^-1"},
		{"%#v", "constant.planckUnits(6.62607015e-34)"},
		{"%s", "%!s(constant.planckUnits=6.62607015e-34 kg m^2 s^-1)"},
	} {
		got := fmt.Sprintf(test.format, Planck)
		if got != test.want {
			t.Errorf("Format %q: got: %q want: %q", test.format, got, test.want)
		}
	}
}
