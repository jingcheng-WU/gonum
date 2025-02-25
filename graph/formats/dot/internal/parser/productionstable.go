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

package parser

import (
	"github.com/jingcheng-WU/gonum/graph/formats/dot/ast"
	"github.com/jingcheng-WU/gonum/graph/formats/dot/internal/astx"
)

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index      int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
	Attrib interface {
	}
)

var productionsTable = ProdTab{
	ProdTabEntry{
		String: `S' : File	<<  >>`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `File : Graph	<< astx.NewFile(X[0]) >>`,
		Id:         "File",
		NTType:     1,
		Index:      1,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewFile(X[0])
		},
	},
	ProdTabEntry{
		String: `File : File Graph	<< astx.AppendGraph(X[0], X[1]) >>`,
		Id:         "File",
		NTType:     1,
		Index:      2,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.AppendGraph(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `Graph : OptStrict DirectedGraph OptID "{" OptStmtList "}"	<< astx.NewGraph(X[0], X[1], X[2], X[4]) >>`,
		Id:         "Graph",
		NTType:     2,
		Index:      3,
		NumSymbols: 6,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewGraph(X[0], X[1], X[2], X[4])
		},
	},
	ProdTabEntry{
		String: `OptStrict : empty	<< false, nil >>`,
		Id:         "OptStrict",
		NTType:     3,
		Index:      4,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return false, nil
		},
	},
	ProdTabEntry{
		String: `OptStrict : strict	<< true, nil >>`,
		Id:         "OptStrict",
		NTType:     3,
		Index:      5,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return true, nil
		},
	},
	ProdTabEntry{
		String: `DirectedGraph : graphx	<< false, nil >>`,
		Id:         "DirectedGraph",
		NTType:     4,
		Index:      6,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return false, nil
		},
	},
	ProdTabEntry{
		String: `DirectedGraph : digraph	<< true, nil >>`,
		Id:         "DirectedGraph",
		NTType:     4,
		Index:      7,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return true, nil
		},
	},
	ProdTabEntry{
		String: `StmtList : Stmt OptSemi	<< astx.NewStmtList(X[0]) >>`,
		Id:         "StmtList",
		NTType:     5,
		Index:      8,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewStmtList(X[0])
		},
	},
	ProdTabEntry{
		String: `StmtList : StmtList Stmt OptSemi	<< astx.AppendStmt(X[0], X[1]) >>`,
		Id:         "StmtList",
		NTType:     5,
		Index:      9,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.AppendStmt(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `OptStmtList : empty	<<  >>`,
		Id:         "OptStmtList",
		NTType:     6,
		Index:      10,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `OptStmtList : StmtList	<<  >>`,
		Id:         "OptStmtList",
		NTType:     6,
		Index:      11,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Stmt : NodeStmt	<<  >>`,
		Id:         "Stmt",
		NTType:     7,
		Index:      12,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Stmt : EdgeStmt	<<  >>`,
		Id:         "Stmt",
		NTType:     7,
		Index:      13,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Stmt : AttrStmt	<<  >>`,
		Id:         "Stmt",
		NTType:     7,
		Index:      14,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Stmt : Attr	<<  >>`,
		Id:         "Stmt",
		NTType:     7,
		Index:      15,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Stmt : Subgraph	<<  >>`,
		Id:         "Stmt",
		NTType:     7,
		Index:      16,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `OptSemi : empty	<<  >>`,
		Id:         "OptSemi",
		NTType:     8,
		Index:      17,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `OptSemi : ";"	<<  >>`,
		Id:         "OptSemi",
		NTType:     8,
		Index:      18,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `NodeStmt : Node OptAttrList	<< astx.NewNodeStmt(X[0], X[1]) >>`,
		Id:         "NodeStmt",
		NTType:     9,
		Index:      19,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewNodeStmt(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `EdgeStmt : Vertex Edge OptAttrList	<< astx.NewEdgeStmt(X[0], X[1], X[2]) >>`,
		Id:         "EdgeStmt",
		NTType:     10,
		Index:      20,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewEdgeStmt(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Edge : DirectedEdge Vertex OptEdge	<< astx.NewEdge(X[0], X[1], X[2]) >>`,
		Id:         "Edge",
		NTType:     11,
		Index:      21,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewEdge(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `DirectedEdge : "--"	<< false, nil >>`,
		Id:         "DirectedEdge",
		NTType:     12,
		Index:      22,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return false, nil
		},
	},
	ProdTabEntry{
		String: `DirectedEdge : "->"	<< true, nil >>`,
		Id:         "DirectedEdge",
		NTType:     12,
		Index:      23,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return true, nil
		},
	},
	ProdTabEntry{
		String: `OptEdge : empty	<<  >>`,
		Id:         "OptEdge",
		NTType:     13,
		Index:      24,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `OptEdge : Edge	<<  >>`,
		Id:         "OptEdge",
		NTType:     13,
		Index:      25,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `AttrStmt : Component AttrList	<< astx.NewAttrStmt(X[0], X[1]) >>`,
		Id:         "AttrStmt",
		NTType:     14,
		Index:      26,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewAttrStmt(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `Component : graphx	<< ast.GraphKind, nil >>`,
		Id:         "Component",
		NTType:     15,
		Index:      27,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.GraphKind, nil
		},
	},
	ProdTabEntry{
		String: `Component : node	<< ast.NodeKind, nil >>`,
		Id:         "Component",
		NTType:     15,
		Index:      28,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NodeKind, nil
		},
	},
	ProdTabEntry{
		String: `Component : edge	<< ast.EdgeKind, nil >>`,
		Id:         "Component",
		NTType:     15,
		Index:      29,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.EdgeKind, nil
		},
	},
	ProdTabEntry{
		String: `AttrList : "[" OptAList "]"	<< X[1], nil >>`,
		Id:         "AttrList",
		NTType:     16,
		Index:      30,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[1], nil
		},
	},
	ProdTabEntry{
		String: `AttrList : AttrList "[" OptAList "]"	<< astx.AppendAttrList(X[0], X[2]) >>`,
		Id:         "AttrList",
		NTType:     16,
		Index:      31,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.AppendAttrList(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `OptAttrList : empty	<<  >>`,
		Id:         "OptAttrList",
		NTType:     17,
		Index:      32,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `OptAttrList : AttrList	<<  >>`,
		Id:         "OptAttrList",
		NTType:     17,
		Index:      33,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `AList : Attr OptSep	<< astx.NewAttrList(X[0]) >>`,
		Id:         "AList",
		NTType:     18,
		Index:      34,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewAttrList(X[0])
		},
	},
	ProdTabEntry{
		String: `AList : AList Attr OptSep	<< astx.AppendAttr(X[0], X[1]) >>`,
		Id:         "AList",
		NTType:     18,
		Index:      35,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.AppendAttr(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `OptAList : empty	<<  >>`,
		Id:         "OptAList",
		NTType:     19,
		Index:      36,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `OptAList : AList	<<  >>`,
		Id:         "OptAList",
		NTType:     19,
		Index:      37,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `OptSep : empty	<<  >>`,
		Id:         "OptSep",
		NTType:     20,
		Index:      38,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `OptSep : ";"	<<  >>`,
		Id:         "OptSep",
		NTType:     20,
		Index:      39,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `OptSep : ","	<<  >>`,
		Id:         "OptSep",
		NTType:     20,
		Index:      40,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Attr : ID "=" ID	<< astx.NewAttr(X[0], X[2]) >>`,
		Id:         "Attr",
		NTType:     21,
		Index:      41,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewAttr(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Subgraph : OptSubgraphID "{" OptStmtList "}"	<< astx.NewSubgraph(X[0], X[2]) >>`,
		Id:         "Subgraph",
		NTType:     22,
		Index:      42,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewSubgraph(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `OptSubgraphID : empty	<<  >>`,
		Id:         "OptSubgraphID",
		NTType:     23,
		Index:      43,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `OptSubgraphID : subgraph OptID	<< X[1], nil >>`,
		Id:         "OptSubgraphID",
		NTType:     23,
		Index:      44,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[1], nil
		},
	},
	ProdTabEntry{
		String: `Vertex : Node	<<  >>`,
		Id:         "Vertex",
		NTType:     24,
		Index:      45,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Vertex : Subgraph	<<  >>`,
		Id:         "Vertex",
		NTType:     24,
		Index:      46,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Node : ID OptPort	<< astx.NewNode(X[0], X[1]) >>`,
		Id:         "Node",
		NTType:     25,
		Index:      47,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewNode(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `Port : ":" ID	<< astx.NewPort(X[1], nil) >>`,
		Id:         "Port",
		NTType:     26,
		Index:      48,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewPort(X[1], nil)
		},
	},
	ProdTabEntry{
		String: `Port : ":" ID ":" ID	<< astx.NewPort(X[1], X[3]) >>`,
		Id:         "Port",
		NTType:     26,
		Index:      49,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewPort(X[1], X[3])
		},
	},
	ProdTabEntry{
		String: `OptPort : empty	<<  >>`,
		Id:         "OptPort",
		NTType:     27,
		Index:      50,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `OptPort : Port	<<  >>`,
		Id:         "OptPort",
		NTType:     27,
		Index:      51,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ID : id	<< astx.NewID(X[0]) >>`,
		Id:         "ID",
		NTType:     28,
		Index:      52,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewID(X[0])
		},
	},
	ProdTabEntry{
		String: `OptID : empty	<< "", nil >>`,
		Id:         "OptID",
		NTType:     29,
		Index:      53,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return "", nil
		},
	},
	ProdTabEntry{
		String: `OptID : ID	<<  >>`,
		Id:         "OptID",
		NTType:     29,
		Index:      54,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
}
