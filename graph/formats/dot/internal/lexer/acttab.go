// Code generated by gocc; DO NOT EDIT.

// This file is dual licensed under CC0 and The Gonum License.
//
// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright ©2017 Robin Eklind.
// This file is made available under a Creative Commons CC0 1.0
// Universal Public Domain Dedication.

package lexer

import (
	"fmt"

	"github.com/jingcheng-WU/gonum/graph/formats/dot/internal/token"
)

type ActionTable [NumStates]ActionRow

type ActionRow struct {
	Accept token.Type
	Ignore string
}

func (a ActionRow) String() string {
	return fmt.Sprintf("Accept=%d, Ignore=%s", a.Accept, a.Ignore)
}

var ActTab = ActionTable{
	ActionRow{ // S0
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S1
		Accept: -1,
		Ignore: "!whitespace",
	},
	ActionRow{ // S2
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S3
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S4
		Accept: 15,
		Ignore: "",
	},
	ActionRow{ // S5
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S6
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S7
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S8
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S9
		Accept: 18,
		Ignore: "",
	},
	ActionRow{ // S10
		Accept: 8,
		Ignore: "",
	},
	ActionRow{ // S11
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S12
		Accept: 16,
		Ignore: "",
	},
	ActionRow{ // S13
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S14
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S15
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S16
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S17
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S18
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S19
		Accept: 13,
		Ignore: "",
	},
	ActionRow{ // S20
		Accept: 14,
		Ignore: "",
	},
	ActionRow{ // S21
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S22
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S23
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S24
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S25
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S26
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S27
		Accept: 2,
		Ignore: "",
	},
	ActionRow{ // S28
		Accept: 3,
		Ignore: "",
	},
	ActionRow{ // S29
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S30
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S31
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S32
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S33
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S34
		Accept: -1,
		Ignore: "!comment",
	},
	ActionRow{ // S35
		Accept: 9,
		Ignore: "",
	},
	ActionRow{ // S36
		Accept: 10,
		Ignore: "",
	},
	ActionRow{ // S37
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S38
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S39
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S40
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S41
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S42
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S43
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S44
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S45
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S46
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S47
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S48
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S49
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S50
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S51
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S52
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S53
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S54
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S55
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S56
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S57
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S58
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S59
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S60
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S61
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S62
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S63
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S64
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S65
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S66
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S67
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S68
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S69
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S70
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S71
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S72
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S73
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S74
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S75
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S76
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S77
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S78
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S79
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S80
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S81
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S82
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S83
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S84
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S85
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S86
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S87
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S88
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S89
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S90
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S91
		Accept: -1,
		Ignore: "!comment",
	},
	ActionRow{ // S92
		Accept: 0,
		Ignore: "",
	},
	ActionRow{ // S93
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S94
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S95
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S96
		Accept: 12,
		Ignore: "",
	},
	ActionRow{ // S97
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S98
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S99
		Accept: 11,
		Ignore: "",
	},
	ActionRow{ // S100
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S101
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S102
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S103
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S104
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S105
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S106
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S107
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S108
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S109
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S110
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S111
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S112
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S113
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S114
		Accept: 6,
		Ignore: "",
	},
	ActionRow{ // S115
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S116
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S117
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S118
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S119
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S120
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S121
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S122
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S123
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S124
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S125
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S126
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S127
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S128
		Accept: 5,
		Ignore: "",
	},
	ActionRow{ // S129
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S130
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S131
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S132
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S133
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S134
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S135
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S136
		Accept: 7,
		Ignore: "",
	},
	ActionRow{ // S137
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S138
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S139
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S140
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S141
		Accept: 19,
		Ignore: "",
	},
	ActionRow{ // S142
		Accept: 17,
		Ignore: "",
	},
}
