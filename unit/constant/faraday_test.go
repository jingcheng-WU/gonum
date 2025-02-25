// Code generated by "go generate github.com/jingcheng-WU/gonum/unit/constant”; DO NOT EDIT.

// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package constant

import (
	"fmt"
	"testing"
)

func TestFaradayFormat(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		format string
		want   string
	}{
		{"%v", "96485.33212 A s mol^-1"},
		{"%.1v", "1e+05 A s mol^-1"},
		{"%50.1v", "                                  1e+05 A s mol^-1"},
		{"%50v", "                            96485.33212 A s mol^-1"},
		{"%1v", "96485.33212 A s mol^-1"},
		{"%#v", "constant.faradayUnits(96485.33212)"},
		{"%s", "%!s(constant.faradayUnits=96485.33212 A s mol^-1)"},
	} {
		got := fmt.Sprintf(test.format, Faraday)
		if got != test.want {
			t.Errorf("Format %q: got: %q want: %q", test.format, got, test.want)
		}
	}
}
