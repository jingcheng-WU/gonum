// Code generated by "go generate github.com/jingcheng-WU/gonum/unit/constant”; DO NOT EDIT.

// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constant

import (
	"fmt"
	"testing"
)

func TestGravitationalFormat(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		format string
		want   string
	}{
		{"%v", "6.6743e-11 m^3 kg^-1 s^-2"},
		{"%.1v", "7e-11 m^3 kg^-1 s^-2"},
		{"%50.1v", "                              7e-11 m^3 kg^-1 s^-2"},
		{"%50v", "                         6.6743e-11 m^3 kg^-1 s^-2"},
		{"%1v", "6.6743e-11 m^3 kg^-1 s^-2"},
		{"%#v", "constant.gravitationalUnits(6.6743e-11)"},
		{"%s", "%!s(constant.gravitationalUnits=6.6743e-11 m^3 kg^-1 s^-2)"},
	} {
		got := fmt.Sprintf(test.format, Gravitational)
		if got != test.want {
			t.Errorf("Format %q: got: %q want: %q", test.format, got, test.want)
		}
	}
}
